package rxtspot

import (
	"context"
	"fmt"
)

// organizationManager implements the OrganizationManager interface
type organizationManager struct {
	client *RackspaceSpotClient
}

// ListOrganizations implements the OrganizationManager interface
func (m *organizationManager) ListOrganizations(ctx context.Context) ([]*Organization, error) {
	orgs, err := m.client.ListOrganizations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	result := make([]*Organization, len(orgs))
	for i, org := range orgs {
		result[i] = &Organization{
			ID:   org.ID,
			Name: org.Name,
		}
	}

	return result, nil
}

// cloudspaceManager implements the CloudspaceManager interface
type cloudspaceManager struct {
	client *RackspaceSpotClient
	org    string
}

// ListCloudspaces implements the CloudspaceManager interface
func (m *cloudspaceManager) ListCloudspaces(ctx context.Context, org string) ([]*CloudSpace, error) {
	list, err := m.client.ListCloudspaces(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("failed to list cloudspaces: %w", err)
	}
	// Convert []CloudSpace to []*CloudSpace
	result := make([]*CloudSpace, len(list.Items))
	for i := range list.Items {
		result[i] = &list.Items[i]
	}
	return result, nil
}

// CreateCloudspace implements the CloudspaceManager interface
func (m *cloudspaceManager) CreateCloudspace(ctx context.Context, cs *CloudSpace) error {
	return m.client.CreateCloudspace(ctx, *cs)
}

// GetCloudspace implements the CloudspaceManager interface
func (m *cloudspaceManager) GetCloudspace(ctx context.Context, org, name string) (*CloudSpace, error) {
	return m.client.GetCloudspace(ctx, org, name)
}

// DeleteCloudspace implements the CloudspaceManager interface
func (m *cloudspaceManager) DeleteCloudspace(ctx context.Context, org, name string) error {
	return m.client.DeleteCloudspace(ctx, org, name)
}

// GetCloudspaceConfig implements the CloudspaceManager interface
func (m *cloudspaceManager) GetCloudspaceConfig(ctx context.Context, org, name string) (string, error) {
	return m.client.GetCloudspaceConfig(ctx, org, name)
}

// regionsManager implements the RegionsManager interface
type regionsManager struct {
	client *RackspaceSpotClient
}

// nodePoolManager implements the NodePoolAPI interface
type spotNodePoolManager struct {
	client     *RackspaceSpotClient
	org        string
	cloudspace string
}

type ondemandNodePoolManager struct {
	client     *RackspaceSpotClient
	org        string
	cloudspace string
}

// ListRegions implements the RegionsManager interface
func (m *regionsManager) ListRegions(ctx context.Context) ([]*Region, error) {
	regions, err := m.client.ListRegions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list regions: %w", err)
	}

	result := make([]*Region, len(regions))
	for i, r := range regions {
		result[i] = &Region{
			Name:        r.Name,
			Description: r.Description,
		}
	}

	return result, nil
}

// GetRegion implements the RegionsManager interface
func (m *regionsManager) GetRegion(ctx context.Context, name string) (*Region, error) {
	return m.client.GetRegion(ctx, name)
}

// serverClassesManager implements the ServerClassesManager interface
type serverClassesManager struct {
	client *RackspaceSpotClient
	region string
}

// ListServerClasses implements the ServerClassesManager interface
func (m *serverClassesManager) ListServerClasses(ctx context.Context, region string) (*ServerClassList, error) {
	list, err := m.client.ListServerClasses(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("failed to list server classes: %w", err)
	}
	return list, nil
}

// GetServerClass implements the ServerClassesManager interface
func (m *serverClassesManager) GetServerClass(ctx context.Context, name string) (*ServerClass, error) {
	return m.client.GetServerClass(ctx, name)
}

// pricingManager implements the PricingManager interface
type pricingManager struct {
	client *RackspaceSpotClient
}

// GetPriceDetailsForServerClass implements the PricingManager interface
func (m *pricingManager) GetPriceDetailsForServerClass(ctx context.Context, serverClass string) (*PriceDetails, error) {
	return m.client.GetPriceDetailsForServerClass(ctx, serverClass)
}

// GetPriceDetails implements the PricingManager interface
func (m *pricingManager) GetPriceDetails(ctx context.Context) ([]*PriceDetails, error) {
	return m.client.GetPriceDetails(ctx)
}

// GetPriceDetailsForRegion implements the PricingManager interface
func (m *pricingManager) GetPriceDetailsForRegion(ctx context.Context, region string) (*PriceDetails, error) {
	return m.client.GetPriceDetailsForRegion(ctx, region)
}
