package config

import "time"

// the target backend server urls
var TargetURLs = []string{
	"http://localhost:8081",
	"http://localhost:8082",
	"http://localhost:8083",
}

// default algorithm is round robin
var Algorithm = AlgorithmLatencyAware

var Default = Config{
	Algorithm:  Algorithm,
	TargetURLs: TargetURLs,
}

const HealthCheckInterval time.Duration = 10 * time.Second
