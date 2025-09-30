package rxtspot

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ListOrganizations retrieves all organizations accessible by the user.
func (c *RackspaceSpotClient) ListOrganizations(ctx context.Context) ([]Organization, error) {
	url := fmt.Sprintf("%s/apis/auth.ngpc.rxt.io/v1/organizations", c.BaseURL) // Correct URL

	// Structure for decoding
	var response struct {
		Organizations []Organization `json:"organizations"`
	}

	// Pass &response to doRequest so it decodes automatically
	_, err := c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &response)
	if err != nil {
		return nil, c.handleAPIError(err, "organization", "", "list")
	}
	return response.Organizations, nil
}

func (c *RackspaceSpotClient) getOrgIDIFExists(ctx context.Context, orgName string) (bool, string, error) {
	url := fmt.Sprintf("%s/apis/auth.ngpc.rxt.io/v1/organizations", c.BaseURL) // Correct URL

	// Structure for decoding
	var response struct {
		Organizations []Organization `json:"organizations"`
	}

	// Pass &response to doRequest so it decodes automatically
	_, err := c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &response)
	if err != nil {
		return false, "", c.handleAPIError(err, "organization", orgName, "find")
	}

	for _, org := range response.Organizations {
		if org.Name == orgName {
			org.ID = strings.ReplaceAll(org.ID, "_", "-")
			org.ID = strings.ToLower(org.ID)
			return true, org.ID, nil
		}
	}
	return false, "", nil
}

func (c *RackspaceSpotClient) getOrgIDIFExistsWithoutNormalizing(ctx context.Context, orgName string) (bool, string, error) {
	url := fmt.Sprintf("%s/apis/auth.ngpc.rxt.io/v1/organizations", c.BaseURL) // Correct URL

	// Structure for decoding
	var response struct {
		Organizations []Organization `json:"organizations"`
	}

	// Pass &response to doRequest so it decodes automatically
	_, err := c.doRequest(ctx, http.MethodGet, url, nil, c.authHeader(), &response)
	if err != nil {
		return false, "", c.handleAPIError(err, "organization", orgName, "find")
	}

	for _, org := range response.Organizations {
		if org.Name == orgName {
			return true, org.ID, nil
		}
	}
	return false, "", nil
}
