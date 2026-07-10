package rxtspot

import "time"

type cloudSpaceGetResponse struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Spec struct {
		BidRequests       []string `json:"bidRequests"`
		Cloud             string   `json:"cloud"`
		CNI               string   `json:"cni"`
		DeploymentType    string   `json:"deploymentType"`
		GpuEnabled        bool     `json:"gpuEnabled"`
		KubernetesVersion string   `json:"kubernetesVersion"`
		Region            string   `json:"region"`
		Type              string   `json:"type"`
		Webhook           string   `json:"webhook"`
	} `json:"spec"`
	Status struct {
		APIServerEndpoint        string                    `json:"APIServerEndpoint"`
		AssignedServers          map[string]AssignedServer `json:"assignedServers"`
		Bids                     map[string]Bid            `json:"bids"`
		CloudspaceClassName      string                    `json:"cloudspaceClassName"`
		CurrentKubernetesVersion string                    `json:"currentKubernetesVersion"`
		FirstReadyTimestamp      time.Time                 `json:"firstReadyTimestamp"`
		Health                   string                    `json:"health"`
		Reason                   string                    `json:"reason"`
		Phase                    string                    `json:"phase"`
	} `json:"status"`
}

// CloudSpaceSpec defines the spec for CloudSpace create requests
type CloudSpaceSpec struct {
	DeploymentType    string `json:"deploymentType"`
	Cloud             string `json:"cloud"`
	Region            string `json:"region"`
	Webhook           string `json:"webhook"`
	CNI               string `json:"cni"`
	KubernetesVersion string `json:"kubernetesVersion"`
	HAControlPlane    bool   `json:"HAControlPlane"`
	GpuEnabled        bool   `json:"gpuEnabled"`
}

// Common read-only response metadata with timestamp
type ResourceMetadataWithTimestamp struct {
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
}

// Spot node pool read-only autoscaling and spec
type SpotNodePoolAutoscalingRO struct {
	Enabled  bool `json:"enabled"`
	MaxNodes int  `json:"maxNodes"`
	MinNodes int  `json:"minNodes"`
}

type SpotNodePoolSpecReadOnly struct {
	Autoscaling       SpotNodePoolAutoscalingRO `json:"autoscaling"`
	BidPrice          string                    `json:"bidPrice"`
	CloudSpace        string                    `json:"cloudSpace"`
	CustomAnnotations map[string]string         `json:"customAnnotations"`
	CustomLabels      map[string]string         `json:"customLabels"`
	CustomTaints      []interface{}             `json:"customTaints,omitempty"`
	Desired           int                       `json:"desired"`
	ServerClass       string                    `json:"serverClass"`
}

type SpotNodePoolStatus struct {
	BidStatus            string   `json:"bidStatus"`
	CustomMetadataStatus struct{} `json:"customMetadataStatus"`
	WonCount             int      `json:"wonCount"`
}

// OnDemand node pool read-only spec and status
type OnDemandNodePoolSpecReadOnly struct {
	CloudSpace        string            `json:"cloudSpace"`
	CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`
	CustomLabels      map[string]string `json:"customLabels,omitempty"`
	CustomTaints      []interface{}     `json:"customTaints,omitempty"`
	Desired           int               `json:"desired"`
	ServerClass       string            `json:"serverClass"`
}

type OnDemandNodePoolStatus struct {
	ReservedCount  int    `json:"reservedCount"`
	ReservedStatus string `json:"reservedStatus"`
}

// Simple metadata with only name (used by server classes and regions)
type ResourceMetadata struct {
	Name string `json:"name"`
}

// Server class reusable spec and status
type ServerClassSpec struct {
	Availability       string `json:"availability"`
	Category           string `json:"category"`
	DisplayName        string `json:"displayName"`
	FlavorType         string `json:"flavorType"`
	MinBidPricePerHour string `json:"minBidPricePerHour,omitempty"`
	OnDemandPricing    struct {
		Cost     string `json:"cost"`
		Interval string `json:"interval"`
	} `json:"onDemandPricing"`
	Provider struct {
		ProviderFlavorID string `json:"providerFlavorID"`
		ProviderType     string `json:"providerType"`
	} `json:"provider"`
	Region    string `json:"region"`
	Resources struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"resources"`
}

type ServerClassStatus struct {
	Available   int `json:"available"`
	Capacity    int `json:"capacity"`
	LastAuction int `json:"lastAuction"`
	Reserved    int `json:"reserved"`
	SpotPricing struct {
		HammerPricePerHour string `json:"hammerPricePerHour"`
		MarketPricePerHour string `json:"marketPricePerHour"`
	} `json:"spotPricing"`
}

// Region reusable spec
type RegionSpec struct {
	Country     string `json:"country"`
	Description string `json:"description"`
	Generation  string `json:"generation"`
	Provider    struct {
		ProviderRegionName string `json:"providerRegionName"`
		ProviderType       string `json:"providerType"`
	} `json:"provider"`
}

