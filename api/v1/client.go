package rxtspot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"errors"

	"k8s.io/klog/v2"
)

// Config holds the configuration for the Rackspace Spot client
type Config struct {
	BaseURL      string
	OAuthURL     string
	RefreshToken string
	HTTPClient   *http.Client
	AccessToken  string
}

// RackspaceSpotClient is the main client for interacting with the Rackspace Spot API
type RackspaceSpotClient struct {
	BaseURL      string
	OAuthURL     string
	HTTPClient   *http.Client
	Token        string
	RefreshToken string
}

// NewSpotClient creates a new RackspaceSpotClient with the given configuration
func NewSpotClient(cfg *Config) (*RackspaceSpotClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	if cfg.OAuthURL == "" {
		cfg.OAuthURL = "https://auth.rackspace.com"
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &RackspaceSpotClient{
		BaseURL:      cfg.BaseURL,
		OAuthURL:     cfg.OAuthURL,
		Token:        cfg.AccessToken,
		HTTPClient:   cfg.HTTPClient,
		RefreshToken: cfg.RefreshToken,
	}, nil
}

// client implements the Client interface
type client struct {
	client *RackspaceSpotClient
}

// Verify that client implements the Client interface
var _ Client = (*client)(nil)

// NewClient creates a new client with the given configuration
func NewClient(cfg *Config) (Client, error) {
	rsc, err := NewSpotClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &client{
		client: rsc,
	}, nil
}

// Authenticate implements the Client interface
func (c *client) Authenticate(ctx context.Context) (string, error) {
	return c.client.Authenticate(ctx)
}

// Organizations implements the Client interface
func (c *client) Organizations() OrganizationManager {
	return &organizationManager{client: c.client}
}

// Cloudspaces implements the Client interface
func (c *client) Cloudspaces(org string) CloudspaceManager {
	return &cloudspaceManager{
		client: c.client,
		org:    org,
	}
}

// NodePools returns a NodePoolAPI for managing node pools in the specified organization and cloudspace
func (c *client) SpotNodePools(org, cloudspace string) SpotNodePoolManager {
	return &spotNodePoolManager{
		client:     c.client,
		org:        org,
		cloudspace: cloudspace,
	}
}

func (c *client) OnDemandNodePools(org, cloudspace string) OnDemandNodePoolManager {
	return &ondemandNodePoolManager{
		client:     c.client,
		org:        org,
		cloudspace: cloudspace,
	}
}

// Regions implements the Client interface
func (c *client) Regions() RegionsManager {
	return &regionsManager{client: c.client}
}

// ServerClasses implements the Client interface
func (c *client) ServerClasses(region string) ServerClassesManager {
	return &serverClassesManager{
		client: c.client,
		region: region,
	}
}

// Pricing implements the Client interface
func (c *client) Pricing() PricingManager {
	return &pricingManager{client: c.client}
}

// ListNodePools implements the NodePoolAPI interface
func (m *spotNodePoolManager) ListNodePools(ctx context.Context, org, cloudspace string, poolType *NodePoolType) ([]*SpotNodePool, error) {
	var pools []*SpotNodePool

	// List spot node pools if requested type is spot or not specified
	if poolType == nil || *poolType == NodePoolTypeSpot {
		spotPools, err := m.client.ListSpotNodePools(ctx, org, cloudspace)
		if err != nil {
			return nil, err
		}

		for _, p := range spotPools {
			if p != nil {
				poolCopy := *p // Create a copy to avoid pointer sharing
				pools = append(pools, &poolCopy)
			}
		}

		// If specific type was requested, return early
		if poolType != nil && *poolType == NodePoolTypeSpot {
			return pools, nil
		}
	}

	// List on-demand node pools if requested type is on-demand or not specified
	if poolType == nil || *poolType == NodePoolTypeOnDemand {
		onDemandPools, err := m.client.ListOnDemandNodePools(ctx, org, cloudspace)
		if err != nil {
			return nil, err
		}

		for _, p := range onDemandPools {
			if p != nil {
				poolCopy := *p // Create a copy to avoid pointer sharing
				pools = append(pools, &onDemandNodePoolWrapper{
					OnDemandNodePool: &poolCopy,
					org:              org,
				})
			}
		}
	}

	return pools, nil
}

// GetNodePool implements the NodePoolAPI interface
func (m *nodePoolManager) GetNodePool(ctx context.Context, org, name string, poolType *NodePoolType) (NodePool, error) {
	if poolType == nil {
		return nil, fmt.Errorf("node pool type must be specified (spot or on-demand)")
	}

	switch *poolType {
	case NodePoolTypeSpot:
		pool, err := m.client.GetSpotNodePool(ctx, org, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get spot node pool %s: %w", name, err)
		}
		return &spotNodePoolWrapper{
			SpotNodePool: pool,
			org:          org,
		}, nil

	case NodePoolTypeOnDemand:
		pool, err := m.client.GetOnDemandNodePool(ctx, org, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get on-demand node pool %s: %w", name, err)
		}
		return &onDemandNodePoolWrapper{
			OnDemandNodePool: pool,
			org:              org,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported node pool type: %s", *poolType)
	}
}

// CreateNodePool implements the NodePoolAPI interface
func (m *nodePoolManager) CreateNodePool(ctx context.Context, org string, pool NodePool) error {
	if pool == nil {
		return fmt.Errorf("node pool cannot be nil")
	}

	switch p := pool.(type) {
	// Handle wrapped types
	case *spotNodePoolWrapper:
		if p.SpotNodePool == nil {
			return fmt.Errorf("spot node pool is nil")
		}
		return m.client.CreateSpotNodePool(ctx, org, *p.SpotNodePool)

	case *onDemandNodePoolWrapper:
		if p.OnDemandNodePool == nil {
			return fmt.Errorf("on-demand node pool is nil")
		}
		return m.client.CreateOnDemandNodePool(ctx, org, *p.OnDemandNodePool)

	// Handle direct types (for backward compatibility with existing code)
	case *SpotNodePool:
		return m.client.CreateSpotNodePool(ctx, org, *p)

	case *OnDemandNodePool:
		return m.client.CreateOnDemandNodePool(ctx, org, *p)

	default:
		return fmt.Errorf("unsupported node pool type: %T - must use NewSpotNodePoolWrapper or NewOnDemandNodePoolWrapper", pool)
	}
}

// UpdateNodePool implements the NodePoolAPI interface
func (m *nodePoolManager) UpdateNodePool(ctx context.Context, org string, pool NodePool) error {
	if pool == nil {
		return fmt.Errorf("node pool cannot be nil")
	}

	switch p := pool.(type) {
	// Handle wrapped types
	case *spotNodePoolWrapper:
		if p.SpotNodePool == nil {
			return fmt.Errorf("spot node pool is nil")
		}
		return m.client.UpdateSpotNodePool(ctx, org, *p.SpotNodePool)

	case *onDemandNodePoolWrapper:
		if p.OnDemandNodePool == nil {
			return fmt.Errorf("on-demand node pool is nil")
		}
		return m.client.UpdateOnDemandNodePool(ctx, org, *p.OnDemandNodePool)

	// Handle direct types (for backward compatibility with existing code)
	case *SpotNodePool:
		return m.client.UpdateSpotNodePool(ctx, org, *p)

	case *OnDemandNodePool:
		return m.client.UpdateOnDemandNodePool(ctx, org, *p)

	default:
		return fmt.Errorf("unsupported node pool type: %T - must use NewSpotNodePoolWrapper or NewOnDemandNodePoolWrapper", pool)
	}
}

// DeleteNodePool implements the NodePoolAPI interface
func (m *nodePoolManager) DeleteNodePool(ctx context.Context, org, name string) error {
	// Try to delete as spot node pool first
	err := m.client.DeleteSpotNodePool(ctx, org, name)
	if err == nil {
		return nil
	}

	// If not found, try to delete as on-demand node pool
	err = m.client.DeleteOnDemandNodePool(ctx, org, name)
	if err != nil {
		return fmt.Errorf("failed to delete node pool %s: %w", name, err)
	}

	return nil
}

// getAuthHeader returns the authorization header with the current token
func (c *RackspaceSpotClient) getAuthHeader() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + c.Token,
	}
}

