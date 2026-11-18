package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	benchmarkpkg "github.com/zaiyan-umer/vectra/app/internal/benchmark"
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/metrics"
)

type Handler struct {
	result     *metrics.BenchmarkRuntimeSnapshot
	repository *benchmarkpkg.Repository
}

func NewHandler(result *metrics.BenchmarkRuntimeSnapshot, db *sql.DB) *Handler {
	return &Handler{result: result, repository: benchmarkpkg.NewRepository(db)}
}

func (h *Handler) StartBenchmark(w http.ResponseWriter, r *http.Request) {
	cfg := config.MustLoad()

	request := metrics.BenchmarkRequest{Algorithm: cfg.Algorithm}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid benchmark request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if request.DurationSeconds <= 0 || request.DurationSeconds > 120 {
		http.Error(w, "durationSeconds must be between 1 and 120", http.StatusBadRequest)
		return
	}
	if request.Connections <= 0 || request.Connections > 1000 {
		http.Error(w, "connections must be between 1 and 1000", http.StatusBadRequest)
		return
	}

	targetURL := "http://localhost:" + cfg.Port + "/health"

	benchmarkID, err := h.repository.CreateBenchmark(request, targetURL)
	if err != nil {
		log.Printf("db.query: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(BenchmarkStartResponse{
		ID:     benchmarkID,
		Status: "queued",
	})
	if err != nil {
		http.Error(w, "failed to encode benchmark", http.StatusInternalServerError)
	}
}

func (h *Handler) GetBenchmarks(w http.ResponseWriter, r *http.Request) {
	benchmarks, err := h.repository.GetBenchmarks()
	if err != nil {
		log.Printf("db.query: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(benchmarks); err != nil {
		log.Printf("json.encode: %v", err)
	}
}

func (h *Handler) GetBenchmark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	benchmark, err := h.repository.GetBenchmark(id)
	if err == sql.ErrNoRows {
		http.Error(w, "benchmark not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("db.query benchmark %q: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(benchmark); err != nil {
		log.Printf("json.encode: %v", err)
	}
}

func (h *Handler) DeleteBenchmark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deleted, err := h.repository.DeleteBenchmark(id)
	if err != nil {
		log.Printf("delete: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !deleted {
		http.Error(w, "benchmark not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
