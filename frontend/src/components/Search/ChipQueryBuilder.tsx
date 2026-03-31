import type { KeyboardEvent } from 'react'

import { INDEX_FIELDS, SERVICE_COLORS, SERVICE_OPTIONS } from '../../constants/schema'
import type { SearchChip, ServiceName } from '../../types/logs'

type ChipQueryBuilderProps = {
  selectedIndex: ServiceName | null
  chips: SearchChip[]
  fieldInput: string
  valueInput: string
  pendingField: string | null
  suggestions: string[]
  onSelectIndex: (service: ServiceName) => void
  onClearIndex: () => void
  onFieldInputChange: (value: string) => void
  onValueInputChange: (value: string) => void
  onSelectField: (field: string) => void
  onCompletePendingField: () => void
  onRemoveChip: (field: string) => void
  onRemoveLastItem: () => void
  onRunSearch: () => void
  canSearch: boolean
}

export function ChipQueryBuilder({
  selectedIndex,
  chips,
  fieldInput,
  valueInput,
  pendingField,
  suggestions,
  onSelectIndex,
  onClearIndex,
  onFieldInputChange,
  onValueInputChange,
  onSelectField,
  onCompletePendingField,
  onRemoveChip,
  onRemoveLastItem,
  onRunSearch,
  canSearch,
}: ChipQueryBuilderProps) {
  const handleKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Backspace') {
      if ((!pendingField && !fieldInput) || (pendingField && !valueInput)) {
        onRemoveLastItem()
      }
      return
    }

    if (event.key !== 'Enter' && event.key !== 'Tab') {
      return
    }

    if (!selectedIndex) {
      return
    }

    if (!pendingField) {
      if (suggestions.length > 0) {
        event.preventDefault()
        onSelectField(suggestions[0])
      }
      return
    }

    if (valueInput.trim()) {
      event.preventDefault()
      onCompletePendingField()
      return
    }

    if (event.key === 'Enter' && canSearch) {
      event.preventDefault()
      onRunSearch()
    }
  }

  return (
    <div className="space-y-4">
      <div className="rounded-2xl border border-border bg-surface p-4">
        <div className="flex min-h-14 flex-wrap items-center gap-2 rounded-xl border border-[#374151] bg-[#1f2937] px-3 py-2">
          {selectedIndex ? (
            <button
              type="button"
              onClick={onClearIndex}
              className="inline-flex items-center gap-2 rounded-full border px-3 py-1.5 font-sans text-sm font-medium transition hover:brightness-110"
              style={{ borderColor: SERVICE_COLORS[selectedIndex], color: SERVICE_COLORS[selectedIndex] }}
            >
              {selectedIndex}
              <span className="text-text-muted">×</span>
            </button>
          ) : null}

          {chips.map((chip) => (
            <div
              key={chip.field}
              className="inline-flex items-center gap-2 rounded-full border border-[#374151] bg-bg px-3 py-1.5 font-mono text-sm"
            >
              <span className="text-level-DEBUG">{chip.field}</span>
              <span className="text-text-muted">:</span>
              <span className="text-text-primary">{chip.value}</span>
              <button
                type="button"
                onClick={() => onRemoveChip(chip.field)}
                className="text-text-muted transition hover:text-level-ERROR"
              >
                ×
              </button>
            </div>
          ))}

          {pendingField ? (
            <div className="inline-flex items-center gap-2 rounded-full border border-[#374151] bg-bg px-3 py-1.5 font-mono text-sm">
              <span className="text-level-DEBUG">{pendingField}</span>
              <span className="text-text-muted">:</span>
              <input
                value={valueInput}
                onChange={(event) => onValueInputChange(event.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="type a value..."
                className="min-w-28 bg-transparent text-text-primary outline-none placeholder:text-text-muted"
              />
            </div>
          ) : selectedIndex ? (
            <input
              value={fieldInput}
              onChange={(event) => onFieldInputChange(event.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="+ type a field..."
              className="min-w-40 flex-1 bg-transparent font-mono text-sm text-text-primary outline-none placeholder:text-text-muted"
            />
          ) : (
            <input
              readOnly
              value=""
              placeholder="Select an index to begin..."
              className="flex-1 bg-transparent font-sans text-sm text-text-primary outline-none placeholder:text-text-muted"
            />
          )}
        </div>

        {selectedIndex ? (
          <div className="mt-3 flex items-center justify-end">
            <button
              type="button"
              onClick={onRunSearch}
              disabled={!canSearch}
              className="rounded-full border border-[#58a6ff]/40 bg-[#58a6ff]/12 px-4 py-2 font-sans text-sm font-medium text-[#58a6ff] transition hover:bg-[#58a6ff]/18 disabled:cursor-not-allowed disabled:opacity-40"
            >
              Run search
            </button>
          </div>
        ) : null}

        {selectedIndex && !pendingField && suggestions.length > 0 ? (
          <div className="mt-3 rounded-xl border border-border bg-bg p-2">
            <div className="flex flex-wrap gap-2">
              {suggestions.map((field) => (
                <button
                  key={field}
                  type="button"
                  onClick={() => onSelectField(field)}
                  className="rounded-full border border-border px-3 py-1.5 font-mono text-xs text-text-muted transition hover:border-[#58a6ff]/40 hover:bg-[#58a6ff]/10 hover:text-[#58a6ff]"
                >
                  {field}
                </button>
              ))}
            </div>
          </div>
        ) : null}
      </div>

      {!selectedIndex ? (
        <div className="flex flex-wrap gap-3">
          {SERVICE_OPTIONS.map((service) => (
            <button
              key={service}
              type="button"
              onClick={() => onSelectIndex(service)}
              className="rounded-full border px-4 py-2 font-sans text-sm font-medium transition hover:shadow-glow"
              style={{
                borderColor: `${SERVICE_COLORS[service]}88`,
                color: SERVICE_COLORS[service],
                backgroundColor: `${SERVICE_COLORS[service]}14`,
              }}
            >
              {service}
            </button>
          ))}
        </div>
      ) : null}

      {selectedIndex ? (
        <p className="font-sans text-xs uppercase tracking-[0.16em] text-text-muted">
          Valid fields: {INDEX_FIELDS[selectedIndex].join(', ')}
        </p>
      ) : null}
    </div>
  )
}
