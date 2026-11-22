package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/zaiyan-umer/vectra/app/internal/backend"
	benchmarkpkg "github.com/zaiyan-umer/vectra/app/internal/benchmark"
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/health"
	"github.com/zaiyan-umer/vectra/app/internal/httpx"
	"github.com/zaiyan-umer/vectra/app/internal/metrics"
)

type Handler struct {
	result     *metrics.BenchmarkRuntimeSnapshot
	repository *benchmarkpkg.Repository
	backends   []*backend.Backend
	health     []health.BackendHealthStatus
}

func NewHandler(result *metrics.BenchmarkRuntimeSnapshot, db *sql.DB, backends []*backend.Backend) *Handler {
	return &Handler{
		result:     result,
		repository: benchmarkpkg.NewRepository(db),
		backends:   backends,
	}
}

func (h *Handler) UpdateHealthStatus(status []health.BackendHealthStatus) {
	h.health = status
}

func (h *Handler) StartBenchmark(w http.ResponseWriter, r *http.Request) {
	cfg := config.MustLoad()

	request := metrics.BenchmarkRequest{Algorithm: cfg.Algorithm}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		slog.Error("invalid benchmark request", "req", request, "err", err)
		httpx.WriteError(w, http.StatusBadRequest, "invalid benchmark request: "+err.Error(), httpx.ErrCodeInvalidJSON)
		return
	}
	if request.DurationSeconds <= 0 || request.DurationSeconds > 120 {
		slog.Error("invalid benchmark request", "req", request, "reason", "durationSeconds must be between 1 and 120")
		httpx.WriteError(w, http.StatusBadRequest, "durationSeconds must be between 1 and 120", httpx.ErrCodeValidationFailed)
		return
	}
	if request.Connections <= 0 || request.Connections > 1000 {
		slog.Error("invalid benchmark request", "req", request, "reason", "connections must be between 1 and 1000")
		httpx.WriteError(w, http.StatusBadRequest, "connections must be between 1 and 1000", httpx.ErrCodeValidationFailed)
		return
	}

	targetURL := "http://localhost:" + cfg.Port + "/health"

	benchmarkID, err := h.repository.CreateBenchmark(request, targetURL)
	if err != nil {
		slog.Error("failed to create benchmark", "req", request, "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error", httpx.ErrCodeInternalError)
		return
	}

	response := BenchmarkStartResponse{
		ID:     benchmarkID,
		Status: "queued",
	}

	err = httpx.WriteJson(w, http.StatusAccepted, response)
	if err != nil {
		slog.Error("failed to encode benchmark start response", "benchmark_id", benchmarkID, "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to encode benchmark", httpx.ErrCodeInternalError)
	}
}

func (h *Handler) GetBenchmarks(w http.ResponseWriter, r *http.Request) {
	benchmarks, err := h.repository.GetBenchmarks()
	if err != nil {
		slog.Error("error while searching benchmarks", "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error", httpx.ErrCodeInternalError)
		return
	}

	err = httpx.WriteJson(w, http.StatusOK, benchmarks)
	if err != nil {
		slog.Error("failed to encode benchmarks response", "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to encode benchmarks", httpx.ErrCodeInternalError)
	}
}

func (h *Handler) GetBenchmark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	benchmark, err := h.repository.GetBenchmark(id)
	if err == sql.ErrNoRows {
		httpx.WriteError(w, http.StatusNotFound, "benchmark not found", httpx.ErrCodeNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to get benchmark", "benchmark_id", id, "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error", httpx.ErrCodeInternalError)
		return
	}

	err = httpx.WriteJson(w, http.StatusOK, benchmark)
	if err != nil {
		slog.Error("failed to encode benchmark response", "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to encode benchmark", httpx.ErrCodeInternalError)
	}
}

