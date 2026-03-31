export type LogLevel = 'INFO' | 'WARN' | 'ERROR' | 'DEBUG'

export type ServiceName = 'auth-service' | 'order-service' | 'payment-service'

export type AuthLog = {
  service: 'auth-service'
  level: LogLevel
  message: string
  request_id: string
  user_id: string
  ip: string
  timestamp: string
}

export type OrderLog = {
  service: 'order-service'
  level: LogLevel
  message: string
  request_id: string
  user_id: string
  order_id: string
  product_id: string
  stock_left: number
  carrier: string
  timestamp: string
}

export type PaymentLog = {
  service: 'payment-service'
  level: LogLevel
  message: string
  request_id: string
  order_id: string
  payment_id: string
  gateway: string
  amount: number
  timestamp: string
}

export type BackendLog = AuthLog | OrderLog | PaymentLog

export type LiveLogEntry = BackendLog & {
  clientId: string
  receivedAt: number
}

export type PaginationInfo = {
  has_more: boolean
  sorted_value: string
}

export type ServiceLogMap = {
  'auth-service': AuthLog
  'order-service': OrderLog
  'payment-service': PaymentLog
}

export type SearchData<TLog> = {
  logs: TLog[]
  base_response: PaginationInfo
}

export type ApiResponse<TData> = {
  statusCode: number
  message: string
  data: TData | null
}

export type SearchResponse<TService extends ServiceName> = ApiResponse<
  SearchData<ServiceLogMap[TService]>
>

export type SearchChip = {
  field: string
  value: string
}

export type CursorStack = Array<string | null>
