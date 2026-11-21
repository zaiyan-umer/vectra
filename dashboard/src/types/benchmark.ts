export type Algorithm = "adaptive" | "round_robin" | "least_conn" | "weighted" | "healthy_round_robin" | "latency_aware";

export type RunStatus = "running" | "completed" | "failed" | "canceled";

export interface BenchmarkRun {
  id: string;
  timestamp: string;
  status: RunStatus;
  algorithm: Algorithm;
  durationSeconds: number;
  connections: number;
  throughput: number;
  avgLatency: number;
  p99: number;
  errors: number;
  errorRate: number;
  totalRequests: number;
  dataTransferredBytes: number;
  httpPerSecond: number;
  latencyPercentiles: {
    p50: number;
    p75: number;
    p90: number;
    p99: number;
    p999: number;
    max: number;
  };
  statusBreakdown: {
    ok2xx: number;
    err5xx: number;
    timeouts: number;
  };
  backendAllocation: {
    url: string;
    requests: number;
    percent: number;
  }[];
  targetUrl?: string;
  createdAt?: string;
  error?: string | null;
}

export interface MetricCardData {
  id: string;
  label: string;
  value: string;
  unit: string;
  delta?: string;
  deltaType?: "positive" | "negative" | "neutral";
  sparklineData: number[];
  sparklineColor: string;
  footerLeft: string;
  footerRight: string;
}
