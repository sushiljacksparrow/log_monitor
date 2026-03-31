import { LEVEL_COLORS, LEVEL_OPTIONS, SERVICE_COLORS, SERVICE_OPTIONS } from '../../constants/schema'
import type { LogLevel, ServiceName } from '../../types/logs'

type FilterBarProps = {
  total: number
  showing: number
  activeLevels: LogLevel[]
  activeServices: ServiceName[]
  onToggleLevel: (level: LogLevel) => void
  onToggleService: (service: ServiceName) => void
}

export function FilterBar({
  total,
  showing,
  activeLevels,
  activeServices,
  onToggleLevel,
  onToggleService,
}: FilterBarProps) {
  return (
    <div className="flex flex-col gap-4 rounded-2xl border border-border bg-surface px-4 py-4">
      <p className="font-sans text-sm text-text-muted">
        Total: <span className="text-text-primary">{total}</span> Showing:{' '}
        <span className="text-text-primary">{showing}</span>
      </p>

      <div className="flex flex-wrap gap-2">
        {LEVEL_OPTIONS.map((level) => {
          const active = activeLevels.includes(level)
          return (
            <button
              key={level}
              type="button"
              onClick={() => onToggleLevel(level)}
              className="rounded-full border px-3 py-1.5 font-sans text-xs font-semibold transition"
              style={{
                borderColor: LEVEL_COLORS[level],
                backgroundColor: active ? `${LEVEL_COLORS[level]}22` : 'transparent',
                color: LEVEL_COLORS[level],
              }}
            >
              {level}
            </button>
          )
        })}
      </div>

      <div className="flex flex-wrap gap-2">
        {SERVICE_OPTIONS.map((service) => {
          const active = activeServices.includes(service)
          return (
            <button
              key={service}
              type="button"
              onClick={() => onToggleService(service)}
              className="rounded-full border px-3 py-1.5 font-sans text-xs font-semibold transition"
              style={{
                borderColor: SERVICE_COLORS[service],
                backgroundColor: active ? `${SERVICE_COLORS[service]}20` : 'transparent',
                color: SERVICE_COLORS[service],
              }}
            >
              {service}
            </button>
          )
        })}
      </div>
    </div>
  )
}
