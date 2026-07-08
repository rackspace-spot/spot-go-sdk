package rxtspot

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newPatchCaptureClient(t *testing.T, captured *[]byte) *RackspaceSpotClient {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/apis/auth.ngpc.rxt.io/v1/organizations":
			w.Write([]byte(`{"organizations":[{"name":"testorg","id":"testorg"}]}`))
		case r.Method == http.MethodPatch:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("failed to read PATCH body: %v", err)
			}
			*captured = body
			w.Write([]byte(`{}`))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)

	return &RackspaceSpotClient{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
		Token:      "test-token",
		RetryConfig: RetryConfig{
			MaxRetries:      1,
			InitialInterval: time.Millisecond,
			Multiplier:      2,
			RetryWaitMax:    time.Second,
		},
	}
}

func specAutoscaling(t *testing.T, body []byte) (map[string]any, bool) {
	t.Helper()
	var parsed struct {
		Spec map[string]json.RawMessage `json:"spec"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("failed to parse update body %s: %v", body, err)
	}
	raw, ok := parsed.Spec["autoscaling"]
	if !ok {
		return nil, false
	}
	var autoscaling map[string]any
	if err := json.Unmarshal(raw, &autoscaling); err != nil {
		t.Fatalf("failed to parse autoscaling %s: %v", raw, err)
	}
	return autoscaling, true
}

func TestUpdateNodePoolAutoscalingPatchBody(t *testing.T) {
	cases := []struct {
		name        string
		autoscaling *Autoscaling
		want        map[string]any
	}{
		{
			name:        "nil autoscaling omitted from patch",
			autoscaling: nil,
			want:        nil,
		},
		{
			name:        "enabled autoscaling sent with bounds",
			autoscaling: &Autoscaling{Enabled: true, MinNodes: 1, MaxNodes: 5},
			want:        map[string]any{"enabled": true, "minNodes": float64(1), "maxNodes": float64(5)},
		},
		{
			name:        "explicit disable sends enabled false",
			autoscaling: &Autoscaling{Enabled: false},
			want:        map[string]any{"enabled": false},
		},
	}

	updaters := map[string]func(c *RackspaceSpotClient, autoscaling *Autoscaling) error{
		"spot": func(c *RackspaceSpotClient, autoscaling *Autoscaling) error {
			return c.UpdateSpotNodePool(context.Background(), "testorg", SpotNodePool{
				Name:        "pool-1",
				Desired:     3,
				Autoscaling: autoscaling,
			})
		},
		"ondemand": func(c *RackspaceSpotClient, autoscaling *Autoscaling) error {
			return c.UpdateOnDemandNodePool(context.Background(), "testorg", OnDemandNodePool{
				Name:        "pool-1",
				Desired:     3,
				Autoscaling: autoscaling,
			})
		},
	}

	for poolType, update := range updaters {
		for _, tc := range cases {
			t.Run(poolType+"/"+tc.name, func(t *testing.T) {
				var captured []byte
				client := newPatchCaptureClient(t, &captured)

				if err := update(client, tc.autoscaling); err != nil {
					t.Fatalf("update failed: %v", err)
				}
				if captured == nil {
					t.Fatal("no PATCH request captured")
				}

				got, present := specAutoscaling(t, captured)
				if tc.want == nil {
					if present {
						t.Errorf("expected autoscaling omitted from patch body, got %s", captured)
					}
					return
				}
				if !present {
					t.Fatalf("expected autoscaling in patch body, got %s", captured)
				}
				if len(got) != len(tc.want) {
					t.Errorf("autoscaling = %v, want %v", got, tc.want)
				}
				for k, v := range tc.want {
					if got[k] != v {
						t.Errorf("autoscaling[%q] = %v, want %v", k, got[k], v)
					}
				}
			})
		}
	}
}
