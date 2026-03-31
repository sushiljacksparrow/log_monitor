import { useMemo, useState } from 'react'

import { LEVEL_OPTIONS, SERVICE_OPTIONS } from '../constants/schema'
import type { LiveLogEntry, LogLevel, ServiceName } from '../types/logs'

export function useLogFilter(logs: LiveLogEntry[]) {
  const [keyword, setKeyword] = useState('')
  const [activeLevels, setActiveLevels] = useState<LogLevel[]>([...LEVEL_OPTIONS])
  const [activeServices, setActiveServices] = useState<ServiceName[]>([...SERVICE_OPTIONS])

  const filteredLogs = useMemo(() => {
    const query = keyword.trim().toLowerCase()

    return logs.filter((log) => {
      if (!activeLevels.includes(log.level) || !activeServices.includes(log.service)) {
        return false
      }

      if (!query) {
        return true
      }

      return Object.values(log).some((value) => String(value).toLowerCase().includes(query))
    })
  }, [activeLevels, activeServices, keyword, logs])

  const toggleLevel = (level: LogLevel) => {
    setActiveLevels((current) =>
      current.includes(level) ? current.filter((item) => item !== level) : [...current, level],
    )
  }

  const toggleService = (service: ServiceName) => {
    setActiveServices((current) =>
      current.includes(service) ? current.filter((item) => item !== service) : [...current, service],
    )
  }

  return {
    keyword,
    setKeyword,
    activeLevels,
    activeServices,
    filteredLogs,
    toggleLevel,
    toggleService,
  }
}
