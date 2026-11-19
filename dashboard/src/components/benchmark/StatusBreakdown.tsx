import { cn } from "@/lib/utils"
import { formatNumber } from "@/utils/format"
import type { BenchmarkRun } from "@/types/benchmark"
import { Badge } from "@/components/ui/badge"

interface StatusBreakdownProps {
  run: BenchmarkRun
}

export function StatusBreakdown({ run }: StatusBreakdownProps) {
  const items = [
    { label: "2xx", value: run.statusBreakdown.ok2xx, className: "badge-healthy" },
    { label: "5xx", value: run.statusBreakdown.err5xx, className: "badge-critical" },
    { label: "Timeouts", value: run.statusBreakdown.timeouts, className: "badge-degraded" },
  ]

  return (
    <div className="flex flex-wrap gap-space-md">
      {items.map((item) => (
        <div key={item.label} className="flex items-center gap-2">
          <Badge className={cn("gap-1.5", item.className)}>
            {item.label}
            <span className="font-code-sm">{formatNumber(item.value)}</span>
          </Badge>
        </div>
      ))}
    </div>
  )
}