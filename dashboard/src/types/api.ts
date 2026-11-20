export interface ApiBenchmarkCreateRequest {
  durationSeconds: number;
  connections: number;
}

export interface ApiBenchmarkCreateResponse {
  id: string;
  status: 'queued';
}

export type ApiBenchmarkStatus = 'queued' | 'running' | 'completed' | 'failed' | 'canceled';

export type ApiAlgorithm = 'adaptive_selection' | 'round_robin' | 'least_conn' | 'weighted';

export interface ApiAutocannonLatency {
  mean: number;
  p1: number;
  p2_5: number;
  p50: number;
  p75: number;
  p90: number;
  p99: number;
  p99_9: number;
  p99_99: number;
  p99_999: number;
  max: number;
  min: number;
}

export interface ApiAutocannonThroughput {
  mean: number;
  max: number;
  min: number;
}

export interface ApiAutocannonRequests {
  total: number;
  mean: number;
  max: number;
  min: number;
}

export interface ApiAutocannonResult {
  latency: ApiAutocannonLatency;
  throughput: ApiAutocannonThroughput;
  requests: ApiAutocannonRequests;
  errors: number;
  non2xx: number;
  timeouts: number;
  statusCodeStats: Record<string, number>;
}

export interface ApiRoutingEntry {
  backend: string;
  port: number;
  assignedRequests: number;
  trafficShare: number;
  averageLatencyMs: number;
}

export interface ApiBenchmarkDetail {
  id: string;
  status: ApiBenchmarkStatus;
  durationSeconds: number;
  connections: number;
  algorithm: ApiAlgorithm;
  targetUrl: string;
  startedAt: string | null;
  finishedAt: string | null;
  error: string | null;
  createdAt: string;
  result: {
    autocannon: ApiAutocannonResult;
    routing: ApiRoutingEntry[];
  } | null;
}