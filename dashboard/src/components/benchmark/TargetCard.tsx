import * as React from "react"
import { cn } from "@/lib/utils"
import { Badge } from "@/components/ui/badge"
import type { BackendTarget } from "@/types/backend"

interface TargetCardProps {
  target: BackendTarget
  className?: string
}

const STATUS_BADGE = {
  healthy: "badge-healthy",
  degraded: "badge-degraded",
  down: "badge-critical",
} as const

export const TargetCard = React.memo(function TargetCard({ target, className }: TargetCardProps) {
  const statusClass = STATUS_BADGE[target.status] || "badge-healthy"

  return (
    <div className={cn("surface-card p-space-md", className)}>
      <div className="flex items-start justify-between gap-space-sm mb-space-sm">
        <div className="flex items-center gap-2">
          <span className={cn(
            "relative flex h-2 w-2 rounded-full",
            target.status === "healthy" && "bg-healthy",
            target.status === "degraded" && "bg-degraded",
            target.status === "down" && "bg-critical"
          )} />
          <span className="font-code-sm text-on-surface">{target.name}</span>
        </div>
        <Badge className={statusClass}>
          {target.status}
        </Badge>
      </div>
      <div className="grid grid-cols-2 gap-space-xs text-body-sm">
        <div>
          <p className="font-label-caps text-on-surface-variant">Latency</p>
          <p className="font-code-sm text-on-surface">{target.avgLatency.toFixed(1)}ms</p>
        </div>
        <div>
          <p className="font-label-caps text-on-surface-variant">Weight</p>
          <p className="font-code-sm text-on-surface">{target.weight}%</p>
        </div>
      </div>
    </div>
  )
})