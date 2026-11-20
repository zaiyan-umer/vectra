import { BackendAllocationBar } from "./BackendAllocationBar"
import type { BenchmarkRun } from "@/types/benchmark"

interface BackendAllocationProps {
  run: BenchmarkRun
}

export function BackendAllocation({ run }: BackendAllocationProps) {
  return (
    <div className="space-y-space-sm">
      <p className="font-label-caps text-on-surface-variant">Backend Allocation</p>
      <div className="space-y-2">
        {run.backendAllocation.map((backend) => (
          <BackendAllocationBar key={backend.url} backend={backend} />
        ))}
      </div>
    </div>
  )
}