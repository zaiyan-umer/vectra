import { useQuery } from '@tanstack/react-query';
import { getConfig } from '@/lib/api';
import type { ConfigResponse } from '@/types/api';

export function useConfig(options?: { enabled?: boolean }) {
  return useQuery<ConfigResponse, Error>({
    queryKey: ['config'],
    queryFn: getConfig,
    enabled: options?.enabled ?? true,
    staleTime: 5 * 60 * 1000,
    refetchInterval: false,
  });
}