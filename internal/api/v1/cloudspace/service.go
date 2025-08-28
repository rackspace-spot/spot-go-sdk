package cloudspace

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rackspace-spot/spot-go-sdk/api/v1/cloudspace"
	intclient "github.com/rackspace-spot/spot-go-sdk/internal/api/v1/client"
)

type service struct {
	client *intclient.HTTPClient
}

// New creates a new CloudSpace service implementation
func New(client *intclient.HTTPClient) cloudspace.Service {
	return &service{
		client: client,
	}
}

func (s *service) Create(ctx context.Context, opts *cloudspace.CreateOptions) (*cloudspace.CloudSpace, error) {
	if opts == nil {
		return nil, fmt.Errorf("create options cannot be nil")
	}

	var created cloudspace.CloudSpace
	if err := s.client.Do(ctx, "POST", "/cloudspaces", opts, &created); err != nil {
		return nil, fmt.Errorf("failed to create cloudspace: %w", err)
	}

	return &created, nil
}

func (s *service) Get(ctx context.Context, name, org string) (*cloudspace.CloudSpace, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if org == "" {
		return nil, fmt.Errorf("org is required")
	}

	var cs cloudspace.CloudSpace
	path := fmt.Sprintf("/orgs/%s/cloudspaces/%s", url.PathEscape(org), url.PathEscape(name))
	if err := s.client.Do(ctx, "GET", path, nil, &cs); err != nil {
		return nil, fmt.Errorf("failed to get cloudspace: %w", err)
	}

	return &cs, nil
}

func (s *service) List(ctx context.Context, org string) ([]*cloudspace.CloudSpace, error) {
	if org == "" {
		return nil, fmt.Errorf("org is required")
	}

	var list struct {
		Items []*cloudspace.CloudSpace `json:"items"`
	}

	path := fmt.Sprintf("/orgs/%s/cloudspaces", url.PathEscape(org))
	if err := s.client.Do(ctx, "GET", path, nil, &list); err != nil {
		return nil, fmt.Errorf("failed to list cloudspaces: %w", err)
	}

	return list.Items, nil
}

func (s *service) Update(ctx context.Context, name string, opts *cloudspace.UpdateOptions) (*cloudspace.CloudSpace, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if opts == nil {
		return nil, fmt.Errorf("update options cannot be nil")
	}
	if opts.Org == "" {
		return nil, fmt.Errorf("org is required in update options")
	}

	path := fmt.Sprintf("/orgs/%s/cloudspaces/%s", url.PathEscape(opts.Org), url.PathEscape(name))
	var updated cloudspace.CloudSpace
	if err := s.client.Do(ctx, "PATCH", path, opts, &updated); err != nil {
		return nil, fmt.Errorf("failed to update cloudspace: %w", err)
	}

	return &updated, nil
}

func (s *service) Delete(ctx context.Context, name, org string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if org == "" {
		return fmt.Errorf("org is required")
	}

	path := fmt.Sprintf("/orgs/%s/cloudspaces/%s", url.PathEscape(org), url.PathEscape(name))
	if err := s.client.Do(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete cloudspace: %w", err)
	}

	return nil
}

func (s *service) GetConfig(ctx context.Context, name, org string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name is required")
	}
	if org == "" {
		return "", fmt.Errorf("org is required")
	}

	var config struct {
		Kubeconfig string `json:"kubeconfig"`
	}

	path := fmt.Sprintf("/orgs/%s/cloudspaces/%s/kubeconfig", url.PathEscape(org), url.PathEscape(name))
	if err := s.client.Do(ctx, "GET", path, nil, &config); err != nil {
		return "", fmt.Errorf("failed to get cloudspace config: %w", err)
	}

	return config.Kubeconfig, nil
}
