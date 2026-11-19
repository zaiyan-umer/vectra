import { cn } from "@/lib/utils"

interface MainLayoutProps {
  children: React.ReactNode
}

export function MainLayout({ children }: MainLayoutProps) {
  return (
    <div className={cn("min-h-screen bg-background font-inter")}>
      {children}
    </div>
  )
}