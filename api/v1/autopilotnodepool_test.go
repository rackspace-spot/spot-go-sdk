package rxtspot

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// captureUpdate runs UpdateAutopilotNodePool against a fake server and returns the PATCH body.
func captureUpdate(t *testing.T, pool AutopilotNodePool) string {
	t.Helper()
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte(`{"organizations":[{"id":"org-1","name":"myorg"}]}`))
		case http.MethodPatch:
			b, _ := io.ReadAll(r.Body)
			got = string(b)
			w.Write([]byte(`{}`))
		}
	}))
	defer srv.Close()

	c := &RackspaceSpotClient{BaseURL: srv.URL, HTTPClient: srv.Client(), Token: "t"}
	if err := c.UpdateAutopilotNodePool(context.Background(), "myorg", pool); err != nil {
		t.Fatalf("update: %v", err)
	}
	if got == "" {
		t.Fatal("no PATCH received")
	}
	return got
}

func TestUpdateAutopilotNodePool_BudgetOnlyOmitsVCPU(t *testing.T) {
	got := captureUpdate(t, AutopilotNodePool{Name: "ap", BudgetPerHour: "0.5"})
	if want := `{"spec":{"vcpuPerNode":{},"budgetPerHour":"0.5"}}`; got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
}

func TestUpdateAutopilotNodePool_VCPUSentWhenSet(t *testing.T) {
	got := captureUpdate(t, AutopilotNodePool{Name: "ap", VCPUTotal: 16})
	if want := `{"spec":{"vcpu":{"total":16},"vcpuPerNode":{}}}`; got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
}
