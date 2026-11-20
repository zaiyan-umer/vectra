import { useQuery } from '@tanstack/react-query';
import { getBackends } from '@/lib/api';
import type { BackendsResponse } from '@/types/api';

const METRICS_INTERVAL = 5000;

export function useBackends(options?: { enabled?: boolean; refetchInterval?: number }) {
  return useQuery<BackendsResponse, Error>({
    queryKey: ['backends'],
    queryFn: getBackends,
    enabled: options?.enabled ?? true,
    refetchInterval: options?.refetchInterval ?? METRICS_INTERVAL,
    staleTime: 0,
  });
}