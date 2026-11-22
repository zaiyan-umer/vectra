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
		http.Error(w, "invalid benchmark request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if request.DurationSeconds <= 0 || request.DurationSeconds > 120 {
		slog.Error("invalid benchmark request", "req", request, "reason", "durationSeconds must be between 1 and 120")
		http.Error(w, "durationSeconds must be between 1 and 120", http.StatusBadRequest)
		return
	}
	if request.Connections <= 0 || request.Connections > 1000 {
		slog.Error("invalid benchmark request", "req", request, "reason", "connections must be between 1 and 1000")
		http.Error(w, "connections must be between 1 and 1000", http.StatusBadRequest)
		return
	}

	targetURL := "http://localhost:" + cfg.Port + "/health"

	benchmarkID, err := h.repository.CreateBenchmark(request, targetURL)
	if err != nil {
		slog.Error("failed to create benchmark", "req", request, "err", err)
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
		slog.Error("failed to encode benchmark start response", "req", request, "benchmark_id", benchmarkID, "err", err)
		http.Error(w, "failed to encode benchmark", http.StatusInternalServerError)
	}
}

func (h *Handler) GetBenchmarks(w http.ResponseWriter, r *http.Request) {
	benchmarks, err := h.repository.GetBenchmarks()
	if err != nil {
		slog.Error("error while searching benchmarks", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(benchmarks); err != nil {
		slog.Error("failed to encode benchmarks response", "err", err)
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
		slog.Error("failed to get benchmark", "benchmark_id", id, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(benchmark); err != nil {
		slog.Error("failed to encode benchmark response", "benchmark_id", id, "err", err)
	}
}

func (h *Handler) DeleteBenchmark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deleted, err := h.repository.DeleteBenchmark(id)
	if err != nil {
		slog.Error("failed to delete benchmark", "benchmark_id", id, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !deleted {
		http.Error(w, "benchmark not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type MetricsResponse struct {
	Pool      metrics.PoolMetrics      `json:"pool"`
	Backends  []metrics.BackendMetrics `json:"backends"`
	Sparkline SparklineData            `json:"sparkline"`
}

type SparklineData struct {
	Timestamps       []string                    `json:"timestamps"`
	PoolRPS          []float64                   `json:"poolRps"`
	PoolErrorRate    []float64                   `json:"poolErrorRate"`
	PoolAvgLatencyMs []float64                   `json:"poolAvgLatencyMs"`
	Backends         map[string]BackendSparkline `json:"backends"`
}

type BackendSparkline struct {
	RPS          []float64 `json:"rps"`
	AvgLatencyMs []float64 `json:"avgLatencyMs"`
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(MetricsResponse{
		Pool:      pool,
		Backends:  backends,
		Sparkline: sparkline,
	}); err != nil {
		slog.Error("failed to encode metrics response", "err", err)
	}
}

type ConfigResponse struct {
	Algorithm                     string                `json:"algorithm"`
	Targets                       []string              `json:"targets"`
	HealthCheckIntervalSeconds    int                   `json:"healthCheckIntervalSeconds"`
	MetricsCaptureIntervalSeconds int                   `json:"metricsCaptureIntervalSeconds"`
	MetricsMaxWindows             int                   `json:"metricsMaxWindows"`
	AdaptiveConfig                config.AdaptiveConfig `json:"adaptiveConfig"`
	AvailableAlgorithms           []string              `json:"availableAlgorithms"`
}

func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := config.MustLoad()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ConfigResponse{
		Algorithm:                     cfg.Algorithm,
		Targets:                       cfg.TargetURLs,
		HealthCheckIntervalSeconds:    int(cfg.HealthCheckInterval.Seconds()),
		MetricsCaptureIntervalSeconds: int(cfg.MetricsCaptureInterval.Seconds()),
		MetricsMaxWindows:             cfg.MetricsMaxWindows,
		AdaptiveConfig:                config.AdaptiveConfiguration,
		AvailableAlgorithms:           config.AvailableAlgorithms,
	}); err != nil {
		slog.Error("failed to encode config response", "err", err)
	}
}

type BackendsResponse struct {
	Backends []BackendDetail `json:"backends"`
}

type BackendDetail struct {
	Index             int     `json:"index"`
	URL               string  `json:"url"`
	Healthy           bool    `json:"healthy"`
	ActiveConnections int64   `json:"activeConnections"`
	RPS               float64 `json:"rps"`
	ErrorRate         float64 `json:"errorRate"`
	AvgLatencyMs      float64 `json:"avgLatencyMs"`
	TrafficShare      float64 `json:"trafficShare"`
	SelectedRequests  uint64  `json:"selectedRequests"`
	CompletedRequests uint64  `json:"completedRequests"`
	Errors            uint64  `json:"errors"`
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(BackendsResponse{Backends: details}); err != nil {
		slog.Error("failed to encode backends response", "err", err)
	}
}

type HealthResponse struct {
	HealthyCount int                          `json:"healthyCount"`
	TotalCount   int                          `json:"totalCount"`
	Backends     []health.BackendHealthStatus `json:"backends"`
}

func (h *Handler) GetBackendsHealth(w http.ResponseWriter, r *http.Request) {
	healthyCount := 0
	for _, b := range h.health {
		if b.Healthy {
			healthyCount++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(HealthResponse{
		HealthyCount: healthyCount,
		TotalCount:   len(h.health),
		Backends:     h.health,
	}); err != nil {
		slog.Error("failed to encode health response", "err", err)
	}
}
