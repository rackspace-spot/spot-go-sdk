// Basic example demonstrating the usage of the Spot Go SDK
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rackspace/spot-go-sdk/api/v1"
)

func main() {
	// Initialize the client
	cfg := v1.ClientConfig{
		BaseURL:    "https://api.spot.io", // Replace with actual API URL
		AuthToken:  os.Getenv("SPOT_AUTH_TOKEN"),
		HTTPClient: nil, // Uses default client if nil
	}

	client, err := v1.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Example: List all regions
	fmt.Println("Listing all regions:")
	regions, err := client.Regions().List(context.Background(), "", nil)
	if err != nil {
		log.Fatalf("Failed to list regions: %v", err)
	}
	for _, r := range regions {
		fmt.Printf("- %s (%s)\n", r.Name, r.Provider)
	}

	// Example: List server classes for a specific region
	if len(regions) > 0 {
		region := regions[0].Name
		fmt.Printf("\nListing server classes for region %s:\n", region)
		
		serverClasses, err := client.ServerClasses().List(context.Background(), region)
		if err != nil {
			log.Fatalf("Failed to list server classes: %v", err)
		}
		
		for _, sc := range serverClasses {
			fmt.Printf("- %s (%s) - %s CPU, %s Memory\n", 
				sc.Name, 
				sc.Category,
				sc.Resources.CPU,
				sc.Resources.Memory,
			)
		}
	}

	// Example: Create a cloudspace (commented out to prevent accidental creation)
	/*
	createOpts := &cloudspace.CreateOptions{
		Name:        "my-cloudspace",
		Org:         "my-org",
		Region:      "us-east-1",
		ServerClass: "general-purpose-1",
	}

	createdCS, err := client.Cloudspaces().Create(context.Background(), createOpts)
	if err != nil {
		log.Fatalf("Failed to create cloudspace: %v", err)
	}
	fmt.Printf("\nCreated cloudspace: %+v\n", createdCS)
	*/

	// Example: List spot node pools for an organization (commented out to prevent API calls)
	/*
	org := "my-org"
	fmt.Printf("\nListing spot node pools for org %s:\n", org)
	
	pools, err := client.SpotNodePools().List(context.Background(), org, "")
	if err != nil {
		log.Fatalf("Failed to list spot node pools: %v", err)
	}
	
	for _, pool := range pools {
		fmt.Printf("- %s (Desired: %d, Status: %s)\n", 
			pool.Name, 
			pool.Desired,
			pool.Status,
		)
	}
	*/
}
