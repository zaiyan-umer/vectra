import { useQuery } from '@tanstack/react-query';
import { getBackends } from '@/lib/api';
import type { BackendsResponse } from '@/types/api';

const pollIntervalMs = Number(import.meta.env.VITE_POLL_INTERVAL_MS)
const METRICS_INTERVAL = Number.isNaN(pollIntervalMs) || pollIntervalMs <= 0 ? 2000 : pollIntervalMs

export function useBackends(options?: { enabled?: boolean; refetchInterval?: number }) {
  return useQuery<BackendsResponse, Error>({
    queryKey: ['backends'],
    queryFn: getBackends,
    enabled: options?.enabled ?? true,
    refetchInterval: options?.refetchInterval ?? METRICS_INTERVAL,
    staleTime: 0,
  });
}