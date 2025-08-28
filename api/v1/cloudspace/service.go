package cloudspace

import (
	"context"
)

// Service defines the interface for cloudspace operations
type Service interface {
	// Create creates a new cloudspace
	Create(ctx context.Context, opts *CreateOptions) (*CloudSpace, error)

	// Get retrieves a cloudspace by name
	Get(ctx context.Context, name, org string) (*CloudSpace, error)

	// List retrieves all cloudspaces, optionally filtered by organization
	List(ctx context.Context, org string) ([]*CloudSpace, error)

	// Update updates an existing cloudspace
	Update(ctx context.Context, name string, opts *UpdateOptions) (*CloudSpace, error)

	// Delete deletes a cloudspace by name
	Delete(ctx context.Context, name, org string) error

	// GetConfig retrieves the kubeconfig for a cloudspace
	GetConfig(ctx context.Context, name, org string) (string, error)
}
