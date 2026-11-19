import { cn } from "@/lib/utils"
import { TargetCard } from "./TargetCard"
import type { BackendTarget } from "@/types/backend"

interface TargetsGridProps {
  targets: BackendTarget[]
  className?: string
}

export function TargetsGrid({ targets, className }: TargetsGridProps) {
  return (
    <div className={cn("grid grid-cols-1 md:grid-cols-3 gap-space-md", className)}>
      {targets.map((target) => (
        <TargetCard key={target.name} target={target} />
      ))}
    </div>
  )
}