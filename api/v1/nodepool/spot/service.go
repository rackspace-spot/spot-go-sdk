package spot

import (
	"context"
)

// Service defines the interface for spot node pool operations
type Service interface {
	// Create creates a new spot node pool
	Create(ctx context.Context, opts *CreateOptions) (*SpotNodePool, error)

	// Get retrieves a spot node pool by name
	Get(ctx context.Context, name, org string) (*SpotNodePool, error)

	// List retrieves all spot node pools, optionally filtered by organization and cloudspace
	List(ctx context.Context, org, cloudspace string) ([]*SpotNodePool, error)

	// Update updates an existing spot node pool
	Update(ctx context.Context, name string, opts *UpdateOptions) (*SpotNodePool, error)

	// Delete deletes a spot node pool by name
	Delete(ctx context.Context, name, org string) error

	// GetBidStatus retrieves the current bid status for a spot node pool
	GetBidStatus(ctx context.Context, name, org string) (*BidStatus, error)
}
