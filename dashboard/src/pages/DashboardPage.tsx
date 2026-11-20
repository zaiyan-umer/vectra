import { useState, useMemo, useCallback, useEffect } from "react"
import { useKeyboardShortcuts } from "@/hooks/useKeyboardShortcuts"
import { useBenchmarkModal } from "@/hooks/useBenchmarkModal"
import { useStartBenchmark } from "@/hooks/useStartBenchmark"
import { useBenchmarkDetail } from "@/hooks/useBenchmarkDetail"
import { getBenchmarks, transformBenchmark } from "@/lib/api"
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
import type { Algorithm, BenchmarkRun } from "@/types/benchmark"
import { toast } from "sonner"

export function DashboardPage() {
  const [selectedRunId, setSelectedRunId] = useState<string | undefined>(undefined)
  const [activeAlgorithm, setActiveAlgorithm] = useState<Algorithm | "all">("all")
  const [searchQuery, setSearchQuery] = useState("")
  const [currentPage, setCurrentPage] = useState(1)
  const [runsData, setRunsData] = useState<BenchmarkRun[]>([])
  const [isLoadingRuns, setIsLoadingRuns] = useState(true)
  const [startedRunId, setStartedRunId] = useState<string | undefined>(undefined)

  const { isOpen, open, close } = useBenchmarkModal()
  const startBenchmark = useStartBenchmark()
  const { data: benchmarkDetail, isLoading: isLoadingDetail } = useBenchmarkDetail(startedRunId)

  const fetchRuns = useCallback(async () => {
    try {
      setIsLoadingRuns(true)
      const data = await getBenchmarks()
      setRunsData(data.map(transformBenchmark))
    } catch (err) {
      console.error("Failed to fetch runs:", err)
      toast.error("Failed to load benchmark runs")
    } finally {
      setIsLoadingRuns(false)
    }
  }, [])

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchRuns()
  }, [fetchRuns])

  const filteredRuns = useMemo(() => {
    let runs = [...runsData]
    if (activeAlgorithm !== "all") {
      runs = runs.filter((r) => r.algorithm === activeAlgorithm)
    }
    if (searchQuery) {
      const q = searchQuery.toLowerCase()
      runs = runs.filter(
        (r) =>
          r.id.toLowerCase().includes(q) ||
          r.algorithm.toLowerCase().includes(q) ||
          (r.status && r.status.toLowerCase().includes(q))
      )
    }
    return runs
  }, [runsData, activeAlgorithm, searchQuery])

  const totalPages = Math.max(1, Math.ceil(filteredRuns.length / 10))

  const selectedRun = benchmarkDetail ?? runsData.find((r) => r.id === selectedRunId)

  const handleRunSelect = (id: string) => {
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

  const handleStartBenchmark = (config: { durationSeconds: number; connections: number }) => {
    startBenchmark.mutate(config, {
      onSuccess: (response) => {
        setStartedRunId(response.id)
        setSelectedRunId(response.id)
      },
      onError: (err) => {
        toast.error(`Failed to start benchmark: ${err.message}`)
      },
    })
  }

  useKeyboardShortcuts({
    onOpenModal: open,
    onCloseModal: close,
    onCloseDrawer: () => setSelectedRunId(undefined),
  })

  const algorithmCounts = useMemo(() => {
    const counts: Record<Algorithm | "all", number> = { all: runsData.length, adaptive: 0, round_robin: 0, least_conn: 0, weighted: 0 }
    runsData.forEach((r) => { counts[r.algorithm]++ })
    return counts
  }, [runsData])

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
          <TargetsGrid targets={[]} />
        </section>

        <section className="mb-space-lg grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-space-md">
          <MetricCard
            data={{
              id: "throughput",
              label: "Throughput",
              value: "—",
              unit: "req/s",
              sparklineData: [],
              sparklineColor: "#2563eb",
              footerLeft: "Live",
              footerRight: "5s",
            }}
          />
          <MetricCard
            data={{
              id: "latency",
              label: "Avg Latency",
              value: "—",
              unit: "ms",
              sparklineData: [],
              sparklineColor: "#10b981",
              footerLeft: "Live",
              footerRight: "5s",
            }}
          />
          <MetricCard
            data={{
              id: "p99",
              label: "P99 Latency",
              value: "—",
              unit: "ms",
              sparklineData: [],
              sparklineColor: "#f59e0b",
              footerLeft: "Live",
              footerRight: "5s",
            }}
          />
          <MetricCard
            data={{
              id: "errors",
              label: "Error Rate",
              value: "—",
              unit: "%",
              sparklineData: [],
              sparklineColor: "#ef4444",
              footerLeft: "Live",
              footerRight: "5s",
            }}
          />
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
              isLoading={isLoadingRuns}
            />
          </div>
          <div className="xl:col-span-5">
            <DetailDrawer
              run={selectedRun}
              isOpen={!!selectedRun}
              onClose={() => setSelectedRunId(undefined)}
              isMobile={false}
              isLoading={isLoadingDetail}
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