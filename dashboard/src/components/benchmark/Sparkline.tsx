import { cn } from "@/lib/utils"
import { generateSparklinePath } from "@/utils/sparkline"

interface SparklineProps {
  data: number[]
  color: string
  className?: string
  width?: number
  height?: number
}

export function Sparkline({ data, color, className, width = 100, height = 24 }: SparklineProps) {
  const path = generateSparklinePath(data, width, height)

  return (
    <svg
      className={cn("w-full h-auto", className)}
      viewBox={`0 0 ${width} ${height}`}
      preserveAspectRatio="none"
      aria-hidden="true"
    >
      <defs>
        <linearGradient id="sparkline-gradient" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor={color} stopOpacity={0.3} />
          <stop offset="100%" stopColor={color} stopOpacity={0} />
        </linearGradient>
      </defs>
      <path
        d={path}
        stroke={color}
        strokeWidth="1.5"
        fill="none"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d={`${path} L${width},${height} L0,${height} Z`}
        fill="url(#sparkline-gradient)"
      />
    </svg>
  )
}