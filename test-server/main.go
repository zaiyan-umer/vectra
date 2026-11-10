package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", serverInfo)
	mux.HandleFunc("GET /health", healthCheck)
	mux.HandleFunc("GET /delay", delayedResponse)
	mux.HandleFunc("GET /error", backendError)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8081"
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("server is listening on %s", httpServer.Addr)

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}

func serverInfo(w http.ResponseWriter, r *http.Request) {
	writeText(w, http.StatusOK, "Test backend")
	logRequest(r)
}

func delayedResponse(w http.ResponseWriter, r *http.Request) {
	const delay = 2 * time.Second

	log.Printf("delaying response for %s", delay)
	time.Sleep(delay)

	writeText(w, http.StatusOK, "Delayed response")
	logRequest(r)
}

func backendError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "simulated backend failure", http.StatusInternalServerError)
	logRequest(r)
	log.Printf("simulated backend failure")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeText(w, http.StatusOK, "Server is healthy.")
	logRequest(r)
}

func writeText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	if _, err := fmt.Fprintln(w, message); err != nil {
		log.Printf("error while writing response: %v", err)
	}
}

func logRequest(r *http.Request) {
	log.Printf(
		"REQ method=%s path=%s remote=%s user_agent=%q",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		r.UserAgent(),
	)
}
