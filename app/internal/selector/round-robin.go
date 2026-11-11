package selector

import (
	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"sync/atomic"
)

type RoundRobin struct {
	backends []*backend.Backend
	next     atomic.Uint64
}

func (r *RoundRobin) Select() int {
	backendLength := len(r.backends)

	if backendLength == 0 {
		return -1
	}

	return NextRoundRobinIndex(backendLength, &r.next)
}

func NextRoundRobinIndex(length int, next *atomic.Uint64) int {
    if length == 0 {
        return -1
    }

    return int((next.Add(1) - 1) % uint64(length))
}