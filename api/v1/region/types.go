package region

// Region represents a cloud region
type Region struct {
	Name        string `json:"name" yaml:"name"`
	DisplayName string `json:"displayName" yaml:"displayName"`
	Provider    string `json:"provider" yaml:"provider"`
	Enabled     bool   `json:"enabled" yaml:"enabled"`
}

// RegionList represents a list of regions
type RegionList struct {
	Items []Region `json:"regions" yaml:"regions"`
}

// ListOptions specifies optional parameters for listing regions
type ListOptions struct {
	Provider string `json:"provider,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

// GetOptions specifies optional parameters for getting a region
type GetOptions struct {
	IncludeZones bool `json:"includeZones,omitempty"`
}

// Zone represents an availability zone within a region
type Zone struct {
	Name    string `json:
ame" yaml:"name"`
	Enabled bool   `json:"enabled" yaml:"enabled"`
}
