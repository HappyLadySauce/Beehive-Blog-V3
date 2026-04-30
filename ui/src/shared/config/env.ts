export type ApiMode = 'mock' | 'live'

// Resolve API mode: only explicit mock opts into in-memory APIs; default is live (gateway).
// 解析 API 模式：仅显式 mock 时使用内存实现；默认 live（走 gateway）。
function readApiMode(value: string | undefined): ApiMode {
  const normalized = value?.trim().toLowerCase()
  return normalized === 'mock' ? 'mock' : 'live'
}

export const appConfig = {
  apiMode: readApiMode(import.meta.env.VITE_API_MODE),
  gatewayBaseUrl: import.meta.env.VITE_GATEWAY_BASE_URL ?? '',
}
