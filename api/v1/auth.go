package rxtspot

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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

	// Check if token is expired (with 60 second leeway for clock skew)
	now := time.Now().Unix()
	return claims.Exp < now-60
}

func GetClientID() string {
	return os.Getenv("RXTSPOT_CLIENT_ID")
}

// Authenticate authenticates with the Rackspace Spot API using the provided credentials.
func (c *RackspaceSpotClient) Authenticate(ctx context.Context) (string, error) {
	// If we have a token and it's not expired, use it
	if c.Token != "" && !isTokenExpired(c.Token) {
		return c.Token, nil
	}
	// No valid token, get a new one
	fmt.Printf("RackspaceSpotClient: %v \n", c)
	fmt.Printf("c.refreshToken: %v \n", c.RefreshToken)
	fmt.Printf("c.OAuthURL: %v \n", c.OAuthURL)
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", GetClientID())
	form.Set("refresh_token", c.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.OAuthURL+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	fmt.Println("Request: ", req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	fmt.Printf("Response: %v", resp)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed: %s", resp.Status)
	}

	var tokenResp struct {
		IDToken string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
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
