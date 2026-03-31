import { useMemo, useRef, useState } from 'react'

import {
  DEFAULT_SEARCH_PAGE_SIZE,
  INDEX_FIELDS,
  MAX_BACKEND_SEARCH_PAGE_SIZE,
  NUMERIC_SEARCH_FIELDS,
} from '../constants/schema'
import axios from 'axios'
import { apiClient } from '../lib/api'
import type { ApiResponse, CursorStack, SearchChip, SearchData, ServiceLogMap, ServiceName } from '../types/logs'

type SearchState = {
  selectedIndex: ServiceName | null
  chips: SearchChip[]
  fieldInput: string
  valueInput: string
  pendingField: string | null
}

type SearchResult = ServiceLogMap[ServiceName]
type SearchRequestBody = Record<string, string | number>

function getErrorMessage(payload: unknown, fallback: string) {
  if (typeof payload === 'string' && payload) {
    return payload
  }

  if (payload && typeof payload === 'object' && 'message' in payload && typeof payload.message === 'string') {
    return payload.message
  }

  return fallback
}

function getAxiosErrorMessage(error: unknown, fallback: string) {
  if (axios.isAxiosError(error)) {
    return getErrorMessage(error.response?.data, error.message || fallback)
  }

  return error instanceof Error ? error.message : fallback
}

function buildSearchBody(service: ServiceName, chips: SearchChip[], cursor: string | null): SearchRequestBody {
  const body: SearchRequestBody = {}
  const numericFields = new Set<string>(NUMERIC_SEARCH_FIELDS[service])

  for (const chip of chips) {
    if (numericFields.has(chip.field)) {
      const numericValue = Number(chip.value.trim())
      if (!Number.isFinite(numericValue)) {
        throw new Error(`${chip.field} must be a valid number`)
      }
      body[chip.field] = numericValue
      continue
    }

    body[chip.field] = chip.value
  }

  if (cursor) {
    body.sorted_value = cursor
  }

  return body
}

