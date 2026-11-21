import { useQuery } from '@tanstack/react-query';
import { getMetrics } from '@/lib/api';
import type { MetricsResponse } from '@/types/api';

const pollIntervalMs = Number(import.meta.env.VITE_POLL_INTERVAL_MS)
const METRICS_INTERVAL = Number.isNaN(pollIntervalMs) || pollIntervalMs <= 0 ? 2000 : pollIntervalMs

export function useMetrics(options?: { enabled?: boolean; refetchInterval?: number }) {
  return useQuery<MetricsResponse, Error>({
    queryKey: ['metrics'],
    queryFn: getMetrics,
    enabled: options?.enabled ?? true,
    refetchInterval: options?.refetchInterval ?? METRICS_INTERVAL,
    staleTime: 0,
  });
}