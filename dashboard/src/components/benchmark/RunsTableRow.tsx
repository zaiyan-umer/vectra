import { cn } from "@/lib/utils"
import { ChevronRight } from "lucide-react"
import { formatTimestamp, formatLatency, formatThroughput, formatPercent } from "@/utils/format"
import type { BenchmarkRun, RunStatus } from "@/types/benchmark"
import { Badge } from "@/components/ui/badge"

interface RunsTableRowProps {
  run: BenchmarkRun
  isSelected: boolean
  onSelect: () => void
}

const STATUS_BADGE_CLASS: Record<RunStatus, "secondary" | "destructive" | "outline"> = {
  running: "secondary",
  completed: "secondary",
  failed: "destructive",
  canceled: "outline",
}

const STATUS_LABELS: Record<RunStatus, string> = {
  running: "Running",
  completed: "Completed",
  failed: "Failed",
  canceled: "Canceled",
}

const ALGORITHM_LABELS: Record<BenchmarkRun["algorithm"], string> = {
  adaptive: "Adaptive",
  round_robin: "Round Robin",
  least_conn: "Least Conn",
  weighted: "Weighted",
}

export function RunsTableRow({ run, isSelected, onSelect }: RunsTableRowProps) {
  const statusClass = STATUS_BADGE_CLASS[run.status]

  return (
    <tr
      className={cn(
        "table-row transition-colors cursor-pointer",
        isSelected && "bg-primary/5"
      )}
      onClick={onSelect}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault()
          onSelect()
        }
      }}
      tabIndex={0}
      role="button"
      aria-selected={isSelected}
    >
      <td className="px-space-md py-space-xs">
        <div className="font-code-sm text-on-surface">#{run.id}</div>
        <div className="font-body-sm text-on-surface-variant">{formatTimestamp(run.timestamp)}</div>
      </td>
      <td className="px-space-md py-space-xs">
        <Badge variant={statusClass} className="gap-1">
          <span className={cn(
            "relative flex h-1.5 w-1.5 rounded-full animate-pulse",
            run.status === "running" && "bg-degraded",
            run.status === "completed" && "bg-healthy",
            run.status === "failed" && "bg-critical",
            run.status === "canceled" && "bg-outline"
          )} />
          {STATUS_LABELS[run.status]}
        </Badge>
      </td>
      <td className="px-space-md py-space-xs">
        <span className="font-label-caps text-on-surface">{ALGORITHM_LABELS[run.algorithm]}</span>
      </td>
      <td className="px-space-md py-space-xs">
        <span className="font-code-sm text-on-surface-variant">
          {run.durationSeconds}s • {run.connections}c • {run.pipelining}p
        </span>
      </td>
      <td className="px-space-md py-space-xs text-right">
        <span className="font-code-sm text-on-surface">{formatThroughput(run.throughput)}</span>
        <span className="font-label-caps text-on-surface-variant ml-1">req/s</span>
      </td>
      <td className="px-space-md py-space-xs text-right">
        <span className="font-code-sm text-on-surface">{formatLatency(run.p99)}</span>
      </td>
      <td className="px-space-md py-space-xs text-right">
        <span className="font-code-sm text-on-surface">{run.errors}</span>
        <span className="font-label-caps text-on-surface-variant ml-1">({formatPercent(run.errorRate * 100, 2)})</span>
      </td>
      <td className="px-space-md py-space-xs text-center">
        <ChevronRight className={cn(
          "mx-auto h-4 w-4 text-on-surface-variant transition-transform",
          isSelected && "rotate-90"
        )} />
      </td>
    </tr>
  )
}