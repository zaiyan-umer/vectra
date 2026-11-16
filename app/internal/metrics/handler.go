package metrics

import (
	"context"
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"time"
	"crypto/rand"
    "encoding/hex"
    "fmt"

	"github.com/zaiyan-umer/vectra/app/internal/config"
)

type BenchmarkConfig struct {
	DurationSeconds int `json:"durationSeconds"`
	Connections     int `json:"connections"`
	Algorithm 		string `json:"algorithm"`
}

type BenchmarkResponse struct {
    ID      string          `json:"id"`
    Status  string          `json:"status"`
    Config  BenchmarkConfig `json:"config"`
    Autocanon json.RawMessage `json:"autocanon"`
    Backends []BackendStats `json:"backends"`
}

type Handler struct {
	result *BenchmarkResult
}

func NewHandler(result *BenchmarkResult) *Handler {
	return &Handler{result: result}
}

func (h *Handler) StartBenchmark(w http.ResponseWriter, r *http.Request) {
	request := BenchmarkConfig{Algorithm: config.Algorithm}
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

	benchmarkID, err := newBenchmarkID()
	if err != nil {
		http.Error(w, "failed to create benchmark ID", http.StatusInternalServerError)
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

	response := BenchmarkResponse{
		ID: benchmarkID,
		Status: "completed",
		Config: request,
		Autocanon: autocannonResult,
		Backends:   report,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode benchmark response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetBenchmark(w http.ResponseWriter, r *http.Request) {

}

func newBenchmarkID() (string, error) {
    var bytes [16]byte

    if _, err := rand.Read(bytes[:]); err != nil {
        return "", fmt.Errorf("generate benchmark ID: %w", err)
    }

    return "benchmark-" + hex.EncodeToString(bytes[:]), nil
}