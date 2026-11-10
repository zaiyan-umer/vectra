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
	// case config.AlgorithmLeastConnections:
	//     return NewLeastConnectionsProxy(targetBackends)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %q", algorithm)
	}
}
