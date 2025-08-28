package ondemand

import (
	"context"
)

// Service defines the interface for on-demand node pool operations
type Service interface {
	// Create creates a new on-demand node pool
	Create(ctx context.Context, opts *CreateOptions) (*OnDemandNodePool, error)

	// Get retrieves an on-demand node pool by name
	Get(ctx context.Context, name, org string) (*OnDemandNodePool, error)

	// List retrieves all on-demand node pools, optionally filtered by organization and cloudspace
	List(ctx context.Context, org, cloudspace string) ([]*OnDemandNodePool, error)

	// Update updates an existing on-demand node pool
	Update(ctx context.Context, name string, opts *UpdateOptions) (*OnDemandNodePool, error)

	// Delete deletes an on-demand node pool by name
	Delete(ctx context.Context, name, org string) error
}
