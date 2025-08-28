package ondemand

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rackspace/spot-go-sdk/api/v1/nodepool/ondemand"
	intclient "github.com/rackspace/spot-go-sdk/internal/api/v1/client"
)

type service struct {
	client *intclient.HTTPClient
}

// New creates a new OnDemandNodePool service implementation
func New(client *intclient.HTTPClient) ondemand.Service {
	return &service{
		client: client,
	}
}

func (s *service) Create(ctx context.Context, opts *ondemand.CreateOptions) (*ondemand.OnDemandNodePool, error) {
	if opts == nil {
		return nil, fmt.Errorf("create options cannot be nil")
	}

	var created ondemand.OnDemandNodePool
	if err := s.client.Do(ctx, "POST", "/ondemand-nodepools", opts, &created); err != nil {
		return nil, fmt.Errorf("failed to create on-demand node pool: %w", err)
	}

	return &created, nil
}

func (s *service) Get(ctx context.Context, name, org string) (*ondemand.OnDemandNodePool, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if org == "" {
		return nil, fmt.Errorf("org is required")
	}

	var pool ondemand.OnDemandNodePool
	path := fmt.Sprintf("/orgs/%s/ondemand-nodepools/%s", url.PathEscape(org), url.PathEscape(name))
	if err := s.client.Do(ctx, "GET", path, nil, &pool); err != nil {
		return nil, fmt.Errorf("failed to get on-demand node pool: %w", err)
	}

	return &pool, nil
}

func (s *service) List(ctx context.Context, org, cloudspace string) ([]*ondemand.OnDemandNodePool, error) {
	if org == "" {
		return nil, fmt.Errorf("org is required")
	}

	path := "/ondemand-nodepools"
	if cloudspace != "" {
		path = fmt.Sprintf("%s?cloudspace=%s", path, url.QueryEscape(cloudspace))
	}

	var list ondemand.OnDemandNodePoolList
	if err := s.client.Do(ctx, "GET", path, nil, &list); err != nil {
		return nil, fmt.Errorf("failed to list on-demand node pools: %w", err)
	}

	return list.Items, nil
}

func (s *service) Update(ctx context.Context, name string, opts *ondemand.UpdateOptions) (*ondemand.OnDemandNodePool, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if opts == nil {
		return nil, fmt.Errorf("update options cannot be nil")
	}
	if opts.Org == "" {
		return nil, fmt.Errorf("org is required in update options")
	}

	path := fmt.Sprintf("/orgs/%s/ondemand-nodepools/%s", url.PathEscape(opts.Org), url.PathEscape(name))
	var updated ondemand.OnDemandNodePool
	if err := s.client.Do(ctx, "PATCH", path, opts, &updated); err != nil {
		return nil, fmt.Errorf("failed to update on-demand node pool: %w", err)
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

	path := fmt.Sprintf("/orgs/%s/ondemand-nodepools/%s", url.PathEscape(org), url.PathEscape(name))
	if err := s.client.Do(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete on-demand node pool: %w", err)
	}

	return nil
}
