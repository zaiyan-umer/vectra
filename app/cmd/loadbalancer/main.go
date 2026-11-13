package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/health"
	"github.com/zaiyan-umer/vectra/app/internal/metrics"
	"github.com/zaiyan-umer/vectra/app/internal/selector"
)

func main() {
	targetBackends := backend.Targets(config.TargetURLs)

	result, err := metrics.NewBenchmarkResult(targetBackends)
	if err != nil {
		log.Fatal(err)
	}

	proxy, err := selector.New(config.Default.Algorithm, targetBackends, result)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := health.Run(ctx, targetBackends, config.HealthCheckInterval); err != nil {
			log.Printf("health checker stopped: %v", err)
			cancel()
		}
	}()

	metricsHandler := metrics.NewHandler(result)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /benchmarks/", metricsHandler.StartBenchmark)
	mux.HandleFunc("GET /benchmarks/", metricsHandler.GetBenchmark)
	mux.Handle("/", proxy)

	httpServer := &http.Server{
		Addr:              config.LoadBalancerURL,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server is listening on %s", httpServer.Addr)

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
