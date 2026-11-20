package metrics

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

const (
	maxTimeWindows = 60 // 5 min at 5s intervals
)

type BackendRuntimeCounters struct {
	BackendIndex      int
	URL               string
	SelectedRequests  atomic.Uint64
	CompletedRequests atomic.Uint64
	Errors            atomic.Uint64
	TotalLatencyNs    atomic.Int64
}

type BenchmarkRuntimeSnapshot struct {
	StartedAt  time.Time
	FinishedAt time.Time
	Backends   []BackendRuntimeCounters

	mu      sync.Mutex
	buffer  []TimeWindow
	cursor  int
	first   bool
}

type BackendBenchmarkStats struct {
	BackendIndex      int     `json:"index"`
	URL               string  `json:"url"`
	SelectedRequests  uint64  `json:"assignedRequests"`
	CompletedRequests uint64  `json:"completedRequests"`
	Errors            uint64  `json:"errors"`
	AverageLatencyMs  float64 `json:"averageLatencyMs"`
	TrafficShare      float64 `json:"trafficShare"`
}

type TimeWindow struct {
	Timestamp          time.Time  `json:"timestamp"`
	PoolRPS            float64    `json:"poolRps"`
	PoolErrorRate      float64    `json:"poolErrorRate"`
	PoolAvgLatencyMs   float64    `json:"poolAvgLatencyMs"`
	Backends           []BackendWindow `json:"backends"`
}

type BackendWindow struct {
	Index           int     `json:"index"`
	RPS             float64 `json:"rps"`
	ErrorRate       float64 `json:"errorRate"`
	AvgLatencyMs    float64 `json:"avgLatencyMs"`
}

type PoolMetrics struct {
	RPS           float64 `json:"rps"`
	ErrorRate     float64 `json:"errorRate"`
	AvgLatencyMs  float64 `json:"avgLatencyMs"`
}

type BackendMetrics struct {
	Index            int     `json:"index"`
	URL              string  `json:"url"`
	Healthy          bool    `json:"healthy"`
	ActiveConnections int64  `json:"activeConnections"`
	RPS              float64 `json:"rps"`
	ErrorRate        float64 `json:"errorRate"`
	AvgLatencyMs     float64 `json:"avgLatencyMs"`
	TrafficShare     float64 `json:"trafficShare"`
	SelectedRequests  uint64  `json:"selectedRequests"`
	CompletedRequests uint64  `json:"completedRequests"`
	Errors           uint64  `json:"errors"`
}

func NewBenchmarkRuntimeSnapshot(backends []*backend.Backend) (*BenchmarkRuntimeSnapshot, error) {
	result := &BenchmarkRuntimeSnapshot{
		StartedAt: time.Now(),
		Backends:  make([]BackendRuntimeCounters, len(backends)),
		buffer:    make([]TimeWindow, 0, maxTimeWindows),
		first:     true,
	}

	for index, target := range backends {
		if target == nil {
			return nil, fmt.Errorf("backend at index %d is nil", index)
		}
		result.Backends[index] = BackendRuntimeCounters{
			BackendIndex: index,
			URL:          target.URL,
		}
	}

	return result, nil
}

func (b *BenchmarkRuntimeSnapshot) validateIndex(index int) error {
	if index < 0 || index >= len(b.Backends) {
		return fmt.Errorf("backend index %d is out of range", index)
	}
	return nil
}

func (b *BenchmarkRuntimeSnapshot) RecordSelected(index int) error {
	if err := b.validateIndex(index); err != nil {
		return err
	}

	b.Backends[index].SelectedRequests.Add(1)
	return nil
}

