package rxtspot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ListOnDemandNodePools retrieves all on-demand node pools in a namespace.
func (c *RackspaceSpotClient) ListOnDemandNodePools(ctx context.Context, org, cloudspaceName string) ([]*OnDemandNodePool, error) {
	if err := ValidateOrgName(org); err != nil {
		return nil, fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(cloudspaceName); err != nil {
		return nil, fmt.Errorf("invalid cloudspace name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIFExists(ctx, org)
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
		"%s/apis/ngpc.rxt.io/v1/namespaces/%s/ondemandnodepools?labelSelector=%s",
		c.BaseURL, orgID, encodedSelector,
	)

	var pool OnDemandNodePoolListResponse
	if _, err := c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &pool); err != nil {
		return nil, c.handleAPIError(err, "ondemand node pool", cloudspaceName, "list")
	}

	var finalList []*OnDemandNodePool
	for _, item := range pool.Items {
		var onDemandPoolcost string
		serverClass, err := c.GetServerClass(ctx, item.Spec.ServerClass)
		if err != nil {
			onDemandPoolcost = "NA"
		}
		onDemandPoolcost = serverClass.OnDemandPricePerHour
		finalList = append(finalList, &OnDemandNodePool{
			Name:                 item.Metadata.Name,
			CreationTimestamp:    item.Metadata.CreationTimestamp,
			CustomAnnotations:    item.Spec.CustomAnnotations,
			CustomLabels:         item.Spec.CustomLabels,
			CustomTaints:         item.Spec.CustomTaints,
			Org:                  org,
			Cloudspace:           item.Spec.CloudSpace,
			ServerClass:          item.Spec.ServerClass,
			Desired:              item.Spec.Desired,
			WonCount:             item.Status.ReservedCount,
			Status:               item.Status.ReservedStatus,
			OnDemandPricePerHour: onDemandPoolcost,
		})
	}
	return finalList, nil
}

// CreateOnDemandNodePool creates a new on-demand node pool in the given namespace.
func (c *RackspaceSpotClient) CreateOnDemandNodePool(ctx context.Context, org string, pool OnDemandNodePool) error {
	if err := ValidateOrgName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(pool.Name); err != nil {
		return fmt.Errorf("invalid node pool name: %w", err)
	}
	serverClass, err := c.GetServerClass(ctx, pool.ServerClass)
	if err != nil {
		return fmt.Errorf("invalid server class: %w", err)
	}
	if err := ValidateServerClass(*serverClass); err != nil {
		return fmt.Errorf("invalid server class: %w", err)
	}
	exists, orgID, err := c.getOrgIDIFExists(ctx, org)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", org)
	}

	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/ondemandnodepools", c.BaseURL, orgID)

	ondemandNodePoolCreateRequestBody := OnDemandNodePoolCreateRequestBody{
		APIVersion: "ngpc.rxt.io/v1",
		Kind:       "OnDemandNodePool",
		Metadata: ObjectMeta{
			Name:      pool.Name,
			Namespace: orgID,
			Labels: map[string]string{
				"ngpc.rxt.io/cloudspace": pool.Cloudspace,
			},
		},
		Spec: OnDemandNodePoolSpec{
			CommonNodePoolSpec: CommonNodePoolSpec{
				ServerClass:       pool.ServerClass,
				Desired:           pool.Desired,
				CloudSpace:        pool.Cloudspace,
				CustomAnnotations: pool.CustomAnnotations,
				CustomLabels:      pool.CustomLabels,
				CustomTaints:      pool.CustomTaints,
			},
			Autoscaling: AutoscalingAny{
				Enabled:  pool.Autoscaling.Enabled,
				MinNodes: pool.Autoscaling.MinNodes,
				MaxNodes: pool.Autoscaling.MaxNodes,
			},
		},
	}

	body, err := json.Marshal(ondemandNodePoolCreateRequestBody)
	if err != nil {
		return err
	}

	_, err = c.doRequest(ctx, http.MethodPost, url, body, c.authHeader(), nil)
	if err != nil {
		return c.handleAPIError(err, "ondemand node pool", pool.Name, "create")
	}
	return nil
}

