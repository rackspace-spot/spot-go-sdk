package serverclass

import (
	"context"
)

// Service defines the interface for server class operations
type Service interface {
	// Get retrieves a server class by name
	Get(ctx context.Context, name string, opts *GetOptions) (*ServerClass, error)

	// List retrieves all server classes, optionally filtered by region
	List(ctx context.Context, region string) ([]*ServerClass, error)

	// GetPricing retrieves pricing information for a server class
	GetPricing(ctx context.Context, name, region string) (*PriceDetails, error)

	// ListPricing retrieves pricing information for all server classes
	ListPricing(ctx context.Context, region string) ([]*PriceDetails, error)
}
