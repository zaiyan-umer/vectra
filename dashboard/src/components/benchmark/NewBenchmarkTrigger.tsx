import { cn } from "@/lib/utils"
import { Plus, Keyboard } from "lucide-react"
import { Button } from "@/components/ui/button"

interface NewBenchmarkTriggerProps {
  onClick: () => void
  className?: string
}

export function NewBenchmarkTrigger({ onClick, className }: NewBenchmarkTriggerProps) {
  return (
    <Button
      onClick={onClick}
      className={cn("gap-2", className)}
    >
      <Plus className="h-4 w-4" />
      <span>New Benchmark</span>
      <Keyboard className="h-3.5 w-3.5 opacity-60 ml-1" />
      <span className="font-label-caps opacity-60">⌘K</span>
    </Button>
  )
}