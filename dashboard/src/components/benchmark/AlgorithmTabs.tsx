import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import type { Algorithm } from "@/types/benchmark"

interface AlgorithmTabsProps {
  activeTab: Algorithm | "all"
  onTabChange: (tab: Algorithm | "all") => void
  counts?: Record<Algorithm | "all", number>
  className?: string
}

const TABS: { id: Algorithm | "all"; label: string }[] = [
  { id: "all", label: "All Runs" },
  { id: "adaptive", label: "Adaptive" },
  { id: "round_robin", label: "Round Robin" },
  { id: "least_conn", label: "Least Conn" },
  { id: "weighted", label: "Weighted" },
]

function CountBadge({ count, isActive }: { count: number; isActive: boolean }) {
  return (
    <span className={cn(
      "ml-1.5 px-1.5 py-0.5 text-xs font-mono rounded-full",
      isActive
        ? "bg-primary-foreground/20 text-primary-foreground"
        : "bg-muted text-muted-foreground"
    )}>
      {count}
    </span>
  )
}

export function AlgorithmTabs({ activeTab, onTabChange, counts, className }: AlgorithmTabsProps) {
  return (
    <div className={cn("flex flex-wrap gap-1.5", className)} role="tablist">
      {TABS.map((tab) => {
        const isActive = activeTab === tab.id
        return (
          <Button
            key={tab.id}
            role="tab"
            aria-selected={isActive}
            variant={isActive ? "default" : "outline"}
            size="sm"
            onClick={() => onTabChange(tab.id)}
            className={cn("cursor-pointer h-7 px-3 text-body-sm font-medium", isActive && "text-white")}
          >
            {tab.label}
            {counts && counts[tab.id] !== undefined && (
              <CountBadge count={counts[tab.id]} isActive={isActive} />
            )}
          </Button>
        )
      })}
    </div>
  )
}