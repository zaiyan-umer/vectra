import { useQuery } from '@tanstack/react-query';
import { getBackendsHealth } from '@/lib/api';
import type { HealthResponse } from '@/types/api';

const pollIntervalMs = Number(import.meta.env.VITE_POLL_INTERVAL_MS)
const METRICS_INTERVAL = Number.isNaN(pollIntervalMs) || pollIntervalMs <= 0 ? 2000 : pollIntervalMs

export function usePoolHealth(options?: { enabled?: boolean; refetchInterval?: number }) {
  return useQuery<HealthResponse, Error>({
    queryKey: ['poolHealth'],
    queryFn: getBackendsHealth,
    enabled: options?.enabled ?? true,
    refetchInterval: options?.refetchInterval ?? METRICS_INTERVAL,
    staleTime: 0,
  });
}