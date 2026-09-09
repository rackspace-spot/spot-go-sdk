package rxtspot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ListAutopilotNodePools retrieves all autopilot node pools in a namespace, optionally
// scoped to a cloudspace.
func (c *RackspaceSpotClient) ListAutopilotNodePools(ctx context.Context, org, cloudspaceName string) ([]*AutopilotNodePool, error) {
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

	reqURL := fmt.Sprintf(
		"%s/apis/ngpc.rxt.io/v1/namespaces/%s/autopilotnodepools?labelSelector=%s",
		c.BaseURL, orgID, encodedSelector,
	)

	var pool AutopilotNodePoolListResponse
	if err := c.doRequest(ctx, http.MethodGet, reqURL, nil, c.authHeader(), &pool); err != nil {
		return nil, c.handleAPIError(err, "autopilot node pool", cloudspaceName, "list")
	}

	var finalList []*AutopilotNodePool
	for _, item := range pool.Items {
		finalList = append(finalList, autopilotNodePoolFromReadResponse(org, item.Metadata, item.Spec, item.Status))
	}
	return finalList, nil
}

// CreateAutopilotNodePool creates a new autopilot node pool in the given namespace.
func (c *RackspaceSpotClient) CreateAutopilotNodePool(ctx context.Context, org string, pool AutopilotNodePool) error {
	if err := ValidateOrgName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(pool.Name); err != nil {
		return fmt.Errorf("invalid autopilot node pool name: %w", err)
	}
	if pool.VCPUTotal <= 0 {
		return fmt.Errorf("vcpu total must be greater than 0")
	}
	if err := ValidateBidPrice(pool.BudgetPerHour); err != nil {
		return fmt.Errorf("invalid budget per hour: %w", err)
	}

	exists, orgID, err := c.getOrgIDIFExists(ctx, org)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", org)
	}
	reqURL := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/autopilotnodepools", c.BaseURL, orgID)

	createRequestBody := AutopilotNodePoolRequestBody{
		APIVersion: "ngpc.rxt.io/v1",
		Kind:       "AutopilotNodePool",
		Metadata: ObjectMeta{
			Name:      pool.Name,
			Namespace: orgID,
			Labels: map[string]string{
				"ngpc.rxt.io/cloudspace": pool.Cloudspace,
			},
		},
		Spec: AutopilotNodePoolSpec{
			Region:             pool.Region,
			CloudSpace:         pool.Cloudspace,
			VCPU:               VCPUTargetRO{Total: pool.VCPUTotal},
			VCPUPerNode:        VCPUPerNodeRangeRO{Min: pool.VCPUPerNodeMin, Max: pool.VCPUPerNodeMax},
			MemoryPerVCPU:      pool.MemoryPerVCPU,
			BudgetPerHour:      pool.BudgetPerHour,
			AllocationStrategy: pool.AllocationStrategy,
			CustomAnnotations:  pool.CustomAnnotations,
			CustomLabels:       pool.CustomLabels,
			CustomTaints:       pool.CustomTaints,
		},
	}

	body, err := json.Marshal(createRequestBody)
	if err != nil {
		return err
	}

	err = c.doRequest(ctx, http.MethodPost, reqURL, body, c.authHeader(), nil)
	return c.handleAPIError(err, "autopilot node pool", pool.Name, "create")
}

// UpdateAutopilotNodePool updates an autopilot node pool in the given namespace. Only the
// mutable fields (region and cloudSpace are immutable) are sent.
func (c *RackspaceSpotClient) UpdateAutopilotNodePool(ctx context.Context, org string, pool AutopilotNodePool) error {
	if err := ValidateOrgName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(pool.Name); err != nil {
		return fmt.Errorf("invalid autopilot node pool name: %w", err)
	}
	if pool.BudgetPerHour != "" {
		if err := ValidateBidPrice(pool.BudgetPerHour); err != nil {
			return fmt.Errorf("invalid budget per hour: %w", err)
		}
	}

	exists, orgID, err := c.getOrgIDIFExists(ctx, org)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", org)
	}
	reqURL := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/autopilotnodepools/%s", c.BaseURL, orgID, pool.Name)

	updateBody := AutopilotNodePoolUpdateRequestBody{
		Spec: AutopilotNodePoolUpdateSpec{
			VCPU:               VCPUTargetRO{Total: pool.VCPUTotal},
			VCPUPerNode:        VCPUPerNodeRangeRO{Min: pool.VCPUPerNodeMin, Max: pool.VCPUPerNodeMax},
			MemoryPerVCPU:      pool.MemoryPerVCPU,
			BudgetPerHour:      pool.BudgetPerHour,
			AllocationStrategy: pool.AllocationStrategy,
			CustomAnnotations:  pool.CustomAnnotations,
			CustomLabels:       pool.CustomLabels,
			CustomTaints:       pool.CustomTaints,
		},
	}

	body, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("failed to marshal update body: %w", err)
	}

	var respBody interface{}
	err = c.doRequest(ctx, http.MethodPatch, reqURL, body, c.authHeader(), &respBody)
	return c.handleAPIError(err, "autopilot node pool", pool.Name, "update")
}

