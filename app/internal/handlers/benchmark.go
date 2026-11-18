package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/config"
	"github.com/zaiyan-umer/vectra/app/internal/metrics"
)

type Handler struct {
	result *metrics.BenchmarkRuntimeSnapshot
	db *sql.DB
}

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

func NewHandler(result *metrics.BenchmarkRuntimeSnapshot, db *sql.DB) *Handler {
	return &Handler{result: result, db: db}
}

func (h *Handler) StartBenchmark(w http.ResponseWriter, r *http.Request) {
	cfg := config.MustLoad()

	request := metrics.BenchmarkRequest{Algorithm: cfg.Algorithm}
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

	targetURL := "http://" + cfg.Port + "/health"
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

	response := metrics.BenchmarkRunResponse{
		ID:         benchmarkID,
		Status:     "completed",
		Config:     request,
		Autocannon: autocannonResult,
		Backends:   report,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode benchmark response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetBenchmarks(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(
		`SELECT id, status, duration_seconds, connections, algorithm,
		target_url, started_at, finished_at, error, created_at
		FROM benchmarks`,
	)
	if err != nil {
		log.Printf("db.query: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	benchmarks := []BenchmarkRecord{}

		for rows.Next() {
			var b BenchmarkRecord
			var startedAt sql.NullTime
			var finishedAt sql.NullTime
			var benchmarkError sql.NullString

			err := rows.Scan(
				&b.ID,
				&b.Status,
				&b.DurationSeconds,
				&b.Connections,
				&b.Algorithm,
				&b.TargetURL,
				&startedAt,
				&finishedAt,
				&benchmarkError,
				&b.CreatedAt,
			)
			if err != nil {
				log.Printf("rows.scan: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			if startedAt.Valid {
				b.StartedAt = &startedAt.Time
			}
			if finishedAt.Valid {
				b.FinishedAt = &finishedAt.Time
			}
			if benchmarkError.Valid {
				b.Error = &benchmarkError.String
			}

			benchmarks = append(benchmarks, b)
		}

		if err := rows.Err(); err != nil {
			log.Printf("rows.next: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(benchmarks); err != nil {
			log.Printf("json.encode: %v", err)
		}
}

func (h *Handler) DeleteBenchmark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	res, err := h.db.Exec(`DELETE FROM benchmarks WHERE id = $1`, id)
	if err != nil {
		log.Printf("delete: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("delete benchmark %q: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "benchmark not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetBenchmark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	row := h.db.QueryRow(`
		SELECT
			b.id,
			b.status,
			b.duration_seconds,
			b.connections,
			b.algorithm,
			b.target_url,
			b.started_at,
			b.finished_at,
			b.error,
			b.created_at,
			r.autocannon_json,
			r.routing_json
		FROM benchmarks b
		LEFT JOIN benchmark_results r ON r.benchmark_id = b.id
		WHERE b.id = $1
	`, id)

	var benchmark BenchmarkDetailRecord
	var startedAt sql.NullTime
	var finishedAt sql.NullTime
	var benchmarkError sql.NullString
	var autocannonJSON []byte
	var routingJSON []byte

	err := row.Scan(
		&benchmark.ID,
		&benchmark.Status,
		&benchmark.DurationSeconds,
		&benchmark.Connections,
		&benchmark.Algorithm,
		&benchmark.TargetURL,
		&startedAt,
		&finishedAt,
		&benchmarkError,
		&benchmark.CreatedAt,
		&autocannonJSON,
		&routingJSON,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "benchmark not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("db.query benchmark %q: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if startedAt.Valid {
		benchmark.StartedAt = &startedAt.Time
	}
	if finishedAt.Valid {
		benchmark.FinishedAt = &finishedAt.Time
	}
	if benchmarkError.Valid {
		benchmark.Error = &benchmarkError.String
	}
	if len(autocannonJSON) > 0 && len(routingJSON) > 0 {
		benchmark.Result = &BenchmarkResultRecord{
			Autocannon: json.RawMessage(autocannonJSON),
			Routing:    json.RawMessage(routingJSON),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(benchmark); err != nil {
		log.Printf("json.encode: %v", err)
	}
}

func newBenchmarkID() (string, error) {
	var bytes [16]byte

	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("generate benchmark ID: %w", err)
	}

	return "benchmark-" + hex.EncodeToString(bytes[:]), nil
}
