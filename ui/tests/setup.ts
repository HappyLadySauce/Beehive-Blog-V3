import { afterEach, beforeEach, vi } from 'vitest'

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
