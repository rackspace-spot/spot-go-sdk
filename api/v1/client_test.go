package rxtspot

import (
	"errors"
	"net/http"
	"testing"
)

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "non-HTTP error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "HTTP 404 error",
			err:      &HTTPStatusError{StatusCode: http.StatusNotFound},
			expected: true,
		},
		{
			name:     "HTTP 403 error",
			err:      &HTTPStatusError{StatusCode: http.StatusForbidden},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("IsNotFound() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsForbidden(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "non-HTTP error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "HTTP 403 error",
			err:      &HTTPStatusError{StatusCode: http.StatusForbidden},
			expected: true,
		},
		{
			name:     "HTTP 404 error",
			err:      &HTTPStatusError{StatusCode: http.StatusNotFound},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsForbidden(tt.err)
			if result != tt.expected {
				t.Errorf("IsForbidden() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsConflict(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "non-HTTP error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "HTTP 409 error",
			err:      &HTTPStatusError{StatusCode: http.StatusConflict},
			expected: true,
		},
		{
			name:     "HTTP 404 error",
			err:      &HTTPStatusError{StatusCode: http.StatusNotFound},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsConflict(tt.err)
			if result != tt.expected {
				t.Errorf("IsConflict() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHandleAPIError(t *testing.T) {
	client := &RackspaceSpotClient{}
	
	tests := []struct {
		name     string
		err      error
		resource string
		id       string
		op       string
		expected string
	}{
		{
			name:     "HTTP 404 error",
			err:      &HTTPStatusError{StatusCode: http.StatusNotFound, Status: "404 Not Found"},
			resource: "cloudspace",
			id:       "test",
			op:       "get",
			expected: "failed to get cloudspace 'test': HTTP error 404: 404 Not Found",
		},
		{
			name:     "HTTP 403 error",
			err:      &HTTPStatusError{StatusCode: http.StatusForbidden, Status: "403 Forbidden"},
			resource: "cloudspace",
			id:       "test",
			op:       "create",
			expected: "failed to create cloudspace 'test': HTTP error 403: 403 Forbidden",
		},
		{
			name:     "HTTP 409 error",
			err:      &HTTPStatusError{StatusCode: http.StatusConflict, Status: "409 Conflict"},
			resource: "cloudspace",
			id:       "test",
			op:       "update",
			expected: "failed to update cloudspace 'test': HTTP error 409: 409 Conflict",
		},
		{
			name:     "non-HTTP error",
			err:      errors.New("some error"),
			resource: "cloudspace",
			id:       "test",
			op:       "delete",
			expected: "failed to delete cloudspace 'test': some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.handleAPIError(tt.err, tt.resource, tt.id, tt.op)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if err.Error() != tt.expected {
				t.Errorf("handleAPIError() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}
