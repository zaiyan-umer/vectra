package selector

import (
	"fmt"
	"net/http"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/metrics"
	"github.com/zaiyan-umer/vectra/app/internal/proxy"
)

func New(algorithm string, targetBackends []*backend.Backend, benchmarkResult *metrics.BenchmarkRuntimeSnapshot) (http.Handler, error) {
	switch algorithm {
	case config.AlgorithmRoundRobin:
		roundRobin := &RoundRobin{
			backends: targetBackends,
		}

		return proxy.NewProxy(roundRobin, targetBackends, benchmarkResult)

	case config.AlgorithmHealthyRoundRobin:
		HealthyRoundRobin := &HealthyRoundRobin{
			backends: targetBackends,
		}

		return proxy.NewProxy(HealthyRoundRobin, targetBackends, benchmarkResult)

	case config.AlgorithmLeastConnections:
		leastConnections := &LeastConnections{
			backends: targetBackends,
		}

		return proxy.NewProxy(leastConnections, targetBackends, benchmarkResult)

	case config.AlgorithmLatencyAware:
		latencyAware := &LatencyAware{
			backends: targetBackends,
		}

		return proxy.NewProxy(latencyAware, targetBackends, benchmarkResult)

	case config.AdaptiveSelection:
		err := config.ValidateAdaptiveConfiguration(config.AdaptiveConfiguration)

		if err != nil {
			return nil, fmt.Errorf("invalid adaptive configuration: %w", err)
		}

		adaptiveSelection := &AdaptiveSelection{
			backends: targetBackends,
		}

		return proxy.NewProxy(adaptiveSelection, targetBackends, benchmarkResult)

	default:
		return nil, fmt.Errorf("unsupported algorithm: %q", algorithm)
	}
}
