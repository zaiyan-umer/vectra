import { cn } from "@/lib/utils"
import { formatLatency } from "@/utils/format"
import type { BenchmarkRun } from "@/types/benchmark"

type PercentileItem = { key: keyof BenchmarkRun["latencyPercentiles"]; label: string; highlight?: boolean }

const PERCENTILES: PercentileItem[] = [
  { key: "p50", label: "P50" },
  { key: "p75", label: "P75" },
  { key: "p90", label: "P90" },
  { key: "p99", label: "P99", highlight: true },
  { key: "p999", label: "P99.9" },
  { key: "max", label: "MAX" },
]

export function LatencyPercentiles({ run }: { run: BenchmarkRun }) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-6 gap-space-md">
      {PERCENTILES.map(({ key, label, highlight }) => (
        <div
          key={key}
          className={cn(
            "space-y-1 text-center p-space-sm rounded-lg",
            highlight && "bg-primary/5 ring-1 ring-primary/20"
          )}
        >
          <p className="font-label-caps text-on-surface-variant">{label}</p>
          <p className={cn(
            "font-code-sm",
            highlight ? "text-primary" : "text-on-surface"
          )}>
            {formatLatency(run.latencyPercentiles[key])}
          </p>
        </div>
      ))}
    </div>
  )
}