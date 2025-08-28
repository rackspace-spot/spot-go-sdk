# Rackspace Spot Go SDK

This package provides an idiomatic Go SDK for interacting with the Rackspace Spot platform. It enables developers and DevOps teams to programmatically manage cloud resources such as cloudspaces (Kubernetes clusters), spot node pools, and on-demand node pools.

## Features

- **Authentication**: Secure authentication with the Rackspace Spot API
- **CloudSpaces**: Create, list, update, and delete Kubernetes clusters
- **Node Pools**: Manage both spot and on-demand node pools
- **Server Classes**: Query available server classes and their specifications
- **Regions**: List available regions and their details
- **Type-Safe**: Strongly typed API for better developer experience
- **Idiomatic Go**: Follows Go best practices and conventions

## Installation

Add the SDK to your Go module:

```bash
go get github.com/rackspace/spot-go-sdk
```

## Usage

### Creating a Client

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/rackspace/spot-go-sdk/api/v1"
)

func main() {
	// Initialize the client
	cfg := v1.ClientConfig{
		BaseURL:    "https://api.spot.io", // Replace with actual API URL
		AuthToken:  os.Getenv("SPOT_AUTH_TOKEN"),
		HTTPClient: nil, // Uses default client if nil
	}

	client, err := v1.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Use the client...
}
```

### Example: List All Regions

```go
regions, err := client.Regions().List(context.Background(), "", nil)
if err != nil {
    log.Fatalf("Failed to list regions: %v", err)
}

for _, r := range regions {
    fmt.Printf("- %s (%s)\n", r.Name, r.Provider)
}
```

### Example: Create a CloudSpace

```go
createOpts := &cloudspace.CreateOptions{
    Name:        "my-cloudspace",
    Org:         "my-org",
    Region:      "us-east-1",
    ServerClass: "general-purpose-1",
}

createdCS, err := client.Cloudspaces().Create(context.Background(), createOpts)
if err != nil {
    log.Fatalf("Failed to create cloudspace: %v", err)
}
fmt.Printf("Created cloudspace: %+v\n", createdCS)
```

## Package Structure

```
api/v1/
├── client.go              # Main client interface and implementation
├── client_interface.go    # Public client interface
├── cloudspace/            # CloudSpace types and service
│   └── types.go
├── nodepool/
│   ├── spot/              # Spot node pool types and service
│   │   ├── types.go
│   │   └── service.go
│   └── ondemand/          # On-demand node pool types and service
│       ├── types.go
│       └── service.go
├── region/                # Region types and service
│   ├── types.go
│   └── service.go
└── serverclass/           # Server class types and service
    ├── types.go
    └── service.go
```

## Authentication

You'll need an authentication token to use the SDK. Set it as an environment variable:

```bash
export SPOT_AUTH_TOKEN=your_auth_token_here
```

## Examples

See the `examples/` directory for more comprehensive examples:

```bash
cd examples/basic
go run main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

See [`examples/main.go`](examples/main.go) for a full example. Here is a minimal usage snippet:

```go
package main

import (
    "context"
    "fmt"
    "os"
    v1 "github.com/rackerlabs/spot-go-sdk/rxtspot/api/v1"
)

func main() {
    refreshToken := os.Getenv("SPOT_REFRESH_TOKEN")
    client := v1.NewClient(refreshToken)
    if err := client.Authenticate(context.Background()); err != nil {
        panic(err)
    }
    orgs, err := client.ListOrganizations(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Println("Organizations:", orgs)
}
```

### 4. Run the Example

```sh
cd examples
export SPOT_REFRESH_TOKEN=your_refresh_token_here
go run main.go
```

This will demonstrate authentication and CRUD operations for all major objects.

---

_See the SDK source and examples for more advanced usage and integration._ 