import { useState, type FormEvent } from "react"
import { cn } from "@/lib/utils"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { DURATION_PRESETS, CONCURRENCY_PRESETS } from "@/utils/constants"

interface BenchmarkModalProps {
  isOpen: boolean
  onClose: () => void
  onStart: (config: {
    durationSeconds: number
    connections: number
  }) => void
}

export function BenchmarkModal({ isOpen, onClose, onStart }: BenchmarkModalProps) {
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

  const isValid = duration > 0 && concurrency > 0

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-lg">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>New Benchmark</DialogTitle>
            <DialogDescription>
              Configure and start a new load test benchmark
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-space-lg">
            <div className="surface-toolbar p-space-md rounded-lg">
              <p className="font-label-caps text-on-surface-variant mb-space-xs">System State</p>
              <div className="grid grid-cols-3 gap-space-sm text-sm">
                <div>
                  <p className="font-label-caps text-on-surface-variant">Targets</p>
                  <p className="font-code-sm text-on-surface">3 healthy</p>
                </div>
                <div>
                  <p className="font-label-caps text-on-surface-variant">Algorithm</p>
                  <p className="font-code-sm text-on-surface">Adaptive</p>
                </div>
                <div>
                  <p className="font-label-caps text-on-surface-variant">Default Config</p>
                  <p className="font-code-sm text-on-surface">30s • 100c</p>
                </div>
              </div>
            </div>

            <div className="space-y-space-md">
              <div>
                <Label htmlFor="duration" className="font-label-caps">Duration (seconds)</Label>
                <div className="flex flex-wrap gap-1.5 mt-1.5">
                  {DURATION_PRESETS.map((preset) => (
                    <button
                      key={preset}
                      type="button"
                      onClick={() => setDuration(preset)}
                      className={cn(
                        "px-3 py-1.5 rounded-lg text-sm font-mono transition-colors",
                        duration === preset
                          ? "bg-primary text-primary-foreground"
                          : "bg-surface-container hover:bg-surface-container-high text-on-surface"
                      )}
                    >
                      {preset}s
                    </button>
                  ))}
                </div>
                <div className="mt-1.5 flex items-center gap-2">
                  <Input
                    id="duration"
                    type="number"
                    min="1"
                    max="3600"
                    value={duration}
                    onChange={(e) => setDuration(Math.max(1, parseInt(e.target.value) || 1))}
                    className="w-24"
                    aria-label="Custom duration"
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
                        "px-3 py-1.5 rounded-lg text-sm font-mono transition-colors",
                        concurrency === preset
                          ? "bg-primary text-primary-foreground"
                          : "bg-surface-container hover:bg-surface-container-high text-on-surface"
                      )}
                    >
                      {preset}c
                    </button>
                  ))}
                </div>
                <div className="mt-1.5 flex items-center gap-2">
                  <Input
                    id="concurrency"
                    type="number"
                    min="1"
                    max="1000"
                    value={concurrency}
                    onChange={(e) => setConcurrency(Math.max(1, parseInt(e.target.value) || 1))}
                    className="w-24"
                    aria-label="Custom concurrency"
                  />
                  {errors.concurrency && (
                    <p className="font-label-caps text-critical">{errors.concurrency}</p>
                  )}
                </div>
              </div>

              </div>
          </div>

          <DialogFooter className="flex justify-end gap-2">
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