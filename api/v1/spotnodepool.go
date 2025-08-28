package rxtspot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ListSpotNodePools retrieves all spot node pools in a namespace.
func (c *RackspaceSpotClient) ListSpotNodePools(ctx context.Context, org, cloudspaceName string) ([]*SpotNodePool, error) {
	if err := ValidateOrgName(org); err != nil {
		return nil, fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(cloudspaceName); err != nil {
		return nil, fmt.Errorf("invalid cloudspace name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIfExists(ctx, org)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("organization '%s' not found", org)
	}

	labelKey := "ngpc.rxt.io/cloudspace"
	labelSelector := fmt.Sprintf("%s=%s", labelKey, cloudspaceName)
	encodedSelector := url.QueryEscape(labelSelector)

	url := fmt.Sprintf(
		"%s/apis/ngpc.rxt.io/v1/namespaces/%s/spotnodepools?labelSelector=%s",
		c.BaseURL, orgID, encodedSelector,
	)

	var pool SpotNodePoolListResponse
	if err := c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &pool); err != nil {
		return nil, c.handleAPIError(err, "spot node pool", cloudspaceName, "list")
	}

	var finalList []*SpotNodePool
	for _, item := range pool.Items {
		finalList = append(finalList, &SpotNodePool{
			Name:              item.Metadata.Name,
			CreationTimestamp: item.Metadata.CreationTimestamp,
			CustomAnnotations: item.Spec.CustomAnnotations,
			CustomLabels:      item.Spec.CustomLabels,
			CustomTaints:      item.Spec.CustomTaints,
			Org:               org,
			Cloudspace:        item.Spec.CloudSpace,
			ServerClass:       item.Spec.ServerClass,
			Desired:           item.Spec.Desired,
			BidPrice:          "$" + item.Spec.BidPrice,
			WonCount:          item.Status.WonCount,
			Status:            item.Status.BidStatus,
		})
	}
	return finalList, nil
}

// CreateSpotNodePool creates a new spot node pool in the given namespace.
func (c *RackspaceSpotClient) CreateSpotNodePool(ctx context.Context, org string, pool SpotNodePool) error {
	if err := ValidateOrgName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(pool.Name); err != nil {
		return fmt.Errorf("invalid node pool name: %w", err)
	}
	if err := ValidateBidPrice(pool.BidPrice); err != nil {
		return fmt.Errorf("invalid bid price: %w", err)
	}

	serverClass, err := c.GetServerClass(ctx, pool.ServerClass)
	if err != nil {
		return fmt.Errorf("invalid server class: %w", err)
	}
	if err := ValidateServerClass(*serverClass); err != nil {
		return fmt.Errorf("invalid server class: %w", err)
	}

	exists, orgID, err := c.getOrgIDIfExists(ctx, org)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", org)
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/spotnodepools", c.BaseURL, orgID)

	spotNodePoolCreateRequestBody := SpotNodePoolRequestBody{
		APIVersion: "ngpc.rxt.io/v1",
		Kind:       "SpotNodePool",
		Metadata: struct {
			Name      string            `json:"name"`
			Namespace string            `json:"namespace"`
			Labels    map[string]string `json:"labels"`
		}{
			Name:      pool.Name,
			Namespace: orgID,
			Labels: map[string]string{
				"ngpc.rxt.io/cloudspace": pool.Cloudspace,
			},
		},
		Spec: struct {
			ServerClass       string            `json:"serverClass"`
			Desired           int               `json:"desired"`
			BidPrice          string            `json:"bidPrice"`
			CloudSpace        string            `json:"cloudSpace"`
			CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`
			CustomLabels      map[string]string `json:"customLabels,omitempty"`
			CustomTaints      []interface{}     `json:"customTaints,omitempty"`
			Autoscaling       struct {
				Enabled  bool  `json:"enabled"`
				MinNodes int64 `json:"minNodes"`
				MaxNodes int64 `json:"maxNodes"`
			} `json:"autoscaling"`
		}{
			ServerClass:       pool.ServerClass,
			Desired:           pool.Desired,
			BidPrice:          pool.BidPrice,
			CloudSpace:        pool.Cloudspace,
			CustomAnnotations: pool.CustomAnnotations,
			CustomLabels:      pool.CustomLabels,
			CustomTaints:      pool.CustomTaints,
			Autoscaling: struct {
				Enabled  bool  `json:"enabled"`
				MinNodes int64 `json:"minNodes"`
				MaxNodes int64 `json:"maxNodes"`
			}{
				Enabled:  pool.Autoscaling.Enabled,
				MinNodes: pool.Autoscaling.MinNodes,
				MaxNodes: pool.Autoscaling.MaxNodes,
			},
		},
	}

	body, err := json.Marshal(spotNodePoolCreateRequestBody)
	if err != nil {
		return err
	}

	err = c.doRequest(ctx, http.MethodPost, url, body, c.authHeader(), nil)
	return c.handleAPIError(err, "spot node pool", pool.Name, "create")
}

// UpdateSpotNodePool updates a spot node pool in the given namespace.
func (c *RackspaceSpotClient) UpdateSpotNodePool(ctx context.Context, org string, pool SpotNodePool) error {
	if err := ValidateOrgName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(pool.Name); err != nil {
		return fmt.Errorf("invalid node pool name: %w", err)
	}
	if pool.BidPrice != "" {
		if err := ValidateBidPrice(pool.BidPrice); err != nil {
			return fmt.Errorf("invalid bid price: %w", err)
		}
	}

	exists, orgID, err := c.getOrgIDIfExists(ctx, org)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", org)
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/spotnodepools/%s", c.BaseURL, orgID, pool.Name)

	// Only include mutable fields in the update request
	updateBody := struct {
		Spec struct {
			Desired           int               `json:"desired,omitempty"`
			BidPrice          string            `json:"bidPrice,omitempty"`
			CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`
			CustomLabels      map[string]string `json:"customLabels,omitempty"`
			CustomTaints      []interface{}     `json:"customTaints,omitempty"`
			Autoscaling       struct {
				Enabled  bool  `json:"enabled"`
				MinNodes int64 `json:"minNodes,omitempty"`
				MaxNodes int64 `json:"maxNodes,omitempty"`
			} `json:"autoscaling"`
		} `json:"spec"`
	}{
		Spec: struct {
			Desired           int               `json:"desired,omitempty"`
			BidPrice          string            `json:"bidPrice,omitempty"`
			CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`
			CustomLabels      map[string]string `json:"customLabels,omitempty"`
			CustomTaints      []interface{}     `json:"customTaints,omitempty"`
			Autoscaling       struct {
				Enabled  bool  `json:"enabled"`
				MinNodes int64 `json:"minNodes,omitempty"`
				MaxNodes int64 `json:"maxNodes,omitempty"`
			} `json:"autoscaling"`
		}{
			Desired:           pool.Desired,
			BidPrice:          pool.BidPrice,
			CustomAnnotations: pool.CustomAnnotations,
			CustomLabels:      pool.CustomLabels,
			CustomTaints:      pool.CustomTaints,
			Autoscaling: struct {
				Enabled  bool  `json:"enabled"`
				MinNodes int64 `json:"minNodes,omitempty"`
				MaxNodes int64 `json:"maxNodes,omitempty"`
			}{
				Enabled:  pool.Autoscaling.Enabled,
				MinNodes: pool.Autoscaling.MinNodes,
				MaxNodes: pool.Autoscaling.MaxNodes,
			},
		},
	}

	body, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("failed to marshal update body: %w", err)
	}

	var respBody interface{}
	err = c.doRequest(ctx, http.MethodPatch, url, body, c.authHeader(), &respBody)
	return c.handleAPIError(err, "spot node pool", pool.Name, "update")
}

// DeleteSpotNodePool deletes a spot node pool by name in the given namespace.
func (c *RackspaceSpotClient) DeleteSpotNodePool(ctx context.Context, org, name string) error {
	if err := ValidateOrgName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(name); err != nil {
		return fmt.Errorf("invalid node pool name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIfExists(ctx, org)
	if err != nil {
		return c.handleAPIError(err, "organization", org, "find")
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", org)
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/spotnodepools/%s", c.BaseURL, orgID, name)

	err = c.doRequest(ctx, http.MethodDelete, url, nil, c.authHeader(), nil)
	return c.handleAPIError(err, "spot node pool", name, "delete")
}

// GetSpotNodePool retrieves a spot node pool by name in the given namespace.
func (c *RackspaceSpotClient) GetSpotNodePool(ctx context.Context, org, name string) (*SpotNodePool, error) {
	if err := ValidateOrgName(org); err != nil {
		return nil, fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(name); err != nil {
		return nil, fmt.Errorf("invalid node pool name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIfExists(ctx, org)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("organization '%s' not found", org)
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/spotnodepools/%s", c.BaseURL, orgID, name)

	var interm SpotNodePoolGetResponse
	if err := c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &interm); err != nil {
		return nil, c.handleAPIError(err, "spot node pool", name, "get")
	}

	labels := interm.Metadata.Labels
	cloudspaceName := labels["ngpc.rxt.io/cloudspace"]

	return &SpotNodePool{
		Name:              interm.Metadata.Name,
		CreationTimestamp: interm.Metadata.CreationTimestamp,
		Org:               org,
		Cloudspace:        cloudspaceName,
		ServerClass:       interm.Spec.ServerClass,
		CustomAnnotations: interm.Spec.CustomAnnotations,
		CustomLabels:      interm.Spec.CustomLabels,
		CustomTaints:      interm.Spec.CustomTaints,
		Desired:           interm.Spec.Desired,
		BidPrice:          "$" + interm.Spec.BidPrice,
		WonCount:          interm.Status.WonCount,
		Status:            interm.Status.BidStatus,
	}, nil
}
