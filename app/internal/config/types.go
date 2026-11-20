package config

import "time"

type Config struct {
    Port string
    DatabaseURL string
	Algorithm  string
	TargetURLs []string
	CORSAllowedOrigins []string
}

type AdaptiveConfig struct {
    MaxExpectedConnections int
    LatencyTarget          time.Duration
    ActiveWeight           float64
    LatencyWeight          float64
    Tolerance              float64
}

const AlgorithmRoundRobin = "round_robin"
const AlgorithmHealthyRoundRobin = "healthy_round_robin"
const AlgorithmLeastConnections = "least_connections"
const AlgorithmLatencyAware = "latency_aware"
const AdaptiveSelection = "adaptive_selection"