import { cn } from "@/lib/utils"
import { Info } from "lucide-react"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import { Badge } from "@/components/ui/badge"

interface AlgorithmSelectorProps {
  algorithm?: string
  className?: string
}

export function AlgorithmSelector({ algorithm = "adaptive", className }: AlgorithmSelectorProps) {
  const labels: Record<string, string> = {
    adaptive: "Adaptive Selection",
    round_robin: "Round Robin",
    least_conn: "Least Connections",
    weighted: "Weighted Response",
  }

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Badge variant="secondary" className={cn("gap-1.5", className)}>
            {labels[algorithm] || algorithm}
            <Info className="h-3.5 w-3.5 shrink-0 opacity-70" />
          </Badge>
        </TooltipTrigger>
        <TooltipContent side="top" align="center">
          <p className="font-body-sm">Automatically selects the best algorithm based on real-time backend health and latency metrics.</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}