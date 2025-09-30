// examples/cloudspace/main_test.go
package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCloudspace(t *testing.T) {
	// Setup mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockSpotAPI(ctrl)

	// Test data
	testCloudspace := v1.CloudSpace{
		Name:        "test-cloudspace",
		DisplayName: "Test CloudSpace",
		Region:      "us-east-1",
	}

	// Set expectations
	mockClient.EXPECT().
		CreateCloudspace(gomock.Any(), testCloudspace).
		Return(nil) // Expect no error

	// Test the function
	err := createExample(mockClient)

	// Verify
	assert.NoError(t, err)
}

func TestListCloudspaces(t *testing.T) {
	// Setup mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockSpotAPI(ctrl)

	// Test data
	expectedList := &v1.CloudSpaceList{
		Items: []v1.CloudSpace{
			{
				Metadata: v1.ObjectMeta{Name: "cs1"},
				Spec:     v1.CloudSpaceSpec{DisplayName: "CloudSpace 1"},
			},
		},
	}

	// Set expectations
	mockClient.EXPECT().
		ListCloudspaces(gomock.Any(), "test-org").
		Return(expectedList, nil)

	// Test the function
	listExample(mockClient)
}
