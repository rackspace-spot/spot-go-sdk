package cloudspaces

import (
	"context"
	"fmt"

	"github.com/rackspace-spot/spot-go-sdk/pkg/validation"
)

type service struct {
	client Client
}

func NewService(client Client) service {
	return service{
		client: client,
	}
}

func (s service) ListCloudspaces(ctx context.Context, org string) (*CloudSpaceList, error) {
	return s.client.ListCloudspaces(ctx, org)
}

func (s service) CreateCloudspace(ctx context.Context, cs CloudSpace) error {
	if err := validation.ValidateOrgName(cs.Org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := validation.ValidateResourceName(cs.Name); err != nil {
		return fmt.Errorf("invalid cloudspace name: %w", err)
	}
	if cs.Region == "" {
		return fmt.Errorf("region is required")
	}
	if cs.KubernetesVersion == "" {
		return fmt.Errorf("kubernetes version is required")
	}

	exists, orgID, err := s.client.getOrgIDIfExists(ctx, cs.Org)
	if err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", cs.Org)
	}
	return s.client.CreateCloudspace(ctx, cs)
}

func (s service) GetCloudspace(ctx context.Context, org, name string) (*CloudSpace, error) {
	return s.client.GetCloudspace(ctx, org, name)
}

func (s service) DeleteCloudspace(ctx context.Context, org, name string) error {
	return s.client.DeleteCloudspace(ctx, org, name)
}

func (s service) GetCloudspaceConfig(ctx context.Context, org, name string) (string, error) {
	return s.client.GetCloudspaceConfig(ctx, org, name)
}
