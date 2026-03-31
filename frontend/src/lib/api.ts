import axios from 'axios'

export function getApiBaseUrl() {
  return import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, '') ?? ''
}

export function getWebSocketUrl() {
  const explicit = import.meta.env.VITE_WS_URL?.replace(/\/$/, '')
  if (explicit) {
    return explicit
  }

  const apiBase = getApiBaseUrl()
  if (apiBase) {
    if (apiBase.startsWith('https://')) {
      return `${apiBase.replace('https://', 'wss://')}/ws`
    }

    return `${apiBase.replace('http://', 'ws://')}/ws`
  }

  if (import.meta.env.DEV) {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    return `${protocol}//${window.location.hostname}:8000/ws`
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
