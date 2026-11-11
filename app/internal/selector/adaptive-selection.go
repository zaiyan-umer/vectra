package selector

import (
	"log"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/health"
)

type AdaptiveSelection struct {
	backends []*backend.Backend
}

func (r *AdaptiveSelection) Select() int {
	healthy := health.GetHealthyBackendIndices(r.backends)

	if len(healthy) == 0 {
		return -1
	}

	scores := map[int]float64{}
	maxExpectedConnections := config.AdaptiveConfiguration.MaxExpectedConnections
	latencyTarget := config.AdaptiveConfiguration.LatencyTarget

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

	for index, score := range scores {
		if score < lowest {
			selected = index
			lowest = score
		}
	}

	return selected
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