func (b *BenchmarkRuntimeSnapshot) RecordCompleted(index int, latency time.Duration) error {
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

func (b *BenchmarkRuntimeSnapshot) RecordError(index int) error {
	if err := b.validateIndex(index); err != nil {
		return err
	}

	b.Backends[index].Errors.Add(1)
	return nil
}

func (b *BenchmarkRuntimeSnapshot) CaptureSnapshot() BenchmarkRuntimeSnapshot {
	snapshot := BenchmarkRuntimeSnapshot{
		StartedAt: b.StartedAt,
		Backends:  make([]BackendRuntimeCounters, len(b.Backends)),
	}

	for i := range b.Backends {
		current := &b.Backends[i]
		copy := BackendRuntimeCounters{
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

func (b *BenchmarkRuntimeSnapshot) Difference(old BenchmarkRuntimeSnapshot) []BackendBenchmarkStats {
	stats := make([]BackendBenchmarkStats, len(b.Backends))
	totalSelected := uint64(0)

	for i := range b.Backends {
		current := &b.Backends[i]
		previous := &old.Backends[i]
		selected := current.SelectedRequests.Load() - previous.SelectedRequests.Load()
		completed := current.CompletedRequests.Load() - previous.CompletedRequests.Load()
		latencyNs := current.TotalLatencyNs.Load() - previous.TotalLatencyNs.Load()

		stats[i] = BackendBenchmarkStats{
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

func (b *BenchmarkRuntimeSnapshot) StartCapture(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var prev BenchmarkRuntimeSnapshot
		for range ticker.C {
			curr := b.CaptureSnapshot()
			if !b.first {
				stats := b.Difference(prev)
				b.addTimeWindow(stats, interval)
			}
			prev = curr
			b.first = false
		}
	}()
}

func (b *BenchmarkRuntimeSnapshot) addTimeWindow(stats []BackendBenchmarkStats, interval time.Duration) {
	var totalRPS float64
	var totalErrors uint64
	var totalCompleted uint64
	var totalLatencyMs float64

	backendWindows := make([]BackendWindow, len(stats))

	for i, s := range stats {
		rps := float64(s.SelectedRequests) / interval.Seconds()
		errorRate := float64(0)
		if s.SelectedRequests > 0 {
			errorRate = float64(s.Errors) / float64(s.SelectedRequests)
		}
		avgLatency := s.AverageLatencyMs

		backendWindows[i] = BackendWindow{
			Index:        s.BackendIndex,
			RPS:          rps,
			ErrorRate:    errorRate,
			AvgLatencyMs: avgLatency,
		}

		totalRPS += rps
		totalErrors += s.Errors
		totalCompleted += s.CompletedRequests
		totalLatencyMs += avgLatency * float64(s.CompletedRequests)
	}

	poolErrorRate := float64(0)
	if totalCompleted > 0 {
		poolErrorRate = float64(totalErrors) / float64(totalCompleted)
	}

	poolAvgLatency := float64(0)
	if totalCompleted > 0 {
		poolAvgLatency = totalLatencyMs / float64(totalCompleted)
	}

	window := TimeWindow{
		Timestamp:         time.Now(),
		PoolRPS:           totalRPS,
		PoolErrorRate:     poolErrorRate,
		PoolAvgLatencyMs:  poolAvgLatency,
		Backends:          backendWindows,
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.buffer) < maxTimeWindows {
		b.buffer = append(b.buffer, window)
	} else {
		b.buffer[b.cursor] = window
		b.cursor = (b.cursor + 1) % maxTimeWindows
	}
}

func (b *BenchmarkRuntimeSnapshot) GetPoolMetrics() PoolMetrics {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.buffer) == 0 {
		return PoolMetrics{}
	}

	latest := b.buffer[(b.cursor-1+len(b.buffer))%len(b.buffer)]
	return PoolMetrics{
		RPS:          latest.PoolRPS,
		ErrorRate:    latest.PoolErrorRate,
		AvgLatencyMs: latest.PoolAvgLatencyMs,
	}
}

func (b *BenchmarkRuntimeSnapshot) GetBackendsMetrics(backends []*backend.Backend) []BackendMetrics {
	b.mu.Lock()
	defer b.mu.Unlock()

	metrics := make([]BackendMetrics, len(b.Backends))

	var latest TimeWindow
	if len(b.buffer) > 0 {
		latest = b.buffer[(b.cursor-1+len(b.buffer))%len(b.buffer)]
	}

	// Calculate total selected requests across all backends
	var totalSelected uint64
	for i := range b.Backends {
		totalSelected += b.Backends[i].SelectedRequests.Load()
	}

	for i := range b.Backends {
		current := &b.Backends[i]
		selected := current.SelectedRequests.Load()
		completed := current.CompletedRequests.Load()
		errors := current.Errors.Load()
		latencyNs := current.TotalLatencyNs.Load()

		var avgLatencyMs float64
		if completed > 0 {
			avgLatencyMs = float64(latencyNs) / float64(completed) / float64(time.Millisecond)
		}

		var rps, errorRate, trafficShare float64
		var bw BackendWindow
		if len(latest.Backends) > i {
			bw = latest.Backends[i]
			rps = bw.RPS
			errorRate = bw.ErrorRate
		}

		if totalSelected > 0 && selected > 0 {
			trafficShare = float64(selected) / float64(totalSelected) * 100
		}

		healthy := false
		var activeConnections int64
		if i < len(backends) && backends[i] != nil {
			healthy = backends[i].IsHealthy()
			activeConnections = backends[i].ActiveConnections()
		}

		metrics[i] = BackendMetrics{
			Index:              i,
			URL:                current.URL,
			Healthy:            healthy,
			ActiveConnections:  activeConnections,
			RPS:                rps,
			ErrorRate:          errorRate,
			AvgLatencyMs:       avgLatencyMs,
			TrafficShare:       trafficShare,
			SelectedRequests:   selected,
			CompletedRequests:  completed,
			Errors:             errors,
		}
	}

	return metrics
}

func (b *BenchmarkRuntimeSnapshot) GetTimeSeries() []TimeWindow {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.buffer) == 0 {
		return []TimeWindow{}
	}

	result := make([]TimeWindow, len(b.buffer))
	for i := range b.buffer {
		idx := (b.cursor + i) % len(b.buffer)
		result[i] = b.buffer[idx]
	}

	return result
}
