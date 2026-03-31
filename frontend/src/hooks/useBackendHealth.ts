import { useEffect, useState } from 'react'

import { apiClient, getApiBaseUrl } from '../lib/api'

const HEALTH_CHECK_INTERVAL_MS = 3000
const HEALTH_CHECK_TIMEOUT_MS = 1500

export type BackendHealthStatus = 'checking' | 'healthy' | 'unhealthy'

export function useBackendHealth() {
  const [isBackendReachable, setIsBackendReachable] = useState(false)
  const [checkedAtLeastOnce, setCheckedAtLeastOnce] = useState(false)

  useEffect(() => {
    let cancelled = false

    const runCheck = async () => {
      try {
        const response = await apiClient.get('/health', {
          timeout: HEALTH_CHECK_TIMEOUT_MS,
        })

        if (!cancelled) {
          setIsBackendReachable(response.status === 200)
          setCheckedAtLeastOnce(true)
        }
      } catch {
        if (!cancelled) {
          setIsBackendReachable(false)
          setCheckedAtLeastOnce(true)
        }
      }
    }

    void runCheck()
    const intervalId = window.setInterval(() => {
      void runCheck()
    }, HEALTH_CHECK_INTERVAL_MS)

    return () => {
      cancelled = true
      window.clearInterval(intervalId)
    }
  }, [])

  const status: BackendHealthStatus = !checkedAtLeastOnce
    ? 'checking'
    : isBackendReachable
      ? 'healthy'
      : 'unhealthy'

  return {
    status,
    isBackendReachable,
    checkedAtLeastOnce,
    backendLabel: getApiBaseUrl() || window.location.origin,
  }
}
