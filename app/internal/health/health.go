package health

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
)

func HealthCheck(ctx context.Context, backends []*backend.Backend) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	now := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(backends))

	for _, target := range backends {
		go func(target *backend.Backend) {
			defer wg.Done()
			url := target.URL + "/health"

			target.SetLastCheck(now)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				target.SetHealthy(false)
				target.IncrementFailures()
				slog.Error("health check request creation failed", "backend_url", target.URL, "url", url, "err", err)
				return
			}

			res, err := client.Do(req)
			if res != nil {
				res.Body.Close()
			}

			if err != nil || res == nil || res.StatusCode < 200 || res.StatusCode > 399 {
				target.SetHealthy(false)
				target.IncrementFailures()
				statusCode := 0
				if res != nil {
					statusCode = res.StatusCode
				}
				slog.Error("health check failed", "backend_url", target.URL, "url", url, "status_code", statusCode, "err", err)
				return
			}

			target.SetHealthy(true)
			target.ResetFailures()
			slog.Info("health check completed", "backend_url", target.URL, "url", url, "status_code", res.StatusCode, "healthy", true)
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

func GetHealthyBackendIndices(backends []*backend.Backend) []int {
	indices := make([]int, 0, len(backends))

	for index, target := range backends {
		if target.IsHealthy() {
			indices = append(indices, index)
		}
	}

	return indices
}

type BackendHealthStatus struct {
	Index               int       `json:"index"`
	URL                 string    `json:"url"`
	Healthy             bool      `json:"healthy"`
	LastCheck           time.Time `json:"lastCheck"`
	ConsecutiveFailures int       `json:"consecutiveFailures"`
}

func GetPoolHealth(backends []*backend.Backend) []BackendHealthStatus {
	statuses := make([]BackendHealthStatus, 0, len(backends))

	for index, target := range backends {
		statuses = append(statuses, BackendHealthStatus{
			Index:               index,
			URL:                 target.URL,
			Healthy:             target.IsHealthy(),
			LastCheck:           target.LastCheck(),
			ConsecutiveFailures: target.ConsecutiveFailures(),
		})
	}

	return statuses
}
