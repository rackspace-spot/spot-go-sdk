package rxtspot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ListCloudspaces retrieves all cloudspaces in a namespace.
func (c *RackspaceSpotClient) ListCloudspaces(ctx context.Context, org string) (*CloudSpaceList, error) {
	if err := ValidateOrgName(org); err != nil {
		return nil, fmt.Errorf("invalid organization name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIfExists(ctx, org)
	if err != nil {
		return nil, c.handleAPIError(err, "organization", org, "find")
	}
	if !exists {
		return nil, fmt.Errorf("organization '%s' not found", org)
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/cloudspaces", c.BaseURL, orgID)

	// Pass &interm to be populated by doRequest JSON decoding
	var interm cloudSpaceListResponse
	err = c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &interm)
	if err != nil {
		return nil, c.handleAPIError(err, "cloudspaces", "", "list")
	}
	var finalList CloudSpaceList
	for _, cs := range interm.Items {
		// Spot node pools
		spotNodePools, err := c.ListSpotNodePools(ctx, org, cs.Metadata.Name)
		if err != nil {
			return nil, c.handleAPIError(err, "spot node pools", "", "list for cloudspace "+cs.Metadata.Name)
		}

		// OnDemand node pools
		onDemandNodePools, err := c.ListOnDemandNodePools(ctx, org, cs.Metadata.Name)
		if err != nil {
			return nil, c.handleAPIError(err, "on-demand node pools", "", "list for cloudspace "+cs.Metadata.Name)
		}
		finalList.Items = append(finalList.Items, CloudSpace{
			Name:                 cs.Metadata.Name,
			Org:                  org,
			CreationTimestamp:    cs.Metadata.CreationTimestamp,
			CNI:                  cs.Spec.CNI,
			DeploymentType:       cs.Spec.DeploymentType,
			GpuEnabled:           cs.Spec.GpuEnabled,
			KubernetesVersion:    cs.Spec.KubernetesVersion,
			Region:               cs.Spec.Region,
			PreemptionWebhookURL: cs.Spec.Webhook,
			APIServerEndpoint:    cs.Status.APIServerEndpoint,
			AssignedServers:      cs.Status.AssignedServers,
			SpotNodepools:        spotNodePools,
			OnDemandNodePools:    onDemandNodePools,
			Status:               cs.Status.Phase,
			Message:              cs.Status.Reason,
		})
	}
	return &finalList, nil
}

// CreateCloudspace creates a new cloudspace in the given namespace.
func (c *RackspaceSpotClient) CreateCloudspace(ctx context.Context, cs CloudSpace) error {
	if err := ValidateOrgName(cs.Org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(cs.Name); err != nil {
		return fmt.Errorf("invalid cloudspace name: %w", err)
	}
	if cs.Region == "" {
		return fmt.Errorf("region is required")
	}
	if cs.KubernetesVersion == "" {
		return fmt.Errorf("kubernetes version is required")
	}

	exists, orgID, err := c.getOrgIDIfExists(ctx, cs.Org)
	if err != nil {
		return c.handleAPIError(err, "organization", cs.Org, "find")
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", cs.Org)
	}

	cloudspaceCreateRequestBody := CloudSpaceCreateRequestBody{
		APIVersion: "ngpc.rxt.io/v1",
		Kind:       "CloudSpace",
		Metadata: struct {
			Name        string `json:"name"`
			Namespace   string `json:"namespace"`
			Annotations struct {
			} `json:"annotations"`
		}{
			Name:      cs.Name,
			Namespace: orgID,
			Annotations: struct {
			}{},
		},
		Spec: struct {
			DeploymentType    string `json:"deploymentType"`
			Cloud             string `json:"cloud"`
			Region            string `json:"region"`
			Webhook           string `json:"webhook"`
			CNI               string `json:"cni"`
			KubernetesVersion string `json:"kubernetesVersion"`
			HAControlPlane    bool   `json:"HAControlPlane"`
			GpuEnabled        bool   `json:"gpuEnabled"`
		}{
			DeploymentType:    "gen2",
			Cloud:             "default",
			Region:            cs.Region,
			Webhook:           cs.PreemptionWebhookURL,
			CNI:               cs.CNI,
			KubernetesVersion: cs.KubernetesVersion,
			HAControlPlane:    false,
			GpuEnabled:        cs.GpuEnabled,
		},
	}

	body, err := json.Marshal(cloudspaceCreateRequestBody)
	if err != nil {
		return c.handleAPIError(err, "cloudspace", cs.Name, "create")
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/cloudspaces", c.BaseURL, orgID)

	if err := c.doRequest(ctx, http.MethodPost, url, body, c.authHeader(), nil); err != nil {
		return c.handleAPIError(err, "cloudspace", cs.Name, "create")
	}
	return nil
}

// DeleteCloudspace deletes a cloudspace by name in the given namespace.
func (c *RackspaceSpotClient) DeleteCloudspace(ctx context.Context, org, name string) error {
	if err := ValidateOrgName(org); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(name); err != nil {
		return fmt.Errorf("invalid cloudspace name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIfExists(ctx, org)
	if err != nil {
		return c.handleAPIError(err, "organization", org, "find")
	}
	if !exists {
		return fmt.Errorf("organization '%s' not found", org)
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/cloudspaces/%s", c.BaseURL, orgID, name)
	err = c.doRequest(ctx, http.MethodDelete, url, nil, c.authHeader(), nil)
	return c.handleAPIError(err, "cloudspace", name, "delete")

}

// GetCloudspace retrieves a cloudspace by name in the given namespace.
func (c *RackspaceSpotClient) GetCloudspace(ctx context.Context, org, name string) (*CloudSpace, error) {
	if err := ValidateOrgName(org); err != nil {
		return nil, fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(name); err != nil {
		return nil, fmt.Errorf("invalid cloudspace name: %w", err)
	}

	exists, orgID, err := c.getOrgIDIfExists(ctx, org)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("organization '%s' not found", org)
	}
	url := fmt.Sprintf("%s/apis/ngpc.rxt.io/v1/namespaces/%s/cloudspaces/%s", c.BaseURL, orgID, name)
	var interm cloudSpaceGetResponse

	// Pass &interm so doRequest will JSON unmarshal into it
	err = c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &interm)
	if err != nil {
		return nil, c.handleAPIError(err, "cloudspace", name, "get")
	}
	// Spot node pools
	spotNodePools, err := c.ListSpotNodePools(ctx, org, interm.Metadata.Name)
	if err != nil {
		return nil, c.handleAPIError(err, "spot node pool", name, "list for cloudspace "+interm.Metadata.Name)
	}

	// OnDemand node pools
	onDemandNodePools, err := c.ListOnDemandNodePools(ctx, org, interm.Metadata.Name)
	if err != nil {
		return nil, c.handleAPIError(err, "on-demand node pool", name, "list for cloudspace "+interm.Metadata.Name)
	}

	finalList := CloudSpace{
		Name:                 interm.Metadata.Name,
		Org:                  org,
		CreationTimestamp:    interm.Metadata.CreationTimestamp,
		CNI:                  interm.Spec.CNI,
		DeploymentType:       interm.Spec.DeploymentType,
		GpuEnabled:           interm.Spec.GpuEnabled,
		KubernetesVersion:    interm.Spec.KubernetesVersion,
		Region:               interm.Spec.Region,
		PreemptionWebhookURL: interm.Spec.Webhook,
		APIServerEndpoint:    interm.Status.APIServerEndpoint,
		AssignedServers:      interm.Status.AssignedServers,
		SpotNodepools:        spotNodePools,
		OnDemandNodePools:    onDemandNodePools,
		Status:               interm.Status.Phase,
		Message:              interm.Status.Reason,
	}

	return &finalList, nil
}

func (c *RackspaceSpotClient) GetCloudspaceConfig(ctx context.Context, namespace, name string) (string, error) {
	if err := ValidateOrgName(namespace); err != nil {
		return "", fmt.Errorf("invalid organization name: %w", err)
	}
	if err := ValidateResourceName(name); err != nil {
		return "", fmt.Errorf("invalid cloudspace name: %w", err)
	}
	if c.RefreshToken == "" {
		return "", fmt.Errorf("refresh token is required")
	}
	url := fmt.Sprintf("%s/apis/auth.ngpc.rxt.io/v1/generate-kubeconfig", c.BaseURL)
	reqBody := struct {
		OrganizationName string `json:"organization_name"`
		CloudspaceName   string `json:"cloudspace_name"`
		RefreshToken     string `json:"refresh_token"`
	}{
		OrganizationName: namespace,
		CloudspaceName:   name,
		RefreshToken:     c.RefreshToken, // actual token
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", c.handleAPIError(err, "cloudspace", name, "get kubeconfig")
	}
	var kubeConfigResponse KubeConfigResponse
	if err := c.doRequest(ctx, http.MethodPost, url, jsonBody, c.authHeader(), &kubeConfigResponse); err != nil {
		return "", c.handleAPIError(err, "cloudspace", name, "get kubeconfig")
	}
	return kubeConfigResponse.Data.Kubeconfig, nil
}
