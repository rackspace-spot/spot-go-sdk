package rxtspot

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

// SpotNodePoolList represents a list of spot node pools
type SpotNodePoolList struct {
	Items []SpotNodePool `json:"spotNodepools" yaml:"spotNodepools"`
}

// SpotNodePool represents a spot node pool configuration
type SpotNodePool struct {
	Name              string            `json:"name" yaml:"name"`
	CreationTimestamp time.Time         `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Org               string            `json:"org,omitempty" yaml:"org,omitempty"`
	Cloudspace        string            `json:"cloudspace,omitempty" yaml:"cloudspace,omitempty"`
	ServerClass       string            `json:"serverClass,omitempty" yaml:"serverClass,omitempty"`
	Desired           int               `json:"desired,omitempty" yaml:"desired,omitempty"`
	WonCount          int               `json:"wonCount,omitempty" yaml:"wonCount,omitempty"`
	CustomAnnotations map[string]string `json:"customAnnotations,omitempty" yaml:"customAnnotations,omitempty"`
	CustomLabels      map[string]string `json:"customLabels,omitempty" yaml:"customLabels,omitempty"`
	CustomTaints      []interface{}     `json:"customTaints,omitempty" yaml:"customTaints,omitempty"`
	Autoscaling       struct {
		Enabled  bool  `json:"enabled" yaml:"enabled"`
		MinNodes int64 `json:"minNodes" yaml:"minNodes"`
		MaxNodes int64 `json:"maxNodes" yaml:"maxNodes"`
	} `json:"autoscaling" yaml:"autoscaling"`
	BidPrice string `json:"bidPrice,omitempty" yaml:"bidPrice,omitempty"`
	Status   string `json:"status,omitempty" yaml:"status,omitempty"`
}

// OnDemandNodePoolList represents a list of on-demand node pools
type OnDemandNodePoolList struct {
	Items []OnDemandNodePool `json:"ondemandNodepools" yaml:"ondemandNodepools"`
}

// OnDemandNodePool represents an on-demand node pool configuration
type OnDemandNodePool struct {
	Name                 string            `json:"name" yaml:"name"`
	CreationTimestamp    time.Time         `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Org                  string            `json:"org,omitempty" yaml:"org,omitempty"`
	Cloudspace           string            `json:"cloudspace,omitempty" yaml:"cloudspace,omitempty"`
	ServerClass          string            `json:"serverClass,omitempty" yaml:"serverClass,omitempty"`
	Desired              int               `json:"desired,omitempty" yaml:"desired,omitempty"`
	WonCount             int               `json:"wonCount,omitempty" yaml:"wonCount,omitempty"`
	CustomAnnotations    map[string]string `json:"customAnnotations,omitempty" yaml:"customAnnotations,omitempty"`
	CustomLabels         map[string]string `json:"customLabels,omitempty" yaml:"customLabels,omitempty"`
	CustomTaints         []interface{}     `json:"customTaints,omitempty" yaml:"customTaints,omitempty"`
	OnDemandPricePerHour string            `json:"onDemandPricePerHour,omitempty" yaml:"onDemandPricePerHour,omitempty"`
	Autoscaling          struct {
		Enabled  bool `json:"enabled" yaml:"enabled"`
		MinNodes int  `json:"minNodes" yaml:"minNodes"`
		MaxNodes int  `json:"maxNodes" yaml:"maxNodes"`
	} `json:"autoscaling" yaml:"autoscaling"`
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
}

type OrganizationList struct {
	Items []Organization `json:"organizations" yaml:"organizations"`
}

type Organization struct {
	Name string `json:"name" yaml:"name"`
	ID   string `json:"id" yaml:"id"`
}

type RegionList struct {
	Items []Region `json:"regions" yaml:"regions"`
}

type Region struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type ServerClassList struct {
	Items []ServerClass `json:"serverClasses" yaml:"serverClasses"`
}

type ServerClass struct {
	Name                      string   `json:"name" yaml:"name"`
	Category                  string   `json:"category,omitempty" yaml:"category,omitempty"`
	Availability              string   `json:"availability,omitempty" yaml:"availability,omitempty"`
	Displayname               string   `json:"displayname,omitempty" yaml:"displayname,omitempty"`
	Region                    string   `json:"region,omitempty" yaml:"region,omitempty"`
	MinBidPricePerHour        string   `json:"minBidPricePerHour,omitempty" yaml:"minBidPricePerHour,omitempty"`
	CurrentMarketPricePerHour string   `json:"currentMarketPricePerHour,omitempty" yaml:"currentMarketPricePerHour,omitempty"`
	OnDemandPricePerHour      string   `json:"onDemandPricePerHour,omitempty" yaml:"onDemandPricePerHour,omitempty"`
	Resources                 Resource `json:"resources,omitempty" yaml:"resources,omitempty"`
}

type Resource struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
	GPU    string `json:"gpu,omitempty" yaml:"gpu,omitempty"`
}

type PriceDetails struct {
	ServerClassName string `json:"serverClassName" yaml:"serverClassName"`
	DisplayName     string `json:"displayName" yaml:"displayName"`
	Category        string `json:"category" yaml:"category"`
	Region          string `json:"region" yaml:"region"`
	MarketPrice     string `json:"currentMarketPrice" yaml:"currentMarketPrice"`
	CPU             string `json:"cpu" yaml:"cpu"`
	Memory          string `json:"memory" yaml:"memory"`
}
