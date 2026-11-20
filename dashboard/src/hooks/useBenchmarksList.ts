import { useQuery } from '@tanstack/react-query';
import { getBenchmarks, transformBenchmark } from '@/lib/api';
import type { ApiBenchmarkDetail } from '@/types/api';
import type { BenchmarkRun } from '@/types/benchmark';

export function useBenchmarksList() {
  return useQuery<ApiBenchmarkDetail[], Error, BenchmarkRun[]>({
    queryKey: ['benchmarks'],
    queryFn: getBenchmarks,
    select: (data) => data.map(transformBenchmark),
    staleTime: 0,
    refetchInterval: 5000,
  });
}