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
	AutopilotNodePools   []*AutopilotNodePool      `json:"autopilotNodepools,omitempty" yaml:"autopilotNodepools,omitempty"`
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

// AutopilotNodePoolList represents a list of autopilot node pools
type AutopilotNodePoolList struct {
	Items []AutopilotNodePool `json:"autopilotNodepools" yaml:"autopilotNodepools"`
}

// AutopilotNodePool represents an autopilot node pool configuration. Unlike SpotNodePool
// (an exact bid for one server class), it takes a total vCPU target and an hourly budget
// ceiling, and the control plane bids across eligible server classes to fill it.
type AutopilotNodePool struct {
	Name               string            `json:"name" yaml:"name"`
	CreationTimestamp  time.Time         `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Org                string            `json:"org,omitempty" yaml:"org,omitempty"`
	Cloudspace         string            `json:"cloudspace,omitempty" yaml:"cloudspace,omitempty"`
	Region             string            `json:"region,omitempty" yaml:"region,omitempty"`
	VCPUTotal          int               `json:"vcpuTotal,omitempty" yaml:"vcpuTotal,omitempty"`
	VCPUPerNodeMin     int               `json:"vcpuPerNodeMin,omitempty" yaml:"vcpuPerNodeMin,omitempty"`
	VCPUPerNodeMax     int               `json:"vcpuPerNodeMax,omitempty" yaml:"vcpuPerNodeMax,omitempty"`
	MemoryPerVCPU      string            `json:"memoryPerVCPU,omitempty" yaml:"memoryPerVCPU,omitempty"`
	BudgetPerHour      string            `json:"budgetPerHour,omitempty" yaml:"budgetPerHour,omitempty"`
	AllocationStrategy string            `json:"allocationStrategy,omitempty" yaml:"allocationStrategy,omitempty"`
	CustomAnnotations  map[string]string `json:"customAnnotations,omitempty" yaml:"customAnnotations,omitempty"`
	CustomLabels       map[string]string `json:"customLabels,omitempty" yaml:"customLabels,omitempty"`
	CustomTaints       []interface{}     `json:"customTaints,omitempty" yaml:"customTaints,omitempty"`

	Phase              string                      `json:"phase,omitempty" yaml:"phase,omitempty"`
	TargetVCPUs        int                         `json:"targetVCPUs,omitempty" yaml:"targetVCPUs,omitempty"`
	ManagedVCPUs       int                         `json:"managedVCPUs,omitempty" yaml:"managedVCPUs,omitempty"`
	ManagedMemoryGB    string                      `json:"managedMemoryGB,omitempty" yaml:"managedMemoryGB,omitempty"`
	Allocations        []AutopilotAllocation       `json:"allocations,omitempty" yaml:"allocations,omitempty"`
	ClassStates        []AutopilotClassState       `json:"classStates,omitempty" yaml:"classStates,omitempty"`
	AllocationsHistory []AutopilotAllocationRecord `json:"allocationsHistory,omitempty" yaml:"allocationsHistory,omitempty"`
}

// AutopilotAllocation is the per-server-class breakdown of what an autopilot pool is managing.
type AutopilotAllocation struct {
	ServerClass        string `json:"serverClass" yaml:"serverClass"`
	SpotNodePool       string `json:"spotNodePool,omitempty" yaml:"spotNodePool,omitempty"`
	MarketPricePerHour string `json:"marketPricePerHour" yaml:"marketPricePerHour"`
	BidPricePerHour    string `json:"bidPricePerHour" yaml:"bidPricePerHour"`
	VCPUPerNode        int    `json:"vcpuPerNode" yaml:"vcpuPerNode"`
	MemoryGBPerNode    string `json:"memoryGBPerNode,omitempty" yaml:"memoryGBPerNode,omitempty"`
	DesiredNodes       int    `json:"desiredNodes" yaml:"desiredNodes"`
	WonNodes           int    `json:"wonNodes" yaml:"wonNodes"`
}

// AutopilotClassState records the most recent bid outcome for one server class.
type AutopilotClassState struct {
	ServerClass     string     `json:"serverClass" yaml:"serverClass"`
	LastOutcome     string     `json:"lastOutcome" yaml:"lastOutcome"`
	LastAttemptTime time.Time  `json:"lastAttemptTime,omitempty" yaml:"lastAttemptTime,omitempty"`
	CooldownUntil   *time.Time `json:"cooldownUntil,omitempty" yaml:"cooldownUntil,omitempty"`
}

// AutopilotAllocationRecord is one entry in an autopilot pool's audit trail.
type AutopilotAllocationRecord struct {
	Time            time.Time `json:"time" yaml:"time"`
	ServerClass     string    `json:"serverClass" yaml:"serverClass"`
	SpotNodePool    string    `json:"spotNodePool,omitempty" yaml:"spotNodePool,omitempty"`
	Event           string    `json:"event" yaml:"event"`
	Reason          string    `json:"reason,omitempty" yaml:"reason,omitempty"`
	DesiredNodes    int       `json:"desiredNodes" yaml:"desiredNodes"`
	WonNodes        int       `json:"wonNodes" yaml:"wonNodes"`
	BidPricePerHour string    `json:"bidPricePerHour,omitempty" yaml:"bidPricePerHour,omitempty"`
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
	Displayname               string   `json:"displayName,omitempty" yaml:"displayName,omitempty"`
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

// VMCloudSpaceList represents a list of VM cloud spaces
type VMCloudSpaceList struct {
	Items []VMCloudSpace `json:"vmCloudSpaces" yaml:"vmCloudSpaces"`
}

// VMSshKeyRef represents a reference to a VM SSH key
type VMSshKeyRef struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

// VMCloudSpace represents a VM cloud space configuration
type VMCloudSpace struct {
	Name              string                      `json:"name" yaml:"name"`
	Org               string                      `json:"org" yaml:"org"`
	CreationTimestamp time.Time                   `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Region            string                      `json:"region" yaml:"region"`
	Webhook           string                      `json:"webhook,omitempty" yaml:"webhook,omitempty"`
	VMSshKeyRef       VMSshKeyRef                 `json:"vmSshKeyRef" yaml:"vmSshKeyRef"`
	VMPools           []*VMPool                   `json:"vmPools,omitempty" yaml:"vmPools,omitempty"`
	AssignedServers   map[string]VMAssignedServer `json:"assignedServers,omitempty" yaml:"assignedServers,omitempty"`
	Status            string                      `json:"status,omitempty" yaml:"status,omitempty"`
	Message           string                      `json:"message,omitempty" yaml:"message,omitempty"`
	Health            string                      `json:"health,omitempty" yaml:"health,omitempty"`
}

