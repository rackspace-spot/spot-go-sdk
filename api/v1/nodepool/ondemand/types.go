package ondemand

import (
	"time"
)

// OnDemandNodePool represents an on-demand node pool configuration
type OnDemandNodePool struct {
	Name                 string                 `json:"name" yaml:"name"`
	CreationTimestamp    time.Time              `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Org                  string                 `json:"org,omitempty" yaml:"org,omitempty"`
	Cloudspace           string                 `json:"cloudspace,omitempty" yaml:"cloudspace,omitempty"`
	ServerClass          string                 `json:"serverClass,omitempty" yaml:"serverClass,omitempty"`
	Desired              int                    `json:"desired,omitempty" yaml:"desired,omitempty"`
	WonCount             int                    `json:"wonCount,omitempty" yaml:"wonCount,omitempty"`
	CustomAnnotations    map[string]string      `json:"customAnnotations,omitempty" yaml:"customAnnotations,omitempty"`
	CustomLabels         map[string]string      `json:"customLabels,omitempty" yaml:"customLabels,omitempty"`
	CustomTaints         []map[string]interface{} `json:"customTaints,omitempty" yaml:"customTaints,omitempty"`
	OnDemandPricePerHour string                 `json:"onDemandPricePerHour,omitempty" yaml:"onDemandPricePerHour,omitempty"`
	Status               string                 `json:"status,omitempty" yaml:"status,omitempty"`
	Autoscaling          Autoscaling            `json:"autoscaling" yaml:"autoscaling"`
}

// OnDemandNodePoolList represents a list of on-demand node pools
type OnDemandNodePoolList struct {
	Items []OnDemandNodePool `json:"ondemandNodepools" yaml:"ondemandNodepools"`
}

// CreateOptions specifies parameters for creating a new on-demand node pool
type CreateOptions struct {
	Name                 string                 `json:"name"`
	Org                  string                 `json:"org"`
	Cloudspace           string                 `json:"cloudspace"`
	ServerClass          string                 `json:"serverClass"`
	Desired              int                    `json:"desired,omitempty"`
	CustomAnnotations    map[string]string      `json:"customAnnotations,omitempty"`
	CustomLabels         map[string]string      `json:"customLabels,omitempty"`
	CustomTaints         []map[string]interface{} `json:"customTaints,omitempty"`
	OnDemandPricePerHour string                 `json:"onDemandPricePerHour,omitempty"`
	Autoscaling          Autoscaling            `json:"autoscaling,omitempty"`
}

// UpdateOptions specifies parameters for updating an on-demand node pool
type UpdateOptions struct {
	Name                 string                 `json:"name"`
	Org                  string                 `json:"org"`
	Desired              int                    `json:"desired,omitempty"`
	CustomAnnotations    map[string]string      `json:"customAnnotations,omitempty"`
	CustomLabels         map[string]string      `json:"customLabels,omitempty"`
	CustomTaints         []map[string]interface{} `json:"customTaints,omitempty"`
	OnDemandPricePerHour string                 `json:"onDemandPricePerHour,omitempty"`
	Autoscaling          *Autoscaling           `json:"autoscaling,omitempty"`
}

// ListOptions specifies optional parameters for listing on-demand node pools
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
