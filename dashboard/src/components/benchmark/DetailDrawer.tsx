import { cn } from "@/lib/utils"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { ExternalLink } from "lucide-react"
import { ScrollArea } from "@/components/ui/scroll-area"
import { ConfigHighlights } from "./ConfigHighlights"
import { AutocannonSummary } from "./AutocannonSummary"
import { LatencyPercentiles } from "./LatencyPercentiles"
import { BackendAllocation } from "./BackendAllocation"
import { StatusBreakdown } from "./StatusBreakdown"
import type { BenchmarkRun } from "@/types/benchmark"

interface DetailDrawerProps {
  run?: BenchmarkRun
  isOpen: boolean
  onClose: () => void
  isMobile?: boolean
  className?: string
}

const ALGORITHM_LABELS = {
  adaptive: "Adaptive Selection",
  round_robin: "Round Robin",
  least_conn: "Least Connections",
  weighted: "Weighted Response",
}

const STATUS_BADGE = {
  running: "secondary",
  completed: "secondary",
  failed: "destructive",
  canceled: "outline",
} as const

export function DetailDrawer({ run, isOpen, onClose, isMobile, className }: DetailDrawerProps) {
  if (!run) {
    return (
      <aside
        className={cn(
          "surface-card",
          isMobile
            ? "fixed inset-y-0 right-0 w-full max-w-md z-50 animate-slide-in"
            : "xl:col-span-5"
        )}
      >
        <div className="flex h-full items-center justify-center">
          <p className="font-body-md text-on-surface-variant">Select a run to view details</p>
        </div>
      </aside>
    )
  }

  const content = (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between px-space-md py-space-sm border-b border-border/50">
        <div className="flex items-center gap-space-sm">
          <span className="font-code-sm text-on-surface">Run #{run.id}</span>
          <Badge variant={STATUS_BADGE[run.status]}>
            {run.status}
          </Badge>
          <span className="font-body-sm text-on-surface-variant">{ALGORITHM_LABELS[run.algorithm]}</span>
        </div>
        <button
          onClick={onClose}
          className="p-1 text-on-surface-variant hover:text-on-surface transition-colors"
          aria-label="Close drawer"
        >
          <ExternalLink className="h-4 w-4" />
        </button>
      </div>

      <ScrollArea className="flex-1">
        <div className="p-space-md space-y-space-lg">
          <ConfigHighlights run={run} />
          <Separator />
          <AutocannonSummary run={run} />
          <Separator />
          <LatencyPercentiles run={run} />
          <Separator />
          <BackendAllocation run={run} />
          <Separator />
          <StatusBreakdown run={run} />
        </div>
      </ScrollArea>
    </div>
  )

  if (isMobile) {
    return (
      <aside className={cn(
        "fixed inset-y-0 right-0 w-full max-w-md z-50 animate-slide-in surface-card flex flex-col",
        isOpen ? "translate-x-0" : "translate-x-full"
      )}>
        {content}
      </aside>
    )
  }

  return (
    <aside className={cn("surface-card xl:col-span-5 flex flex-col", className)}>
      {content}
    </aside>
  )
}