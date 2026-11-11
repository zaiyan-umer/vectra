package selector

import (
	"testing"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/config"
)

func setAdaptiveTestConfiguration(t *testing.T) {
	t.Helper()

	previous := config.AdaptiveConfiguration
	config.AdaptiveConfiguration = config.AdaptiveConfig{
		MaxExpectedConnections: 10,
		LatencyTarget:          100 * time.Millisecond,
		ActiveWeight:           0.5,
		LatencyWeight:          0.5,
	}
	t.Cleanup(func() {
		config.AdaptiveConfiguration = previous
	})
}

func TestAdaptiveSelectionChoosesLowestCombinedScore(t *testing.T) {
	setAdaptiveTestConfiguration(t)

	backends := backend.Targets([]string{"http://one", "http://two", "http://three"})
	for _, target := range backends {
		target.SetHealthy(true)
	}

	backends[0].UpdateResponseTime(100*time.Millisecond, 1)
	backends[1].UpdateResponseTime(20*time.Millisecond, 1)
	backends[2].UpdateResponseTime(80*time.Millisecond, 1)
	for range 4 {
		backends[1].Acquire()
	}
	for range 2 {
		backends[2].Acquire()
	}

	adaptive := &AdaptiveSelection{backends: backends}

	if got := adaptive.Select(); got != 1 {
		t.Fatalf("selection = %d, want 1", got)
	}
}

func TestAdaptiveSelectionIgnoresUnhealthyBackends(t *testing.T) {
	setAdaptiveTestConfiguration(t)

	backends := backend.Targets([]string{"http://one", "http://two"})
	backends[0].SetHealthy(true)
	backends[0].UpdateResponseTime(100*time.Millisecond, 1)
	backends[1].UpdateResponseTime(1*time.Millisecond, 1)

	adaptive := &AdaptiveSelection{backends: backends}

	if got := adaptive.Select(); got != 0 {
		t.Fatalf("selection = %d, want healthy backend 0", got)
	}
}

func TestAdaptiveSelectionReturnsUnavailableWhenNoBackendIsHealthy(t *testing.T) {
	setAdaptiveTestConfiguration(t)

	adaptive := &AdaptiveSelection{
		backends: backend.Targets([]string{"http://one", "http://two"}),
	}

	if got := adaptive.Select(); got != -1 {
		t.Fatalf("selection = %d, want -1", got)
	}
}

func TestAdaptiveSelectionWithNoBackends(t *testing.T) {
	setAdaptiveTestConfiguration(t)

	adaptive := &AdaptiveSelection{}

	if got := adaptive.Select(); got != -1 {
		t.Fatalf("selection = %d, want -1", got)
	}
}

func TestAdaptiveSelectionAccountsForActiveConnections(t *testing.T) {
	setAdaptiveTestConfiguration(t)

	backends := backend.Targets([]string{"http://one", "http://two"})
	for _, target := range backends {
		target.SetHealthy(true)
	}
	backends[0].UpdateResponseTime(10*time.Millisecond, 1)
	backends[1].UpdateResponseTime(100*time.Millisecond, 1)
	for range 8 {
		backends[0].Acquire()
	}

	adaptive := &AdaptiveSelection{backends: backends}

	if got := adaptive.Select(); got != 0 {
		t.Fatalf("selection = %d, want 0", got)
	}
}
