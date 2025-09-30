package rxtspot

import (
	"context"
	"fmt"
	"net/http"
)

// ListRegions retrieves all available regions.
func (c *RackspaceSpotClient) ListRegions(ctx context.Context) ([]Region, error) {
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/regions", c.BaseURL)

	var regions ListRegionsResponse
	if _, err := c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &regions); err != nil {
		return nil, c.handleAPIError(err, "region", "", "list")
	}
	var regionList []Region
	for _, item := range regions.Items {
		regionList = append(regionList, Region{
			Name:        item.Metadata.Name,
			Description: item.Spec.Description,
		})
	}
	return regionList, nil
}

// GetRegion retrieves a region by name.
func (c *RackspaceSpotClient) GetRegion(ctx context.Context, name string) (*Region, error) {
	regions, err := c.ListRegions(ctx)
	if err != nil {
		return nil, c.handleAPIError(err, "region", name, "get")
	}
	///fmt.Printf("Regions: %v\n", regions)
	for _, region := range regions {
		if region.Name == name {
			return &region, nil
		}
	}
	return nil, c.handleAPIError(fmt.Errorf("region '%s' not found", name), "region", name, "get")
}

func (c *RackspaceSpotClient) checkIfRegionExists(ctx context.Context, name string) (bool, error) {
	regions, err := c.ListRegions(ctx)
	if err != nil {
		return false, c.handleAPIError(err, "region", name, "check")
	}
	for _, region := range regions {
		if region.Name == name {
			return true, nil
		}
	}
	return false, nil
}