// doRequest performs an HTTP request and decodes the JSON response
/*
func (c *RackspaceSpotClient) doRequest(ctx context.Context, method, url string, body []byte, headers map[string]string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Set content type if not set
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}*/

// isAlphanumeric checks if a rune is alphanumeric
func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
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

// RequestOptions holds optional parameters for the HTTP request
type RequestOptions struct {
	// Headers to be included in the request
	Headers map[string]string
	// ExpectedStatus is the expected HTTP status code (default: 200-299)
	ExpectedStatus int
	// RetryPolicy defines the retry strategy for failed requests
	RetryPolicy RetryPolicy
}

// RetryPolicy defines the retry strategy for failed requests
type RetryPolicy struct {
	// MaxRetries is the maximum number of retries (default: 3)
	MaxRetries int
	// RetryableStatusCodes are the HTTP status codes that should be retried
	RetryableStatusCodes []int
	// Backoff is the initial backoff duration (default: 100ms)
	Backoff time.Duration
}

// DefaultRetryPolicy returns a sensible default retry policy
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxRetries: 3,
		RetryableStatusCodes: []int{
			http.StatusRequestTimeout,      // 408
			http.StatusTooManyRequests,     // 429
			http.StatusInternalServerError, // 500
			http.StatusBadGateway,          // 502
			http.StatusServiceUnavailable,  // 503
			http.StatusGatewayTimeout,      // 504
		},
		Backoff: 100 * time.Millisecond,
	}
}

