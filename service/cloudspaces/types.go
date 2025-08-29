package cloudspaces

import "time"

type CloudSpaceList struct {
	Items []CloudSpace `json:"cloudspaces" yaml:"cloudspaces"`
}

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
	SpotNodepools        []*SpotNodePool           `json:"spotNodepools,omitempty" yaml:"spotNodepools,omitempty"`
	OnDemandNodePools    []*OnDemandNodePool       `json:"ondemandNodepools,omitempty" yaml:"ondemandNodepools,omitempty"`
	Status               string                    `json:"status,omitempty" yaml:"status,omitempty"`
	Message              string                    `json:"message,omitempty"`
}

type AssignedServer struct {
	IP              string `json:"IP" yaml:"IP"`
	ClusterRole     string `json:"clusterRole" yaml:"clusterRole"`
	ServerClassName string `json:"serverClassName" yaml:"serverClassName"`
	State           string `json:"state" yaml:"state"`
}