// VMAssignedServer represents a server assigned to a VM cloud space
type VMAssignedServer struct {
	IPAddress       string `json:"ipAddress" yaml:"ipAddress"`
	CPU             string `json:"cpu" yaml:"cpu"`
	Memory          string `json:"memory" yaml:"memory"`
	DisplayName     string `json:"displayName" yaml:"displayName"`
	ServerClassName string `json:"serverClassName" yaml:"serverClassName"`
	ServerName      string `json:"serverName" yaml:"serverName"`
	ServerType      string `json:"serverType" yaml:"serverType"`
	NodePoolName    string `json:"nodePoolName" yaml:"nodePoolName"`
	State           string `json:"state" yaml:"state"`
	Error           string `json:"error,omitempty" yaml:"error,omitempty"`
}

// VMPoolList represents a list of VM pools
type VMPoolList struct {
	Items []VMPool `json:"vmPools" yaml:"vmPools"`
}

// VMPool represents a VM pool configuration
type VMPool struct {
	Name              string    `json:"name" yaml:"name"`
	CreationTimestamp time.Time `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Org               string    `json:"org,omitempty" yaml:"org,omitempty"`
	VMCloudSpace      string    `json:"vmCloudSpace,omitempty" yaml:"vmCloudSpace,omitempty"`
	ServerClass       string    `json:"serverClass,omitempty" yaml:"serverClass,omitempty"`
	Desired           int       `json:"desired,omitempty" yaml:"desired,omitempty"`
	BidPrice          string    `json:"bidPrice,omitempty" yaml:"bidPrice,omitempty"`
	PoolType          string    `json:"poolType,omitempty" yaml:"poolType,omitempty"`
	VMImage           string    `json:"vmImage,omitempty" yaml:"vmImage,omitempty"`
	VMUserData        string    `json:"vmUserData,omitempty" yaml:"vmUserData,omitempty"`
	WonCount          int       `json:"wonCount,omitempty" yaml:"wonCount,omitempty"`
	BidStatus         string    `json:"bidStatus,omitempty" yaml:"bidStatus,omitempty"`
}

// VMSSHKeyList represents a list of VM SSH keys
type VMSSHKeyList struct {
	Items []VMSSHKey `json:"vmSshKeys" yaml:"vmSshKeys"`
}

// VMSSHKey represents a VM SSH key
type VMSSHKey struct {
	Name              string    `json:"name" yaml:"name"`
	CreationTimestamp time.Time `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Org               string    `json:"org,omitempty" yaml:"org,omitempty"`
	PublicKey         string    `json:"publicKey" yaml:"publicKey"`
	Description       string    `json:"description,omitempty" yaml:"description,omitempty"`
	Fingerprint       string    `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	Validated         bool      `json:"validated,omitempty" yaml:"validated,omitempty"`
}