// doRequest performs an HTTP request with the given method, URL, body, and options.
// It handles retries, rate limiting, and error handling.
//
// For backward compatibility, it also supports the old signature:
//
//	doRequest(ctx, method, url string, body []byte, headers map[string]string, out interface{}) error
//
// The new signature is:
//
//	doRequest(ctx, method, url string, body []byte, opts ...func(*RequestOptions)) ([]byte, error)
func (c *RackspaceSpotClient) doRequest(
	ctx context.Context,
	method string,
	urlStr string,
	body []byte,
	args ...interface{},
) ([]byte, error) {
	// Handle backward compatibility with old signature
	var out interface{}
	reqOpts := &RequestOptions{
		Headers:     make(map[string]string),
		RetryPolicy: DefaultRetryPolicy(),
	}

	// Check if using old signature (headers map, out interface{})
	if len(args) > 0 {
		switch v := args[0].(type) {
		case map[string]string:
			// Old signature: headers map[string]string
			for k, val := range v {
				reqOpts.Headers[k] = val
			}
			if len(args) > 1 {
				// out interface{}
				out = args[1]
			}
		case func(*RequestOptions):
			// New signature: opts ...func(*RequestOptions)
			for _, opt := range args {
				if fn, ok := opt.(func(*RequestOptions)); ok {
					fn(reqOpts)
				}
			}
		}
	}

	// Validate URL
	if err := validateURL(urlStr); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	// fmt.Printf("method, url - %+v , %+v \n", method, urlStr)

	var lastErr error
	var resp *http.Response
	var respBody []byte

	// Prepare request outside retry loop to avoid recreating it
	req, err := http.NewRequestWithContext(ctx, method, urlStr, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set default headers
	req.Header.Set("User-Agent", "spot-go-sdk/1.0")
	req.Header.Set("Accept", "application/json")
	if method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/merge-patch+json")
	} else if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add custom headers
	for key, value := range reqOpts.Headers {
		req.Header.Set(key, value)
	}

	// Log request
	klog.V(2).Infof("[%s] %s", method, urlStr)
	if klog.V(3).Enabled() {
		headers := make(map[string]string, len(req.Header))
		for k, v := range req.Header {
			headers[k] = strings.Join(v, ", ")
		}
		klog.V(3).Infof("Request headers: %+v", headers)
	}
	if len(body) > 0 && klog.V(4).Enabled() {
		klog.V(4).Infof("Request body: %s", string(body))
	}

	// Execute request with retries
	for attempt := 0; attempt <= reqOpts.RetryPolicy.MaxRetries; attempt++ {
		if attempt > 0 {
			// Apply exponential backoff
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * reqOpts.RetryPolicy.Backoff
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
				// Continue with the next attempt
			}

			// Reset the request body for retries
			if body != nil {
				req.Body = io.NopCloser(bytes.NewReader(body))
			}
		}

		// Perform the request
		start := time.Now()
		resp, err = c.HTTPClient.Do(req)
		duration := time.Since(start)
		//fmt.Printf("resp -- %+v \n", resp.Body)
		// Handle request errors
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed after %v: %w", duration, err)
			klog.Errorf("Attempt %d/%d failed: %v", attempt+1, reqOpts.RetryPolicy.MaxRetries+1, lastErr)
			continue
		}

		// Read the response body
		respBody, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		//fmt.Printf("resp body - %+v \n", string(respBody))
		if err != nil {
			lastErr = fmt.Errorf("read response body: %w", err)
			klog.Errorf("Failed to read response body: %v", lastErr)
			continue
		}

		klog.V(2).Infof("Response status: %d %s (attempt %d/%d, duration: %v)",
			resp.StatusCode, http.StatusText(resp.StatusCode),
			attempt+1, reqOpts.RetryPolicy.MaxRetries+1, duration)

		// Check if we should retry based on status code
		if shouldRetry(resp.StatusCode, reqOpts.RetryPolicy.RetryableStatusCodes) {
			err := &HTTPStatusError{
				StatusCode: resp.StatusCode,
				Status:     resp.Status,
				Body:       string(respBody),
			}
			klog.Warningf("Retryable status code %d: %s", resp.StatusCode, string(respBody))
			lastErr = err
			continue
		}

		// Check for expected status code
		if reqOpts.ExpectedStatus > 0 && resp.StatusCode != reqOpts.ExpectedStatus {
			err := &HTTPStatusError{
				StatusCode: resp.StatusCode,
				Status:     resp.Status,
				Body:       string(respBody),
			}
			klog.Errorf("Unexpected status code %d: %s", resp.StatusCode, string(respBody))
			return nil, err
		}

		// Success case
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Handle old signature where out parameter is provided
			if out != nil && len(respBody) > 0 {
				if err := json.Unmarshal(respBody, out); err != nil {
					return nil, fmt.Errorf("unmarshal response: %w", err)
				}
			}
			return respBody, nil
		}

		// Non-retryable error
		err := &HTTPStatusError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(respBody),
		}
		klog.Errorf("Request failed with status %d: %s", resp.StatusCode, string(respBody))
		return nil, err
	}

	return nil, fmt.Errorf("max retries exceeded, last error: %w", lastErr)
}

