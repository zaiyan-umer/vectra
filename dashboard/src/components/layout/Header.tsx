import { cn } from "@/lib/utils"

interface HeaderProps {
  className?: string
}

export function Header({ className }: HeaderProps) {
  return (
    <header className={cn("border-b border-border/50", className)}>
      <div className="mx-margin-desktop px-gutter-desktop py-space-md">
        <div className="surface-toolbar p-space-md flex items-center justify-between">
          <span className="font-headline-md text-on-surface">Vectra</span>
          <span className="inline-flex items-center gap-1.5 rounded-full bg-healthy-bg border border-healthy-border px-2 py-0.5">
            <span className="relative flex size-1.5 rounded-full bg-healthy animate-pulse" />
            <span className="font-label-caps text-healthy-text">Engine Active</span>
          </span>
        </div>
      </div>
    </header>
  )
}