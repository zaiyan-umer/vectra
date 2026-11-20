import { useQuery } from '@tanstack/react-query';
import { getBenchmark, transformBenchmark } from '@/lib/api';
import type { ApiBenchmarkDetail, ApiBenchmarkStatus } from '@/types/api';
import type { BenchmarkRun } from '@/types/benchmark';

const ACTIVE_STATUSES: ApiBenchmarkStatus[] = ['queued', 'running'];

export function useBenchmarkDetail(id: string | undefined, options?: { enabled?: boolean }) {
  const query = useQuery<ApiBenchmarkDetail, Error, BenchmarkRun>({
    queryKey: ['benchmark', id],
    queryFn: () => getBenchmark(id!),
    enabled: !!id && (options?.enabled ?? true),
    select: transformBenchmark,
    refetchInterval: (query) => {
      const data = query.state.data;
      if (!data) return false;
      return ACTIVE_STATUSES.includes(data.status) ? 2000 : false;
    },
  });

  return query;
}