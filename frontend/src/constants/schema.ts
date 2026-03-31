export const SERVICE_COLORS = {
  'auth-service': '#a371f7',
  'order-service': '#f0883e',
  'payment-service': '#39d353',
} as const

export const LEVEL_COLORS = {
  INFO: '#3fb950',
  WARN: '#d29922',
  ERROR: '#f85149',
  DEBUG: '#58a6ff',
} as const

export const INDEX_FIELDS = {
  'auth-service': ['service', 'level', 'message', 'request_id', 'user_id', 'ip', 'timestamp'],
  'order-service': [
    'service',
    'level',
    'message',
    'request_id',
    'user_id',
    'order_id',
    'product_id',
    'stock_left',
    'carrier',
    'timestamp',
  ],
  'payment-service': [
    'service',
    'level',
    'message',
    'request_id',
    'order_id',
    'payment_id',
    'gateway',
    'amount',
    'timestamp',
  ],
} as const

export const TABLE_COLUMNS = {
  'auth-service': ['timestamp', 'level', 'service', 'message', 'request_id', 'user_id', 'ip'],
  'order-service': [
    'timestamp',
    'level',
    'service',
    'message',
    'request_id',
    'user_id',
    'order_id',
    'product_id',
    'stock_left',
    'carrier',
  ],
  'payment-service': [
    'timestamp',
    'level',
    'service',
    'message',
    'request_id',
    'order_id',
    'payment_id',
    'gateway',
    'amount',
  ],
} as const

export const DEFAULT_SEARCH_PAGE_SIZE = 25
export const MAX_BACKEND_SEARCH_PAGE_SIZE = 100

export const NUMERIC_SEARCH_FIELDS = {
  'auth-service': [],
  'order-service': [],
  'payment-service': ['amount'],
} as const

export const SERVICE_OPTIONS = Object.keys(SERVICE_COLORS) as Array<keyof typeof SERVICE_COLORS>

export const LEVEL_OPTIONS = Object.keys(LEVEL_COLORS) as Array<keyof typeof LEVEL_COLORS>
