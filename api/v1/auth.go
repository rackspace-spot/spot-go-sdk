package rxtspot

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"k8s.io/klog/v2"
)

// JWTClaims represents the standard JWT claims we care about
type JWTClaims struct {
	Exp int64 `json:"exp"` // Expiration time (Unix timestamp)
}

// isTokenExpired checks if a JWT token is expired
func isTokenExpired(tokenString string) bool {
	if tokenString == "" {
		return true
	}

	// Split the token into parts (header.payload.signature)
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return true // Invalid token format, consider it expired
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return true // Invalid base64, consider it expired
	}

	// Parse the claims
	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return true // Invalid claims, consider it expired
	}

	// Check if token is expired with a 60s early-expiry leeway for clock skew.
	// Treat the token as expired if current time is after (exp - 60s).
	expTime := time.Unix(claims.Exp, 0)
	return time.Now().After(expTime.Add(-60 * time.Second))
}

func GetClientID() string {
	clientID := os.Getenv("RXTSPOT_CLIENT_ID")
	if clientID == "" {
		clientID = "mwG3lUMV8KyeMqHe4fJ5Bb3nM1vBvRNa"
	}
	return clientID
}

// Authenticate authenticates with the Rackspace Spot API using the provided credentials.
func (c *RackspaceSpotClient) Authenticate(ctx context.Context) (string, error) {
	// If we have a token and it's not expired, use it
	if c.Token != "" && !isTokenExpired(c.Token) {
		return c.Token, nil
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", GetClientID())
	form.Set("refresh_token", c.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.OAuthURL+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fmt.Printf("req: %+v\n", req)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	fmt.Printf("resp: %+v\n", resp)
	defer resp.Body.Close()
	
	// Read the response body for error details
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		klog.Errorf("Authentication failed with status %d: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("authentication failed: %s - %s", resp.Status, string(body))
	}
	
	// For successful responses, we need to use the body for token parsing
	// So we'll create a new reader for the body since we've already read it
	resp.Body = io.NopCloser(bytes.NewReader(body))

	var tokenResp struct {
		IDToken string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	fmt.Printf("tokenResp: %+v\n", tokenResp)
	if tokenResp.IDToken == "" {
		return "", errors.New("no id_token in authentication response")
	}

	c.Token = tokenResp.IDToken
	return tokenResp.IDToken, nil
}

// authHeader returns the Authorization header for the current access token.
func (c *RackspaceSpotClient) authHeader() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.Token),
		"Content-Type":  "application/json",
	}
}
