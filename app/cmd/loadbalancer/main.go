package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/health"
	"github.com/zaiyan-umer/vectra/app/internal/selector"
)

func main() {
	targetBackends := backend.Targets(config.TargetURLs)

	proxy, err := selector.New(config.Default.Algorithm, targetBackends)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go health.Run(ctx, targetBackends, config.HealthCheckInterval)

	mux := http.NewServeMux()
	mux.Handle("/", proxy)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("server is listening on %s", httpServer.Addr)

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
