package region

import (
	"context"
)

// Service defines the interface for region operations
type Service interface {
	// Get retrieves a region by name
	Get(ctx context.Context, name string, opts *GetOptions) (*Region, error)

	// List retrieves all regions, optionally filtered by provider and enabled status
	List(ctx context.Context, provider string, enabled *bool) ([]*Region, error)
}
