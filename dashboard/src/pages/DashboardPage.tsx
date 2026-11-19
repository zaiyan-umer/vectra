import { useState, useMemo } from "react"
import { useKeyboardShortcuts } from "@/hooks/useKeyboardShortcuts"
import { useBenchmarkModal } from "@/hooks/useBenchmarkModal"
import { MainLayout } from "@/components/layout/MainLayout"
import { Header } from "@/components/layout/Header"
import { Footer } from "@/components/layout/Footer"
import { AlgorithmSelector } from "@/components/benchmark/AlgorithmSelector"
import { PoolHealthIndicator } from "@/components/benchmark/PoolHealthIndicator"
import { NewBenchmarkTrigger } from "@/components/benchmark/NewBenchmarkTrigger"
import { TargetsGrid } from "@/components/benchmark/TargetsGrid"
import { MetricCard } from "@/components/benchmark/MetricCard"
import { AlgorithmTabs } from "@/components/benchmark/AlgorithmTabs"
import { SearchInput } from "@/components/benchmark/SearchInput"
import { RunsTable } from "@/components/benchmark/RunsTable"
import { DetailDrawer } from "@/components/benchmark/DetailDrawer"
import { BenchmarkModal } from "@/components/modal/BenchmarkModal"
import { METRIC_CARDS_DATA, MOCK_RUNS, MOCK_BACKENDS } from "@/lib/mockData"
import type { Algorithm } from "@/types/benchmark"

export function DashboardPage() {
  const [selectedRunId, setSelectedRunId] = useState<number | undefined>(MOCK_RUNS[0]?.id)
  const [activeAlgorithm, setActiveAlgorithm] = useState<Algorithm | "all">("all")
  const [searchQuery, setSearchQuery] = useState("")
  const [currentPage, setCurrentPage] = useState(1)

  const { isOpen, open, close } = useBenchmarkModal()

  const filteredRuns = useMemo(() => {
    let runs = [...MOCK_RUNS]
    if (activeAlgorithm !== "all") {
      runs = runs.filter((r) => r.algorithm === activeAlgorithm)
    }
    if (searchQuery) {
      const q = searchQuery.toLowerCase()
      runs = runs.filter(
        (r) =>
          r.id.toString().includes(q) ||
          r.algorithm.toLowerCase().includes(q) ||
          (r.status && r.status.toLowerCase().includes(q))
      )
    }
    return runs
  }, [activeAlgorithm, searchQuery])

  const totalPages = Math.max(1, Math.ceil(filteredRuns.length / 10))

  const selectedRun = MOCK_RUNS.find((r) => r.id === selectedRunId)

  const handleRunSelect = (id: number) => {
    setSelectedRunId(id)
    setCurrentPage(1)
  }

  const handleAlgorithmChange = (algo: Algorithm | "all") => {
    setActiveAlgorithm(algo)
    setCurrentPage(1)
  }

  const handleSearchChange = (value: string) => {
    setSearchQuery(value)
    setCurrentPage(1)
  }

  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page)
    }
  }

  const handleStartBenchmark = (config: { duration: number; concurrency: number; pipelining: number }) => {
    console.log("Start benchmark with config:", config)
  }

  useKeyboardShortcuts({
    onOpenModal: open,
    onCloseModal: close,
    onCloseDrawer: () => setSelectedRunId(undefined),
  })

  const algorithmCounts = useMemo(() => {
    const counts: Record<Algorithm | "all", number> = { all: MOCK_RUNS.length, adaptive: 0, round_robin: 0, least_conn: 0, weighted: 0 }
    MOCK_RUNS.forEach((r) => { counts[r.algorithm]++ })
    return counts
  }, [])

  return (
    <MainLayout>
      <Header />
      <main className="mx-margin-desktop px-gutter-desktop py-space-lg">
        <section className="mb-space-lg surface-card p-space-lg">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-space-md mb-space-lg">
            <div className="flex items-center gap-space-md">
              <AlgorithmSelector algorithm="adaptive" />
              <PoolHealthIndicator />
            </div>
            <NewBenchmarkTrigger onClick={open} />
          </div>
          <TargetsGrid targets={MOCK_BACKENDS} />
        </section>

        <section className="mb-space-lg grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-space-md">
          {METRIC_CARDS_DATA.map((data) => (
            <MetricCard key={data.id} data={data} />
          ))}
        </section>

        <section className="grid grid-cols-1 xl:grid-cols-12 gap-space-lg">
          <div className="xl:col-span-7 space-y-space-md">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-space-md">
              <AlgorithmTabs
                activeTab={activeAlgorithm}
                onTabChange={handleAlgorithmChange}
                counts={algorithmCounts}
              />
              <SearchInput value={searchQuery} onChange={handleSearchChange} />
            </div>
            <RunsTable
              runs={filteredRuns}
              selectedRunId={selectedRunId}
              onSelectRun={handleRunSelect}
              onPageChange={handlePageChange}
              currentPage={currentPage}
              totalPages={totalPages}
            />
          </div>
          <div className="xl:col-span-5">
            <DetailDrawer
              run={selectedRun}
              isOpen={!!selectedRun}
              onClose={() => setSelectedRunId(undefined)}
              isMobile={false}
            />
          </div>
        </section>
      </main>
      <Footer />
      <BenchmarkModal
        isOpen={isOpen}
        onClose={close}
        onStart={handleStartBenchmark}
      />
    </MainLayout>
  )
}

export default DashboardPage