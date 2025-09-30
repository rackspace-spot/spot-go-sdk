package rxtspot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/klog/v2"
)

type ServerData struct {
	Regions map[string]RegionPricingDetails `json:"regions"`
}

type RegionPricingDetails struct {
	Generation    string                               `json:"generation"`
	ServerClasses map[string]ServerClassPricingDetails `json:"serverclasses"`
}

type ServerClassPricingDetails struct {
	Percentile20 float64 `json:"20_percentile"`
	Percentile50 float64 `json:"50_percentile"`
	Percentile80 float64 `json:"80_percentile"`
	MarketPrice  string  `json:"market_price"`
	CPU          string  `json:"cpu"`
	Memory       string  `json:"memory"`
	DisplayName  string  `json:"display_name"`
	Category     string  `json:"category"`
	Description  string  `json:"description"`
}

// GetPriceDetails retrieves the price details for a server class.
func (c *RackspaceSpotClient) GetPriceDetails(ctx context.Context) ([]*PriceDetails, error) {

	var serverData ServerData
	respBody, err := c.doRequest(ctx, http.MethodGet, PriceDetailsURL, nil, nil, &serverData)
	if err != nil {
		return nil, c.handleAPIError(err, "server class", "", "get price details")
	}
	if err := json.Unmarshal(respBody, &serverData); err != nil {
		klog.Errorf("Failed to unmarshal pricing data: %v", err)
		klog.V(4).Infof("Response that failed to unmarshal: %s", string(respBody))
		return nil, fmt.Errorf("failed to parse pricing data: %w", err)
	}
	var completePriceDetails []*PriceDetails

	for region, details := range serverData.Regions {
		for serverClassName, pricingDetails := range details.ServerClasses {
			completePriceDetails = append(completePriceDetails, &PriceDetails{
				ServerClassName: serverClassName,
				Region:          region,
				MarketPrice:     "$" + pricingDetails.MarketPrice,
				CPU:             pricingDetails.CPU,
				Memory:          pricingDetails.Memory,
				DisplayName:     pricingDetails.DisplayName,
				Category:        pricingDetails.Category,
			})
		}
	}
	return completePriceDetails, nil
}

func (c *RackspaceSpotClient) GetPriceDetailsForServerClass(ctx context.Context, serverClass string) (*PriceDetails, error) {

	var serverData ServerData
	respBody, err := c.doRequest(ctx, http.MethodGet, PriceDetailsURL, nil, nil, &serverData)
	if err != nil {
		return nil, c.handleAPIError(err, "server class", serverClass, "get price details")
	}
	if err := json.Unmarshal(respBody, &serverData); err != nil {
		klog.Errorf("Failed to unmarshal pricing data: %v", err)
		klog.V(4).Infof("Response that failed to unmarshal: %s", string(respBody))
		return nil, fmt.Errorf("failed to parse pricing data: %w", err)
	}

	if serverData.Regions == nil {
		klog.Error("serverData.Regions is nil - check if the API response format matches the expected structure")
		return nil, fmt.Errorf("invalid pricing data format: no regions data")
	}

	var priceDetails PriceDetails

	for region, details := range serverData.Regions {
		for serverClassName, pricingDetails := range details.ServerClasses {
			if serverClassName == serverClass {
				priceDetails = PriceDetails{
					ServerClassName: serverClassName,
					Region:          region,
					MarketPrice:     "$" + pricingDetails.MarketPrice,
					CPU:             pricingDetails.CPU,
					Memory:          pricingDetails.Memory,
					DisplayName:     pricingDetails.DisplayName,
					Category:        pricingDetails.Category,
				}
				return &priceDetails, nil

			}
		}
	}
	return nil, fmt.Errorf("server class '%s' not found", serverClass)
}

func (c *RackspaceSpotClient) GetPriceDetailsForRegion(ctx context.Context, regionName string) (*PriceDetails, error) {
	var serverData ServerData
	respBody, err := c.doRequest(ctx, http.MethodGet, PriceDetailsURL, nil, nil, &serverData)
	if err != nil {
		return nil, c.handleAPIError(err, "region", regionName, "get price details")
	}
	if err := json.Unmarshal(respBody, &serverData); err != nil {
		klog.Errorf("Failed to unmarshal pricing data: %v", err)
		klog.V(4).Infof("Response that failed to unmarshal: %s", string(respBody))
		return nil, fmt.Errorf("failed to parse pricing data: %w", err)
	}

	for region, details := range serverData.Regions {
		if region == regionName {
			// Return the first server class found for the region
			for serverClassName, pricingDetails := range details.ServerClasses {
				priceDetails := &PriceDetails{
					ServerClassName: serverClassName,
					Region:          region,
					MarketPrice:     "$" + pricingDetails.MarketPrice,
					CPU:             pricingDetails.CPU,
					Memory:          pricingDetails.Memory,
					DisplayName:     pricingDetails.DisplayName,
					Category:        pricingDetails.Category,
				}
				return priceDetails, nil
			}
			return nil, fmt.Errorf("no server classes found for region '%s'", regionName)
		}
	}
	return nil, fmt.Errorf("region '%s' not found", regionName)
}

func (c *RackspaceSpotClient) GetMarketPriceForServerClass(ctx context.Context, serverClass string) (string, error) {
	var serverData ServerData

	// First, make the request and get the raw response
	respBody, err := c.doRequest(ctx, http.MethodGet, PriceDetailsURL, nil, nil, nil)
	if err != nil {
		klog.Errorf("Failed to fetch price details: %v", err)
		return "", c.handleAPIError(err, "server class", serverClass, "get market price")
	}
	// Try to unmarshal the response manually
	if err := json.Unmarshal(respBody, &serverData); err != nil {
		klog.Errorf("Failed to unmarshal pricing data: %v", err)
		return "", fmt.Errorf("failed to parse pricing data: %w", err)
	}

	// Check if regions is nil or empty
	if serverData.Regions == nil {
		klog.Error("serverData.Regions is nil - check if the API response format matches the expected structure")
		return "", fmt.Errorf("invalid pricing data format: no regions data")
	}

	for _, details := range serverData.Regions {
		for className, pricing := range details.ServerClasses {
			if className == serverClass {
				return "$" + pricing.MarketPrice, nil
			}
		}
	}

	return "", fmt.Errorf("server class '%s' not found", serverClass)
}

func (c *RackspaceSpotClient) GetMinimumBidPriceForServerClass(ctx context.Context, serverClass string) (string, error) {
	ServerClassDetails, _ := c.GetServerClass(ctx, serverClass)
	return ServerClassDetails.MinBidPricePerHour, nil
}
