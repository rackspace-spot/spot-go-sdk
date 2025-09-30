// Package rxtspot provides a client for interacting with the Rackspace Spot API.
//
// This version (v1) introduces a new unified interface for managing node pools
// and other resources. The old interfaces have been removed in favor of a more
// consistent and type-safe API.
package rxtspot

import (
	"context"
)

// NodePoolType represents the type of node pool
type NodePoolType string

// const (
// 	// NodePoolTypeSpot represents a spot node pool
// 	NodePoolTypeSpot NodePoolType = "spot"
// 	// NodePoolTypeOnDemand represents an on-demand node pool
// 	NodePoolTypeOnDemand NodePoolType = "on-demand"
// )

// NodePool represents the common interface for all node pool types
// type NodePool interface {
// 	// GetType returns the type of the node pool (spot or on-demand)
// 	GetType() NodePoolType
// 	// GetName returns the name of the node pool
// 	GetName() string
// 	// GetOrg returns the organization ID this node pool belongs to
// 	GetOrg() string
// }

// NodePoolAPI defines the unified interface for managing node pools
// type NodePoolAPI interface {
// 	// ListNodePools lists all node pools, optionally filtered by type
// 	ListNodePools(ctx context.Context, org, cloudspace string, poolType *NodePoolType) ([]NodePool, error)
// 	// GetNodePool retrieves a specific node pool by name and type.
// 	// If poolType is nil, it will return an error as the type must be specified.
// 	GetNodePool(ctx context.Context, org, name string, poolType *NodePoolType) (NodePool, error)
// 	// CreateNodePool creates a new node pool
// 	CreateNodePool(ctx context.Context, org string, pool NodePool) error
// 	// UpdateNodePool updates an existing node pool
// 	UpdateNodePool(ctx context.Context, org string, pool NodePool) error
// 	// DeleteNodePool deletes a node pool by name
// 	DeleteNodePool(ctx context.Context, org, name string) error
// }

type SpotNodePoolManager interface {
	// ListNodePools lists all node pools
	ListNodePools(ctx context.Context, org, cloudspace string) ([]SpotNodePool, error)
	// GetNodePool retrieves a specific node pool by name and type.
	// If poolType is nil, it will return an error as the type must be specified.
	GetNodePool(ctx context.Context, org, name string) (SpotNodePool, error)
	// CreateNodePool creates a new node pool
	CreateNodePool(ctx context.Context, org string, pool SpotNodePool) error
	// UpdateNodePool updates an existing node pool
	UpdateNodePool(ctx context.Context, org string, pool SpotNodePool) error
	// DeleteNodePool deletes a node pool by name
	DeleteNodePool(ctx context.Context, org, name string) error
}

type OnDemandNodePoolManager interface {
	// ListNodePools lists all node pools
	ListNodePools(ctx context.Context, org, cloudspace string) ([]OnDemandNodePool, error)
	// GetNodePool retrieves a specific node pool by name and type.
	// If poolType is nil, it will return an error as the type must be specified.
	GetNodePool(ctx context.Context, org, name string) (OnDemandNodePool, error)
	// CreateNodePool creates a new node pool
	CreateNodePool(ctx context.Context, org string, pool OnDemandNodePool) error
	// UpdateNodePool updates an existing node pool
	UpdateNodePool(ctx context.Context, org string, pool OnDemandNodePool) error
	// DeleteNodePool deletes a node pool by name
	DeleteNodePool(ctx context.Context, org, name string) error
}

// CloudspaceManager defines the interface for managing cloudspaces
type CloudspaceManager interface {
	// ListCloudspaces lists all cloudspaces in an organization
	ListCloudspaces(ctx context.Context, org string) ([]*CloudSpace, error)
	// CreateCloudspace creates a new cloudspace
	CreateCloudspace(ctx context.Context, cs *CloudSpace) error
	// GetCloudspace retrieves a cloudspace by name
	GetCloudspace(ctx context.Context, org, name string) (*CloudSpace, error)
	// DeleteCloudspace deletes a cloudspace by name
	DeleteCloudspace(ctx context.Context, org, name string) error
	// GetCloudspaceConfig gets the kubeconfig for a cloudspace
	GetCloudspaceConfig(ctx context.Context, org, name string) (string, error)
}

// OrganizationManager defines the interface for managing organizations
type OrganizationManager interface {
	// ListOrganizations lists all organizations the authenticated user has access to
	ListOrganizations(ctx context.Context) ([]*Organization, error)
}

// RegionsManager defines the interface for managing regions
type RegionsManager interface {
	// ListRegions lists all available regions
	ListRegions(ctx context.Context) ([]*Region, error)
	// GetRegion gets a specific region by name
	GetRegion(ctx context.Context, name string) (*Region, error)
}

// ServerClassesManager defines the interface for managing server classes
type ServerClassesManager interface {
	// ListServerClasses lists all server classes for a region
	ListServerClasses(ctx context.Context, region string) (*ServerClassList, error)
	// GetServerClass gets a specific server class by name
	GetServerClass(ctx context.Context, name string) (*ServerClass, error)
}

// PricingManager defines the interface for managing pricing information
type PricingManager interface {
	// GetPriceDetailsForServerClass gets pricing details for a specific server class
	GetPriceDetailsForServerClass(ctx context.Context, serverClass string) (*PriceDetails, error)
	// GetPriceDetails gets all pricing details
	GetPriceDetails(ctx context.Context) ([]*PriceDetails, error)
	// GetPriceDetailsForRegion gets pricing details for a specific region
	GetPriceDetailsForRegion(ctx context.Context, region string) (*PriceDetails, error)
}

// Client is the main interface for the Rackspace Spot SDK
type Client interface {
	// Authenticate authenticates the client and returns an access token
	Authenticate(ctx context.Context) (string, error)

	// Organizations returns the organization manager
	Organizations() OrganizationManager

	// Cloudspaces returns a cloudspace manager for a specific organization
	Cloudspaces(org string) CloudspaceManager

	// SpotNodePools returns a node pool manager for a specific organization and cloudspace
	SpotNodepools(org, cloudspace string) SpotNodePoolManager

	// OnDemandNodePools returns a node pool manager for a specific organization and cloudspace
	OnDemandNodePools(org, cloudspace string) ondemandNodePoolManager

	// Regions returns the regions manager
	Regions() RegionsManager

	// ServerClasses returns the server classes manager for a specific region
	ServerClasses(region string) ServerClassesManager

	// Pricing returns the pricing manager
	Pricing() PricingManager
}

// Ensure compatibility with the old interface
type SpotAPI = Client
