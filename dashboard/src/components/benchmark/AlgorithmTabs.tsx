import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { API_ALGORITHM_LABELS } from "@/utils/constants"

interface AlgorithmTabsProps {
  activeTab: string | "all"
  onTabChange: (tab: string | "all") => void
  availableAlgorithms: string[]
  counts?: Record<string, number>
  className?: string
}

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

export function AlgorithmTabs({ activeTab, onTabChange, availableAlgorithms, counts, className }: AlgorithmTabsProps) {
  const tabs = [
    { id: "all", label: "All Runs" },
    ...availableAlgorithms.map((algo) => ({
      id: algo,
      label: API_ALGORITHM_LABELS[algo] ?? algo,
    })),
  ]

  return (
    <div className={cn("flex flex-wrap gap-1.5", className)} role="tablist">
      {tabs.map((tab) => {
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