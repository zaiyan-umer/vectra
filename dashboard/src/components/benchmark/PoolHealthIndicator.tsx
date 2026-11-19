import { cn } from "@/lib/utils"

interface PoolHealthIndicatorProps {
  healthy?: number
  total?: number
  lastCheck?: string
  className?: string
}

export function PoolHealthIndicator({ healthy = 3, total = 3, lastCheck = "2s ago", className }: PoolHealthIndicatorProps) {
  return (
    <div className={cn("flex items-center gap-space-sm font-body-sm", className)}>
      <span className={cn(
        "relative flex h-2 w-2 rounded-full bg-healthy animate-pulse",
        healthy === total ? "bg-healthy" : healthy > 0 ? "bg-degraded" : "bg-critical"
      )} />
      <span className="font-code-sm text-on-surface">
        {healthy}/{total} Healthy
      </span>
      <span className="text-on-surface-variant">• Last check: {lastCheck}</span>
    </div>
  )
}