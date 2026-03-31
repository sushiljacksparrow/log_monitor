import { useMemo, useState } from 'react'

import { LEVEL_COLORS, SERVICE_COLORS } from '../../constants/schema'
import type { LiveLogEntry } from '../../types/logs'

type LogRowProps = {
  log: LiveLogEntry
  keyword: string
}

const PRIMARY_FIELDS = new Set(['clientId', 'receivedAt', 'timestamp', 'level', 'service', 'message'])

function highlightText(text: string, query: string) {
  const trimmed = query.trim()
  if (!trimmed) {
    return text
  }

  const escaped = trimmed.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const matcher = new RegExp(`(${escaped})`, 'ig')
  const parts = text.split(matcher)

  return parts.map((part, index) =>
    part.toLowerCase() === trimmed.toLowerCase() ? (
      <mark key={`${part}-${index}`} className="rounded bg-level-DEBUG/20 px-0.5 text-level-DEBUG">
        {part}
      </mark>
    ) : (
      <span key={`${part}-${index}`}>{part}</span>
    ),
  )
}

export function LogRow({ log, keyword }: LogRowProps) {
  const [expanded, setExpanded] = useState(false)
  const detailEntries = useMemo(
    () => Object.entries(log).filter(([key]) => !['clientId', 'receivedAt'].includes(key)),
    [log],
  )
  const inlineEntries = useMemo(
    () => Object.entries(log).filter(([key]) => !PRIMARY_FIELDS.has(key)),
    [log],
  )

  return (
    <div className="rounded-xl border border-transparent bg-surface/55 transition hover:border-border hover:bg-[#1b2330]">
      <button
        type="button"
        onClick={() => setExpanded((current) => !current)}
        className="grid w-full grid-cols-1 gap-3 px-4 py-3 text-left md:grid-cols-[170px_68px_132px_minmax(0,1fr)] md:items-start"
      >
        <span className="font-mono text-xs text-text-muted">{log.timestamp}</span>
        <span
          className="inline-flex w-fit rounded-full border px-2 py-1 font-sans text-[11px] font-bold"
          style={{
            borderColor: LEVEL_COLORS[log.level],
            color: LEVEL_COLORS[log.level],
            backgroundColor: `${LEVEL_COLORS[log.level]}18`,
          }}
        >
          {log.level}
        </span>
        <span
          className="inline-flex w-fit rounded-full border px-2 py-1 font-sans text-[11px] font-semibold"
          style={{
            borderColor: SERVICE_COLORS[log.service],
            color: SERVICE_COLORS[log.service],
            backgroundColor: `${SERVICE_COLORS[log.service]}18`,
          }}
        >
          {log.service}
        </span>
        <div className="min-w-0 space-y-3">
          <div className="break-words font-mono text-sm leading-6 text-text-primary">
            {highlightText(log.message, keyword)}
          </div>
          {inlineEntries.length > 0 ? (
            <div className="flex flex-wrap gap-2">
              {inlineEntries.map(([key, value]) => (
                <span
                  key={key}
                  className="inline-flex max-w-full items-center gap-2 rounded-full border px-2.5 py-1 font-mono text-xs"
                  style={{
                    borderColor: `${SERVICE_COLORS[log.service]}44`,
                    backgroundColor: `${SERVICE_COLORS[log.service]}12`,
                    color: '#e6edf3',
                  }}
                >
                  <span className="text-text-muted">{key}</span>
                  <span className="truncate">{String(value)}</span>
                </span>
              ))}
            </div>
          ) : null}
        </div>
      </button>

      <div
        className={`grid overflow-hidden transition-all duration-150 ease-in-out ${
          expanded ? 'grid-rows-[1fr] border-t border-border' : 'grid-rows-[0fr]'
        }`}
      >
        <div className="min-h-0">
          <div className="grid gap-3 px-4 py-4 sm:grid-cols-2 xl:grid-cols-3">
            {detailEntries.map(([key, value]) => (
              <div
                key={key}
                className="rounded-xl border p-3"
                style={{
                  borderColor: `${SERVICE_COLORS[log.service]}44`,
                  backgroundColor: `${SERVICE_COLORS[log.service]}10`,
                }}
              >
                <p className="mb-1 font-sans text-xs uppercase tracking-[0.18em] text-text-muted">{key}</p>
                <p className="break-words font-mono text-sm text-text-primary">{String(value)}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
