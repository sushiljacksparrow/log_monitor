import { useEffect, useRef, useState } from 'react'

import { getWebSocketUrl } from '../lib/api'
import type { BackendLog, LiveLogEntry } from '../types/logs'

const MAX_LOGS = 500
const RECONNECT_DELAY_MS = 3000

function isObject(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null
}

function isBackendLog(value: unknown): value is BackendLog {
  if (!isObject(value) || typeof value.service !== 'string' || typeof value.level !== 'string') {
    return false
  }

  if (
    typeof value.message !== 'string' ||
    typeof value.request_id !== 'string' ||
    typeof value.timestamp !== 'string'
  ) {
    return false
  }

  switch (value.service) {
    case 'auth-service':
      return typeof value.user_id === 'string' && typeof value.ip === 'string'
    case 'order-service':
      return (
        typeof value.user_id === 'string' &&
        typeof value.order_id === 'string' &&
        typeof value.product_id === 'string' &&
        typeof value.stock_left === 'number' &&
        typeof value.carrier === 'string'
      )
    case 'payment-service':
      return (
        typeof value.order_id === 'string' &&
        typeof value.payment_id === 'string' &&
        typeof value.gateway === 'string' &&
        (typeof value.amount === 'number' || typeof value.amount === 'string')
      )
    default:
      return false
  }
}

export function useWebSocket() {
  const [logs, setLogs] = useState<LiveLogEntry[]>([])
  const [isConnected, setIsConnected] = useState(false)
  const [isReconnecting, setIsReconnecting] = useState(false)
  const reconnectTimerRef = useRef<number | null>(null)
  const socketRef = useRef<WebSocket | null>(null)
  const idCounterRef = useRef(0)

  useEffect(() => {
    let active = true

    const clearTimer = () => {
      if (reconnectTimerRef.current !== null) {
        window.clearTimeout(reconnectTimerRef.current)
        reconnectTimerRef.current = null
      }
    }

    const scheduleReconnect = () => {
      if (!active || reconnectTimerRef.current !== null) {
        return
      }
      setIsReconnecting(true)
      reconnectTimerRef.current = window.setTimeout(() => {
        reconnectTimerRef.current = null
        connect()
      }, RECONNECT_DELAY_MS)
    }

    const handleChunk = (chunk: string) => {
      const trimmed = chunk.trim()
      if (!trimmed) {
        return
      }

      try {
        const parsed = JSON.parse(trimmed) as unknown
        if (!isBackendLog(parsed)) {
          return
        }

        const entry: LiveLogEntry = {
          ...parsed,
          clientId: `live-${Date.now()}-${idCounterRef.current++}`,
          receivedAt: Date.now(),
        }

        setLogs((current) => {
          const next = [...current, entry]
          return next.length > MAX_LOGS ? next.slice(next.length - MAX_LOGS) : next
        })
      } catch {
        // Ignore malformed chunks from the stream.
      }
    }

    const connect = () => {
      clearTimer()

      try {
        const socket = new WebSocket(getWebSocketUrl())
        socketRef.current = socket

        socket.addEventListener('open', () => {
          if (!active) {
            socket.close()
            return
          }

          setIsConnected(true)
          setIsReconnecting(false)
        })

        socket.addEventListener('message', (event) => {
          const payload = typeof event.data === 'string' ? event.data : ''
          payload.split('\n').forEach(handleChunk)
        })

        socket.addEventListener('close', () => {
          setIsConnected(false)
          scheduleReconnect()
        })

        socket.addEventListener('error', () => {
          socket.close()
        })
      } catch {
        setIsConnected(false)
        scheduleReconnect()
      }
    }

    connect()

    return () => {
      active = false
      clearTimer()
      socketRef.current?.close()
    }
  }, [])

  return {
    logs,
    isConnected,
    isReconnecting,
  }
}
