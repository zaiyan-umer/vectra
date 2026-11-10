package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

func TestHealthCheckMarksBackendsHealthyAndUnhealthy(t *testing.T) {
	healthyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("health-check path = %q, want /health", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer healthyServer.Close()

	unhealthyServer := httptest.NewServer(http.NotFoundHandler())
	defer unhealthyServer.Close()

	backends := []*backend.Backend{
		{URL: healthyServer.URL},
		{URL: unhealthyServer.URL},
		{URL: "://invalid"},
	}

	HealthCheck(context.Background(), backends)

	if !backends[0].IsHealthy() {
		t.Error("healthy backend was marked unhealthy")
	}
	if backends[1].IsHealthy() {
		t.Error("unhealthy backend was marked healthy")
	}
	if backends[2].IsHealthy() {
		t.Error("invalid backend was marked healthy")
	}
}

func TestHealthCheckHonorsCanceledContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("health check started after context cancellation")
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	backends := []*backend.Backend{{URL: server.URL}}
	HealthCheck(ctx, backends)

	if backends[0].IsHealthy() {
		t.Error("canceled health check marked backend healthy")
	}
}
