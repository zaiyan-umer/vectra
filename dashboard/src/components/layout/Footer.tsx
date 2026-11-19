import { cn } from "@/lib/utils"
import { Cpu, Server, Shield } from "lucide-react"

interface FooterProps {
  className?: string
}

export function Footer({ className }: FooterProps) {
  return (
    <footer className={cn("border-t border-border/50 bg-surface-container-low", className)}>
      <div className="mx-margin-desktop px-gutter-desktop py-space-lg">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-space-lg text-center md:text-left">
          <div className="flex items-center justify-center md:justify-start gap-space-sm">
            <Cpu className="h-5 w-5 text-primary" />
            <div>
              <p className="font-label-caps text-primary">Vectra Core Engine</p>
              <p className="font-code-sm text-on-surface-variant">v2.1.0</p>
            </div>
          </div>
          <div className="flex items-center justify-center md:justify-center gap-space-sm">
            <Server className="h-5 w-5 text-secondary" />
            <div>
              <p className="font-label-caps text-secondary">Target Host</p>
              <p className="font-code-sm text-on-surface-variant">vectra-loadbalancer:8080</p>
            </div>
          </div>
          <div className="flex items-center justify-center md:justify-end gap-space-sm">
            <Shield className="h-5 w-5 text-tertiary" />
            <div>
              <p className="font-label-caps text-tertiary">Status</p>
              <p className="font-code-sm text-healthy">All Systems Operational</p>
            </div>
          </div>
        </div>
        <div className="mt-space-lg pt-space-lg border-t border-border/50 text-center">
          <p className="font-body-sm text-on-surface-variant">
            © 2025 Vectra. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  )
}