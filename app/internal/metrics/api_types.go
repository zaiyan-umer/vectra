package metrics

import "time"

type BenchmarkRequest struct {
	DurationSeconds int    `json:"durationSeconds"`
	Connections     int    `json:"connections"`
	Algorithm       string `json:"algorithm"`
}

type TimeWindow struct {
	Timestamp          time.Time  `json:"timestamp"`
	PoolRPS            float64    `json:"poolRps"`
	PoolErrorRate      float64    `json:"poolErrorRate"`
	PoolAvgLatencyMs   float64    `json:"poolAvgLatencyMs"`
	Backends           []BackendWindow `json:"backends"`
}

type BackendWindow struct {
	Index           int     `json:"index"`
	RPS             float64 `json:"rps"`
	ErrorRate       float64 `json:"errorRate"`
	AvgLatencyMs    float64 `json:"avgLatencyMs"`
}

type PoolMetrics struct {
	RPS           float64 `json:"rps"`
	ErrorRate     float64 `json:"errorRate"`
	AvgLatencyMs  float64 `json:"avgLatencyMs"`
}

type BackendMetrics struct {
	Index            int     `json:"index"`
	URL              string  `json:"url"`
	Healthy          bool    `json:"healthy"`
	ActiveConnections int64  `json:"activeConnections"`
	RPS              float64 `json:"rps"`
	ErrorRate        float64 `json:"errorRate"`
	AvgLatencyMs     float64 `json:"avgLatencyMs"`
	TrafficShare     float64 `json:"trafficShare"`
	SelectedRequests  uint64  `json:"selectedRequests"`
	CompletedRequests uint64  `json:"completedRequests"`
	Errors           uint64  `json:"errors"`
}