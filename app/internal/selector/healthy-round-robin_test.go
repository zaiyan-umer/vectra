package selector

import (
	"sync"
	"testing"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

func TestHealthyRoundRobinCyclesThroughHealthyBackends(t *testing.T) {
	backends := backend.Targets([]string{"http://one", "http://two", "http://three"})
	backends[0].SetHealthy(true)
	backends[2].SetHealthy(true)

	roundRobin := &HealthyRoundRobin{backends: backends}
	want := []int{0, 2, 0, 2}

	for i, expected := range want {
		if got := roundRobin.Acquire(); got != expected {
			t.Errorf("selection %d = %d, want %d", i, got, expected)
		}
	}
}

func TestHealthyRoundRobinReturnsUnavailableWhenAllBackendsAreUnhealthy(t *testing.T) {
	roundRobin := &HealthyRoundRobin{
		backends: backend.Targets([]string{"http://one", "http://two"}),
	}

	if got := roundRobin.Acquire(); got != -1 {
		t.Fatalf("selection = %d, want -1", got)
	}
}

func TestHealthyRoundRobinWithNoBackends(t *testing.T) {
	roundRobin := &HealthyRoundRobin{}

	if got := roundRobin.Acquire(); got != -1 {
		t.Fatalf("selection = %d, want -1", got)
	}
}

func TestHealthyRoundRobinConcurrentSelection(t *testing.T) {
	const selections = 1000

	backends := backend.Targets([]string{"http://one", "http://two"})
	for _, target := range backends {
		target.SetHealthy(true)
	}

	roundRobin := &HealthyRoundRobin{backends: backends}
	counts := make([]int, len(backends))
	var countsMu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(selections)

	for range selections {
		go func() {
			defer wg.Done()
			index := roundRobin.Acquire()
			countsMu.Lock()
			counts[index]++
			countsMu.Unlock()
		}()
	}

	wg.Wait()

	for index, count := range counts {
		if count != selections/len(backends) {
			t.Errorf("backend %d received %d selections, want %d", index, count, selections/len(backends))
		}
	}
}
