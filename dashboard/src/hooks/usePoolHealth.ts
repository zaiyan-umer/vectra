import { useQuery } from '@tanstack/react-query';
import { getBackendsHealth } from '@/lib/api';
import type { HealthResponse } from '@/types/api';

const METRICS_INTERVAL = 5000;

export function usePoolHealth(options?: { enabled?: boolean; refetchInterval?: number }) {
  return useQuery<HealthResponse, Error>({
    queryKey: ['poolHealth'],
    queryFn: getBackendsHealth,
    enabled: options?.enabled ?? true,
    refetchInterval: options?.refetchInterval ?? METRICS_INTERVAL,
    staleTime: 0,
  });
}