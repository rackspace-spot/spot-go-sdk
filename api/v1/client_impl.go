package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rackspace-spot/spot-go-sdk/api/v1/cloudspace"
	cloudspacesvc "github.com/rackspace-spot/spot-go-sdk/internal/api/v1/cloudspace"
	ondemand "github.com/rackspace-spot/spot-go-sdk/internal/api/v1/nodepool/ondemand"
	svcspot "github.com/rackspace-spot/spot-go-sdk/internal/api/v1/nodepool/spot"
	region "github.com/rackspace-spot/spot-go-sdk/internal/api/v1/region"
	serverclass "github.com/rackspace-spot/spot-go-sdk/internal/api/v1/serverclass"
	intclient "github.com/rackspace-spot/spot-go-sdk/internal/api/v1/client"
)

// ClientConfig holds configuration for creating a new Client
type ClientConfig struct {
	// BaseURL is the base URL of the Spot API
	BaseURL string

	// HTTPClient is an optional HTTP client to use for requests.
	// If nil, a default client with a 30-second timeout will be used.
	HTTPClient *http.Client

	// AuthToken is the authentication token to use for requests
	AuthToken string
}

// client is the concrete implementation of the Client interface
type client struct {
	httpClient *intclient.HTTPClient
}

// NewClient creates a new Spot API client with the given configuration
func NewClient(cfg ClientConfig) (Client, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	httpClient := intclient.New(cfg.BaseURL, cfg.HTTPClient)
	if cfg.AuthToken != "" {
		httpClient.SetHeader("Authorization", "Bearer "+cfg.AuthToken)
	}

	return &client{
		httpClient: httpClient,
	}, nil
}

// Cloudspaces implements the Client interface
func (c *client) Cloudspaces() cloudspace.Service {
	return cloudspacesvc.New(c.httpClient)
}

// SpotNodePools implements the Client interface
func (c *client) SpotNodePools() spot.Service {
	return svcspot.New(c.httpClient)
}

// OnDemandNodePools implements the Client interface
func (c *client) OnDemandNodePools() ondemand.Service {
	return ondemand.New(c.httpClient)
}

// ServerClasses implements the Client interface
func (c *client) ServerClasses() serverclass.Service {
	return serverclass.New(c.httpClient)
}

// Regions implements the Client interface
func (c *client) Regions() region.Service {
	return region.New(c.httpClient)
}

// Close releases any resources held by the client
func (c *client) Close() error {
	// Currently nothing to clean up, but keeping the method for future use
	return nil
}
