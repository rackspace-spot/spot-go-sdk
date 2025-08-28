package region

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rackspace/spot-go-sdk/api/v1/region"
	intclient "github.com/rackspace/spot-go-sdk/internal/api/v1/client"
)

type service struct {
	client *intclient.HTTPClient
}

// New creates a new Region service implementation
func New(client *intclient.HTTPClient) region.Service {
	return &service{
		client: client,
	}
}

func (s *service) Get(ctx context.Context, name string, opts *region.GetOptions) (*region.Region, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	path := fmt.Sprintf("/regions/%s", url.PathEscape(name))
	if opts != nil && opts.IncludeZones {
		path = fmt.Sprintf("%s?includeZones=true", path)
	}

	var r region.Region
	if err := s.client.Do(ctx, "GET", path, nil, &r); err != nil {
		return nil, fmt.Errorf("failed to get region: %w", err)
	}

	return &r, nil
}

func (s *service) List(ctx context.Context, provider string, enabled *bool) ([]*region.Region, error) {
	path := "/regions"
	params := make(url.Values)

	if provider != "" {
		params.Add("provider", provider)
	}
	if enabled != nil {
		params.Add("enabled", fmt.Sprintf("%t", *enabled))
	}

	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	var list region.RegionList
	if err := s.client.Do(ctx, "GET", path, nil, &list); err != nil {
		return nil, fmt.Errorf("failed to list regions: %w", err)
	}

	return list.Items, nil
}
