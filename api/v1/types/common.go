package types

import (
	"fmt"
	"time"
)

// Resource represents compute resources (CPU, Memory)
type Resource struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

// PriceDetails represents pricing information for a resource
type PriceDetails struct {
	Region     string    `json:"region" yaml:"region"`
	ServerClass string   `json:"serverClass" yaml:"serverClass"`
	MarketPrice string   `json:"marketPrice" yaml:"marketPrice"`
	UpdatedAt  time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// Region represents a cloud region
type Region struct {
	Name        string `json:"name" yaml:"name"`
	DisplayName string `json:"displayName" yaml:"displayName"`
	Provider    string `json:"provider" yaml:"provider"`
	Enabled     bool   `json:"enabled" yaml:"enabled"`
}

// Common fields used across multiple types
type Metadata struct {
	Name              string            `json:"name" yaml:"name"`
	Namespace         string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	CreationTimestamp time.Time         `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Labels            map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// Autoscaling configuration
type Autoscaling struct {
	Enabled  bool  `json:"enabled" yaml:"enabled"`
	MinNodes int64 `json:"minNodes,omitempty" yaml:"minNodes,omitempty"`
	MaxNodes int64 `json:"maxNodes,omitempty" yaml:"maxNodes,omitempty"`
}

// Common error types
var (
	ErrNotFound = NewAPIError("not found", "The requested resource was not found", 404)
	ErrInvalidInput = NewAPIError("invalid input", "The provided input is invalid", 400)
)

// APIError represents an API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewAPIError(code, message string, status int) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}
