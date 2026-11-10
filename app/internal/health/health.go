package health

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

func HealthCheck(ctx context.Context, backends []*backend.Backend) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(len(backends))

	for _, target := range backends {
		go func(target *backend.Backend) {
			defer wg.Done()
			url := target.URL + "/health"

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				target.SetHealthy(false)
				return
			}

			res, err := client.Do(req)
			if res != nil {
				res.Body.Close()
			}

			if err != nil || res == nil || res.StatusCode < 200 || res.StatusCode > 399 {
				target.SetHealthy(false)
				log.Printf("HEALTH CHECK FAILED - %s\n", target.URL)
				return
			}

			target.SetHealthy(true)
			log.Printf("HEALTH CHECK - %s, Healthy: %t\n", target.URL, target.IsHealthy())
		}(target)
	}
	wg.Wait()
}

func Run(ctx context.Context, backends []*backend.Backend, interval time.Duration) error {
	if interval <= 0 {
		return fmt.Errorf("health-check interval must be positive: %s", interval)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	HealthCheck(ctx, backends)

	for {
		select {
		case <-ticker.C:
			HealthCheck(ctx, backends)
		case <-ctx.Done():
			return nil
		}
	}
}
