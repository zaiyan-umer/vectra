export interface ApiBenchmarkCreateRequest {
  durationSeconds: number;
  connections: number;
}

export interface ApiBenchmarkCreateResponse {
  id: string;
  status: 'queued';
}

export type ApiBenchmarkStatus = 'queued' | 'running' | 'completed' | 'failed' | 'canceled';

export type ApiAlgorithm = 'adaptive_selection' | 'round_robin' | 'healthy_round_robin' | 'least_connections' | 'latency_aware' | 'weighted';

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
  statusCodeStats: Record<string, number | { count: number }>;
  '2xx': number | { count: number };
  '5xx': number | { count: number };
}

export interface ApiRoutingEntry {
  backend: string;
  port: number;
  assignedRequests: number;
  trafficShare: number;
  averageLatencyMs: number;
  url: string;
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

export interface BackendSparkline {
  rps: number[];
  avgLatencyMs: number[];
}

export interface SparklineData {
  timestamps: string[];
  poolRps: number[];
  poolErrorRate: number[];
  poolAvgLatencyMs: number[];
  backends: Record<string, BackendSparkline>;
}

export interface PoolMetrics {
  rps: number;
  errorRate: number;
  avgLatencyMs: number;
}

export interface BackendMetrics {
  index: number;
  url: string;
  healthy: boolean;
  activeConnections: number;
  rps: number;
  errorRate: number;
  avgLatencyMs: number;
  trafficShare: number;
  selectedRequests: number;
  completedRequests: number;
  errors: number;
}

export interface MetricsResponse {
  pool: PoolMetrics;
  backends: BackendMetrics[];
  sparkline: SparklineData;
}

export interface AdaptiveConfig {
  maxExpectedConnections: number;
  latencyTargetMs: number;
  activeWeight: number;
  latencyWeight: number;
  tolerance: number;
}

export interface ConfigResponse {
  algorithm: string;
  targets: string[];
  healthCheckIntervalSeconds: number;
  adaptiveConfig: AdaptiveConfig;
  availableAlgorithms: string[];
}

export interface BackendDetail {
  index: number;
  url: string;
  healthy: boolean;
  activeConnections: number;
  rps: number;
  errorRate: number;
  avgLatencyMs: number;
  trafficShare: number;
  selectedRequests: number;
  completedRequests: number;
  errors: number;
}

export interface BackendsResponse {
  backends: BackendDetail[];
}

export interface BackendHealthStatus {
  index: number;
  url: string;
  healthy: boolean;
  lastCheck: string;
  consecutiveFailures: number;
}

export interface HealthResponse {
  healthyCount: number;
  totalCount: number;
  backends: BackendHealthStatus[];
}