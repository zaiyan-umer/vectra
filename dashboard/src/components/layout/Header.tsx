import { cn } from "@/lib/utils"

interface HeaderProps {
  className?: string
}

export function Header({ className }: HeaderProps) {
  return (
    <header className={cn("border-b border-border/50", className)}>
      <div className="mx-margin-desktop px-gutter-desktop py-space-md">
        <div className="surface-toolbar p-space-md">
          <div className="flex items-center gap-space-sm text-body-sm text-on-surface-variant">
            <span className="font-label-caps">Alert</span>
            <span>No active alerts</span>
          </div>
        </div>
      </div>
    </header>
  )
}