package rxtspot

import (
	"context"
	"fmt"
	"net/http"
)

// ListServerClasses retrieves all available server classes.
func (c *RackspaceSpotClient) ListServerClasses(ctx context.Context, region string) (*ServerClassList, error) {
	if region != "" {
		exists, err := c.checkIfRegionExists(ctx, region)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("region '%s' not found", region)
		}
	}

	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/serverclasses", c.BaseURL)

	var interm ListServerClassesResponse
	if err := c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &interm); err != nil {
		return nil, c.handleAPIError(err, "server class", "", "list")
	}
	var serverclasses []ServerClass
	if region != "" {
		for _, item := range interm.Items {
			if item.Spec.Region == region {
				marketPrice, err := c.GetMarketPriceForServerClass(ctx, item.Metadata.Name)
				if err != nil {
					// If market price is not found, set to "N/A" and continue
					marketPrice = "N/A"
				}
				serverclasses = append(serverclasses, ServerClass{
					Availability:              item.Spec.Availability,
					Name:                      item.Metadata.Name,
					Category:                  item.Spec.Category,
					Region:                    item.Spec.Region,
					CurrentMarketPricePerHour: marketPrice,
					MinBidPricePerHour:        "$" + item.Spec.MinBidPricePerHour,
					OnDemandPricePerHour:      "$" + item.Spec.OnDemandPricing.Cost,
					Resources: Resource{
						CPU:    item.Spec.Resources.CPU,
						Memory: item.Spec.Resources.Memory,
					},
				})
			}
		}
		return &ServerClassList{Items: serverclasses}, nil
	}

	for _, item := range interm.Items {
		marketPrice, err := c.GetMarketPriceForServerClass(ctx, item.Metadata.Name)
		if err != nil {
			// If market price is not found, set to "N/A" and continue
			marketPrice = "N/A"
		}
		serverclasses = append(serverclasses, ServerClass{
			Availability:              item.Spec.Availability,
			Name:                      item.Metadata.Name,
			Category:                  item.Spec.Category,
			Region:                    item.Spec.Region,
			CurrentMarketPricePerHour: marketPrice,
			MinBidPricePerHour:        "$" + item.Spec.MinBidPricePerHour,
			OnDemandPricePerHour:      "$" + item.Spec.OnDemandPricing.Cost,
			Resources: Resource{
				CPU:    item.Spec.Resources.CPU,
				Memory: item.Spec.Resources.Memory,
			},
		})
	}
	return &ServerClassList{Items: serverclasses}, nil
}

// GetServerClass retrieves a server class by name.
func (c *RackspaceSpotClient) GetServerClass(ctx context.Context, name string) (*ServerClass, error) {
	if err := ValidateResourceName(name); err != nil {
		return nil, fmt.Errorf("invalid server class name: %w", err)
	}

	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/serverclasses/%s", c.BaseURL, name)

	var interm GetServerClassResponse
	if err := c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &interm); err != nil {
		return nil, c.handleAPIError(err, "server class", name, "get")
	}
	marketPrice, err := c.GetMarketPriceForServerClass(ctx, name)
	if err != nil {
		return nil, err
	}

	serverclass := ServerClass{
		Availability:              interm.Spec.Availability,
		Name:                      interm.Metadata.Name,
		Category:                  interm.Spec.Category,
		Region:                    interm.Spec.Region,
		MinBidPricePerHour:        "$" + interm.Spec.MinBidPricePerHour,
		CurrentMarketPricePerHour: marketPrice,
		OnDemandPricePerHour:      "$" + interm.Spec.OnDemandPricing.Cost,
		Resources: Resource{
			CPU:    interm.Spec.Resources.CPU,
			Memory: interm.Spec.Resources.Memory,
		},
	}
	return &serverclass, nil
}
