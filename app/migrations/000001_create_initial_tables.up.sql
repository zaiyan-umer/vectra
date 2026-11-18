CREATE TABLE benchmarks (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL DEFAULT 'queued'
        CHECK (status IN ('queued', 'running', 'completed', 'failed', 'canceled')),
    duration_seconds INTEGER NOT NULL CHECK (duration_seconds > 0),
    connections INTEGER NOT NULL CHECK (connections > 0),
    algorithm TEXT NOT NULL,
    target_url TEXT NOT NULL,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_benchmarks_status ON benchmarks (status);
CREATE INDEX idx_benchmarks_created_at ON benchmarks (created_at DESC);

CREATE TABLE benchmark_results (
    benchmark_id TEXT PRIMARY KEY REFERENCES benchmarks(id) ON DELETE CASCADE,
    autocannon_json JSONB NOT NULL,
    routing_json JSONB NOT NULL
);

