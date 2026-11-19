import { formatNumber } from "@/utils/format"
import type { BenchmarkRun } from "@/types/benchmark"

interface ConfigHighlightsProps {
  run: BenchmarkRun
}

const ALGORITHM_LABELS = {
  adaptive: "Adaptive",
  round_robin: "Round Robin",
  least_conn: "Least Conn",
  weighted: "Weighted",
}

export function ConfigHighlights({ run }: ConfigHighlightsProps) {
  return (
    <div className="grid grid-cols-3 gap-space-md">
      <div className="space-y-1">
        <p className="font-label-caps text-on-surface-variant">Connections</p>
        <p className="font-code-sm text-on-surface">{formatNumber(run.connections)}</p>
      </div>
      <div className="space-y-1">
        <p className="font-label-caps text-on-surface-variant">Pipelining</p>
        <p className="font-code-sm text-on-surface">{run.pipelining}</p>
      </div>
      <div className="space-y-1">
        <p className="font-label-caps text-on-surface-variant">Engine</p>
        <p className="font-label-caps text-on-surface">{ALGORITHM_LABELS[run.algorithm]}</p>
      </div>
    </div>
  )
}