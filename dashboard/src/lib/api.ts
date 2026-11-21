import type {
  ApiBenchmarkCreateRequest,
  ApiBenchmarkCreateResponse,
  ApiBenchmarkDetail,
  ApiBenchmarkStatus,
  ApiAlgorithm,
  ApiAutocannonResult,
  MetricsResponse,
  ConfigResponse,
  BackendsResponse,
  HealthResponse,
} from '@/types/api';
import type { BenchmarkRun } from '@/types/benchmark';

const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080';

const STATUS_MAP: Record<ApiBenchmarkStatus, BenchmarkRun['status']> = {
  queued: 'running',
  running: 'running',
  completed: 'completed',
  failed: 'failed',
  canceled: 'canceled',
};

const ALGORITHM_MAP: Record<ApiAlgorithm, BenchmarkRun['algorithm']> = {
  adaptive_selection: 'adaptive',
  round_robin: 'round_robin',
  healthy_round_robin: 'healthy_round_robin',
  least_connections: 'least_conn',
  latency_aware: 'latency_aware',
  weighted: 'weighted',
};

function getNumericValue(value: number | { count: number } | undefined): number {
  if (value === undefined) return 0;
  if (typeof value === 'number') return value;
  if (typeof value === 'object' && value !== null && 'count' in value) {
    return value.count;
  }
  return 0;
}

function aggregateStatusCodes(autocannon: ApiAutocannonResult) {
  const ok2xx = getNumericValue(autocannon['2xx']) ?? 0;
  const err5xx = getNumericValue(autocannon['5xx']) ?? 0;

  return { ok2xx, err5xx };
}

function transformBenchmark(api: ApiBenchmarkDetail): BenchmarkRun {
  const autocannon = api.result?.autocannon;
  const routing = api.result?.routing ?? [];

  const { ok2xx, err5xx } = autocannon
    ? aggregateStatusCodes(autocannon)
    : { ok2xx: 0, err5xx: 0 };

  return {
    id: api.id,
    timestamp: api.startedAt ?? api.createdAt,
    status: STATUS_MAP[api.status],
    algorithm: ALGORITHM_MAP[api.algorithm] ?? 'adaptive',
    durationSeconds: api.durationSeconds,
    connections: api.connections,
    throughput: autocannon?.throughput?.mean ?? 0,
    avgLatency: autocannon?.latency?.mean ?? 0,
    p99: autocannon?.latency?.p99 ?? 0,
    errors: (autocannon?.errors ?? 0) + (autocannon?.non2xx ?? 0),
    errorRate: autocannon?.requests?.total
      ? (((autocannon?.errors ?? 0) + (autocannon?.non2xx ?? 0)) / autocannon.requests.total) * 100
      : 0,
    totalRequests: autocannon?.requests?.total ?? 0,
    dataTransferredBytes: 0,
    httpPerSecond: autocannon?.throughput?.mean ?? 0,
    latencyPercentiles: {
      p50: autocannon?.latency?.p50 ?? 0,
      p75: autocannon?.latency?.p75 ?? 0,
      p90: autocannon?.latency?.p90 ?? 0,
      p99: autocannon?.latency?.p99 ?? 0,
      p999: autocannon?.latency?.p99_9 ?? 0,
      max: autocannon?.latency?.max ?? 0,
    },
    statusBreakdown: {
      ok2xx,
      err5xx,
      timeouts: autocannon?.timeouts ?? 0,
    },
    backendAllocation: routing.map((r) => ({
      url: r.url,
      requests: r.assignedRequests,
      percent: r.trafficShare,
    })),
    targetUrl: api.targetUrl,
    createdAt: api.createdAt,
    error: api.error,
  };
}

async function fetchJson<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ message: res.statusText }));
    throw new Error(err.message ?? `HTTP ${res.status}`);
  }

  return res.json();
}

export async function postBenchmark(
  config: ApiBenchmarkCreateRequest
): Promise<ApiBenchmarkCreateResponse> {
  return fetchJson<ApiBenchmarkCreateResponse>('/benchmarks', {
    method: 'POST',
    body: JSON.stringify(config),
  });
}

export async function getBenchmark(id: string): Promise<ApiBenchmarkDetail> {
  return fetchJson<ApiBenchmarkDetail>(`/benchmarks/${id}`);
}

export async function getBenchmarks(): Promise<ApiBenchmarkDetail[]> {
  return fetchJson<ApiBenchmarkDetail[]>('/benchmarks');
}

export async function getMetrics(): Promise<MetricsResponse> {
  return fetchJson<MetricsResponse>('/metrics');
}

export async function getConfig(): Promise<ConfigResponse> {
  return fetchJson<ConfigResponse>('/config');
}

export async function getBackends(): Promise<BackendsResponse> {
  return fetchJson<BackendsResponse>('/backends');
}

export async function getBackendsHealth(): Promise<HealthResponse> {
  return fetchJson<HealthResponse>('/backends/health');
}

export { transformBenchmark, STATUS_MAP, ALGORITHM_MAP };