type SpotNodePoolGetResponse struct {
	APIVersion string                        `json:"apiVersion"`
	Kind       string                        `json:"kind"`
	Metadata   ResourceMetadataWithTimestamp `json:"metadata"`
	Spec       SpotNodePoolSpecReadOnly      `json:"spec"`
	Status     SpotNodePoolStatus            `json:"status"`
}

type Bid struct {
	BidName  string `json:"bidName"`
	Type     string `json:"type"`
	WonCount int    `json:"wonCount"`
}

type cloudSpaceListResponse struct {
	Items []cloudSpaceGetResponse `json:"items"`
}

type CloudSpaceCreateRequestBody struct {
	APIVersion string                    `json:"apiVersion"`
	Kind       string                    `json:"kind"`
	Metadata   ObjectMetaWithAnnotations `json:"metadata"`
	Spec       CloudSpaceSpec            `json:"spec"`
}

type ObjectMeta struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
}

// ObjectMetaWithAnnotations is used for resources that carry annotations instead of labels in metadata
type ObjectMetaWithAnnotations struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Annotations map[string]string `json:"annotations"`
}

type CommonNodePoolSpec struct {
	ServerClass       string            `json:"serverClass"`
	Desired           int               `json:"desired"`
	CloudSpace        string            `json:"cloudSpace"`
	CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`
	CustomLabels      map[string]string `json:"customLabels,omitempty"`
	CustomTaints      []interface{}     `json:"customTaints,omitempty"`
}

type AutoscalingInt64 struct {
	Enabled  bool  `json:"enabled"`
	MinNodes int64 `json:"minNodes"`
	MaxNodes int64 `json:"maxNodes"`
}

type SpotNodePoolSpec struct {
	CommonNodePoolSpec
	BidPrice    string           `json:"bidPrice"`
	Autoscaling AutoscalingInt64 `json:"autoscaling"`
}

type SpotNodePoolRequestBody struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   ObjectMeta       `json:"metadata"`
	Spec       SpotNodePoolSpec `json:"spec"`
}

type KubeConfigResponse struct {
	Data struct {
		Kubeconfig string `json:"kubeconfig"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type SpotNodePoolListResponse struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string                        `json:"apiVersion"`
		Kind       string                        `json:"kind"`
		Metadata   ResourceMetadataWithTimestamp `json:"metadata"`
		Spec       SpotNodePoolSpecReadOnly      `json:"spec"`
		Status     SpotNodePoolStatus            `json:"status"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
		Continue        string `json:"continue"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
}

type OnDemandNodePoolSpec struct {
	CommonNodePoolSpec
}

type OnDemandNodePoolCreateRequestBody struct {
	APIVersion string               `json:"apiVersion"`
	Kind       string               `json:"kind"`
	Metadata   ObjectMeta           `json:"metadata"`
	Spec       OnDemandNodePoolSpec `json:"spec"`
}

type AutoscalingInt64Update struct {
	Enabled  bool  `json:"enabled"`
	MinNodes int64 `json:"minNodes,omitempty"`
	MaxNodes int64 `json:"maxNodes,omitempty"`
}

type SpotNodePoolUpdateSpec struct {
	Desired           int                     `json:"desired,omitempty"`
	BidPrice          string                  `json:"bidPrice,omitempty"`
	CustomAnnotations map[string]string       `json:"customAnnotations,omitempty"`
	CustomLabels      map[string]string       `json:"customLabels,omitempty"`
	CustomTaints      []interface{}           `json:"customTaints,omitempty"`
	Autoscaling       *AutoscalingInt64Update `json:"autoscaling,omitempty"`
}

