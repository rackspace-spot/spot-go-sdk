package rxtspot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"k8s.io/klog/v2"
)

// Config holds the configuration for the Rackspace Spot API client
type Config struct {
	BaseURL      string
	OAuthURL     string
	HTTPClient   *http.Client
	AccessToken  string
	RefreshToken string
}

type RackspaceSpotClient struct {
	BaseURL      string
	OAuthURL     string
	HTTPClient   *http.Client
	Token        string
	RefreshToken string
}

type HTTPStatusError struct {
	StatusCode int
	Body       string
}

// NewSpotClient creates a new Rackspace Spot API client with secure defaults
func NewSpotClient(cfg *Config) (*RackspaceSpotClient, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	return &RackspaceSpotClient{
		BaseURL:      cfg.BaseURL,
		OAuthURL:     cfg.OAuthURL,
		HTTPClient:   cfg.HTTPClient,
		Token:        cfg.AccessToken,
		RefreshToken: cfg.RefreshToken,
	}, nil
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

// handleAPIError processes HTTP errors and returns user-friendly error messages
func (c *RackspaceSpotClient) handleAPIError(err error, resourceType, resourceName string, operation string) error {
	if err == nil {
		return nil
	}

	httpErr, ok := err.(*HTTPStatusError)
	if !ok {
		return fmt.Errorf("failed to %s %s: %w", operation, resourceType, err)
	}

	switch httpErr.StatusCode {
	case http.StatusForbidden, http.StatusUnauthorized:
		// Try to extract the detailed error message from the response body
		var apiErr struct {
			Message string `json:"message"`
		}
		if httpErr.Body != "" {
			if json.Unmarshal([]byte(httpErr.Body), &apiErr) == nil && apiErr.Message != "" {
				return fmt.Errorf("access denied: %s", apiErr.Message)
			}
		}
		return fmt.Errorf("access denied: you do not have permission to %s the %s '%s'", operation, resourceType, resourceName)

	case http.StatusNotFound:
		return fmt.Errorf("%s '%s' not found", resourceType, resourceName)

	case http.StatusBadRequest, http.StatusConflict:
		if httpErr.Body != "" {
			return fmt.Errorf("invalid request: %s", httpErr.Body)
		}
		return fmt.Errorf("invalid request: failed to %s %s", operation, resourceType)

	default:
		if httpErr.Body != "" {
			return fmt.Errorf("API error (HTTP %d): %s", httpErr.StatusCode, httpErr.Body)
		}
		return fmt.Errorf("API error (HTTP %d): failed to %s %s", httpErr.StatusCode, operation, resourceType)
	}
}

// validateURL performs basic validation of a URL
func validateURL(urlStr string) error {
	// Check for empty URL
	if urlStr == "" {
		return errors.New("URL cannot be empty")
	}

	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Check scheme (http or https)
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("invalid URL scheme: %s, must be http or https", scheme)
	}

	// Check hostname
	if parsedURL.Host == "" {
		return errors.New("URL must contain a host")
	}

	// Validate hostname
	hostname := parsedURL.Hostname()
	if !isValidHostname(hostname) {
		return fmt.Errorf("invalid hostname: %s", hostname)
	}

	// Check for invalid characters in path
	if strings.ContainsAny(parsedURL.Path, "\x00\r\n") {
		return errors.New("URL path contains invalid characters")
	}

	// Validate port if present
	if port := parsedURL.Port(); port != "" {
		portNum, err := strconv.Atoi(port)
		if err != nil || portNum < 1 || portNum > 65535 {
			return fmt.Errorf("invalid port number")
		}
	}

	return nil
}

// doRequest performs an HTTP request with the given method, URL, body, and headers.
// It handles authentication, rate limiting, and error handling.
func (c *RackspaceSpotClient) doRequest(ctx context.Context, method, url string, body []byte, headers map[string]string, out interface{}) error {
	if err := validateURL(url); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		klog.Errorf("doRequest: failed to create request: %v", err)
		return err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/merge-patch+json")
	}

	// ----- Request logging -----
	klog.V(1).Infof("[%s] %s", method, url)
	if len(headers) > 0 {
		klog.V(2).Infof("Request headers: %+v", headers)
	}
	if len(body) > 0 {
		klog.V(3).Infof("Request body: %s", string(body))
	}

	// ----- Perform HTTP request -----
	start := time.Now()
	resp, err := c.HTTPClient.Do(req)
	duration := time.Since(start)
	if err != nil {
		klog.Errorf("HTTP request failed after %v: %v", duration, err)
		return err
	}
	defer resp.Body.Close()

	klog.V(2).Infof("Response status: %d %s (duration: %v)", resp.StatusCode, http.StatusText(resp.StatusCode), duration)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			klog.Errorf("Failed to read response body: %v", err)
			return fmt.Errorf("read response body: %w", err)
		}
		return &HTTPStatusError{StatusCode: resp.StatusCode, Body: string(b)}
	}

	// ----- Response logging -----
	if klog.V(3).Enabled() {
		klog.V(3).Infof("Response headers: %+v", resp.Header)
	}
	if klog.V(4).Enabled() {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			klog.Errorf("Failed to read response body: %v", err)
			return fmt.Errorf("read response body: %w", err)
		}
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
		klog.V(4).Infof("Response body: %s", string(respBody))
	}

	// ----- Decode if out != nil -----
	if out != nil {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(out); err != nil && err != io.EOF {
			klog.Errorf("Failed to decode JSON: %v", err)
			return fmt.Errorf("decode json: %w", err)
		}
		klog.V(4).Infof("Decoded object: %+v", out)
	}

	return nil
}

// isValidHostname checks if a string is a valid hostname (supports both domains and IPs)
func isValidHostname(hostname string) bool {
	// First try to parse as IP address
	if ip := net.ParseIP(hostname); ip != nil {
		return true
	}

	// Then check as domain name (simplified check)
	// This is a basic check - for production you might want more comprehensive validation
	if len(hostname) > 253 || len(hostname) == 0 {
		return false
	}

	// Check each label in the hostname
	for _, label := range strings.Split(hostname, ".") {
		if len(label) > 63 || len(label) == 0 {
			return false
		}

		// Labels must start and end with alphanumeric characters
		if len(label) > 0 && !isAlphanumeric(rune(label[0])) || !isAlphanumeric(rune(label[len(label)-1])) {
			return false
		}

		// Labels can contain alphanumeric characters and hyphens
		for _, r := range label {
			if !isAlphanumeric(r) && r != '-' {
				return false
			}
		}
	}

	return true
}

// isAlphanumeric checks if a rune is alphanumeric
func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

// APIError represents an error response from the API
type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// Helper matchers for consumers (like your CLI)
func IsNotFound(err error) bool {
	var e *HTTPStatusError
	if errors.As(err, &e) {
		return e.StatusCode == http.StatusNotFound
	}
	return false
}

func IsForbidden(err error) bool {
	var e *HTTPStatusError
	if errors.As(err, &e) {
		return e.StatusCode == http.StatusForbidden
	}
	return false
}

func IsConflict(err error) bool {
	var e *HTTPStatusError
	if errors.As(err, &e) {
		return e.StatusCode == http.StatusConflict
	}
	return false
}
