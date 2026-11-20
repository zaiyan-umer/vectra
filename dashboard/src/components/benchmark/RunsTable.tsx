import { cn } from "@/lib/utils"
import type { BenchmarkRun } from "@/types/benchmark"
import { RunsTableRow } from "./RunsTableRow"
import { Button } from "@/components/ui/button"
import { ChevronRight } from "lucide-react"
import { Loader2 } from "lucide-react"

interface RunsTableProps {
  runs: BenchmarkRun[]
  selectedRunId?: string
  onSelectRun: (id: string) => void
  onPageChange?: (page: number) => void
  currentPage?: number
  totalPages?: number
  className?: string
  isLoading?: boolean
}

const ROWS_PER_PAGE = 10

export function RunsTable({
  runs,
  selectedRunId,
  onSelectRun,
  onPageChange,
  currentPage = 1,
  totalPages = 1,
  className,
  isLoading = false,
}: RunsTableProps) {
  const paginatedRuns = runs.slice((currentPage - 1) * ROWS_PER_PAGE, currentPage * ROWS_PER_PAGE)

  if (isLoading && runs.length === 0) {
    return (
      <div className={cn("surface-card overflow-hidden", className)}>
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <span className="ml-3 font-body-sm text-on-surface-variant">Loading runs...</span>
        </div>
      </div>
    )
  }

  return (
    <div className={cn("surface-card overflow-hidden", className)}>
      <div className="overflow-x-auto">
        <table className="w-full" role="table">
          <thead>
            <tr className="border-b border-border/50">
              <th className="table-header px-space-md text-left">Run ID / Time</th>
              <th className="table-header px-space-md text-left">Status</th>
              <th className="table-header px-space-md text-left">Algorithm</th>
              <th className="table-header px-space-md text-left">Load Config</th>
              <th className="table-header px-space-md text-right">Throughput</th>
              <th className="table-header px-space-md text-right">P99</th>
              <th className="table-header px-space-md text-right">Errors</th>
              <th className="table-header px-space-md text-center">Inspect</th>
            </tr>
          </thead>
          <tbody>
            {paginatedRuns.map((run) => (
              <RunsTableRow
                key={run.id}
                run={run}
                isSelected={selectedRunId === run.id}
                onSelect={() => onSelectRun(run.id)}
              />
            ))}
            {paginatedRuns.length === 0 && !isLoading && (
              <tr>
                <td colSpan={8} className="px-space-md py-space-lg text-center text-on-surface-variant">
                  No runs found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      {totalPages > 1 && (
        <div className="flex items-center justify-between px-space-md py-space-sm border-t border-border/50">
          <p className="font-body-sm text-on-surface-variant">
            Page {currentPage} of {totalPages}
          </p>
          <div className="flex items-center gap-1">
            <Button
              variant="outline"
              size="sm"
              onClick={() => onPageChange?.(currentPage - 1)}
              disabled={currentPage <= 1}
              aria-label="Previous page"
            >
              <ChevronRight className="h-4 w-4 rotate-180" />
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => onPageChange?.(currentPage + 1)}
              disabled={currentPage >= totalPages}
              aria-label="Next page"
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}