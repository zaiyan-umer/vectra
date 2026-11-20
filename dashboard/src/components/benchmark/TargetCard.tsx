import * as React from "react"
import { cn } from "@/lib/utils"
import { Badge } from "@/components/ui/badge"
import type { BackendDetail } from "@/types/api"
import { formatLatency, formatRPS, formatErrorRate } from "@/utils/format"

interface TargetCardProps {
  target: BackendDetail
  className?: string
}

export const TargetCard = React.memo(function TargetCard({ target, className }: TargetCardProps) {
  const statusClass = target.healthy ? "badge-healthy" : "badge-critical"

  const url = new URL(target.url);
  const name = url.host;

  return (
    <div className={cn("surface-card p-space-md", className)}>
      <div className="flex items-start justify-between gap-space-sm mb-space-sm">
        <div className="flex items-center gap-2">
          <span className={cn(
            "relative flex h-2 w-2 rounded-full",
            target.healthy && "bg-healthy",
            !target.healthy && "bg-critical"
          )} />
          <span className="font-code-sm text-on-surface">{name}</span>
        </div>
        <Badge className={statusClass}>
          {target.healthy ? "healthy" : "down"}
        </Badge>
      </div>
      <div className="grid grid-cols-2 gap-space-xs text-body-sm mb-space-sm">
        <div>
          <p className="font-label-caps text-on-surface-variant">Latency</p>
          <p className="font-code-sm text-on-surface">{formatLatency(target.avgLatencyMs)}</p>
        </div>
        <div>
          <p className="font-label-caps text-on-surface-variant">Traffic</p>
          <p className="font-code-sm text-on-surface">{target.trafficShare.toFixed(1)}%</p>
        </div>
      </div>
      <div className="grid grid-cols-2 gap-space-xs text-body-sm">
        <div>
          <p className="font-label-caps text-on-surface-variant">RPS</p>
          <p className="font-code-sm text-on-surface">{formatRPS(target.rps)}</p>
        </div>
        <div>
          <p className="font-label-caps text-on-surface-variant">Error Rate</p>
          <p className="font-code-sm text-on-surface">{formatErrorRate(target.errorRate)}%</p>
        </div>
      </div>
    </div>
  )
})