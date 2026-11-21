import { cn } from "@/lib/utils"

function Skeleton({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn("animate-pulse bg-surface-container-high rounded", className)}
      {...props}
    />
  )
}

export { Skeleton }
