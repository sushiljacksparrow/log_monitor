import { useEffect, useRef, useState } from 'react'

import type { LiveLogEntry } from '../../types/logs'
import { LogRow } from './LogRow'

type LogListProps = {
  logs: LiveLogEntry[]
  keyword: string
}

export function LogList({ logs, keyword }: LogListProps) {
  const containerRef = useRef<HTMLDivElement | null>(null)
  const [autoScroll, setAutoScroll] = useState(true)

  useEffect(() => {
    if (!autoScroll || !containerRef.current) {
      return
    }

    containerRef.current.scrollTop = containerRef.current.scrollHeight
  }, [autoScroll, logs])

  const handleScroll = () => {
    const node = containerRef.current
    if (!node) {
      return
    }

    const threshold = 48
    const isAtBottom = node.scrollHeight - node.scrollTop - node.clientHeight <= threshold
    setAutoScroll(isAtBottom)
  }

  const jumpToLatest = () => {
    if (!containerRef.current) {
      return
    }

    containerRef.current.scrollTop = containerRef.current.scrollHeight
    setAutoScroll(true)
  }

  return (
    <div className="relative rounded-2xl border border-border bg-surface">
      <div ref={containerRef} onScroll={handleScroll} className="max-h-[58vh] space-y-3 overflow-y-auto p-3">
        {logs.length === 0 ? (
          <div className="flex min-h-40 items-center justify-center rounded-xl border border-dashed border-border text-sm text-text-muted">
            No live logs match the current filters
          </div>
        ) : (
          logs.map((log) => <LogRow key={log.clientId} log={log} keyword={keyword} />)
        )}
      </div>

      {!autoScroll && logs.length > 0 ? (
        <button
          type="button"
          onClick={jumpToLatest}
          className="absolute bottom-4 right-4 rounded-full border border-[#58a6ff]/50 bg-[#58a6ff]/12 px-4 py-2 font-sans text-sm font-medium text-[#58a6ff] shadow-glow transition hover:bg-[#58a6ff]/18"
        >
          ↓ Jump to latest
        </button>
      ) : null}
    </div>
  )
}
