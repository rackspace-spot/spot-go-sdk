package spot

import (
	"time"
)

// SpotNodePool represents a spot node pool configuration
type SpotNodePool struct {
	Name              string                 `json:"name" yaml:"name"`
	CreationTimestamp time.Time              `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Org               string                 `json:"org,omitempty" yaml:"org,omitempty"`
	Cloudspace        string                 `json:"cloudspace,omitempty" yaml:"cloudspace,omitempty"`
	ServerClass       string                 `json:"serverClass,omitempty" yaml:"serverClass,omitempty"`
	Desired           int                    `json:"desired,omitempty" yaml:"desired,omitempty"`
	WonCount          int                    `json:"wonCount,omitempty" yaml:"wonCount,omitempty"`
	CustomAnnotations map[string]string      `json:"customAnnotations,omitempty" yaml:"customAnnotations,omitempty"`
	CustomLabels      map[string]string      `json:"customLabels,omitempty" yaml:"customLabels,omitempty"`
	CustomTaints      []map[string]interface{} `json:"customTaints,omitempty" yaml:"customTaints,omitempty"`
	BidPrice          string                 `json:"bidPrice,omitempty" yaml:"bidPrice,omitempty"`
	Status            string                 `json:"status,omitempty" yaml:"status,omitempty"`
	Autoscaling       Autoscaling            `json:"autoscaling" yaml:"autoscaling"`
}

// SpotNodePoolList represents a list of spot node pools
type SpotNodePoolList struct {
	Items []SpotNodePool `json:"spotNodepools" yaml:"spotNodepools"`
}

// CreateOptions specifies parameters for creating a new spot node pool
type CreateOptions struct {
	Name              string                 `json:"name"`
	Org               string                 `json:"org"`
	Cloudspace        string                 `json:"cloudspace"`
	ServerClass       string                 `json:"serverClass"`
	Desired           int                    `json:"desired,omitempty"`
	BidPrice          string                 `json:"bidPrice"`
	CustomAnnotations map[string]string      `json:"customAnnotations,omitempty"`
	CustomLabels      map[string]string      `json:"customLabels,omitempty"`
	CustomTaints      []map[string]interface{} `json:"customTaints,omitempty"`
	Autoscaling       Autoscaling            `json:"autoscaling,omitempty"`
}

// UpdateOptions specifies parameters for updating a spot node pool
type UpdateOptions struct {
	Name              string                 `json:"name"`
	Org               string                 `json:"org"`
	Desired           int                    `json:"desired,omitempty"`
	BidPrice          string                 `json:"bidPrice,omitempty"`
	CustomAnnotations map[string]string      `json:"customAnnotations,omitempty"`
	CustomLabels      map[string]string      `json:"customLabels,omitempty"`
	CustomTaints      []map[string]interface{} `json:"customTaints,omitempty"`
	Autoscaling       *Autoscaling           `json:"autoscaling,omitempty"`
}

// ListOptions specifies optional parameters for listing spot node pools
type ListOptions struct {
	Org        string            `json:"org,omitempty"`
	Cloudspace string            `json:"cloudspace,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

// Autoscaling represents autoscaling configuration for a node pool
type Autoscaling struct {
	Enabled  bool  `json:"enabled" yaml:"enabled"`
	MinNodes int64 `json:"minNodes,omitempty" yaml:"minNodes,omitempty"`
	MaxNodes int64 `json:"maxNodes,omitempty" yaml:"maxNodes,omitempty"`
}

// BidStatus represents the current bid status for a spot node pool
type BidStatus struct {
	// CurrentBidPrice is the current bid price for the node pool
	CurrentBidPrice string `json:"currentBidPrice" yaml:"currentBidPrice"`
	// Status is the current status of the bid (e.g., "active", "outbid", etc.)
	Status string `json:"status" yaml:"status"`
	// LastUpdated is when the bid status was last updated
	LastUpdated time.Time `json:"lastUpdated" yaml:"lastUpdated"`
}
