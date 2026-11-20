import { formatNumber, formatLatency, formatThroughput } from "@/utils/format"
import type { BenchmarkRun } from "@/types/benchmark"

interface AutocannonSummaryProps {
  run: BenchmarkRun
}

export function AutocannonSummary({ run }: AutocannonSummaryProps) {
  return (
    <div className="flex flex-col md:flex-row gap-space-md">
      <div className="flex-1 space-y-1">
        <p className="font-label-caps text-on-surface-variant">Total Requests</p>
        <p className="font-code-sm text-on-surface">{formatNumber(run.totalRequests)}</p>
      </div>
      <div className="flex-1 space-y-1">
        <p className="font-label-caps text-on-surface-variant">Throughput</p>
        <p className="font-code-sm text-on-surface">{formatThroughput(run.throughput)} req/s</p>
      </div>
      <div className="flex-1 space-y-1">
        <p className="font-label-caps text-on-surface-variant">HTTP Time/Sec</p>
        <p className="font-code-sm text-on-surface">{formatLatency(run.avgLatency)}</p>
      </div>
    </div>
  )
}