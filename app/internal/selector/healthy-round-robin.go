package selector

import (
	"sync/atomic"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/health"
)

type HealthyRoundRobin struct {
	backends []*backend.Backend
	next     atomic.Uint64
}

func (r *HealthyRoundRobin) Select() int {
	healthy := health.GetHealthyBackendIndices(r.backends)

	if len(healthy) == 0 {
		return -1
	}

	position := NextRoundRobinIndex(len(healthy), &r.next)
	return healthy[position]
}
