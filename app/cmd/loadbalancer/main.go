package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/db"
	"github.com/zaiyan-umer/vectra/app/internal/handlers"
	"github.com/zaiyan-umer/vectra/app/internal/health"
	"github.com/zaiyan-umer/vectra/app/internal/metrics"
	"github.com/zaiyan-umer/vectra/app/internal/selector"
)

func main() {
	cfg := config.MustLoad()

	targetBackends := backend.Targets(cfg.TargetURLs)

	result, err := metrics.NewBenchmarkRuntimeSnapshot(targetBackends)
	if err != nil {
		log.Fatalf("main.metrics.NewBenchmarkRuntimeSnapshot: %v", err)
	}

	proxy, err := selector.New(cfg.Algorithm, targetBackends, result)
	if err != nil {
		log.Fatalf("main.selector.New: %v", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatalf("main.os.Getenv: DATABASE_URL is required.")
	}

	dbConnection, err := db.Connect(dbUrl)
	if err != nil {
		log.Fatalf("main.db.Connect: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := health.Run(ctx, targetBackends, config.HealthCheckInterval); err != nil {
			log.Printf("health checker stopped: %v", err)
			cancel()
		}
	}()

	metricsHandler := handlers.NewHandler(result)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /benchmarks/", metricsHandler.StartBenchmark)
	mux.HandleFunc("GET /benchmarks/", metricsHandler.GetBenchmark(dbConnection))
	mux.Handle("/", proxy)

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
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
