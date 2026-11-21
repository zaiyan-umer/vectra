package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var AdaptiveConfiguration = AdaptiveConfig{
    MaxExpectedConnections: 100,
    LatencyTarget:          500 * time.Millisecond,
    ActiveWeight:           0.5,
    LatencyWeight:          0.5,
	Tolerance: 				0.00001,
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required")
	}

	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		panic("DATABASE_URL is required")
	}

	algorithm := os.Getenv("ALGORITHM")
	if algorithm == "" {
		panic("ALGORITHM is required")
	}

	parsedAlgo, err := ParseAlgorithm(algorithm)
	if err != nil {
		panic(err)
	}

	rawURLs := os.Getenv("BACKEND_URLS")
	if rawURLs == "" {
		panic("BACKEND_URLS is required")
	}

	targetURLs := makeTargetUrls(rawURLs)

	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	var corsAllowedOrigins []string
	if corsOrigins != "" {
		for _, o := range strings.Split(corsOrigins, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				corsAllowedOrigins = append(corsAllowedOrigins, trimmed)
			}
		}
	} else {
		corsAllowedOrigins = []string{"http://localhost:5173"}
	}

	healthCheckInterval := 10 * time.Second
	if v := os.Getenv("HEALTH_CHECK_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			healthCheckInterval = d
		}
	}

	metricsCaptureInterval := 2 * time.Second
	if v := os.Getenv("METRICS_CAPTURE_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			metricsCaptureInterval = d
		}
	}

	metricsMaxWindows := 180
	if v := os.Getenv("METRICS_MAX_WINDOWS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 10 {
			metricsMaxWindows = n
		}
	}

	benchmarkWorkerPollInterval := 2 * time.Second
	if v := os.Getenv("BENCHMARK_WORKER_POLL_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d >= 500*time.Millisecond && d <= 30*time.Second {
			benchmarkWorkerPollInterval = d
		}
	}

	return Config{
		Port:                       port,
		DatabaseURL:                db_url,
		Algorithm:                  parsedAlgo,
		TargetURLs:                 targetURLs,
		CORSAllowedOrigins:         corsAllowedOrigins,
		HealthCheckInterval:        healthCheckInterval,
		MetricsCaptureInterval:     metricsCaptureInterval,
		MetricsMaxWindows:          metricsMaxWindows,
		BenchmarkWorkerPollInterval: benchmarkWorkerPollInterval,
	}
}

func makeTargetUrls(rawURLs string) []string {
	targetURLs := make([]string, 0)
	for _, rawURL := range strings.Split(rawURLs, ",") {
		url := strings.TrimSpace(rawURL)
		if url != "" {
			targetURLs = append(targetURLs, url)
		}
	}

	return targetURLs
}