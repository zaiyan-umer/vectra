import { QueryClient } from '@tanstack/react-query';

const staleTimeMs = Number(import.meta.env.VITE_STALE_TIME_MS)
const STALE_TIME = Number.isNaN(staleTimeMs) || staleTimeMs < 0 ? 1000 : staleTimeMs

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: STALE_TIME,
      refetchOnWindowFocus: false,
      retry: 1,
    },
    mutations: {
      retry: 0,
    },
  },
});