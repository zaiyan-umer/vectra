package selector

import (
	"log"
	"math"
	"sync"
	"sync/atomic"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/health"
)

type AdaptiveSelection struct {
	backends []*backend.Backend
	next     atomic.Uint64
	mu       sync.Mutex
}

func (r *AdaptiveSelection) Select() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.selectLocked()
}

func (r *AdaptiveSelection) SelectAndAcquire() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	selected := r.selectLocked()
	if selected >= 0 {
		r.backends[selected].Acquire()
	}

	return selected
}

func (r *AdaptiveSelection) selectLocked() int {
	healthy := health.GetHealthyBackendIndices(r.backends)

	if len(healthy) == 0 {
		return -1
	}

	scores := map[int]float64{}
	maxExpectedConnections := config.AdaptiveConfiguration.MaxExpectedConnections
	latencyTarget := config.AdaptiveConfiguration.LatencyTarget
	tolerance := config.AdaptiveConfiguration.Tolerance

	for _, index := range healthy {
		activeConnections := r.backends[index].ActiveConnections()
		latency := r.backends[index].ResponseTime()

		normalizedActive := clamp(
			float64(activeConnections)/float64(maxExpectedConnections),
			0, 1,
		)

		normalizedLatency := clamp(
			float64(latency)/float64(latencyTarget),
			0, 1,
		)

		score := config.AdaptiveConfiguration.ActiveWeight*normalizedActive + config.AdaptiveConfiguration.LatencyWeight*normalizedLatency

		scores[index] = score
		log.Printf(
			"adaptive backend=%d active_connections=%d latency=%s score=%.4f",
			index+1,
			activeConnections,
			latency,
			score,
		)
	}

	selected := healthy[0]
	lowest := scores[selected]

	for _, index := range healthy {
		if scores[index] < lowest {
			selected = index
			lowest = scores[index]
		}
	}

	tied := make([]int, 0)

	for _, index := range healthy {
		if math.Abs(scores[index]-lowest) <= tolerance {
			tied = append(tied, index)
		}
	}

	if len(tied) == 1 {
		return selected
	}

	position := NextRoundRobinIndex(len(tied), &r.next)
	return tied[position]
}

func clamp(value, lower, upper float64) float64 {
	if value < lower {
		return lower
	}
	if value > upper {
		return upper
	}
	return value
}
