import { cn } from "@/lib/utils"
import { Sparkline } from "./Sparkline"
import type { MetricCardData } from "@/types/benchmark"

const DELTA_STYLE = {
  positive: "text-healthy",
  negative: "text-critical",
  neutral: "text-on-surface-variant",
} as const

export function MetricCard({ data, className }: { data: MetricCardData; className?: string }) {
  return (
    <div className={cn("telemetry-card", className)}>
      <div className="flex items-start justify-between mb-space-xs">
        <div className="flex items-center gap-2">
          <span className="font-label-caps text-on-surface-variant">{data.label}</span>
        </div>
      </div>
      <div className="flex items-baseline gap-1 mb-space-xs">
        <span className="font-headline-lg font-code-lg text-on-surface">{data.value}</span>
        <span className="font-body-sm text-on-surface-variant">{data.unit}</span>
      </div>
      <div className="flex items-center justify-between mb-space-md">
        <span className={cn("font-label-caps", DELTA_STYLE[data.deltaType || "neutral"])}>
          {data.delta ? (data.deltaType === "positive" ? "+" : "") + data.delta : "—"}
        </span>
        <span className="font-label-caps text-on-surface-variant">{data.footerRight}</span>
      </div>
      <Sparkline data={data.sparklineData} color={data.sparklineColor} />
      <p className="mt-space-xs font-label-caps text-on-surface-variant">{data.footerLeft}</p>
    </div>
  )
}