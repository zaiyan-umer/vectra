import { useState, useMemo } from "react"
import { useQueryClient } from "@tanstack/react-query"
import { useKeyboardShortcuts } from "@/hooks/useKeyboardShortcuts"
import { useBenchmarkModal } from "@/hooks/useBenchmarkModal"
import { useStartBenchmark } from "@/hooks/useStartBenchmark"
import { useBenchmarkDetail } from "@/hooks/useBenchmarkDetail"
import { useBenchmarksList } from "@/hooks/useBenchmarksList"
import { useMetrics } from "@/hooks/useMetrics"
import { useConfig } from "@/hooks/useConfig"
import { useBackends } from "@/hooks/useBackends"
import { usePoolHealth } from "@/hooks/usePoolHealth"
import { MainLayout } from "@/components/layout/MainLayout"
import { Header } from "@/components/layout/Header"
import { AlgorithmSelector } from "@/components/benchmark/AlgorithmSelector"
import { PoolHealthIndicator } from "@/components/benchmark/PoolHealthIndicator"
import { NewBenchmarkTrigger } from "@/components/benchmark/NewBenchmarkTrigger"
import { TargetsGrid } from "@/components/benchmark/TargetsGrid"
import { MetricCard } from "@/components/benchmark/MetricCard"
import { AlgorithmTabs } from "@/components/benchmark/AlgorithmTabs"
import { RunsTable } from "@/components/benchmark/RunsTable"
import { DetailDrawer } from "@/components/benchmark/DetailDrawer"
import { BenchmarkModal } from "@/components/modal/BenchmarkModal"
import type { Algorithm } from "@/types/benchmark"
import type { BackendDetail } from "@/types/api"
import { formatRPS, formatLatency, formatErrorRate, formatRelativeTime } from "@/utils/format"
import { CHART_COLORS } from "@/utils/constants"
import { toast } from "sonner"

