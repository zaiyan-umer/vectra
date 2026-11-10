package backend

import (
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL string

	mu sync.RWMutex
	healthy bool

	activeConnections atomic.Int64
	responseTime atomic.Int64
}

func Targets(targetURLs []string) []*Backend {
	backends := make([]*Backend, 0, len(targetURLs))

	for _, targetURL := range targetURLs {
		backends = append(backends, &Backend{URL: targetURL, healthy: false})
	}

	return backends
}

func (b *Backend) IsHealthy() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.healthy
}

func (b *Backend) SetHealthy(healthy bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.healthy = healthy
}

func (b *Backend) ActiveConnections() int64 {
	return b.activeConnections.Load()
}

func (b *Backend) Acquire() {
	b.activeConnections.Add(1)
}

func (b *Backend) Release() {
	b.activeConnections.Add(-1)
}

func (b *Backend) ResponseTime() time.Duration {
	return time.Duration(b.responseTime.Load())
}

func (b *Backend) UpdateResponseTime(sample time.Duration, alpha float64) {
    for {
        previous := b.responseTime.Load()

        var next int64
        if previous == 0 {
            next = int64(sample)
        } else {
            next = int64(alpha*float64(sample) + (1-alpha)*float64(previous))
        }

        if b.responseTime.CompareAndSwap(previous, next) {
            return
        }
    }
}