type OnDemandNodePoolUpdateSpec struct {
	Desired           int               `json:"desired,omitempty"`
	CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`
	CustomLabels      map[string]string `json:"customLabels,omitempty"`
	CustomTaints      []interface{}     `json:"customTaints,omitempty"`
}

type SpotNodePoolUpdateRequestBody struct {
	Spec SpotNodePoolUpdateSpec `json:"spec"`
}

type OnDemandNodePoolUpdateRequestBody struct {
	Spec OnDemandNodePoolUpdateSpec `json:"spec"`
}

type OnDemandNodePoolGetResponse struct {
	APIVersion string                        `json:"apiVersion"`
	Kind       string                        `json:"kind"`
	Metadata   ResourceMetadataWithTimestamp `json:"metadata"`
	Spec       OnDemandNodePoolSpecReadOnly  `json:"spec"`
	Status     OnDemandNodePoolStatus        `json:"status"`
}

type OnDemandNodePoolListResponse struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string                        `json:"apiVersion"`
		Kind       string                        `json:"kind"`
		Metadata   ResourceMetadataWithTimestamp `json:"metadata"`
		Spec       OnDemandNodePoolSpecReadOnly  `json:"spec"`
		Status     OnDemandNodePoolStatus        `json:"status"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
		Continue        string `json:"continue"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
}

type GetServerClassResponse struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   ResourceMetadata  `json:"metadata"`
	Spec       ServerClassSpec   `json:"spec"`
	Status     ServerClassStatus `json:"status"`
}

type ListServerClassesResponse struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string            `json:"apiVersion"`
		Kind       string            `json:"kind"`
		Metadata   ResourceMetadata  `json:"metadata"`
		Spec       ServerClassSpec   `json:"spec"`
		Status     ServerClassStatus `json:"status"`
	} `json:"items"`
}

type ListRegionsResponse struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string           `json:"apiVersion"`
		Kind       string           `json:"kind"`
		Metadata   ResourceMetadata `json:"metadata"`
		Spec       RegionSpec       `json:"spec"`
	} `json:"items"`
}

// --- VM Backend Types ---

// VMCloudSpace API response types
type vmCloudSpaceGetResponse struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Spec struct {
		BidRequests []string `json:"bidRequests"`
		Region      string   `json:"region"`
		Webhook     string   `json:"webhook"`
		VMSshKeyRef struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"vmSshKeyRef"`
	} `json:"spec"`
	Status struct {
		AssignedServers map[string]VMAssignedServer `json:"assignedServers"`
		Bids            map[string]Bid              `json:"bids"`
		Phase           string                      `json:"phase"`
		Reason          string                      `json:"reason"`
		Health          string                      `json:"health"`
	} `json:"status"`
}

type vmCloudSpaceListResponse struct {
	Items []vmCloudSpaceGetResponse `json:"items"`
}

type VMCloudSpaceCreateRequestBody struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Region      string `json:"region"`
		Webhook     string `json:"webhook,omitempty"`
		VMSshKeyRef struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"vmSshKeyRef"`
	} `json:"spec"`
}

// VMCloudSpace update request body - only Webhook field is allowed to be modified
type VMCloudSpaceUpdateRequestBody struct {
	Spec struct {
		Webhook string `json:"webhook"`
	} `json:"spec"`
}

// VMPool API response types
type vmPoolGetResponse struct {
	APIVersion string                        `json:"apiVersion"`
	Kind       string                        `json:"kind"`
	Metadata   ResourceMetadataWithTimestamp `json:"metadata"`
	Spec       struct {
		BidPrice     string `json:"bidPrice"`
		Desired      int    `json:"desired"`
		PoolType     string `json:"poolType"`
		ServerClass  string `json:"serverClass"`
		VMCloudSpace string `json:"vmCloudSpace"`
		VMImage      string `json:"vmImage"`
		VMUserData   string `json:"vmUserData,omitempty"`
	} `json:"spec"`
	Status struct {
		BidStatus string `json:"bidStatus"`
		WonCount  int    `json:"wonCount"`
	} `json:"status"`
}

type vmPoolListResponse struct {
	APIVersion string              `json:"apiVersion"`
	Items      []vmPoolGetResponse `json:"items"`
}

type VMPoolCreateRequestBody struct {
	APIVersion string     `json:"apiVersion"`
	Kind       string     `json:"kind"`
	Metadata   ObjectMeta `json:"metadata"`
	Spec       struct {
		BidPrice     string `json:"bidPrice"`
		Desired      int    `json:"desired"`
		PoolType     string `json:"poolType,omitempty"`
		ServerClass  string `json:"serverClass"`
		VMCloudSpace string `json:"vmCloudSpace"`
		VMImage      string `json:"vmImage,omitempty"`
		VMUserData   string `json:"vmUserData,omitempty"`
	} `json:"spec"`
}

type VMPoolUpdateRequestBody struct {
	Spec struct {
		Desired  int    `json:"desired,omitempty"`
		BidPrice string `json:"bidPrice,omitempty"`
	} `json:"spec"`
}

// VMSSHKey API response types
type vmSSHKeyGetResponse struct {
	APIVersion string                        `json:"apiVersion"`
	Kind       string                        `json:"kind"`
	Metadata   ResourceMetadataWithTimestamp `json:"metadata"`
	Spec       struct {
		PublicKey   string `json:"publicKey"`
		Description string `json:"description"`
	} `json:"spec"`
	Status struct {
		Fingerprint string `json:"fingerprint"`
		Validated   bool   `json:"validated"`
	} `json:"status"`
}

type vmSSHKeyListResponse struct {
	APIVersion string                `json:"apiVersion"`
	Items      []vmSSHKeyGetResponse `json:"items"`
}

type VMSSHKeyCreateRequestBody struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		PublicKey   string `json:"publicKey"`
		Description string `json:"description,omitempty"`
	} `json:"spec"`
}
