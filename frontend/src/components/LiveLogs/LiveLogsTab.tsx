import { useEffect } from 'react'

import { useLogFilter } from '../../hooks/useLogFilter'
import { useWebSocket } from '../../hooks/useWebSocket'
import { FilterBar } from './FilterBar'
import { LogList } from './LogList'

type LiveLogsTabProps = {
  onConnectionChange: (connected: boolean) => void
}

export function LiveLogsTab({ onConnectionChange }: LiveLogsTabProps) {
  const { logs, isConnected, isReconnecting } = useWebSocket()
  const { keyword, setKeyword, activeLevels, activeServices, filteredLogs, toggleLevel, toggleService } =
    useLogFilter(logs)

  useEffect(() => {
    onConnectionChange(isConnected)
  }, [isConnected, onConnectionChange])

  return (
    <section className="space-y-4">
      {!isConnected && isReconnecting ? (
        <div className="rounded-2xl border border-level-ERROR/40 bg-level-ERROR/10 px-4 py-3 font-sans text-sm text-level-ERROR">
          WebSocket disconnected — retrying in 3s...
        </div>
      ) : null}

      <div className="rounded-2xl border border-border bg-surface px-4 py-4">
        <input
          value={keyword}
          onChange={(event) => setKeyword(event.target.value)}
          placeholder="Search logs by keyword, service, ID..."
          className="w-full rounded-xl border border-border bg-[#1f2937] px-4 py-3 font-sans text-sm text-text-primary outline-none transition placeholder:text-text-muted focus:border-[#58a6ff]/60 focus:shadow-glow"
        />
      </div>

      <FilterBar
        total={logs.length}
        showing={filteredLogs.length}
        activeLevels={activeLevels}
        activeServices={activeServices}
        onToggleLevel={toggleLevel}
        onToggleService={toggleService}
      />

      <LogList logs={filteredLogs} keyword={keyword} />
    </section>
  )
}
