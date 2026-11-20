import { formatNumber, formatPercent } from "@/utils/format"

interface BackendAllocationBarProps {
  backend: {
    url: string
    requests: number
    percent: number
  }
}

export function BackendAllocationBar({ backend }: BackendAllocationBarProps) {
  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between">
        <span className="font-code-sm text-on-surface">{backend.url}</span>
        <div className="flex items-center gap-2 text-right">
          <span className="font-code-sm text-on-surface">{formatNumber(backend.requests)}</span>
          <span className="font-label-caps text-on-surface-variant">{formatPercent(backend.percent)}</span>
        </div>
      </div>
      <div className="h-2 bg-surface-container overflow-hidden rounded-full">
        <div
          className="h-full bg-primary transition-all duration-300"
          style={{ width: `${backend.percent}%` }}
        />
      </div>
    </div>
  )
}