import { cn } from "@/lib/utils"
import { Command, Plus } from "lucide-react"
import { Button } from "@/components/ui/button"

interface NewBenchmarkTriggerProps {
  onClick: () => void
  className?: string
}

export function NewBenchmarkTrigger({ onClick, className }: NewBenchmarkTriggerProps) {
  return (
    <Button
      onClick={onClick}
      className={cn("cursor-pointer gap-2", className)}
    >
      <Plus className="h-4 w-4" />
      <span>New Benchmark</span>
      <kbd
        aria-label="Keyboard shortcut: Command K"
        className="ml-1 inline-flex items-center gap-0.5 rounded-md border border-primary-foreground/25 bg-primary-foreground/10 px-1.5 py-0.5 font-mono text-sm font-semibold text-primary-foreground shadow-sm"
      >
        <Command aria-hidden="true" className="size-4" />
        <span className="leading-4">K</span>
      </kbd>
    </Button>
  )
}