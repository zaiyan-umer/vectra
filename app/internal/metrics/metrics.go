package metrics

import (
	"fmt"
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
