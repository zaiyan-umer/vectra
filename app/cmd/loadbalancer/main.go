package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/benchmark"
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

	dbConnection, err := db.Connect(cfg.DatabaseURL)
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

	benchmarkRepository := benchmark.NewRepository(dbConnection)
	benchmarkWorker := benchmark.NewWorker(benchmarkRepository, result)
	go benchmarkWorker.Run(ctx)

	metricsHandler := handlers.NewHandler(result, dbConnection)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /benchmarks/", metricsHandler.StartBenchmark)
	mux.HandleFunc("GET /benchmarks/{id}", metricsHandler.GetBenchmark)
	mux.HandleFunc("GET /benchmarks/", metricsHandler.GetBenchmarks)
	mux.HandleFunc("DELETE /benchmarks/{id}", metricsHandler.DeleteBenchmark)
	mux.Handle("/", proxy)

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
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
