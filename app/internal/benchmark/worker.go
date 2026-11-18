package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/zaiyan-umer/vectra/app/internal/metrics"
)

type Worker struct {
	repository      *Repository
	runtimeSnapshot *metrics.BenchmarkRuntimeSnapshot
}

func NewWorker(repo *Repository, sh *metrics.BenchmarkRuntimeSnapshot) *Worker {
	return &Worker{
		repository:      repo,
		runtimeSnapshot: sh,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
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
		fmt.Printf("worker could not claim the job: %v", err)
		return
	}

	w.execute(ctx, job)
}

func (w *Worker) execute(ctx context.Context, job *BenchmarkJob) {
	before := w.runtimeSnapshot.CaptureSnapshot()

	bechmarkCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
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
		log.Printf("exec.CommandContext: %v", err)
		if err := w.repository.FailJob(ctx, job.ID, err.Error()); err != nil {
			log.Printf("fail benchmark %q: %v", job.ID, err)
		}
		return
	}

	after := w.runtimeSnapshot.CaptureSnapshot()
	report := after.Difference(before)

	var autocannonResult json.RawMessage
	if err := json.Unmarshal(output, &autocannonResult); err != nil {
		log.Printf("json.unmarshal: %v", err)
		if err := w.repository.FailJob(ctx, job.ID, err.Error()); err != nil {
			log.Printf("fail benchmark %q: %v", job.ID, err)
		}
		return
	}

	routingJSON, err := json.Marshal(report)
	if err != nil {
		log.Printf("json.marshal routing: %v", err)
		if err := w.repository.FailJob(ctx, job.ID, err.Error()); err != nil {
			log.Printf("fail benchmark %q: %v", job.ID, err)
		}
		return
	}

	if err := w.repository.CompleteJob(ctx, job.ID, autocannonResult, routingJSON); err != nil {
		log.Printf("complete benchmark %q: %v", job.ID, err)
	}
}
