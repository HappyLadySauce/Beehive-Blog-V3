import { afterEach, beforeEach, vi } from 'vitest'

const viteEnv = import.meta.env as Record<string, string>
viteEnv.VITE_API_MODE = 'mock'
viteEnv.VITE_GATEWAY_BASE_URL = ''

beforeEach(() => {
  window.localStorage.clear()
  Object.defineProperty(window, 'scrollTo', {
    configurable: true,
    value: vi.fn(),
  })
})

afterEach(() => {
  window.localStorage.clear()
})
