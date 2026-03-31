import { useSearch } from '../../hooks/useSearch'
import { ChipQueryBuilder } from './ChipQueryBuilder'
import { Pagination } from './Pagination'
import { ResultsTable } from './ResultsTable'

export function SearchTab() {
  const search = useSearch()

  return (
    <section className="space-y-4">
      <ChipQueryBuilder
        selectedIndex={search.selectedIndex}
        chips={search.chips}
        fieldInput={search.fieldInput}
        valueInput={search.valueInput}
        pendingField={search.pendingField}
        suggestions={search.suggestions}
        onSelectIndex={search.selectIndex}
        onClearIndex={search.clearIndex}
        onFieldInputChange={search.setFieldInput}
        onValueInputChange={search.setValueInput}
        onSelectField={search.selectField}
        onCompletePendingField={search.completePendingField}
        onRemoveChip={search.removeChip}
        onRemoveLastItem={search.removeLastItem}
        onRunSearch={search.runSearch}
        canSearch={search.canSearch}
      />

      {search.error ? (
        <div className="rounded-2xl border border-level-ERROR/40 bg-level-ERROR/10 px-4 py-3 font-sans text-sm text-level-ERROR">
          {search.error}
        </div>
      ) : null}

      {search.selectedIndex ? (
        <div className="flex flex-wrap items-center justify-between gap-3 font-sans text-sm text-text-muted">
          <p>
            Showing <span className="text-text-primary">{search.results.length}</span> results
          </p>
          <p className="text-xs">
            Backend limit: <span className="text-text-primary">100</span> per request
          </p>
        </div>
      ) : null}

      {search.selectedIndex ? (
        <div className="flex flex-col gap-4">
          {search.hasSearched ? (
            <ResultsTable selectedIndex={search.selectedIndex} results={search.results} loading={search.loading} />
          ) : null}

          <div className="flex justify-end">
            <div className="w-full max-w-xl">
              <Pagination
                page={search.currentPageIndex + 1}
                canGoPrevious={search.canGoPrevious}
                canGoNext={search.canGoNext}
                pageSizeInput={search.pageSizeInput}
                pageSizeError={search.pageSizeError}
                maxPageSize={search.maxPageSize}
                onPageSizeInputChange={search.setPageSizeInput}
                onPrevious={search.goToPreviousPage}
                onNext={search.goToNextPage}
              />
            </div>
          </div>
        </div>
      ) : null}
    </section>
  )
}
