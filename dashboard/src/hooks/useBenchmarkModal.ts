import { create } from "zustand"

interface BenchmarkModalState {
  isOpen: boolean
  open: () => void
  close: () => void
  formData: {
    duration: number
    concurrency: number
  }
  setFormData: (data: Partial<BenchmarkModalState["formData"]>) => void
  resetForm: () => void
}

const DEFAULT_FORM_DATA = {
  duration: 30,
  concurrency: 100,
}

export const useBenchmarkModal = create<BenchmarkModalState>((set) => ({
  isOpen: false,
  open: () => set({ isOpen: true }),
  close: () => set({ isOpen: false }),
  formData: DEFAULT_FORM_DATA,
  setFormData: (data) => set((state) => ({ formData: { ...state.formData, ...data } })),
  resetForm: () => set({ formData: DEFAULT_FORM_DATA }),
}))