package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient is the internal HTTP client that implements the v1.Client interface
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// New creates a new HTTP client that implements the v1.Client interface
func New(baseURL string, httpClient *http.Client) *HTTPClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &HTTPClient{
		baseURL:    baseURL,
		httpClient: httpClient,
		headers:    make(map[string]string),
	}
}

// SetHeader sets a header that will be included in all requests
func (c *HTTPClient) SetHeader(key, value string) {
	c.headers[key] = value
}

// Do sends an HTTP request and decodes the JSON response into out
func (c *HTTPClient) Do(ctx context.Context, method, path string, body interface{}, out interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check for error responses
	if resp.StatusCode >= 400 {
		return parseErrorResponse(resp)
	}

	// If no out parameter was provided, don't try to decode the response
	if out == nil {
		return nil
	}

	// Decode successful response
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

// APIError represents an error response from the API
type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("API request failed with status code: %d", e.StatusCode)
}

// parseErrorResponse parses an error response from the API
func parseErrorResponse(resp *http.Response) error {
	var apiErr APIError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return fmt.Errorf("failed to decode error response: %v", err)
	}
	apiErr.StatusCode = resp.StatusCode
	return &apiErr
}
