package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestRunRejectsNonPositiveInterval(t *testing.T) {
	for _, interval := range []time.Duration{0, -time.Second} {
		t.Run(interval.String(), func(t *testing.T) {
			err := Run(context.Background(), nil, interval)
			if err == nil {
				t.Fatal("Run returned nil error for non-positive interval")
			}
			if !strings.Contains(err.Error(), "health-check interval must be positive") {
				t.Fatalf("error = %q, want interval validation error", err)
			}
		})
	}
}

func TestRunStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- Run(ctx, nil, time.Millisecond)
	}()

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after context cancellation")
	}
}
