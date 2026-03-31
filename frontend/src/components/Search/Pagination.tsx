type PaginationProps = {
  page: number
  canGoPrevious: boolean
  canGoNext: boolean
  pageSizeInput: string
  pageSizeError: string | null
  maxPageSize: number
  onPageSizeInputChange: (pageSize: string) => void
  onPrevious: () => void
  onNext: () => void
}

export function Pagination({
  page,
  canGoPrevious,
  canGoNext,
  pageSizeInput,
  pageSizeError,
  maxPageSize,
  onPageSizeInputChange,
  onPrevious,
  onNext,
}: PaginationProps) {
  return (
    <div className="flex flex-col gap-3 rounded-2xl border border-border bg-surface px-4 py-3">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <label className="flex items-center gap-3 font-sans text-sm text-text-muted">
          Per page
          <input
            type="number"
            min={1}
            max={maxPageSize}
            value={pageSizeInput}
            onChange={(event) => onPageSizeInputChange(event.target.value)}
            className="w-24 rounded-lg border border-border bg-bg px-3 py-2 text-text-primary outline-none transition focus:border-[#58a6ff]/60"
          />
          <span className="text-xs">max {maxPageSize}</span>
        </label>

        <div className="flex items-center justify-between gap-3 sm:justify-end">
          <button
            type="button"
            disabled={!canGoPrevious}
            onClick={onPrevious}
            className="rounded-full border border-border px-4 py-2 font-sans text-sm text-text-primary transition disabled:cursor-not-allowed disabled:opacity-40"
            aria-label="Previous page"
          >
            &lt;
          </button>
          <span className="rounded-full border border-border px-4 py-2 font-sans text-sm text-text-primary">
            Page {page}
          </span>
          <button
            type="button"
            disabled={!canGoNext}
            onClick={onNext}
            className="rounded-full border border-border px-4 py-2 font-sans text-sm text-text-primary transition disabled:cursor-not-allowed disabled:opacity-40"
            aria-label="Next page"
          >
            &gt;
          </button>
        </div>
      </div>

      {pageSizeError ? <p className="font-sans text-xs text-level-ERROR">{pageSizeError}</p> : null}
    </div>
  )
}
