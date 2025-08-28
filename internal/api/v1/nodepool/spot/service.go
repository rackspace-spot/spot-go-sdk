package spot

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rackspace-spot/spot-go-sdk/api/v1/nodepool/spot"
	intclient "github.com/rackspace-spot/spot-go-sdk/internal/api/v1/client"
)

type service struct {
	client *intclient.HTTPClient
}

// New creates a new SpotNodePool service implementation
func New(client *intclient.HTTPClient) spot.Service {
	return &service{
		client: client,
	}
}

func (s *service) Create(ctx context.Context, opts *spot.CreateOptions) (*spot.SpotNodePool, error) {
	if opts == nil {
		return nil, fmt.Errorf("create options cannot be nil")
	}

	var created spot.SpotNodePool
	if err := s.client.Do(ctx, "POST", "/spot-nodepools", opts, &created); err != nil {
		return nil, fmt.Errorf("failed to create spot node pool: %w", err)
	}

	return &created, nil
}

func (s *service) Get(ctx context.Context, name, org string) (*spot.SpotNodePool, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if org == "" {
		return nil, fmt.Errorf("org is required")
	}

	var pool spot.SpotNodePool
	path := fmt.Sprintf("/orgs/%s/spot-nodepools/%s", url.PathEscape(org), url.PathEscape(name))
	if err := s.client.Do(ctx, "GET", path, nil, &pool); err != nil {
		return nil, fmt.Errorf("failed to get spot node pool: %w", err)
	}

	return &pool, nil
}

func (s *service) List(ctx context.Context, org, cloudspace string) ([]*spot.SpotNodePool, error) {
	if org == "" {
		return nil, fmt.Errorf("org is required")
	}

	var list struct {
		Items []*spot.SpotNodePool `json:"spotNodepools"`
	}

	path := fmt.Sprintf("/orgs/%s/spot-nodepools", url.PathEscape(org))
	if cloudspace != "" {
		path = fmt.Sprintf("%s?cloudspace=%s", path, url.QueryEscape(cloudspace))
	}

	if err := s.client.Do(ctx, "GET", path, nil, &list); err != nil {
		return nil, fmt.Errorf("failed to list spot node pools: %w", err)
	}

	return list.Items, nil
}

func (s *service) Update(ctx context.Context, name string, opts *spot.UpdateOptions) (*spot.SpotNodePool, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if opts == nil {
		return nil, fmt.Errorf("update options cannot be nil")
	}
	if opts.Org == "" {
		return nil, fmt.Errorf("org is required in update options")
	}

	path := fmt.Sprintf("/orgs/%s/spot-nodepools/%s", url.PathEscape(opts.Org), url.PathEscape(name))
	var updated spot.SpotNodePool
	if err := s.client.Do(ctx, "PATCH", path, opts, &updated); err != nil {
		return nil, fmt.Errorf("failed to update spot node pool: %w", err)
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

	path := fmt.Sprintf("/orgs/%s/spot-nodepools/%s", url.PathEscape(org), url.PathEscape(name))
	if err := s.client.Do(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete spot node pool: %w", err)
	}

	return nil
}

func (s *service) GetBidStatus(ctx context.Context, name, org string) (*spot.BidStatus, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if org == "" {
		return nil, fmt.Errorf("org is required")
	}

	var status spot.BidStatus
	path := fmt.Sprintf("/orgs/%s/spot-nodepools/%s/bid-status", url.PathEscape(org), url.PathEscape(name))
	if err := s.client.Do(ctx, "GET", path, nil, &status); err != nil {
		return nil, fmt.Errorf("failed to get bid status: %w", err)
	}

	return &status, nil
}
