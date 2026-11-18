package metrics

// import "encoding/json"

type BenchmarkRequest struct {
	DurationSeconds int    `json:"durationSeconds"`
	Connections     int    `json:"connections"`
	Algorithm       string `json:"algorithm"`
}

// type BenchmarkRunResponse struct {
// 	ID         string                  `json:"id"`
// 	Status     string                  `json:"status"`
// 	Config     BenchmarkRequest        `json:"config"`
// 	Autocannon json.RawMessage         `json:"autocannon"`
// 	Backends   []BackendBenchmarkStats `json:"backends"`
// }
