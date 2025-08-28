package cloudspace

import (
	"time"

	"github.com/rackspace/spot-go-sdk/api/v1/types"
)

// CloudSpace represents a cloudspace configuration
type CloudSpace struct {
	Name                 string                    `json:"name" yaml:"name"`
	Org                  string                    `json:"org" yaml:"org"`
	CreationTimestamp    time.Time                 `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	CNI                  string                    `json:"cni,omitempty" yaml:"cni,omitempty"`
	DeploymentType       string                    `json:"deploymentType,omitempty" yaml:"deploymentType,omitempty"`
	GpuEnabled           bool                      `json:"gpuEnabled,omitempty" yaml:"gpuEnabled,omitempty"`
	KubernetesVersion    string                    `json:"kubernetesVersion,omitempty" yaml:"kubernetesVersion,omitempty"`
	Region               string                    `json:"region,omitempty" yaml:"region,omitempty"`
	PreemptionWebhookURL string                    `json:"preEmptionWebhookURL,omitempty" yaml:"preEmptionWebhookURL,omitempty"`
	APIServerEndpoint    string                    `json:"apiServerEndpoint,omitempty" yaml:"apiServerEndpoint,omitempty"`
	AssignedServers      map[string]AssignedServer `json:"assignedServers,omitempty" yaml:"assignedServers,omitempty"`
	SpotNodepools        []*nodepool.Spot          `json:"spotNodepools,omitempty" yaml:"spotNodepools,omitempty"`
	OnDemandNodePools    []*nodepool.OnDemand      `json:"ondemandNodepools,omitempty" yaml:"ondemandNodepools,omitempty"`
	Status               string                    `json:"status,omitempty" yaml:"status,omitempty"`
	Message              string                    `json:"message,omitempty"`
}

// AssignedServer represents a server assigned to a cloudspace
type AssignedServer struct {
	IP              string `json:"IP" yaml:"IP"`
	ClusterRole     string `json:"clusterRole" yaml:"clusterRole"`
	ServerClassName string `json:"serverClassName" yaml:"serverClassName"`
	State           string `json:"state" yaml:"state"`
}

// CloudSpaceList represents a list of cloudspaces
type CloudSpaceList struct {
	Items []CloudSpace `json:"cloudspaces" yaml:"cloudspaces"`
}

// CreateOptions specifies parameters for creating a new cloudspace
type CreateOptions struct {
	Name          string            `json:"name"`
	Org           string            `json:"org"`
	Region        string            `json:"region"`
	CNI           string            `json:"cni,omitempty"`
	GpuEnabled    bool              `json:"gpuEnabled,omitempty"`
	K8sVersion    string            `json:"kubernetesVersion,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`
}

// UpdateOptions specifies parameters for updating a cloudspace
type UpdateOptions struct {
	Name          string            `json:"name"`
	Org           string            `json:"org"`
	Labels        map[string]string `json:"labels,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`n}

// ListOptions specifies optional parameters for listing cloudspaces
type ListOptions struct {
	Org    string            `json:"org,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}
