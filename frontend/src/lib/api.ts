import axios from 'axios'

export function getApiBaseUrl() {
  return import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, '') ?? ''
}

function stripApiSuffix(url: string) {
  return url.replace(/\/api$/i, '')
}

export function getWebSocketUrl() {
  const explicit = import.meta.env.VITE_WS_URL?.replace(/\/$/, '')
  if (explicit) {
    return explicit
  }

  const apiBase = getApiBaseUrl()
  if (apiBase) {
    if (apiBase.startsWith('/')) {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      return `${protocol}//${window.location.host}/ws`
    }

    const wsBase = stripApiSuffix(apiBase)
    if (apiBase.startsWith('https://')) {
      return `${wsBase.replace('https://', 'wss://')}/ws`
    }

    return `${wsBase.replace('http://', 'ws://')}/ws`
  }

  if (import.meta.env.DEV) {
    return '/ws'
  }

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}/ws`
}

export const apiClient = axios.create({
  baseURL: getApiBaseUrl() || undefined,
  headers: {
    'Content-Type': 'application/json',
  },
})
