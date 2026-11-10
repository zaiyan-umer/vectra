package selector

import (
	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/health"
)

type LatencyAware struct {
	backends []*backend.Backend
}

func (r *LatencyAware) Select() int {
	healthy := health.GetHealthyBackendIndices(r.backends)

	if len(healthy) == 0 {
		return -1
	}

	selected := healthy[0]
	lowest := r.backends[selected].ResponseTime()

	for _, index := range healthy[1:] {
		connections := r.backends[index].ResponseTime()
		if connections < lowest {
			selected = index
			lowest = connections
		}
	}

	return selected
}