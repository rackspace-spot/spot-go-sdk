package cloudspaces

import "context"

type Client interface {
	ListCloudspaces(ctx context.Context, org string) (*CloudSpaceList, error)
	CreateCloudspace(ctx context.Context, cs CloudSpace) error
	GetCloudspace(ctx context.Context, org, name string) (*CloudSpace, error)
	DeleteCloudspace(ctx context.Context, org, name string) error
	GetCloudspaceConfig(ctx context.Context, org, name string) (string, error)
}
