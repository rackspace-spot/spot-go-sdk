package serverclass

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rackspace/spot-go-sdk/api/v1/serverclass"
	intclient "github.com/rackspace/spot-go-sdk/internal/api/v1/client"
)

type service struct {
	client *intclient.HTTPClient
}

// New creates a new ServerClass service implementation
func New(client *intclient.HTTPClient) serverclass.Service {
	return &service{
		client: client,
	}
}

func (s *service) Get(ctx context.Context, name string, opts *serverclass.GetOptions) (*serverclass.ServerClass, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	path := fmt.Sprintf("/server-classes/%s", url.PathEscape(name))
	if opts != nil && opts.IncludePricing {
		path = fmt.Sprintf("%s?includePricing=true", path)
	}

	var sc serverclass.ServerClass
	if err := s.client.Do(ctx, "GET", path, nil, &sc); err != nil {
		return nil, fmt.Errorf("failed to get server class: %w", err)
	}

	return &sc, nil
}

func (s *service) List(ctx context.Context, region string) ([]*serverclass.ServerClass, error) {
	path := "/server-classes"
	if region != "" {
		path = fmt.Sprintf("%s?region=%s", path, url.QueryEscape(region))
	}

	var list serverclass.ServerClassList
	if err := s.client.Do(ctx, "GET", path, nil, &list); err != nil {
		return nil, fmt.Errorf("failed to list server classes: %w", err)
	}

	return list.Items, nil
}

func (s *service) GetPricing(ctx context.Context, name, region string) (*serverclass.PriceDetails, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}

	path := fmt.Sprintf("/server-classes/%s/pricing?region=%s", 
		url.PathEscape(name), 
		url.QueryEscape(region))

	var pricing serverclass.PriceDetails
	if err := s.client.Do(ctx, "GET", path, nil, &pricing); err != nil {
		return nil, fmt.Errorf("failed to get server class pricing: %w", err)
	}

	return &pricing, nil
}

func (s *service) ListPricing(ctx context.Context, region string) ([]*serverclass.PriceDetails, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}

	path := fmt.Sprintf("/server-classes/pricing?region=%s", url.QueryEscape(region))

	var pricingList []*serverclass.PriceDetails
	if err := s.client.Do(ctx, "GET", path, nil, &pricingList); err != nil {
		return nil, fmt.Errorf("failed to list server class pricing: %w", err)
	}

	return pricingList, nil
}
