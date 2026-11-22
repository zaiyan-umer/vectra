import { cn } from "@/lib/utils"
import { Info } from "lucide-react"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import { Badge } from "@/components/ui/badge"
import { API_ALGORITHM_LABELS } from "@/utils/constants"

interface AlgorithmSelectorProps {
  algorithm?: string
  className?: string
}

export function AlgorithmSelector({ algorithm = "adaptive_selection", className }: AlgorithmSelectorProps) {
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Badge variant="secondary" className={cn("gap-1.5", className)}>
            {API_ALGORITHM_LABELS[algorithm] || algorithm}
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