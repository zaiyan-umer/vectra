import { useState, type FormEvent } from "react"
import { cn } from "@/lib/utils"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Skeleton } from "@/components/ui/skeleton"
import { DURATION_PRESETS, CONCURRENCY_PRESETS } from "@/utils/constants"
import { ALGORITHM_LABELS } from "@/utils/constants"
import type { Algorithm } from "@/types/benchmark"

interface SystemConfig {
  algorithm: string
  targets: string[]
  healthyCount: number
  totalCount: number
  isLoading: boolean
}

interface BenchmarkModalProps {
  isOpen: boolean
  onClose: () => void
  onStart: (config: {
    durationSeconds: number
    connections: number
  }) => void
  systemConfig: SystemConfig
}

export function BenchmarkModal({ isOpen, onClose, onStart, systemConfig }: BenchmarkModalProps) {
  const [duration, setDuration] = useState(30)
  const [concurrency, setConcurrency] = useState(100)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validate = () => {
    const newErrors: Record<string, string> = {}
    if (!duration || duration < 1) newErrors.duration = "Duration must be at least 1 second"
    if (!concurrency || concurrency < 1) newErrors.concurrency = "Concurrency must be at least 1"
    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (validate()) {
      onStart({ durationSeconds: duration, connections: concurrency })
      onClose()
    }
  }

  const isValid = duration > 0 && concurrency > 0 && !systemConfig.isLoading

  const algorithmLabel = systemConfig.isLoading
    ? "—"
    : ALGORITHM_LABELS[systemConfig.algorithm as Algorithm] ?? systemConfig.algorithm

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <form onSubmit={handleSubmit}>
          <DialogHeader className="pb-space-xl border-b border-border/50 mb-space-xl">
            <DialogTitle className="mb-space-md font-bold">New Benchmark</DialogTitle>
          </DialogHeader>

          <div className="surface-toolbar p-space-lg rounded-lg mb-space-xl" aria-live="polite">
            <p className="font-label-caps text-on-surface-variant mb-space-sm">System State</p>
            <div className="grid grid-cols-3 gap-space-md text-sm">
              <div>
                <p className="font-label-caps text-on-surface-variant">Targets</p>
                {systemConfig.isLoading ? (
                  <Skeleton className="h-4 w-16" />
                ) : (
                  <p className="font-code-sm text-on-surface">
                    {systemConfig.healthyCount}/{systemConfig.totalCount} healthy
                  </p>
                )}
              </div>
              <div className="min-w-0">
                <p className="font-label-caps text-on-surface-variant">Algorithm</p>
                {systemConfig.isLoading ? (
                  <Skeleton className="h-4 w-24" />
                ) : (
                  <p className="font-code-sm text-on-surface whitespace-normal break-words">
                    {algorithmLabel}
                  </p>
                )}
              </div>
              <div className="text-right">
                <p className="font-label-caps text-on-surface-variant">Defaults</p>
                {systemConfig.isLoading ? (
                  <Skeleton className="h-4 w-20 mx-auto" />
                ) : (
                  <p className="font-code-sm text-on-surface">
                    {DURATION_PRESETS[1]}s • {CONCURRENCY_PRESETS[2]}c
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className="space-y-space-lg pb-space-md mb-space-xl">
            <div>
              <Label htmlFor="duration" className="font-label-caps">Duration (seconds)</Label>
              <div className="flex flex-wrap gap-1.5 mt-1.5">
                {DURATION_PRESETS.map((preset) => (
                  <button
                    key={preset}
                    type="button"
                    onClick={() => setDuration(preset)}
                    className={cn(
                      "px-3 py-1.5 rounded-lg text-sm font-mono transition-colors cursor-pointer",
                      duration === preset
                        ? "bg-primary text-primary-foreground"
                        : "bg-surface-container hover:bg-surface-container-high text-on-surface"
                    )}
                  >
                    {preset}s
                  </button>
                ))}
              </div>
              <div className="mt-2.5 flex items-center gap-2">
                <Input
                  id="duration"
                  type="number"
                  min="1"
                  max="3600"
                  value={duration}
                  onChange={(e) => setDuration(Math.max(1, parseInt(e.target.value) || 1))}
                  className="w-24"
                  aria-label="Custom duration"
                  disabled={systemConfig.isLoading}
                />
                {errors.duration && (
                  <p className="font-label-caps text-critical">{errors.duration}</p>
                )}
              </div>
            </div>

            <div>
              <Label htmlFor="concurrency" className="font-label-caps">Concurrency</Label>
              <div className="flex flex-wrap gap-1.5 mt-1.5">
                {CONCURRENCY_PRESETS.map((preset) => (
                  <button
                    key={preset}
                    type="button"
                    onClick={() => setConcurrency(preset)}
                    className={cn(
                      "px-3 py-1.5 rounded-lg text-sm font-mono transition-colors cursor-pointer",
                      concurrency === preset
                        ? "bg-primary text-primary-foreground"
                        : "bg-surface-container hover:bg-surface-container-high text-on-surface"
                    )}
                  >
                    {preset}c
                  </button>
                ))}
              </div>
              <div className="mt-2.5 flex items-center gap-2">
                <Input
                  id="concurrency"
                  type="number"
                  min="1"
                  max="1000"
                  value={concurrency}
                  onChange={(e) => setConcurrency(Math.max(1, parseInt(e.target.value) || 1))}
                  className="w-24"
                  aria-label="Custom concurrency"
                  disabled={systemConfig.isLoading}
                />
                {errors.concurrency && (
                  <p className="font-label-caps text-critical">{errors.concurrency}</p>
                )}
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={!isValid}>
              Start Benchmark
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}