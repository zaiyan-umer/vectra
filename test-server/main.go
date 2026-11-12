package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var serverName string

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

	serverName = os.Getenv("SERVER_NAME")
	if serverName == "" {
		serverName = addr
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("server %s is listening on %s", serverName, httpServer.Addr)

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}

func serverInfo(w http.ResponseWriter, r *http.Request) {
	writeText(w, http.StatusOK, "Test backend")
	logRequest(serverName, r)
}

func delayedResponse(w http.ResponseWriter, r *http.Request) {
	const delay = 2 * time.Second

	log.Printf("delaying response for %s", delay)
	time.Sleep(delay)

	writeText(w, http.StatusOK, "Delayed response")
	logRequest(serverName, r)
}

func backendError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "simulated backend failure on "+serverName, http.StatusInternalServerError)
	logRequest(serverName, r)
	log.Printf("server=%s simulated backend failure", serverName)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeText(w, http.StatusOK, "Server is healthy.")
	logRequest(serverName, r)
}

func writeText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "server=%s %s\n", serverName, message); err != nil {
		log.Printf("error while writing response: %v", err)
	}
}

func logRequest(serverName string, r *http.Request) {
	if r.URL.Path == "/health" {
        return
    }
	log.Printf(
		"REQ server=%s method=%s path=%s remote=%s user_agent=%q",
		serverName,
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		r.UserAgent(),
	)
}