// DeleteAutopilotNodePool deletes an autopilot node pool by name in the given namespace.
func (c *RackspaceSpotClient) DeleteAutopilotNodePool(ctx context.Context, org, name string) error {
	if err := ValidateOrgName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(name); err != nil {
		return fmt.Errorf("invalid autopilot node pool name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIFExists(ctx, org)
	if err != nil {
		return c.handleAPIError(err, "organization", org, "find")
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", org)
	}
	reqURL := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/autopilotnodepools/%s", c.BaseURL, orgID, name)

	err = c.doRequest(ctx, http.MethodDelete, reqURL, nil, c.authHeader(), nil)
	return c.handleAPIError(err, "autopilot node pool", name, "delete")
}

// GetAutopilotNodePool retrieves an autopilot node pool by name in the given namespace.
func (c *RackspaceSpotClient) GetAutopilotNodePool(ctx context.Context, org, name string) (*AutopilotNodePool, error) {
	if err := ValidateOrgName(org); err != nil {
		return nil, fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(name); err != nil {
		return nil, fmt.Errorf("invalid autopilot node pool name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIFExists(ctx, org)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("organization '%s' not found", org)
	}
	reqURL := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/autopilotnodepools/%s", c.BaseURL, orgID, name)

	var interm AutopilotNodePoolGetResponse
	if err := c.doRequest(ctx, http.MethodGet, reqURL, nil, c.authHeader(), &interm); err != nil {
		return nil, c.handleAPIError(err, "autopilot node pool", name, "get")
	}

	return autopilotNodePoolFromReadResponse(org, interm.Metadata, interm.Spec, interm.Status), nil
}

// autopilotNodePoolFromReadResponse maps the wire read shape (shared by Get and List) into the
// public AutopilotNodePool struct.
func autopilotNodePoolFromReadResponse(org string, metadata ResourceMetadataWithTimestamp, spec AutopilotNodePoolSpecReadOnly, status AutopilotNodePoolStatus) *AutopilotNodePool {
	labels := metadata.Labels
	cloudspaceName := spec.CloudSpace
	if cloudspaceName == "" {
		cloudspaceName = labels["ngpc.rxt.io/cloudspace"]
	}

	allocations := make([]AutopilotAllocation, 0, len(status.Allocations))
	for _, a := range status.Allocations {
		allocations = append(allocations, AutopilotAllocation{
			ServerClass:        a.ServerClass,
			SpotNodePool:       a.SpotNodePool,
			MarketPricePerHour: a.MarketPricePerHour,
			BidPricePerHour:    a.BidPricePerHour,
			VCPUPerNode:        a.VCPUPerNode,
			MemoryGBPerNode:    a.MemoryGBPerNode,
			DesiredNodes:       a.DesiredNodes,
			WonNodes:           a.WonNodes,
		})
	}

	classStates := make([]AutopilotClassState, 0, len(status.ClassStates))
	for _, s := range status.ClassStates {
		classStates = append(classStates, AutopilotClassState{
			ServerClass:     s.ServerClass,
			LastOutcome:     s.LastOutcome,
			LastAttemptTime: s.LastAttemptTime,
			CooldownUntil:   s.CooldownUntil,
		})
	}

	history := make([]AutopilotAllocationRecord, 0, len(status.AllocationsHistory))
	for _, r := range status.AllocationsHistory {
		history = append(history, AutopilotAllocationRecord{
			Time:            r.Time,
			ServerClass:     r.ServerClass,
			SpotNodePool:    r.SpotNodePool,
			Event:           r.Event,
			Reason:          r.Reason,
			DesiredNodes:    r.DesiredNodes,
			WonNodes:        r.WonNodes,
			BidPricePerHour: r.BidPricePerHour,
		})
	}

	return &AutopilotNodePool{
		Name:               metadata.Name,
		CreationTimestamp:  metadata.CreationTimestamp,
		Org:                org,
		Cloudspace:         cloudspaceName,
		Region:             spec.Region,
		VCPUTotal:          spec.VCPU.Total,
		VCPUPerNodeMin:     spec.VCPUPerNode.Min,
		VCPUPerNodeMax:     spec.VCPUPerNode.Max,
		MemoryPerVCPU:      spec.MemoryPerVCPU,
		BudgetPerHour:      spec.BudgetPerHour,
		AllocationStrategy: spec.AllocationStrategy,
		CustomAnnotations:  spec.CustomAnnotations,
		CustomLabels:       spec.CustomLabels,
		CustomTaints:       spec.CustomTaints,
		Phase:              status.Phase,
		TargetVCPUs:        status.TargetVCPUs,
		ManagedVCPUs:       status.ManagedVCPUs,
		ManagedMemoryGB:    status.ManagedMemoryGB,
		Allocations:        allocations,
		ClassStates:        classStates,
		AllocationsHistory: history,
	}
}
