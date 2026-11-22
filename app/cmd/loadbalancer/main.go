package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
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

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	targetBackends := backend.Targets(cfg.TargetURLs)

	result, err := metrics.NewBenchmarkRuntimeSnapshot(targetBackends, cfg.MetricsMaxWindows)
	if err != nil {
		slog.Error("failed to initialize benchmark runtime metrics", "err", err)
		os.Exit(1)
	}

	proxy, err := selector.New(cfg.Algorithm, targetBackends, result)
	if err != nil {
		slog.Error("failed to initialize load-balancer selector", "algorithm", cfg.Algorithm, "err", err)
		os.Exit(1)
	}

	dbConnection, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := health.Run(ctx, targetBackends, cfg.HealthCheckInterval); err != nil {
			slog.Error("health checker stopped", "interval", cfg.HealthCheckInterval, "err", err)
			cancel()
		}
	}()

	benchmarkRepository := benchmark.NewRepository(dbConnection)
	benchmarkWorker := benchmark.NewWorker(benchmarkRepository, result, cfg.BenchmarkWorkerPollInterval)
	go benchmarkWorker.Run(ctx)

	metricsHandler := handlers.NewHandler(result, dbConnection, targetBackends)

	result.StartCapture(cfg.MetricsCaptureInterval)

	go func() {
		ticker := time.NewTicker(cfg.HealthCheckInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				healthStatus := health.GetPoolHealth(targetBackends)
				metricsHandler.UpdateHealthStatus(healthStatus)
			case <-ctx.Done():
				return
			}
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /benchmarks/", metricsHandler.StartBenchmark)
	mux.HandleFunc("GET /benchmarks/{id}", metricsHandler.GetBenchmark)
	mux.HandleFunc("GET /benchmarks/", metricsHandler.GetBenchmarks)
	mux.HandleFunc("DELETE /benchmarks/{id}", metricsHandler.DeleteBenchmark)
	mux.HandleFunc("GET /metrics", metricsHandler.GetMetrics)
	mux.HandleFunc("GET /config", metricsHandler.GetConfig)
	mux.HandleFunc("GET /backends", metricsHandler.GetBackends)
	mux.HandleFunc("GET /backends/health", metricsHandler.GetBackendsHealth)
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

	slog.Info("server is listening", "address", httpServer.Addr)

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server failed", "address", httpServer.Addr, "err", err)
		os.Exit(1)
	}
}
