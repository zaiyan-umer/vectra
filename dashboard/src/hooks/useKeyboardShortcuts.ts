import { useEffect } from "react"

interface KeyboardShortcutsProps {
  onOpenModal?: () => void
  onCloseModal?: () => void
  onCloseDrawer?: () => void
}

export function useKeyboardShortcuts({
  onOpenModal,
  onCloseModal,
  onCloseDrawer,
}: KeyboardShortcutsProps) {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const isMeta = e.metaKey || e.ctrlKey

      if (isMeta && e.key.toLowerCase() === "k") {
        e.preventDefault()
        onOpenModal?.()
      }

      if (e.key === "Escape") {
        if (onCloseModal) onCloseModal()
        if (onCloseDrawer) onCloseDrawer()
      }
    }

    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [onOpenModal, onCloseModal, onCloseDrawer])
}