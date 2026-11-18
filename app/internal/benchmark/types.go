package benchmark

import (
	"encoding/json"
	"time"
)

type BenchmarkRecord struct {
	ID              string     `json:"id"`
	Status          string     `json:"status"`
	DurationSeconds int        `json:"durationSeconds"`
	Connections     int        `json:"connections"`
	Algorithm       string     `json:"algorithm"`
	TargetURL       string     `json:"targetUrl"`
	StartedAt       *time.Time `json:"startedAt"`
	FinishedAt      *time.Time `json:"finishedAt"`
	Error           *string    `json:"error"`
	CreatedAt       time.Time  `json:"createdAt"`
}

type BenchmarkResultRecord struct {
	Autocannon json.RawMessage `json:"autocannon"`
	Routing    json.RawMessage `json:"routing"`
}

type BenchmarkDetailRecord struct {
	BenchmarkRecord
	Result *BenchmarkResultRecord `json:"result,omitempty"`
}
