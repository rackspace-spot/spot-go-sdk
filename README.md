# Rackspace Spot Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/rackspace-spot/spot-go-sdk.svg)](https://pkg.go.dev/github.com/rackspace-spot/spot-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/rackspace-spot/spot-go-sdk)](https://goreportcard.com/report/github.com/rackspace-spot/spot-go-sdk)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Go client library for the Rackspace Spot API, providing a convenient way to interact with Rackspace Spot services programmatically.

## Features

- **Unified API**: Single client interface for all Rackspace Spot services
- **Type Safety**: Strongly typed API for better developer experience
- **Comprehensive Coverage**: Supports all major Rackspace Spot services
- **Concurrency Safe**: Designed for use in concurrent applications
- **Extensible**: Easy to extend with new features and services

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [Usage Examples](#usage-examples)
  - [Managing Node Pools](#managing-node-pools)
  - [Working with Cloudspaces](#working-with-cloudspaces)
  - [Managing Organizations](#managing-organizations)
- [Advanced Configuration](#advanced-configuration)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Installation

To install the Rackspace Spot Go SDK, use `go get`:

```bash
go get github.com/rackspace-spot/spot-go-sdk
```

## Quick Start

Here's a quick example of how to use the SDK to list all organizations:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rackspace-spot/spot-go-sdk/api/v1"
)

func main() {
	// Initialize client with default configuration
	client, err := v1.NewClient(&v1.Config{
		BaseURL:      "https://api.spot.io",
		OAuthURL:     "https://auth.spot.io",
		RefreshToken: os.Getenv("SPOT_REFRESH_TOKEN"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// List all organizations
	orgs, err := client.Organizations().ListOrganizations(context.Background())
	if err != nil {
		log.Fatalf("Failed to list organizations: %v", err)
	}

	// Print organization names
	for _, org := range orgs {
		fmt.Printf("Organization: %s (ID: %s)\n", org.Name, org.ID)
	}
}
```

## Authentication

The SDK uses OAuth2 for authentication. You'll need a valid Rackspace Spot refresh token to authenticate.

### Getting a Refresh Token

1. Log in to the Rackspace Spot Console
2. Navigate to Settings > API Tokens
3. Generate a new API token
4. Set it as an environment variable:

```bash
export SPOT_REFRESH_TOKEN=your_refresh_token_here
```

### Configuring the Client

When creating a new client, you can customize the configuration:

```go
config := &v1.Config{
    BaseURL:      "https://api.spot.io",      // API base URL
    OAuthURL:     "https://auth.spot.io",     // OAuth2 server URL
    RefreshToken: "your_refresh_token_here",  // Your refresh token
    HTTPClient:   &http.Client{                // Custom HTTP client (optional)
        Timeout: 30 * time.Second,
    },
}

client, err := v1.NewClient(config)
```

## Usage Examples

### Managing Node Pools

```go
// Create a new node pool manager
nodePools := client.NodePools("your-org-id", "your-cloudspace-id")

// List all node pools
pools, err := nodePools.ListNodePools(context.Background(), "your-org-id", "your-cloudspace-id", nil)
if err != nil {
    log.Fatalf("Failed to list node pools: %v", err)
}
for _, pool := range pools {
    fmt.Printf("Node Pool: %s (Type: %s)\n", pool.GetName(), pool.GetType())
}

// Get a specific node pool
pool, err := nodePools.GetNodePool(context.Background(), "your-org-id", "your-pool-name")
if err != nil {
    log.Fatalf("Failed to get node pool: %v", err)
}
fmt.Printf("Found node pool: %+v\n", pool)
```

### Working with Cloudspaces

```go
// Get a cloudspace manager
cloudspaces := client.Cloudspaces("your-org-id")

// List all cloudspaces
cloudspaceList, err := cloudspaces.ListCloudspaces(context.Background(), "your-org-id")
if err != nil {
    log.Fatalf("Failed to list cloudspaces: %v", err)
}
for _, cs := range cloudspaceList {
    fmt.Printf("Cloudspace: %s (ID: %s)\n", cs.Name, cs.ID)
}

// Get kubeconfig for a cloudspace
kubeconfig, err := cloudspaces.GetCloudspaceConfig(context.Background(), "your-org-id", "your-cloudspace-name")
if err != nil {
    log.Fatalf("Failed to get kubeconfig: %v", err)
}
fmt.Printf("Kubeconfig: %s\n", kubeconfig)
```

### Managing Organizations

```go
// Get the organization manager
orgManager := client.Organizations()

// List all organizations
orgs, err := orgManager.ListOrganizations(context.Background())
if err != nil {
    log.Fatalf("Failed to list organizations: %v", err)
}
for _, org := range orgs {
    fmt.Printf("Organization: %s (ID: %s)\n", org.Name, org.ID)
}
```

## Advanced Configuration

### Custom HTTP Client

You can provide a custom HTTP client with your own configuration:

```go
import (
    "crypto/tls"
    "net/http"
    "time"

    "github.com/rackspace-spot/spot-go-sdk/api/v1"
)

// Create a custom HTTP client with timeouts and TLS configuration
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: false, // Set to true only for testing
        },
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

// Create a new client with the custom HTTP client
client, err := v1.NewClient(&v1.Config{
    BaseURL:      "https://api.spot.io",
    OAuthURL:     "https://auth.spot.io",
    RefreshToken: os.Getenv("SPOT_REFRESH_TOKEN"),
    HTTPClient:   httpClient,
})
```

### Request Timeouts

You can set timeouts at the context level for individual requests:

```go
// Create a context with a timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Use the context with a request
orgs, err := client.Organizations().ListOrganizations(ctx)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        log.Fatal("Request timed out")
    }
    log.Fatalf("Request failed: %v", err)
}
```

## ⚠️ Breaking Changes in v1.0.0

Version 1.0.0 introduces several breaking changes to improve the API design and consistency:

### Major Changes

1. **Unified Node Pool Interface**
   - Combined `SpotNodePoolAPI` and `OnDemandNodePoolAPI` into a single `NodePoolAPI`
   - Added `NodePool` interface with common methods for both pool types
   - Added `NodePoolType` to distinguish between spot and on-demand pools

2. **Simplified Client Interface**
   - New `Client` interface replaces the old `SpotAPI`
   - Resource-specific methods grouped into manager interfaces
   - More consistent method signatures and return types

3. **Improved Error Handling**
   - Better error types and messages
   - More consistent error handling patterns

### Migration Guide

#### Before

```go
// Old way
client := rxtspot.NewSpotClient(cfg)

// List spot node pools
spotPools, err := client.ListSpotNodePools(ctx, org, cloudspace)

// List on-demand node pools
onDemandPools, err := client.ListOnDemandNodePools(ctx, org, cloudspace)
```

#### After

```go
// New way
client, err := rxtspot.NewClient(cfg)

// List all node pools
allPools, err := client.NodePools(org, cloudspace).ListNodePools(ctx, org, cloudspace, nil)

// List only spot node pools
spotPools, err := client.NodePools(org, cloudspace).ListNodePools(ctx, org, cloudspace, rxtspot.NodePoolTypeSpot)

// List only on-demand node pools
onDemandPools, err := client.NodePools(org, cloudspace).ListNodePools(ctx, org, cloudspace, rxtspot.NodePoolTypeOnDemand)
```

For more details, see the [API Reference](https://pkg.go.dev/github.com/rackspace-spot/spot-go-sdk/api/v1).

---

## Features

- Authenticate with Rackspace Spot using OAuth2 refresh tokens
- Create, list, and manage cloudspaces (Kubernetes clusters)
- Unified interface for managing both spot and on-demand node pools
- Query available regions, server classes, and price history
- Comprehensive error handling and logging
- Thread-safe client implementation

## Features (Planned)
- Authenticate with Rackspace Spot using OAuth2 refresh tokens
- Create, list, and delete cloudspaces
- Manage spot and on-demand node pools
- Query available regions, server classes, and price history
- Example CLI for resource management
- Comprehensive documentation and usage examples

## Roadmap
1. Core SDK: Authentication, cloudspace management
2. Node pool management (spot/on-demand)
3. Utility methods (regions, server classes, price history)
4. Example CLI tool
5. Tests and documentation


## Installation

### 1. Install the SDK

Clone this repository and use Go modules to import the SDK in your project:

```sh
git clone https://github.com/rackerlabs/spot-go-sdk.git
cd spot-go-sdk/rxtspot
```

Or add to your Go project:

```go
import v1 "github.com/rackerlabs/spot-go-sdk/rxtspot/api/v1"
```

### 2. Authentication

You need a Rackspace Spot refresh token. Set it as an environment variable:

```sh
export SPOT_REFRESH_TOKEN=your_refresh_token_here
```

## Example Usage

### Creating a Client

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rackspace-spot/spot-go-sdk/api/v1"
)

func main() {
	// Initialize client with configuration
	cfg := &v1.Config{
		BaseURL:      "https://api.spot.io",
		RefreshToken: os.Getenv("SPOT_REFRESH_TOKEN"),
	}

	client, err := v1.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Authenticate (this will be called automatically when needed)
	token, err := client.Authenticate(context.Background())
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	log.Printf("Authenticated with token: %s\n", token)

	// Example: List organizations
	orgs, err := client.Organizations().ListOrganizations(context.Background())
	if err != nil {
		log.Fatalf("Failed to list organizations: %v", err)
	}

	fmt.Println("Organizations:")
	for _, org := range orgs {
		fmt.Printf("- %s (ID: %s)\n", org.Name, org.ID)

		// Example: List cloudspaces for each organization
		cloudspaces, err := client.Cloudspaces(org.ID).ListCloudspaces(context.Background(), org.ID)
		if err != nil {
			log.Printf("Warning: Failed to list cloudspaces for org %s: %v\n", org.ID, err)
			continue
		}

		for _, cs := range cloudspaces {
			fmt.Printf("  - Cloudspace: %s\n", cs.Name)

			// Example: List node pools
			nodePools, err := client.NodePools(org.ID, cs.Name).ListNodePools(
				context.Background(),
				org.ID,
				cs.Name,
				nil, // nil means get all types of node pools
			)
			if err != nil {
				log.Printf("Warning: Failed to list node pools: %v\n", err)
				continue
			}

			for _, pool := range nodePools {
				fmt.Printf("    - %s (%s)\n", pool.GetName(), pool.GetType())
			}
		}
	}
}
```

### Managing Node Pools

```go
// Create a new spot node pool
spotPool := &v1.SpotNodePool{
    Name:        "my-spot-pool",
    InstanceType: "t3.medium",
    MinSize:     1,
    MaxSize:     5,
    TargetSize:  3,
}

// Add the node pool to a cloudspace
err = client.NodePools("my-org", "my-cloudspace").CreateNodePool(
    context.Background(),
    "my-org",
    spotPool,
)
if err != nil {
    log.Fatalf("Failed to create node pool: %v", err)
}

// List only spot node pools
spotPools, err := client.NodePools("my-org", "my-cloudspace").ListNodePools(
    context.Background(),
    "my-org",
    "my-cloudspace",
    v1.NodePoolTypeSpot, // Filter for spot pools only
)
```

### Error Handling

The SDK provides detailed error information through the `APIError` type:

```go
pool, err := client.NodePools("my-org", "my-cloudspace").GetNodePool(context.Background(), "my-org", "nonexistent-pool")
if err != nil {
    var apiErr *v1.APIError
    if errors.As(err, &apiErr) {
        log.Printf("API Error (Code %d): %s\n", apiErr.Code, apiErr.Message)
        // Handle specific error codes
        if apiErr.Code == http.StatusNotFound {
            log.Println("The requested resource was not found")
        }
    } else {
        log.Printf("Unexpected error: %v\n", err)

---

_See the SDK source and examples for more advanced usage and integration._ 