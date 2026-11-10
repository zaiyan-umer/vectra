package selector

import (
	"testing"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

func TestLatencyAwareSelectsHealthyBackendWithLowestLatency(t *testing.T) {
	backends := backend.Targets([]string{"http://one", "http://two", "http://three"})
	for _, target := range backends {
		target.SetHealthy(true)
	}
	backends[0].UpdateResponseTime(300*time.Millisecond, 1)
	backends[1].UpdateResponseTime(100*time.Millisecond, 1)
	backends[2].UpdateResponseTime(200*time.Millisecond, 1)

	latencyAware := &LatencyAware{backends: backends}

	if got := latencyAware.Select(); got != 1 {
		t.Fatalf("selection = %d, want 1", got)
	}
}

func TestLatencyAwareIgnoresUnhealthyBackend(t *testing.T) {
	backends := backend.Targets([]string{"http://one", "http://two"})
	backends[0].SetHealthy(true)
	backends[0].UpdateResponseTime(500*time.Millisecond, 1)
	backends[1].UpdateResponseTime(100*time.Millisecond, 1)

	latencyAware := &LatencyAware{backends: backends}

	if got := latencyAware.Select(); got != 0 {
		t.Fatalf("selection = %d, want healthy backend 0", got)
	}
}

func TestLatencyAwareReturnsUnavailableWhenNoBackendIsHealthy(t *testing.T) {
	latencyAware := &LatencyAware{
		backends: backend.Targets([]string{"http://one", "http://two"}),
	}

	if got := latencyAware.Select(); got != -1 {
		t.Fatalf("selection = %d, want -1", got)
	}
}

func TestLatencyAwareWithNoBackends(t *testing.T) {
	latencyAware := &LatencyAware{}

	if got := latencyAware.Select(); got != -1 {
		t.Fatalf("selection = %d, want -1", got)
	}
}

func TestLatencyAwareSelectsFirstBackendWhenLatenciesAreEqual(t *testing.T) {
	backends := backend.Targets([]string{"http://one", "http://two"})
	for _, target := range backends {
		target.SetHealthy(true)
		target.UpdateResponseTime(100*time.Millisecond, 1)
	}

	latencyAware := &LatencyAware{backends: backends}

	if got := latencyAware.Select(); got != 0 {
		t.Fatalf("selection = %d, want 0", got)
	}
}
