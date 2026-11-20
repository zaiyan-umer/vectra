import { useMutation } from '@tanstack/react-query';
import { postBenchmark } from '@/lib/api';
import type { ApiBenchmarkCreateRequest, ApiBenchmarkCreateResponse } from '@/types/api';

export function useStartBenchmark() {
  return useMutation<ApiBenchmarkCreateResponse, Error, ApiBenchmarkCreateRequest>({
    mutationFn: postBenchmark,
  });
}