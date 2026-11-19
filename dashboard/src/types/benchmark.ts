export type Algorithm = "adaptive" | "round_robin" | "least_conn" | "weighted";

export type RunStatus = "running" | "completed" | "failed" | "canceled";

export interface BenchmarkRun {
  id: number;
  timestamp: string;
  status: RunStatus;
  algorithm: Algorithm;
  durationSeconds: number;
  connections: number;
  pipelining: number;
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
    name: string;
    port: number;
    requests: number;
    percent: number;
  }[];
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