// shouldRetry checks if a status code should be retried
func shouldRetry(statusCode int, retryableStatuses []int) bool {
	// Always retry on 5xx status codes unless it's in the excluded list
	if statusCode >= 500 && statusCode < 600 {
		return true
	}

	// Check against the list of retryable status codes
	for _, code := range retryableStatuses {
		if statusCode == code {
			return true
		}
	}

	return false
}

// HTTPStatusError represents an HTTP error with a status code
type HTTPStatusError struct {
	StatusCode int
	Status     string
	Body       string
}

// Error implements the error interface
func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("HTTP error %d: %s", e.StatusCode, e.Status)
}

// doRequestJSON is a convenience wrapper around doRequest that handles JSON serialization/deserialization
func (c *RackspaceSpotClient) doRequestJSON(
	ctx context.Context,
	method string,
	urlStr string,
	request, response interface{},
	args ...interface{},
) error {
	// Marshal request body if provided
	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
	}

	// Execute the request
	respBody, err := c.doRequest(ctx, method, urlStr, body, args...)
	if err != nil {
		return err
	}

	// Unmarshal response if a target is provided
	if response != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

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

// IsNotFound checks if the error represents a "not found" condition (HTTP 404).
// It checks if the error is an HTTPStatusError with StatusCode http.StatusNotFound.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	var httpErr *HTTPStatusError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == http.StatusNotFound
	}

	return false
}

// IsForbidden checks if the error represents a "forbidden" condition (HTTP 403).
// It checks if the error is an HTTPStatusError with StatusCode http.StatusForbidden.
func IsForbidden(err error) bool {
	if err == nil {
		return false
	}

	var httpErr *HTTPStatusError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == http.StatusForbidden
	}

	return false
}

// IsConflict checks if the error represents a "conflict" condition (HTTP 409).
// It checks if the error is an HTTPStatusError with StatusCode http.StatusConflict.
func IsConflict(err error) bool {
	if err == nil {
		return false
	}

	var httpErr *HTTPStatusError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == http.StatusConflict
	}

	return false
}

// spotNodePoolWrapper wraps SpotNodePool to implement the NodePool interface
type spotNodePoolWrapper struct {
	*SpotNodePool
	org string
}

// GetType returns the type of the node pool (spot)
func (w *spotNodePoolWrapper) GetType() NodePoolType { return NodePoolTypeSpot }

// GetName returns the name of the node pool
func (w *spotNodePoolWrapper) GetName() string {
	if w.SpotNodePool == nil {
		return ""
	}
	return w.SpotNodePool.Name
}

// GetOrg returns the organization ID of the node pool
func (w *spotNodePoolWrapper) GetOrg() string { return w.org }

// NewSpotNodePoolWrapper creates a new NodePool wrapper for a SpotNodePool
func NewSpotNodePoolWrapper(pool *SpotNodePool, org string) NodePool {
	return &spotNodePoolWrapper{
		SpotNodePool: pool,
		org:          org,
	}
}

// NewOnDemandNodePoolWrapper creates a new NodePool wrapper for an OnDemandNodePool
func NewOnDemandNodePoolWrapper(pool *OnDemandNodePool, org string) NodePool {
	return &onDemandNodePoolWrapper{
		OnDemandNodePool: pool,
		org:              org,
	}
}

// onDemandNodePoolWrapper wraps OnDemandNodePool to implement the NodePool interface
type onDemandNodePoolWrapper struct {
	*OnDemandNodePool
	org string
}

// GetType returns the type of the node pool (on-demand)
func (w *onDemandNodePoolWrapper) GetType() NodePoolType { return NodePoolTypeOnDemand }

// GetName returns the name of the node pool
func (w *onDemandNodePoolWrapper) GetName() string {
	if w.OnDemandNodePool == nil {
		return ""
	}
	return w.OnDemandNodePool.Name
}

// GetOrg returns the organization ID of the node pool
func (w *onDemandNodePoolWrapper) GetOrg() string { return w.org }
