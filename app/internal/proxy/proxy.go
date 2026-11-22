package proxy

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/httpx"
	"github.com/zaiyan-umer/vectra/app/internal/metrics"
)

// Any type that has a Select() method satisfies BackendSelector interface.
// So each algorithm must implement Select() so it's Select() method can be used here
type BackendSelector interface {
	Select() int
}

type AcquiringBackendSelector interface {
	SelectAndAcquire() int
}

type Proxy struct {
	selector         BackendSelector
	backends         []*backend.Backend
	proxies          []*httputil.ReverseProxy
	benchmarkResults *metrics.BenchmarkRuntimeSnapshot
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	acquiringSelector, acquires := p.selector.(AcquiringBackendSelector)

	var index int
	if acquires {
		index = acquiringSelector.SelectAndAcquire()
	} else {
		index = p.selector.Select()
	}

	if index < 0 || index >= len(p.proxies) {
		slog.Error("no backend available", "backend_index", index, "backend_count", len(p.proxies))
		httpx.WriteError(w, http.StatusServiceUnavailable, "no backend is available", httpx.ErrCodeUnavailable)
		return
	}

	target := p.backends[index]
	p.benchmarkResults.RecordSelected(index)

	if !acquires {
		target.Acquire()
	}
	defer target.Release()

	start := time.Now()
	defer func() {
		target.UpdateResponseTime(time.Since(start), 0.25)
	}()

	p.proxies[index].ServeHTTP(w, r)

	latency := time.Since(start)
	p.benchmarkResults.RecordCompleted(index, latency)

	backendMetrics := &p.benchmarkResults.Backends[index]

	completed := backendMetrics.CompletedRequests.Load()
	averageLatency := time.Duration(0)

	if completed > 0 {
		averageLatency = time.Duration(
			backendMetrics.TotalLatencyNs.Load() / int64(completed),
		)
	}

	slog.Info("proxy request completed",
		"backend_index", index,
		"selected_requests", backendMetrics.SelectedRequests.Load(),
		"completed_requests", completed,
		"errors", backendMetrics.Errors.Load(),
		"average_latency", averageLatency,
	)
}

func NewProxy(selector BackendSelector, targetBackends []*backend.Backend, benchmarkResults ...*metrics.BenchmarkRuntimeSnapshot) (*Proxy, error) {
	if selector == nil {
		return nil, fmt.Errorf("backend selector must not be nil")
	}

	benchmarkResult, err := metrics.NewBenchmarkRuntimeSnapshot(targetBackends, 180)
	if len(benchmarkResults) > 0 && benchmarkResults[0] != nil {
		benchmarkResult = benchmarkResults[0]
	}
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,

		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		IdleConnTimeout:       90 * time.Second,
	}

	p := &Proxy{
		selector:         selector,
		backends:         targetBackends,
		proxies:          make([]*httputil.ReverseProxy, 0, len(targetBackends)),
		benchmarkResults: benchmarkResult,
	}

	for index, rawTarget := range targetBackends {
		target, err := config.ParseBackendURL(rawTarget.URL)
		if err != nil {
			return nil, err
		}

		reverseProxy := httputil.NewSingleHostReverseProxy(target)
		reverseProxy.Transport = transport

		backendIndex := index
		reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			if !errors.Is(err, context.Canceled) {
				p.benchmarkResults.RecordError(backendIndex)
			}

			handleUnhealthyBackendResponse(w, r, err, backendIndex, rawTarget)
		}

		p.proxies = append(p.proxies, reverseProxy)
	}

	return p, nil
}

func handleUnhealthyBackendResponse(w http.ResponseWriter, r *http.Request, err error, backendIndex int, rawTarget *backend.Backend) {
	if errors.Is(err, context.Canceled) {
		slog.Info("proxy request canceled",
			"backend_index", backendIndex,
			"backend_url", rawTarget.URL,
			"method", r.Method,
			"path", r.URL.Path,
			"err", err,
		)
		return
	}

	rawTarget.SetHealthy(false)
	slog.Error("backend unavailable",
		"backend_index", backendIndex,
		"backend_url", rawTarget.URL,
		"method", r.Method,
		"path", r.URL.Path,
		"err", err,
	)
	httpx.WriteError(w, http.StatusBadGateway, "backend unavailable", httpx.ErrCodeBadGateway)
}
