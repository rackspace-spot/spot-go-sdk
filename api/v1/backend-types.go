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

type SpotNodePoolGetResponse struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		CreationTimestamp time.Time         `json:"creationTimestamp"`
		Labels            map[string]string `json:"labels"`
		Name              string            `json:"name"`
		Namespace         string            `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Autoscaling struct {
			Enabled  bool `json:"enabled"`
			MaxNodes int  `json:"maxNodes"`
			MinNodes int  `json:"minNodes"`
		} `json:"autoscaling"`
		BidPrice          string            `json:"bidPrice"`
		CloudSpace        string            `json:"cloudSpace"`
		CustomAnnotations map[string]string `json:"customAnnotations"`
		CustomLabels      map[string]string `json:"customLabels"`
		CustomTaints      []interface{}     `json:"customTaints,omitempty"`
		Desired           int               `json:"desired"`
		ServerClass       string            `json:"serverClass"`
	} `json:"spec"`
	Status struct {
		BidStatus            string `json:"bidStatus"`
		CustomMetadataStatus struct {
		} `json:"customMetadataStatus"`
		WonCount int `json:"wonCount"`
	} `json:"status"`
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
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name        string            `json:"name"`
		Namespace   string            `json:"namespace"`
		Annotations map[string]string `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		DeploymentType    string `json:"deploymentType"`
		Cloud             string `json:"cloud"`
		Region            string `json:"region"`
		Webhook           string `json:"webhook"`
		CNI               string `json:"cni"`
		KubernetesVersion string `json:"kubernetesVersion"`
		HAControlPlane    bool   `json:"HAControlPlane"`
		GpuEnabled        bool   `json:"gpuEnabled"`
	} `json:"spec"`
}

type SpotNodePoolRequestBody struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string            `json:"name"`
		Namespace string            `json:"namespace"`
		Labels    map[string]string `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		ServerClass       string            `json:"serverClass"`
		Desired           int               `json:"desired"`
		BidPrice          string            `json:"bidPrice"`
		CloudSpace        string            `json:"cloudSpace"`
		CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`
		CustomLabels      map[string]string `json:"customLabels,omitempty"`
		CustomTaints      []interface{}     `json:"customTaints,omitempty"`
		Autoscaling       struct {
			Enabled  bool  `json:"enabled"`
			MinNodes int64 `json:"minNodes"`
			MaxNodes int64 `json:"maxNodes"`
		} `json:"autoscaling"`
	} `json:"spec"`
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
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			CreationTimestamp time.Time         `json:"creationTimestamp"`
			Name              string            `json:"name"`
			Namespace         string            `json:"namespace"`
			Labels            map[string]string `json:"labels"`
		} `json:"metadata"`
		Spec struct {
			Autoscaling struct {
				Enabled  bool `json:"enabled"`
				MaxNodes int  `json:"maxNodes"`
				MinNodes int  `json:"minNodes"`
			} `json:"autoscaling"`
			BidPrice          string            `json:"bidPrice"`
			CloudSpace        string            `json:"cloudSpace"`
			CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`
			CustomLabels      map[string]string `json:"customLabels,omitempty"`
			CustomTaints      []interface{}     `json:"customTaints,omitempty"`
			Desired           int               `json:"desired"`
			ServerClass       string            `json:"serverClass"`
		} `json:"spec"`
		Status struct {
			BidStatus            string `json:"bidStatus"`
			CustomMetadataStatus struct {
			} `json:"customMetadataStatus"`
			WonCount int `json:"wonCount"`
		} `json:"status"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
		Continue        string `json:"continue"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
}

type OnDemandNodePoolCreateRequestBody struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string            `json:"name"`
		Namespace string            `json:"namespace"`
		Labels    map[string]string `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		ServerClass       string            `json:"serverClass"`
		Desired           int               `json:"desired"`
		CloudSpace        string            `json:"cloudSpace"`
		CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`
		CustomLabels      map[string]string `json:"customLabels,omitempty"`
		CustomTaints      []interface{}     `json:"customTaints,omitempty"`
		Autoscaling       struct {
			Enabled  bool `json:"enabled"`
			MinNodes any  `json:"minNodes"`
			MaxNodes any  `json:"maxNodes"`
		} `json:"autoscaling"`
	} `json:"spec"`
}

type OnDemandNodePoolGetResponse struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		CreationTimestamp time.Time         `json:"creationTimestamp"`
		Labels            map[string]string `json:"labels"`
		Name              string            `json:"name"`
		Namespace         string            `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		CloudSpace        string            `json:"cloudSpace"`
		CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`
		CustomLabels      map[string]string `json:"customLabels,omitempty"`
		CustomTaints      []interface{}     `json:"customTaints,omitempty"`
		Desired           int               `json:"desired"`
		ServerClass       string            `json:"serverClass"`
	} `json:"spec"`
	Status struct {
		ReservedCount  int    `json:"reservedCount"`
		ReservedStatus string `json:"reservedStatus"`
	} `json:"status"`
}

type OnDemandNodePoolListResponse struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			CreationTimestamp time.Time         `json:"creationTimestamp"`
			Labels            map[string]string `json:"labels"`
			Name              string            `json:"name"`
			Namespace         string            `json:"namespace"`
		} `json:"metadata"`
		Spec struct {
			CloudSpace        string            `json:"cloudSpace"`
			CustomAnnotations map[string]string `json:"customAnnotations"`
			CustomLabels      map[string]string `json:"customLabels"`
			CustomTaints      []interface{}     `json:"customTaints,omitempty"`
			Desired           int               `json:"desired"`
			ServerClass       string            `json:"serverClass"`
		} `json:"spec"`
		Status struct {
			ReservedCount  int    `json:"reservedCount"`
			ReservedStatus string `json:"reservedStatus"`
		} `json:"status"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
		Continue        string `json:"continue"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
}

type GetServerClassResponse struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
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
	} `json:"spec"`
	Status struct {
		Available   int `json:"available"`
		Capacity    int `json:"capacity"`
		LastAuction int `json:"lastAuction"`
		Reserved    int `json:"reserved"`
		SpotPricing struct {
			HammerPricePerHour string `json:"hammerPricePerHour"`
			MarketPricePerHour string `json:"marketPricePerHour"`
		} `json:"spotPricing"`
	} `json:"status"`
}

type ListServerClassesResponse struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Spec struct {
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
		} `json:"spec"`
		Status struct {
			Available   int `json:"available"`
			Capacity    int `json:"capacity"`
			LastAuction int `json:"lastAuction"`
			Reserved    int `json:"reserved"`
			SpotPricing struct {
				HammerPricePerHour string `json:"hammerPricePerHour"`
				MarketPricePerHour string `json:"marketPricePerHour"`
			} `json:"spotPricing"`
		} `json:"status"`
	} `json:"items"`
}

type ListRegionsResponse struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Spec struct {
			Country     string `json:"country"`
			Description string `json:"description"`
			Generation  string `json:"generation"`
			Provider    struct {
				ProviderRegionName string `json:"providerRegionName"`
				ProviderType       string `json:"providerType"`
			} `json:"provider"`
		} `json:"spec"`
	} `json:"items"`
}
