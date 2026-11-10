package selector

import (
	"sync"
	"testing"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

func TestRoundRobinCyclesThroughBackends(t *testing.T) {
	roundRobin := &RoundRobin{
		backends: backend.Targets([]string{"http://one", "http://two", "http://three"}),
	}

	want := []int{0, 1, 2, 0, 1}
	for i, expected := range want {
		if got := roundRobin.Acquire(); got != expected {
			t.Errorf("selection %d = %d, want %d", i, got, expected)
		}
	}
}

func TestRoundRobinWithNoBackends(t *testing.T) {
	roundRobin := &RoundRobin{}

	if got := roundRobin.Acquire(); got != -1 {
		t.Fatalf("got %d, want -1", got)
	}
}

func TestRoundRobinIsHealthUnaware(t *testing.T) {
	roundRobin := &RoundRobin{
		backends: backend.Targets([]string{"http://one", "http://two"}),
	}

	if got := roundRobin.Acquire(); got != 0 {
		t.Fatalf("selection = %d, want 0", got)
	}
}

func TestRoundRobinConcurrentSelection(t *testing.T) {
	const (
		backendCount = 4
		selections   = 1000
	)

	roundRobin := &RoundRobin{
		backends: backend.Targets([]string{"http://one", "http://two", "http://three", "http://four"}),
	}
	counts := make([]int, backendCount)
	var countsMu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(selections)

	for range selections {
		go func() {
			defer wg.Done()
			selection := roundRobin.Acquire()
			if selection >= 0 && selection < backendCount {
				countsMu.Lock()
				counts[selection]++
				countsMu.Unlock()
			}
		}()
	}

	wg.Wait()

	for i, count := range counts {
		if count != selections/backendCount {
			t.Errorf("backend %d received %d selections, want %d", i, count, selections/backendCount)
		}
	}
}
