package v1

import (
	"context"

	"github.com/rackspace-spot/spot-go-sdk/api/v1/cloudspace"
	"github.com/rackspace-spot/spot-go-sdk/api/v1/nodepool/ondemand"
	"github.com/rackspace-spot/spot-go-sdk/api/v1/nodepool/spot"
	"github.com/rackspace-spot/spot-go-sdk/api/v1/region"
	"github.com/rackspace-spot/spot-go-sdk/api/v1/serverclass"
)

// Client defines the main interface for the Spot SDK client.
// It provides access to all the services available in the SDK.
type Client interface {
	// Cloudspaces provides access to cloudspace management operations
	Cloudspaces() cloudspace.Service

	// SpotNodePools provides access to spot node pool management operations
	SpotNodePools() spot.Service

	// OnDemandNodePools provides access to on-demand node pool management operations
	OnDemandNodePools() ondemand.Service

	// ServerClasses provides access to server class information
	ServerClasses() serverclass.Service

	// Regions provides access to region information
	Regions() region.Service

	// Close releases any resources held by the client
	Close() error
}
