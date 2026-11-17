package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const HealthCheckInterval time.Duration = 10 * time.Second

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

	return Config{
		Port: port,
		DatabaseURL: db_url,
		Algorithm: parsedAlgo,
		TargetURLs: targetURLs,
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