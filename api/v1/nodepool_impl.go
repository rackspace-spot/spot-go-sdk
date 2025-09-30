package rxtspot

// GetType returns the type of the node pool (spot)
func (s *SpotNodePool) GetType() NodePoolType {
	return NodePoolTypeSpot
}

// GetName returns the name of the node pool
func (s *SpotNodePool) GetName() string {
	return s.Name
}

// GetOrg returns the organization ID this node pool belongs to
func (s *SpotNodePool) GetOrg() string {
	return s.Org
}

// GetName returns the name of the node pool
func (o *OnDemandNodePool) GetName() string {
	return o.Name
}

// GetType returns the type of the node pool (spot)
func (o *OnDemandNodePool) GetType() NodePoolType {
	return NodePoolTypeOnDemand
}

// GetOrg returns the organization ID this node pool belongs to
func (o *OnDemandNodePool) GetOrg() string {
	return o.Org
}

// Ensure SpotNodePool implements NodePool interface
var _ NodePool = (*SpotNodePool)(nil)

// Ensure OnDemandNodePool implements NodePool interface
var _ NodePool = (*OnDemandNodePool)(nil)
