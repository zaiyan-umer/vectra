package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/selector"
)

func main() {
	proxy, err := selector.New(config.Default.Algorithm, config.Default.TargetURLs)
	if err != nil {
		log.Fatal(err)
	}

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
