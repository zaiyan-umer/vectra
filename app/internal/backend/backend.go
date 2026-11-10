package backend

import "sync"

type Backend struct {
	URL string

	mu sync.RWMutex
	healthy bool
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
