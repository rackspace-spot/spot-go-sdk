package serverclass

import (
	"time"

	"github.com/rackspace/spot-go-sdk/api/v1/types"
)

// ServerClass represents a server class configuration
type ServerClass struct {
	Name                      string   `json:"name" yaml:"name"`
	Category                  string   `json:"category,omitempty" yaml:"category,omitempty"`
	Availability              string   `json:"availability,omitempty" yaml:"availability,omitempty"`
	DisplayName               string   `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Region                    string   `json:"region,omitempty" yaml:"region,omitempty"`
	MinBidPricePerHour        string   `json:"minBidPricePerHour,omitempty" yaml:"minBidPricePerHour,omitempty"`
	CurrentMarketPricePerHour string   `json:"currentMarketPricePerHour,omitempty" yaml:"currentMarketPricePerHour,omitempty"`
	OnDemandPricePerHour      string   `json:"onDemandPricePerHour,omitempty" yaml:"onDemandPricePerHour,omitempty"`
	Resources                 Resource `json:"resources,omitempty" yaml:"resources,omitempty"`
}

// ServerClassList represents a list of server classes
type ServerClassList struct {
	Items []ServerClass `json:"serverClasses" yaml:"serverClasses"`
}

// Resource represents compute resources (CPU, Memory)
type Resource struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

// ListOptions specifies optional parameters for listing server classes
type ListOptions struct {
	Region string `json:"region,omitempty"`
}

// GetOptions specifies optional parameters for getting a server class
type GetOptions struct {
	IncludePricing bool `json:"includePricing,omitempty"`
}

// PriceDetails represents pricing information for a server class
type PriceDetails struct {
	ServerClass    string    `json:"serverClass" yaml:"serverClass"`
	Region         string    `json:"region" yaml:"region"`
	MarketPrice    string    `json:"marketPrice" yaml:"marketPrice"`
	OnDemandPrice  string    `json:"onDemandPrice" yaml:"onDemandPrice"`
	MinBidPrice    string    `json:"minBidPrice" yaml:"minBidPrice"`
	UpdatedAt      time.Time `json:"updatedAt" yaml:"updatedAt"`
}
