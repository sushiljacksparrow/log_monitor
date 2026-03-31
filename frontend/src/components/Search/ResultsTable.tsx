import { LEVEL_COLORS, SERVICE_COLORS, TABLE_COLUMNS } from '../../constants/schema'
import type { ServiceLogMap, ServiceName } from '../../types/logs'

type ResultsTableProps<TService extends ServiceName> = {
  selectedIndex: TService | null
  results: ServiceLogMap[TService][]
  loading: boolean
}

export function ResultsTable<TService extends ServiceName>({
  selectedIndex,
  results,
  loading,
}: ResultsTableProps<TService>) {
  if (!selectedIndex) {
    return null
  }

  if (loading) {
    return (
      <div className="flex min-h-64 items-center justify-center rounded-2xl border border-border bg-surface">
        <div className="h-10 w-10 animate-spin rounded-full border-2 border-border border-t-[#58a6ff]" />
      </div>
    )
  }

  if (results.length === 0) {
    return (
      <div className="flex min-h-64 items-center justify-center rounded-2xl border border-border bg-surface font-sans text-sm text-text-muted">
        No results found
      </div>
    )
  }

  const columns = TABLE_COLUMNS[selectedIndex]

  return (
    <div className="overflow-hidden rounded-2xl border border-border bg-surface">
      <div className="overflow-x-auto">
        <table className="min-w-full border-collapse">
          <thead>
            <tr className="border-b border-border bg-bg/65">
              {columns.map((column) => (
                <th
                  key={column}
                  className="px-4 py-3 text-left font-sans text-xs uppercase tracking-[0.18em] text-text-muted"
                >
                  {column}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {results.map((row, index) => (
              <tr
                key={`${row.request_id}-${index}`}
                className="border-b border-border/60 last:border-b-0 hover:bg-[#1b2330]"
              >
                {columns.map((column) => {
                  const value = row[column as keyof typeof row]
                  if (column === 'level' && typeof value === 'string') {
                    const color = LEVEL_COLORS[value as keyof typeof LEVEL_COLORS]
                    return (
                      <td key={column} className="px-4 py-3">
                        <span
                          className="inline-flex rounded-full border px-2 py-1 font-sans text-[11px] font-bold"
                          style={{ borderColor: color, color, backgroundColor: `${color}18` }}
                        >
                          {value}
                        </span>
                      </td>
                    )
                  }

                  if (column === 'service' && typeof value === 'string') {
                    const color = SERVICE_COLORS[value as keyof typeof SERVICE_COLORS]
                    return (
                      <td key={column} className="px-4 py-3">
                        <span
                          className="inline-flex rounded-full border px-2 py-1 font-sans text-[11px] font-semibold"
                          style={{ borderColor: color, color, backgroundColor: `${color}18` }}
                        >
                          {value}
                        </span>
                      </td>
                    )
                  }

                  return (
                    <td key={column} className="max-w-[320px] px-4 py-3 font-mono text-sm text-text-primary">
                      <div className="break-words">{String(value ?? '')}</div>
                    </td>
                  )
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
