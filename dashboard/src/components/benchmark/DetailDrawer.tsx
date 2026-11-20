import { cn } from "@/lib/utils"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { Loader2, AlertCircle, ExternalLink } from "lucide-react"
import { ScrollArea } from "@/components/ui/scroll-area"
import { ConfigHighlights } from "./ConfigHighlights"
import { AutocannonSummary } from "./AutocannonSummary"
import { LatencyPercentiles } from "./LatencyPercentiles"
import { BackendAllocation } from "./BackendAllocation"
import { StatusBreakdown } from "./StatusBreakdown"
import type { BenchmarkRun, RunStatus } from "@/types/benchmark"
import { API_ALGORITHM_LABELS } from "@/utils/constants"
import { formatDateTime } from "@/utils/format"

const STATUS_BADGE_VARIANT: Record<RunStatus, "default" | "secondary" | "destructive" | "outline" | "ghost" | "healthy"> = {
  completed: "healthy",
  running: "secondary",
  failed: "destructive",
  canceled: "outline",
};

interface DetailDrawerProps {
  run?: BenchmarkRun
  isOpen: boolean
  isMobile?: boolean
  isLoading?: boolean
  className?: string
}

export function DetailDrawer({ run, isOpen, isMobile, isLoading = false, className }: DetailDrawerProps) {
  if (!run) {
    return (
      <aside
        className={cn(
          "surface-card",
          isMobile
            ? "fixed inset-y-0 right-0 w-full max-w-md z-50 animate-slide-in"
            : "xl:col-span-7"
        )}
      >
        <div className="flex h-full items-center justify-center">
          <p className="font-body-md text-on-surface-variant">Select a run to view details</p>
        </div>
      </aside>
    )
  }

  if (isLoading) {
    return (
      <aside className={cn(
        "surface-card flex flex-col",
        isMobile
          ? "fixed inset-y-0 right-0 w-full max-w-md z-50 animate-slide-in"
          : "xl:col-span-7"
      )}>
        <div className="flex items-center justify-between px-space-md py-space-sm border-b border-border/50">
          <div className="flex items-center gap-space-sm">
            <Loader2 className="h-4 w-4 animate-spin text-primary" />
            <span className="font-code-sm text-on-surface">Loading...</span>
          </div>
        </div>
        <div className="flex-1 flex items-center justify-center">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      </aside>
    )
  }

  const content = (
    <div className="flex flex-col h-full">
      <div className="bg-primary/5 px-space-md py-space-md border-b border-border/50">
        <div className="flex items-center justify-between flex-wrap gap-space-md">
          <div className="flex items-center gap-space-sm flex-wrap">
            <span className="font-headline-md font-bold text-on-surface">
              Run #{run.id.slice(0, 8)} Details
            </span>
          </div>
          <div className="flex items-center gap-space-md flex-wrap">
            <span className="font-body-sm text-on-surface-variant">
              {formatDateTime(run.timestamp)}
            </span>
            <Badge variant={STATUS_BADGE_VARIANT[run.status]}>
              {run.status}
            </Badge>
            {run.error && (
              <Badge variant="destructive" className="gap-1">
                <AlertCircle className="h-3 w-3" />
                Failed: {run.error}
              </Badge>
            )}
            <span className="font-body-sm text-on-surface-variant">
              {API_ALGORITHM_LABELS[run.algorithm] ?? run.algorithm}
            </span>
          </div>
        </div>
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
          {run.targetUrl && (
            <>
              <Separator />
              <div className="p-3 bg-surface-container rounded-lg border border-border/50">
                <div className="font-label-caps text-on-surface-variant mb-1">Target URL</div>
                <a
                  href={run.targetUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="font-code-sm text-primary hover:underline flex items-center gap-1"
                >
                  {run.targetUrl}
                  <ExternalLink className="h-3 w-3" />
                </a>
              </div>
            </>
          )}
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
    <aside className={cn("surface-card xl:col-span-7 flex flex-col", className)}>
      {content}
    </aside>
  )
}