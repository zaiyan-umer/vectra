package benchmark

import (
	"database/sql"
	"encoding/json"

	"github.com/zaiyan-umer/vectra/app/internal/metrics"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (d *Repository) CreateBenchmark(request metrics.BenchmarkRequest, targetURL string) (string, error) {
	var benchmarkID string
	err := d.db.QueryRow(`
		INSERT INTO benchmarks (
			status, duration_seconds, connections, algorithm, target_url
		)
		VALUES ('queued', $1, $2, $3, $4)
		RETURNING id
	`,
		request.DurationSeconds,
		request.Connections,
		request.Algorithm,
		targetURL,
	).Scan(&benchmarkID)

	return benchmarkID, err
}

func (d *Repository) GetBenchmarks() ([]BenchmarkRecord, error) {
	rows, err := d.db.Query(`
		SELECT id, status, duration_seconds, connections, algorithm,
			target_url, started_at, finished_at, error, created_at
		FROM benchmarks
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	benchmarks := []BenchmarkRecord{}

	for rows.Next() {
		var benchmark BenchmarkRecord
		var startedAt sql.NullTime
		var finishedAt sql.NullTime
		var benchmarkError sql.NullString

		err := rows.Scan(
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
		)
		if err != nil {
			return nil, err
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

		benchmarks = append(benchmarks, benchmark)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return benchmarks, nil
}

func (d *Repository) GetBenchmark(id string) (*BenchmarkDetailRecord, error) {
	row := d.db.QueryRow(`
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
	if err != nil {
		return nil, err
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

	return &benchmark, nil
}

func (d *Repository) DeleteBenchmark(id string) (bool, error) {
	result, err := d.db.Exec(`DELETE FROM benchmarks WHERE id = $1`, id)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}
