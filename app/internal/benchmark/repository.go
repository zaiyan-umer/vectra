package benchmark

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

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

func (d *Repository) ClaimNext(ctx context.Context) (*BenchmarkJob, bool, error) {
	row := d.db.QueryRowContext(ctx, `
        WITH next_job AS (
            SELECT id
            FROM benchmarks
            WHERE status = 'queued'
            ORDER BY created_at
            FOR UPDATE SKIP LOCKED
            LIMIT 1
        )
        UPDATE benchmarks b
        SET status = 'running', started_at = NOW()
        FROM next_job n
        WHERE b.id = n.id
        RETURNING
            b.id,
            b.status,
            b.duration_seconds,
            b.connections,
            b.algorithm,
            b.target_url
    `)

	job := &BenchmarkJob{}
	err := row.Scan(
		&job.ID,
		&job.Status,
		&job.DurationSeconds,
		&job.Connections,
		&job.Algorithm,
		&job.TargetURL,
	)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("claim benchmark: %w", err)
	}
	return job, true, nil
}

func (d *Repository) CompleteJob(
	ctx context.Context,
	id string,
	autocannon json.RawMessage,
	routing json.RawMessage,
) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin complete benchmark %q: %w", id, err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		INSERT INTO benchmark_results (
			benchmark_id,
			autocannon_json,
			routing_json
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (benchmark_id) DO UPDATE
		SET autocannon_json = EXCLUDED.autocannon_json,
			routing_json = EXCLUDED.routing_json
	`, id, autocannon, routing)
	if err != nil {
		return fmt.Errorf("store benchmark %q result: %w", id, err)
	}

	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("check benchmark %q result: %w", id, err)
	} else if affected == 0 {
		return fmt.Errorf("benchmark %q result was not stored", id)
	}

	result, err = tx.ExecContext(ctx, `
		UPDATE benchmarks
		SET status = 'completed',
			finished_at = NOW(),
			error = NULL
		WHERE id = $1
		  AND status = 'running'
	`, id)
	if err != nil {
		return fmt.Errorf("complete benchmark %q: %w", id, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check completed benchmark %q: %w", id, err)
	}
	if affected == 0 {
		return fmt.Errorf("benchmark %q is not running", id)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit benchmark %q completion: %w", id, err)
	}

	return nil
}

func (d *Repository) FailJob(ctx context.Context, id string, errorMessage string) error {
	result, err := d.db.ExecContext(ctx, `
        UPDATE benchmarks
        SET status = 'failed',
            finished_at = NOW(),
            error = $2
        WHERE id = $1
          AND status = 'running'
    `, id, errorMessage)
	if err != nil {
		return fmt.Errorf("fail benchmark %q: %w", id, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check failed benchmark %q: %w", id, err)
	}
	if affected == 0 {
		return fmt.Errorf("benchmark %q is not running", id)
	}

	return nil
}