export function useSearch() {
  const [state, setState] = useState<SearchState>({
    selectedIndex: null,
    chips: [],
    fieldInput: '',
    valueInput: '',
    pendingField: null,
  })
  const [results, setResults] = useState<SearchResult[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [message, setMessage] = useState<string | null>(null)
  const [hasMore, setHasMore] = useState(false)
  const [, setCursorStack] = useState<CursorStack>([null])
  const [currentPageIndex, setCurrentPageIndex] = useState(0)
  const [pageSize, setPageSize] = useState(DEFAULT_SEARCH_PAGE_SIZE)
  const [pageSizeInput, setPageSizeInput] = useState(String(DEFAULT_SEARCH_PAGE_SIZE))
  const [pageSizeError, setPageSizeError] = useState<string | null>(null)
  const [hasSearched, setHasSearched] = useState(false)
  const cursorStackRef = useRef<CursorStack>([null])

  const suggestions = useMemo(() => {
    if (!state.selectedIndex || state.pendingField) {
      return []
    }

    const usedFields = new Set(state.chips.map((chip) => chip.field))
    const query = state.fieldInput.trim().toLowerCase()

    return INDEX_FIELDS[state.selectedIndex].filter((field) => {
      if (usedFields.has(field)) {
        return false
      }

      return !query || field.toLowerCase().includes(query)
    })
  }, [state.chips, state.fieldInput, state.pendingField, state.selectedIndex])

  const resetResults = () => {
    setResults([])
    setError(null)
    setMessage(null)
    setHasMore(false)
    setHasSearched(false)
    setCursorStack([null])
    setCurrentPageIndex(0)
    cursorStackRef.current = [null]
  }

  const selectIndex = (service: ServiceName) => {
    setState({
      selectedIndex: service,
      chips: [],
      fieldInput: '',
      valueInput: '',
      pendingField: null,
    })
    resetResults()
  }

  const clearIndex = () => {
    setState({
      selectedIndex: null,
      chips: [],
      fieldInput: '',
      valueInput: '',
      pendingField: null,
    })
    resetResults()
  }

  const selectField = (field: string) => {
    setState((current) => ({
      ...current,
      pendingField: field,
      fieldInput: '',
      valueInput: '',
    }))
  }

  const completePendingField = () => {
    const value = state.valueInput.trim()
    if (!state.pendingField || !value) {
      return
    }

    setState((current) => ({
      ...current,
      chips: [...current.chips, { field: current.pendingField as string, value }],
      pendingField: null,
      valueInput: '',
      fieldInput: '',
    }))
  }

  const removeChip = (field: string) => {
    setState((current) => ({
      ...current,
      chips: current.chips.filter((chip) => chip.field !== field),
    }))
  }

  const removeLastItem = () => {
    setState((current) => {
      if (current.pendingField) {
        return { ...current, pendingField: null, valueInput: '' }
      }

      if (current.chips.length === 0) {
        return current
      }

      return { ...current, chips: current.chips.slice(0, -1) }
    })
  }

  const fetchPage = async (pageIndex: number, cursor: string | null, shouldPushCursor: boolean) => {
    if (!state.selectedIndex) {
      return
    }

    setLoading(true)
    setError(null)
    setMessage(null)

    try {
      const body = buildSearchBody(state.selectedIndex, state.chips, cursor)
      const response = await apiClient.post<ApiResponse<SearchData<SearchResult>>>(
        `/api/search/${state.selectedIndex}`,
        body,
        {
          params: { size: pageSize },
        },
      )

      const payload = response.data

      if (!payload.data) {
        throw new Error(getErrorMessage(payload, 'Search response was empty'))
      }

      const nextCursor = payload.data.base_response.sorted_value || null

      setResults(payload.data.logs)
      setHasMore(payload.data.base_response.has_more)
      setMessage(payload.message)
      setHasSearched(true)

      if (shouldPushCursor) {
        setCursorStack((current) => {
          const next = current.slice(0, pageIndex + 1)
          next[pageIndex + 1] = nextCursor
          cursorStackRef.current = next
          return next
        })
      }

      setCurrentPageIndex(pageIndex)
    } catch (caught) {
      const nextMessage = getAxiosErrorMessage(caught, 'Search request failed')
      setResults([])
      setHasMore(false)
      setError(nextMessage)
      setHasSearched(true)
    } finally {
      setLoading(false)
    }
  }

  const runSearch = async () => {
    if (pageSizeError) {
      setError(pageSizeError)
      return
    }
    setCursorStack([null])
    setCurrentPageIndex(0)
    cursorStackRef.current = [null]
    await fetchPage(0, null, true)
  }

  const goToNextPage = async () => {
    if (pageSizeError) {
      setError(pageSizeError)
      return
    }
    const cursor = cursorStackRef.current[currentPageIndex + 1] ?? null
    await fetchPage(currentPageIndex + 1, cursor, true)
  }

  const goToPreviousPage = async () => {
    if (pageSizeError) {
      setError(pageSizeError)
      return
    }
    const prevPageIndex = currentPageIndex - 1
    const cursor = cursorStackRef.current[prevPageIndex] ?? null
    await fetchPage(prevPageIndex, cursor, false)
  }

  return {
    ...state,
    suggestions,
    results,
    loading,
    error,
    message,
    hasMore,
    hasSearched,
    currentPageIndex,
    pageSize,
    pageSizeInput,
    pageSizeError,
    maxPageSize: MAX_BACKEND_SEARCH_PAGE_SIZE,
    canSearch: Boolean(state.selectedIndex && !state.pendingField && !pageSizeError),
    canGoPrevious: currentPageIndex > 0 && !loading && !pageSizeError,
    canGoNext: hasMore && !loading && !pageSizeError,
    setPageSizeInput: (value: string) => {
      setPageSizeInput(value)

      const trimmed = value.trim()
      if (!trimmed) {
        setPageSizeError('Page size is required')
        return
      }

      const numericValue = Number(trimmed)
      if (!Number.isInteger(numericValue) || numericValue < 1) {
        setPageSizeError('Page size must be a whole number between 1 and 100')
        return
      }

      if (numericValue > MAX_BACKEND_SEARCH_PAGE_SIZE) {
        setPageSizeError('Page size cannot be more than 100')
        return
      }

      setPageSizeError(null)
      setPageSize(numericValue)
      setCursorStack([null])
      cursorStackRef.current = [null]
      setCurrentPageIndex(0)
      setResults([])
      setHasMore(false)
      setError(null)
      setMessage(null)
      setHasSearched(false)
    },
    setFieldInput: (fieldInput: string) => setState((current) => ({ ...current, fieldInput })),
    setValueInput: (valueInput: string) => setState((current) => ({ ...current, valueInput })),
    selectIndex,
    clearIndex,
    selectField,
    completePendingField,
    removeChip,
    removeLastItem,
    runSearch,
    goToNextPage,
    goToPreviousPage,
  }
}
