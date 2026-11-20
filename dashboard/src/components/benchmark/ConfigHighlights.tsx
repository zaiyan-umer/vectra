import { formatNumber } from "@/utils/format"
import type { BenchmarkRun } from "@/types/benchmark"
import { ALGORITHM_LABELS } from "@/utils/constants"

interface ConfigHighlightsProps {
  run: BenchmarkRun
}

export function ConfigHighlights({ run }: ConfigHighlightsProps) {
  return (
    <div className="grid grid-cols-3 gap-space-md">
      <div className="space-y-1">
        <p className="font-label-caps text-on-surface-variant">Duration</p>
        <p className="font-code-sm text-on-surface">{run.durationSeconds}s</p>
      </div>
      <div className="space-y-1">
        <p className="font-label-caps text-on-surface-variant">Connections</p>
        <p className="font-code-sm text-on-surface">{formatNumber(run.connections)}</p>
      </div>
      <div className="space-y-1">
        <p className="font-label-caps text-on-surface-variant">Engine</p>
        <p className="font-label-caps text-on-surface">{ALGORITHM_LABELS[run.algorithm]}</p>
      </div>
    </div>
  )
}