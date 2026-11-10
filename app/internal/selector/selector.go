package selector

import (
	"fmt"
	"net/http"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/proxy"
)

func New(algorithm string, targetBackends []*backend.Backend) (http.Handler, error) {
	switch algorithm {
	case config.AlgorithmRoundRobin:
		roundRobin := &RoundRobin{
			backends: targetBackends,
		}

		return proxy.NewProxy(roundRobin, targetBackends)

	case config.AlgorithmHealthyRoundRobin:
		HealthyRoundRobin := &HealthyRoundRobin{
			backends: targetBackends,
		}

		return proxy.NewProxy(HealthyRoundRobin, targetBackends)

	case config.AlgorithmLeastConnections:
		leastConnections := &LeastConnections{
			backends: targetBackends,
		}

		return proxy.NewProxy(leastConnections, targetBackends)

	case config.AlgorithmLatencyAware:
		latencyAware := &LatencyAware{
			backends: targetBackends,
		}

		return proxy.NewProxy(latencyAware, targetBackends)
		
	default:
		return nil, fmt.Errorf("unsupported algorithm: %q", algorithm)
	}
}
