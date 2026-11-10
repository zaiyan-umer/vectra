package selector

import (
	"testing"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

func TestLeastConnectionsSelectsHealthyBackendWithFewestConnections(t *testing.T) {
	backends := backend.Targets([]string{"http://one", "http://two", "http://three"})
	for _, target := range backends {
		target.SetHealthy(true)
	}

	backends[0].Acquire()
	backends[0].Acquire()
	backends[1].Acquire()

	leastConnections := &LeastConnections{backends: backends}

	if got := leastConnections.Select(); got != 2 {
		t.Fatalf("selection = %d, want 2", got)
	}
}

func TestLeastConnectionsIgnoresUnhealthyBackends(t *testing.T) {
	backends := backend.Targets([]string{"http://one", "http://two"})
	backends[0].SetHealthy(true)
	backends[0].Acquire()
	backends[1].Acquire()

	leastConnections := &LeastConnections{backends: backends}

	if got := leastConnections.Select(); got != 0 {
		t.Fatalf("selection = %d, want healthy backend 0", got)
	}
}

func TestLeastConnectionsReturnsUnavailableWhenNoBackendIsHealthy(t *testing.T) {
	leastConnections := &LeastConnections{
		backends: backend.Targets([]string{"http://one", "http://two"}),
	}

	if got := leastConnections.Select(); got != -1 {
		t.Fatalf("selection = %d, want -1", got)
	}
}

func TestLeastConnectionsWithNoBackends(t *testing.T) {
	leastConnections := &LeastConnections{}

	if got := leastConnections.Select(); got != -1 {
		t.Fatalf("selection = %d, want -1", got)
	}
}