// DeleteOnDemandNodePool deletes an on-demand node pool by name in the given namespace.
func (c *RackspaceSpotClient) DeleteOnDemandNodePool(ctx context.Context, org, name string) error {
	if err := ValidateOrgName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(name); err != nil {
		return fmt.Errorf("invalid node pool name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIFExists(ctx, org)
	if err != nil {
		return c.handleAPIError(err, "organization", org, "find")
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", org)
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/ondemandnodepools/%s", c.BaseURL, orgID, name)

	_, err = c.doRequest(ctx, http.MethodDelete, url, nil, c.authHeader(), nil)
	if err != nil {
		return c.handleAPIError(err, "ondemand node pool", name, "delete")
	}

	return nil
}

// GetOnDemandNodePool retrieves an on-demand node pool by name in the given namespace.
func (c *RackspaceSpotClient) GetOnDemandNodePool(ctx context.Context, org, name string) (*OnDemandNodePool, error) {
	if err := ValidateOrgName(org); err != nil {
		return nil, fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(name); err != nil {
		return nil, fmt.Errorf("invalid node pool name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIFExists(ctx, org)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("organization '%s' not found", org)
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/ondemandnodepools/%s", c.BaseURL, orgID, name)

	var interm OnDemandNodePoolGetResponse
	if err := c.doRequestJSON(ctx, http.MethodGet, url, nil, &interm, c.authHeader()); err != nil {
		return nil, c.handleAPIError(err, "ondemand node pool", name, "get")
	}

	serverClass, err := c.GetServerClass(ctx, interm.Spec.ServerClass)
	if err != nil {
		return nil, err
	}

	labels := interm.Metadata.Labels
	cloudspaceName := labels["ngpc.rxt.io/cloudspace"]

	return &OnDemandNodePool{
		Name:                 interm.Metadata.Name,
		CreationTimestamp:    interm.Metadata.CreationTimestamp,
		CustomAnnotations:    interm.Spec.CustomAnnotations,
		CustomLabels:         interm.Spec.CustomLabels,
		CustomTaints:         interm.Spec.CustomTaints,
		Org:                  org,
		Cloudspace:           cloudspaceName,
		ServerClass:          interm.Spec.ServerClass,
		Desired:              interm.Spec.Desired,
		WonCount:             interm.Status.ReservedCount,
		Status:               interm.Status.ReservedStatus,
		OnDemandPricePerHour: serverClass.OnDemandPricePerHour,
	}, nil
}

func (c *RackspaceSpotClient) UpdateOnDemandNodePool(ctx context.Context, org string, pool OnDemandNodePool) error {
	if err := ValidateOrgName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(pool.Name); err != nil {
		return fmt.Errorf("invalid node pool name: %w", err)
	}

	if pool.Name == "" {
		return fmt.Errorf("name must be provided")
	}

	exists, orgID, err := c.getOrgIDIFExists(ctx, org)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", org)
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/ondemandnodepools/%s", c.BaseURL, orgID, pool.Name)

	// Only include mutable fields in the update request
	updateBody := OnDemandNodePoolUpdateRequestBody{
		Spec: OnDemandNodePoolUpdateSpec{
			Desired:           pool.Desired,
			CustomAnnotations: pool.CustomAnnotations,
			CustomLabels:      pool.CustomLabels,
			CustomTaints:      pool.CustomTaints,
			Autoscaling: AutoscalingInt64Update{
				Enabled:  pool.Autoscaling.Enabled,
				MinNodes: int64(pool.Autoscaling.MinNodes),
				MaxNodes: int64(pool.Autoscaling.MaxNodes),
			},
		},
	}

	body, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("failed to marshal update body: %w", err)
	}

	var respBody interface{}
	_, err = c.doRequest(ctx, http.MethodPatch, url, body, c.authHeader(), &respBody)
	return c.handleAPIError(err, "ondemand node pool", pool.Name, "update")

}
