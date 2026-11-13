package metrics

import (
	"context"
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/config"
)

type BenchmarkRequest struct {
	DurationSeconds int `json:"durationSeconds"`
	Connections     int `json:"connections"`
}

type Handler struct {
	result *BenchmarkResult
}

func NewHandler(result *BenchmarkResult) *Handler {
	return &Handler{result: result}
}

func (h *Handler) StartBenchmark(w http.ResponseWriter, r *http.Request) {
	var request BenchmarkRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid benchmark request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if request.DurationSeconds <= 0 || request.DurationSeconds > 120 {
		http.Error(w, "durationSeconds must be between 1 and 120", http.StatusBadRequest)
		return
	}
	if request.Connections <= 0 || request.Connections > 1000 {
		http.Error(w, "connections must be between 1 and 1000", http.StatusBadRequest)
		return
	}

	before := h.result.CaptureSnapshot()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	targetURL := "http://" + config.LoadBalancerURL + "/health"
	cmd := exec.CommandContext(
		ctx,
		"npx",
		"--yes",
		"autocannon",
		"--json",
		"-c", strconv.Itoa(request.Connections),
		"-d", strconv.Itoa(request.DurationSeconds),
		targetURL,
	)

	output, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	after := h.result.CaptureSnapshot()
	report := after.Difference(before)

	var autocannonResult json.RawMessage
	if err := json.Unmarshal(output, &autocannonResult); err != nil {
		http.Error(w, "invalid autocannon output", http.StatusInternalServerError)
		return
	}

	response := struct {
		Autocannon json.RawMessage `json:"autocannon"`
		Backends   any             `json:"backends"`
	}{
		Autocannon: autocannonResult,
		Backends:   report,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode benchmark response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetBenchmark(w http.ResponseWriter, r *http.Request) {

}