func (h *Handler) DeleteBenchmark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deleted, err := h.repository.DeleteBenchmark(id)
	if err != nil {
		slog.Error("failed to delete benchmark", "benchmark_id", id, "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error", httpx.ErrCodeInternalError)
		return
	}

	if !deleted {
		httpx.WriteError(w, http.StatusNotFound, "benchmark not found", httpx.ErrCodeNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	pool := h.result.GetPoolMetrics()
	backends := h.result.GetBackendsMetrics(h.backends)
	series := h.result.GetTimeSeries()

	sparkline := SparklineData{
		Timestamps:       make([]string, 0, len(series)),
		PoolRPS:          make([]float64, 0, len(series)),
		PoolErrorRate:    make([]float64, 0, len(series)),
		PoolAvgLatencyMs: make([]float64, 0, len(series)),
		Backends:         make(map[string]BackendSparkline),
	}

	for _, w := range series {
		sparkline.Timestamps = append(sparkline.Timestamps, w.Timestamp.Format("15:04:05"))
		sparkline.PoolRPS = append(sparkline.PoolRPS, w.PoolRPS)
		sparkline.PoolErrorRate = append(sparkline.PoolErrorRate, w.PoolErrorRate)
		sparkline.PoolAvgLatencyMs = append(sparkline.PoolAvgLatencyMs, w.PoolAvgLatencyMs)

		for _, bw := range w.Backends {
			key := string(rune('0' + bw.Index))
			if _, ok := sparkline.Backends[key]; !ok {
				sparkline.Backends[key] = BackendSparkline{
					RPS:          make([]float64, 0, len(series)),
					AvgLatencyMs: make([]float64, 0, len(series)),
				}
			}
			s := sparkline.Backends[key]
			s.RPS = append(s.RPS, bw.RPS)
			s.AvgLatencyMs = append(s.AvgLatencyMs, bw.AvgLatencyMs)
			sparkline.Backends[key] = s
		}
	}

	metrics := MetricsResponse{
		Pool:      pool,
		Backends:  backends,
		Sparkline: sparkline,
	}

	err := httpx.WriteJson(w, http.StatusOK, metrics)
	if err != nil {
		slog.Error("failed to encode metrics response", "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to encode metrics response", httpx.ErrCodeInternalError)
	}
}

func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := config.MustLoad()

	response := ConfigResponse{
		Algorithm:                     cfg.Algorithm,
		Targets:                       cfg.TargetURLs,
		HealthCheckIntervalSeconds:    int(cfg.HealthCheckInterval.Seconds()),
		MetricsCaptureIntervalSeconds: int(cfg.MetricsCaptureInterval.Seconds()),
		MetricsMaxWindows:             cfg.MetricsMaxWindows,
		AdaptiveConfig:                config.AdaptiveConfiguration,
		AvailableAlgorithms:           config.AvailableAlgorithms,
	}

	if err := httpx.WriteJson(w, http.StatusOK, response); err != nil {
		slog.Error("failed to encode config response", "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to encode config response", httpx.ErrCodeInternalError)
	}
}

func (h *Handler) GetBackends(w http.ResponseWriter, r *http.Request) {
	backendMetrics := h.result.GetBackendsMetrics(h.backends)

	details := make([]BackendDetail, len(backendMetrics))
	for i, bm := range backendMetrics {
		details[i] = BackendDetail{
			Index:             bm.Index,
			URL:               bm.URL,
			Healthy:           bm.Healthy,
			ActiveConnections: bm.ActiveConnections,
			RPS:               bm.RPS,
			ErrorRate:         bm.ErrorRate,
			AvgLatencyMs:      bm.AvgLatencyMs,
			TrafficShare:      bm.TrafficShare,
			SelectedRequests:  bm.SelectedRequests,
			CompletedRequests: bm.CompletedRequests,
			Errors:            bm.Errors,
		}
	}

	response := BackendsResponse{Backends: details}
	if err := httpx.WriteJson(w, http.StatusOK, response); err != nil {
		slog.Error("failed to encode backends response", "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to encode backends response", httpx.ErrCodeInternalError)
	}
}

func (h *Handler) GetBackendsHealth(w http.ResponseWriter, r *http.Request) {
	healthyCount := 0
	for _, b := range h.health {
		if b.Healthy {
			healthyCount++
		}
	}

	response := HealthResponse{
		HealthyCount: healthyCount,
		TotalCount:   len(h.health),
		Backends:     h.health,
	}

	if err := httpx.WriteJson(w, http.StatusOK, response); err != nil {
		slog.Error("failed to encode health response", "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to encode health response", httpx.ErrCodeInternalError)
	}
}
