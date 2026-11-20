import { useQuery } from '@tanstack/react-query';
import { getMetrics } from '@/lib/api';
import type { MetricsResponse } from '@/types/api';

const METRICS_INTERVAL = 5000;

export function useMetrics(options?: { enabled?: boolean; refetchInterval?: number }) {
  return useQuery<MetricsResponse, Error>({
    queryKey: ['metrics'],
    queryFn: getMetrics,
    enabled: options?.enabled ?? true,
    refetchInterval: options?.refetchInterval ?? METRICS_INTERVAL,
    staleTime: 0,
  });
}