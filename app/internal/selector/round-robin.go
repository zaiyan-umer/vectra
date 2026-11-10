package selector

import (
	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"sync/atomic"
)

type RoundRobin struct {
	backends []backend.Backend
	next     atomic.Uint64
}

func (r *RoundRobin) Select() int {
	backendLength := len(r.backends)

	if backendLength == 0 {
		return -1
	}

	index := (r.next.Add(1) - 1) % uint64(backendLength)
	return int(index)
}
