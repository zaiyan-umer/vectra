import type { Algorithm, RunStatus } from "@/types/benchmark";

export const ALGORITHMS: { value: Algorithm; label: string; short: string }[] = [
  { value: "adaptive", label: "Adaptive Selection", short: "Adaptive" },
  { value: "round_robin", label: "Round Robin", short: "Round Robin" },
  { value: "least_conn", label: "Least Connections", short: "Least Conn" },
  { value: "weighted", label: "Weighted Response", short: "Weighted" },
  { value: "healthy_round_robin", label: "Healthy Round Robin", short: "Healthy RR" },
  { value: "latency_aware", label: "Latency Aware", short: "Latency" },
];

export const ALGORITHM_LABELS: Record<Algorithm, string> = {
  adaptive: "Adaptive",
  round_robin: "Round Robin",
  least_conn: "Least Conn",
  weighted: "Weighted",
  healthy_round_robin: "Healthy Round Robin",
  latency_aware: "Latency Aware",
};

export const API_ALGORITHM_LABELS: Record<string, string> = {
  adaptive_selection: "Adaptive Selection",
  round_robin: "Round Robin",
  healthy_round_robin: "Healthy Round Robin",
  least_connections: "Least Connections",
  latency_aware: "Latency Aware",
  weighted: "Weighted Response",
};

export const STATUS_LABELS: Record<RunStatus, string> = {
  running: "Running",
  completed: "Completed",
  failed: "Failed",
  canceled: "Canceled",
};

export const API_STATUS_LABELS: Record<string, string> = {
  queued: "Queued",
  running: "Running",
  completed: "Completed",
  failed: "Failed",
  canceled: "Canceled",
};

export const STATUS_BADGE_CLASS: Record<RunStatus, string> = {
  running: "badge-degraded",
  completed: "badge-healthy",
  failed: "badge-critical",
  canceled: "badge-degraded",
};

export const API_STATUS_BADGE_CLASS: Record<string, string> = {
  queued: "badge-degraded",
  running: "badge-degraded",
  completed: "badge-healthy",
  failed: "badge-critical",
  canceled: "badge-degraded",
};

export const STATUS_DOT_COLOR: Record<RunStatus, string> = {
  running: "bg-[var(--color-degraded)]",
  completed: "bg-[var(--color-healthy)]",
  failed: "bg-[var(--color-critical)]",
  canceled: "bg-[var(--color-outline)]",
};

export const BACKEND_STATUS_BADGE_CLASS = {
  healthy: "badge-healthy",
  degraded: "badge-degraded",
  down: "badge-critical",
} as const;

export const CHART_COLORS = {
  primary: "#2563eb",
  healthy: "#10b981",
  degraded: "#f59e0b",
  critical: "#ef4444",
  neutral: "#737686",
} as const;

export const DURATION_PRESETS = [10, 30, 60, 120] as const;

export const CONCURRENCY_PRESETS = [10, 25, 50, 100, 250] as const;
