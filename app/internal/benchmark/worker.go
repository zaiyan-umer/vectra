package benchmark

import (
	"context"
	"encoding/json"
	"log/slog"
	"os/exec"
	"strconv"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/metrics"
)

type Worker struct {
	repository                      *Repository
	runtimeSnapshot                 *metrics.BenchmarkRuntimeSnapshot
	benchmarkWorkerPollInterval     time.Duration
}

func NewWorker(repo *Repository, sh *metrics.BenchmarkRuntimeSnapshot, pollInterval time.Duration) *Worker {
	if pollInterval <= 0 {
		pollInterval = 2 * time.Second
	}
	return &Worker{
		repository:                     repo,
		runtimeSnapshot:                sh,
		benchmarkWorkerPollInterval:    pollInterval,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.benchmarkWorkerPollInterval)
	defer ticker.Stop()

	for {
		w.processNext(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *Worker) processNext(ctx context.Context) {
	job, found, err := w.repository.ClaimNext(ctx)
	if !found {
		return
	}
	if err != nil {
		slog.Error("worker could not claim the job", "job_id", job.ID, "err", err)
		return
	}

	w.execute(ctx, job)
}

func (w *Worker) execute(ctx context.Context, job *BenchmarkJob) {
	before := w.runtimeSnapshot.CaptureSnapshot()

	timeout := time.Duration(job.DurationSeconds+30) * time.Second
	bechmarkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(
		bechmarkCtx,
		"npx",
		"--yes",
		"autocannon",
		"--json",
		"-c", strconv.Itoa(job.Connections),
		"-d", strconv.Itoa(job.DurationSeconds),
		job.TargetURL,
	)

	output, err := cmd.Output()
	if err != nil {
		slog.Error("error on cmd.output", "job_id", job.ID, "err", err)
		if err := w.repository.FailJob(ctx, job.ID, err.Error()); err != nil {
			slog.Error("benchmark failed", "job_id", job.ID, "err", err)
		}
		return
	}

	after := w.runtimeSnapshot.CaptureSnapshot()
	report := after.Difference(before)

	var autocannonResult json.RawMessage
	if err := json.Unmarshal(output, &autocannonResult); err != nil {
		slog.Error("error on json.unmarshal", "job_id", job.ID, "err", err)
		if err := w.repository.FailJob(ctx, job.ID, err.Error()); err != nil {
			slog.Error("benchmark failed", "job_id", job.ID, "err", err)
		}
		return
	}

	routingJSON, err := json.Marshal(report)
	if err != nil {
		slog.Error("error on json.marshal routing", "job_id", job.ID, "err", err)
		if err := w.repository.FailJob(ctx, job.ID, err.Error()); err != nil {
			slog.Error("benchmark failed", "job_id", job.ID, "err", err)
		}
		return
	}

	if err := w.repository.CompleteJob(ctx, job.ID, autocannonResult, routingJSON); err != nil {
		slog.Info("benchmark completed", "job_id", job.ID)
	}
}