export function DashboardPage() {
  const [selectedRunId, setSelectedRunId] = useState<string | undefined>(undefined)
  const [activeAlgorithm, setActiveAlgorithm] = useState<Algorithm | "all">("all")
  const [currentPage, setCurrentPage] = useState(1)
  const [startedRunId, setStartedRunId] = useState<string | undefined>(undefined)

  const queryClient = useQueryClient()

  const { isOpen, open, close } = useBenchmarkModal()
  const startBenchmark = useStartBenchmark()
  const { data: benchmarkDetail, isLoading: isLoadingDetail } = useBenchmarkDetail(startedRunId)
  const { data: selectedRunDetail, isLoading: isLoadingSelectedDetail } = useBenchmarkDetail(selectedRunId, {
    enabled: !!selectedRunId && selectedRunId !== startedRunId,
  })
  const { data: runsData = [], isLoading: isLoadingRuns } = useBenchmarksList()
  const { data: metrics } = useMetrics()
  const { data: config } = useConfig()
  const { data: backendsData } = useBackends()
  const { data: healthData } = usePoolHealth()

  const filteredRuns = useMemo(() => {
    let runs = [...runsData]
    if (activeAlgorithm !== "all") {
      runs = runs.filter((r) => r.algorithm === activeAlgorithm)
    }
    return runs
  }, [runsData, activeAlgorithm])

  const totalPages = Math.max(1, Math.ceil(filteredRuns.length / 10))

  const selectedRun = benchmarkDetail ?? selectedRunDetail ?? runsData.find((r) => r.id === selectedRunId)

  const handleRunSelect = (id: string) => {
    setSelectedRunId(id)
    setCurrentPage(1)
  }

  const handleAlgorithmChange = (algo: Algorithm | "all") => {
    setActiveAlgorithm(algo)
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
        queryClient.invalidateQueries({ queryKey: ['benchmarks'] })
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

  const metricCards = useMemo(() => {
    const pool = metrics?.pool
    const sparkline = metrics?.sparkline

    if (!pool || !sparkline) {
      return [
        { id: "throughput", label: "Throughput", value: "\u2014", unit: "req/s", sparklineData: [], sparklineColor: CHART_COLORS.primary, footerLeft: "Live", footerRight: "5s" },
        { id: "latency", label: "Avg Latency", value: "\u2014", unit: "ms", sparklineData: [], sparklineColor: CHART_COLORS.healthy, footerLeft: "Live", footerRight: "5s" },
        { id: "errors", label: "Error Rate", value: "\u2014", unit: "%", sparklineData: [], sparklineColor: CHART_COLORS.critical, footerLeft: "Live", footerRight: "5s" },
      ]
    }

    return [
      { id: "throughput", label: "Throughput", value: formatRPS(pool.rps), unit: "req/s", sparklineData: sparkline.poolRps, sparklineColor: CHART_COLORS.primary, footerLeft: "Live", footerRight: "5s" },
      { id: "latency", label: "Avg Latency", value: formatLatency(pool.avgLatencyMs), unit: "ms", sparklineData: sparkline.poolAvgLatencyMs, sparklineColor: CHART_COLORS.healthy, footerLeft: "Live", footerRight: "5s" },
      { id: "errors", label: "Error Rate", value: `${formatErrorRate(pool.errorRate)}%`, unit: "", sparklineData: sparkline.poolErrorRate, sparklineColor: CHART_COLORS.critical, footerLeft: "Live", footerRight: "5s" },
    ]
  }, [metrics])

  const targets = useMemo((): BackendDetail[] => {
    return backendsData?.backends ?? []
  }, [backendsData])

  const poolHealthy = healthData?.healthyCount ?? 0
  const poolTotal = healthData?.totalCount ?? 0
  const poolLastCheck = healthData?.backends?.[0]?.lastCheck
    ? formatRelativeTime(healthData.backends[0].lastCheck)
    : "\u2014"

  const currentAlgorithm = config?.algorithm ?? "adaptive"

  const isSystemConfigLoading = config === undefined || healthData === undefined

  const systemConfig = {
    algorithm: config?.algorithm ?? "adaptive_selection",
    targets: config?.targets ?? [],
    healthyCount: healthData?.healthyCount ?? 0,
    totalCount: healthData?.totalCount ?? 0,
    isLoading: isSystemConfigLoading,
  }

  return (
    <MainLayout>
      <Header />
      <main className="mx-margin-desktop px-gutter-desktop py-space-lg">
        <section className="mb-space-lg surface-card p-space-lg">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-space-md mb-space-lg">
            <div className="flex items-center gap-space-md">
              <AlgorithmSelector algorithm={currentAlgorithm} />
              <PoolHealthIndicator healthy={poolHealthy} total={poolTotal} lastCheck={poolLastCheck} />
            </div>
            <NewBenchmarkTrigger onClick={open} />
          </div>
          <TargetsGrid targets={targets} />
        </section>

        <section className="mb-space-lg grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-space-md">
          {metricCards.map((card) => (
            <MetricCard key={card.id} data={card} />
          ))}
        </section>

        <section className="grid grid-cols-1 xl:grid-cols-12 gap-space-lg">
          <div className="xl:col-span-5 space-y-space-md">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-space-md">
              <AlgorithmTabs
                activeTab={activeAlgorithm}
                onTabChange={handleAlgorithmChange}
                counts={algorithmCounts}
              />
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
          <div className="xl:col-span-7">
            <DetailDrawer
              run={selectedRun}
              isOpen={!!selectedRun}
              isMobile={false}
              isLoading={isLoadingDetail || isLoadingSelectedDetail}
            />
          </div>
        </section>
      </main>
      <BenchmarkModal
        isOpen={isOpen}
        onClose={close}
        onStart={handleStartBenchmark}
        systemConfig={systemConfig}
      />
    </MainLayout>
  )
}

export default DashboardPage
