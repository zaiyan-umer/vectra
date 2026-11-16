package metrics

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

type BackendBenchmark struct {
	BackendIndex      int
	URL               string
	SelectedRequests  atomic.Uint64
	CompletedRequests atomic.Uint64
	Errors            atomic.Uint64
	TotalLatencyNs    atomic.Int64
}

type BenchmarkResult struct {
	StartedAt  time.Time
	FinishedAt time.Time
	Backends   []BackendBenchmark
}

type BackendStats struct {
	BackendIndex      int     `json:"index"`
	URL               string  `json:"url"`
	SelectedRequests  uint64  `json:"assignedRequests"`
	CompletedRequests uint64  `json:"completedRequests"`
	Errors            uint64  `json:"errors"`
	AverageLatencyMs  float64 `json:"averageLatencyMs"`
	TrafficShare      float64 `json:"trafficShare"`
}

func NewBenchmarkResult(backends []*backend.Backend) (*BenchmarkResult, error) {
	result := &BenchmarkResult{
		StartedAt: time.Now(),
		Backends:  make([]BackendBenchmark, len(backends)),
	}

	for index, target := range backends {
		if target == nil {
			return nil, fmt.Errorf("backend at index %d is nil", index)
		}
		result.Backends[index] = BackendBenchmark{
			BackendIndex: index,
			URL:          target.URL,
		}
	}

	return result, nil
}

func (b *BenchmarkResult) validateIndex(index int) error {
	if index < 0 || index >= len(b.Backends) {
		return fmt.Errorf("backend index %d is out of range", index)
	}
	return nil
}

func (b *BenchmarkResult) RecordSelected(index int) error {
	if err := b.validateIndex(index); err != nil {
		return err
	}

	b.Backends[index].SelectedRequests.Add(1)
	return nil
}

func (b *BenchmarkResult) RecordCompleted(index int, latency time.Duration) error {
	if err := b.validateIndex(index); err != nil {
		return err
	}
	if latency < 0 {
		return fmt.Errorf("latency must not be negative: %s", latency)
	}

	b.Backends[index].CompletedRequests.Add(1)
	b.Backends[index].TotalLatencyNs.Add(int64(latency))
	return nil
}

func (b *BenchmarkResult) RecordError(index int) error {
	if err := b.validateIndex(index); err != nil {
		return err
	}

	b.Backends[index].Errors.Add(1)
	return nil
}

func (b *BenchmarkResult) CaptureSnapshot() BenchmarkResult {
	snapshot := BenchmarkResult{
		StartedAt: b.StartedAt,
		Backends:  make([]BackendBenchmark, len(b.Backends)),
	}

	for i := range b.Backends {
		current := &b.Backends[i]
		copy := BackendBenchmark{
			BackendIndex: current.BackendIndex,
			URL:          current.URL,
		}

		copy.SelectedRequests.Store(current.SelectedRequests.Load())
		copy.CompletedRequests.Store(current.CompletedRequests.Load())
		copy.Errors.Store(current.Errors.Load())
		copy.TotalLatencyNs.Store(current.TotalLatencyNs.Load())

		snapshot.Backends[i] = copy
	}

	return snapshot
}

func (b *BenchmarkResult) Difference(old BenchmarkResult) []BackendStats {
	stats := make([]BackendStats, len(b.Backends))
	totalSelected := uint64(0)

	for i := range b.Backends {
		current := &b.Backends[i]
		previous := &old.Backends[i]
		selected := current.SelectedRequests.Load() - previous.SelectedRequests.Load()
		completed := current.CompletedRequests.Load() - previous.CompletedRequests.Load()
		latencyNs := current.TotalLatencyNs.Load() - previous.TotalLatencyNs.Load()

		stats[i] = BackendStats{
			BackendIndex:      current.BackendIndex,
			URL:               current.URL,
			SelectedRequests:  selected,
			CompletedRequests: completed,
			Errors:            current.Errors.Load() - previous.Errors.Load(),
		}
		if completed > 0 {
			stats[i].AverageLatencyMs = math.Round(float64(latencyNs)/float64(completed)/float64(time.Millisecond)*100) / 100
		}
		totalSelected += selected
	}

	if totalSelected > 0 {
		for i := range stats {
			stats[i].TrafficShare = math.Round(float64(stats[i].SelectedRequests)/float64(totalSelected)*100*100) / 100
		}
	}

	return stats
}
