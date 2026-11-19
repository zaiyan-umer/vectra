import { useState, type FormEvent } from "react"
import { cn } from "@/lib/utils"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Slider } from "@/components/ui/slider"
import { Label } from "@/components/ui/label"
import { DURATION_PRESETS, CONCURRENCY_PRESETS, PIPELINING_MIN, PIPELINING_MAX, PIPELINING_DEFAULT } from "@/utils/constants"

interface BenchmarkModalProps {
  isOpen: boolean
  onClose: () => void
  onStart: (config: {
    duration: number
    concurrency: number
    pipelining: number
  }) => void
}

export function BenchmarkModal({ isOpen, onClose, onStart }: BenchmarkModalProps) {
  const [duration, setDuration] = useState(30)
  const [concurrency, setConcurrency] = useState(100)
  const [pipelining, setPipelining] = useState(PIPELINING_DEFAULT)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validate = () => {
    const newErrors: Record<string, string> = {}
    if (!duration || duration < 1) newErrors.duration = "Duration must be at least 1 second"
    if (!concurrency || concurrency < 1) newErrors.concurrency = "Concurrency must be at least 1"
    if (pipelining < PIPELINING_MIN || pipelining > PIPELINING_MAX) {
      newErrors.pipelining = `Pipelining must be between ${PIPELINING_MIN} and ${PIPELINING_MAX}`
    }
    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (validate()) {
      onStart({ duration, concurrency, pipelining })
      onClose()
    }
  }

  const isValid = duration > 0 && concurrency > 0 && pipelining >= PIPELINING_MIN && pipelining <= PIPELINING_MAX

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

              <div>
                <div className="flex items-center justify-between mb-1.5">
                  <Label htmlFor="pipelining" className="font-label-caps mb-0">
                    Pipelining Factor
                  </Label>
                  <span className="font-code-sm text-on-surface-variant">{pipelining}</span>
                </div>
                <Slider
                  id="pipelining"
                  min={PIPELINING_MIN}
                  max={PIPELINING_MAX}
                  step={1}
                  value={[pipelining]}
                  onValueChange={([val]) => setPipelining(val)}
                  className="w-full"
                />
                <div className="flex justify-between text-xs text-on-surface-variant mt-1">
                  <span>{PIPELINING_MIN}</span>
                  <span>{PIPELINING_MAX}</span>
                </div>
                {errors.pipelining && (
                  <p className="font-label-caps text-critical mt-1">{errors.pipelining}</p>
                )}
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