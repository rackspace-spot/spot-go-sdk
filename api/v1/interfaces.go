package rxtspot

import "context"

//go:generate mockgen -source interfaces.go -destination mocks/mock_interfaces.go -package mocks

// OrganizationAPI defines organization-related methods.
type OrganizationAPI interface {
	ListOrganizations(ctx context.Context) ([]Organization, error)
}

// CloudspaceAPI defines cloudspace-related methods.
type CloudspaceAPI interface {
	ListCloudspaces(ctx context.Context, org string) (*CloudSpaceList, error)
	CreateCloudspace(ctx context.Context, cs CloudSpace) error
	GetCloudspace(ctx context.Context, org, name string) (*CloudSpace, error)
	DeleteCloudspace(ctx context.Context, org, name string) error
	GetCloudspaceConfig(ctx context.Context, org, name string) (string, error)
}

// SpotNodePoolAPI defines spot node pool methods.
type SpotNodePoolAPI interface {
	ListSpotNodePools(ctx context.Context, org string, cloudspace string) ([]*SpotNodePool, error)
	CreateSpotNodePool(ctx context.Context, org string, pool SpotNodePool) error
	UpdateSpotNodePool(ctx context.Context, org string, pool SpotNodePool) error
	GetSpotNodePool(ctx context.Context, org, name string) (*SpotNodePool, error)
	DeleteSpotNodePool(ctx context.Context, org, name string) error
}

// AutopilotNodePoolAPI defines autopilot node pool methods.
type AutopilotNodePoolAPI interface {
	ListAutopilotNodePools(ctx context.Context, org string, cloudspace string) ([]*AutopilotNodePool, error)
	CreateAutopilotNodePool(ctx context.Context, org string, pool AutopilotNodePool) error
	UpdateAutopilotNodePool(ctx context.Context, org string, pool AutopilotNodePool) error
	GetAutopilotNodePool(ctx context.Context, org, name string) (*AutopilotNodePool, error)
	DeleteAutopilotNodePool(ctx context.Context, org, name string) error
}

// OnDemandNodePoolAPI defines on-demand node pool methods.
type OnDemandNodePoolAPI interface {
	ListOnDemandNodePools(ctx context.Context, org string, cloudspace string) ([]*OnDemandNodePool, error)
	CreateOnDemandNodePool(ctx context.Context, org string, pool OnDemandNodePool) error
	UpdateOnDemandNodePool(ctx context.Context, org string, pool OnDemandNodePool) error
	GetOnDemandNodePool(ctx context.Context, org, name string) (*OnDemandNodePool, error)
	DeleteOnDemandNodePool(ctx context.Context, org, name string) error
}

type SpotRegionsAPI interface {
	ListRegions(ctx context.Context) ([]Region, error)
	GetRegion(ctx context.Context, name string) (*Region, error)
}

type SpotServerClassesAPI interface {
	ListServerClasses(ctx context.Context, region string) (*ServerClassList, error)
	GetServerClass(ctx context.Context, name string) (*ServerClass, error)
}

type SpotPricingAPI interface {
	GetPriceDetailsForServerClass(ctx context.Context, serverClass string) (*PriceDetails, error)
	GetPriceDetails(ctx context.Context) ([]*PriceDetails, error)
	GetPriceDetailsForRegion(ctx context.Context, region string) (*PriceDetails, error)
	GetMarketPriceForServerClass(ctx context.Context, serverClassStatus *ServerClassStatus) string
	GetMinimumBidPriceForServerClass(ctx context.Context, serverClassSpec *ServerClassSpec) string
}

// VMCloudSpaceAPI defines VM cloud space methods.
type VMCloudSpaceAPI interface {
	ListVMCloudSpaces(ctx context.Context, org string) (*VMCloudSpaceList, error)
	CreateVMCloudSpace(ctx context.Context, vmcs VMCloudSpace) error
	GetVMCloudSpace(ctx context.Context, org, name string) (*VMCloudSpace, error)
	UpdateVMCloudSpace(ctx context.Context, org string, vmcs VMCloudSpace) error
	DeleteVMCloudSpace(ctx context.Context, org, name string) error
}

// VMPoolAPI defines VM pool methods.
type VMPoolAPI interface {
	ListVMPools(ctx context.Context, org string, vmCloudSpace string) ([]*VMPool, error)
	CreateVMPool(ctx context.Context, org string, pool VMPool) error
	UpdateVMPool(ctx context.Context, org string, pool VMPool) error
	GetVMPool(ctx context.Context, org, name string) (*VMPool, error)
	DeleteVMPool(ctx context.Context, org, name string) error
}

// VMSSHKeyAPI defines VM SSH key methods.
type VMSSHKeyAPI interface {
	ListVMSSHKeys(ctx context.Context, org string) (*VMSSHKeyList, error)
	CreateVMSSHKey(ctx context.Context, key VMSSHKey) error
	GetVMSSHKey(ctx context.Context, org, name string) (*VMSSHKey, error)
	DeleteVMSSHKey(ctx context.Context, org, name string) error
}

// SpotAPI defines the complete interface for the Rackspace Spot SDK client.
// It embeds all the specific APIs to provide a unified interface.
type SpotAPI interface {
	Authenticate(ctx context.Context) (string, error)
	OrganizationAPI
	CloudspaceAPI
	SpotNodePoolAPI
	OnDemandNodePoolAPI
	AutopilotNodePoolAPI
	SpotRegionsAPI
	SpotServerClassesAPI
	SpotPricingAPI
	VMCloudSpaceAPI
	VMPoolAPI
	VMSSHKeyAPI
}
