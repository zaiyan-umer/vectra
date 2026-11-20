import * as React from "react"
import { cn } from "@/lib/utils"
import { formatTimestamp } from "@/utils/format"
import type { BenchmarkRun } from "@/types/benchmark"
import { Badge } from "@/components/ui/badge"
import { STATUS_BADGE_CLASS, STATUS_LABELS, ALGORITHM_LABELS } from "@/utils/constants"

export const RunsTableRow = React.memo(function RunsTableRow({ run, isSelected, onSelect }: {
  run: BenchmarkRun
  isSelected: boolean
  onSelect: () => void
}) {
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
        <div className="font-code-sm text-on-surface">{run.id.slice(0, 8)}</div>
        <div className="font-body-sm text-on-surface-variant">{formatTimestamp(run.timestamp)}</div>
      </td>
      <td className="px-space-sm py-space-xs">
        <Badge variant={statusClass as "default" | "link" | "secondary" | "destructive" | "outline" | "ghost"} className="gap-1">
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
      <td className="px-space-sm py-space-xs">
        <span className="font-label-caps text-on-surface">{ALGORITHM_LABELS[run.algorithm]}</span>
      </td>
      <td className="px-space-sm py-space-xs">
        <span className="font-code-sm text-on-surface-variant">
          {run.durationSeconds}s • {run.connections}c
        </span>
      </td>
    </tr>
  )
})