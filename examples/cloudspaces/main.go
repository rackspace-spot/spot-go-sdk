// examples/cloudspace/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	v1 "github.com/rackspace-spot/spot-go-sdk/api/v1"
)

func main() {
	// Initialize spot client

	os.Setenv("RXTSPOT_CLIENT_ID", "BSQn6mv3OhI3jxQj0cvaZAUW4FBvLctu")

	spotClient, err := v1.NewSpotClient(&v1.Config{
		RefreshToken: "gTLEWRNYgwI3qNMLCgPC9vaL_S0n9MeZZy4spo5sj0djr",
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("spotClient: %+v\n", spotClient)
	regions, err := spotClient.ListRegions(context.Background())
	if err != nil {
		log.Fatalf("Failed to list regions: %v", err)
	}

	fmt.Println("Regions:")
	for _, region := range regions {
		fmt.Printf("- %s (%s)\n", region.Name, region.Name)
	}
}

/*
func createCloudspace(client *v1.RackspaceSpotClient) {
	ctx := context.Background()

	// Create a new cloudspace
	cloudspace := v1.CloudSpace{
		Name:        "example-cloudspace",
		DisplayName: "Example CloudSpace",
		Region:      "us-east-1",
		// Add other required fields
	}

	err := client.CreateCloudspace(ctx, cloudspace)
	if err != nil {
		log.Fatalf("Failed to create cloudspace: %v", err)
	}

	fmt.Println("Successfully created cloudspace:", cloudspace.Name)
}

func listCloudspaces(client *v1.RackspaceSpotClient) {
	ctx := context.Background()

	// List all cloudspaces in the organization
	cloudspaces, err := client.ListCloudspaces(ctx, "your-org-id")
	if err != nil {
		log.Fatalf("Failed to list cloudspaces: %v", err)
	}

	fmt.Println("CloudSpaces:")
	for _, cs := range cloudspaces.Items {
		fmt.Printf("- %s (%s)\n", cs.Metadata.Name, cs.Spec.DisplayName)
	}
}
*/
