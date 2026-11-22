package handlers

import (
	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/health"
	"github.com/zaiyan-umer/vectra/app/internal/metrics"
)

type BenchmarkStartResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
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

type ConfigResponse struct {
	Algorithm                     string                `json:"algorithm"`
	Targets                       []string              `json:"targets"`
	HealthCheckIntervalSeconds    int                   `json:"healthCheckIntervalSeconds"`
	MetricsCaptureIntervalSeconds int                   `json:"metricsCaptureIntervalSeconds"`
	MetricsMaxWindows             int                   `json:"metricsMaxWindows"`
	AdaptiveConfig                config.AdaptiveConfig `json:"adaptiveConfig"`
	AvailableAlgorithms           []string              `json:"availableAlgorithms"`
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

type HealthResponse struct {
	HealthyCount int                          `json:"healthyCount"`
	TotalCount   int                          `json:"totalCount"`
	Backends     []health.BackendHealthStatus `json:"backends"`
}