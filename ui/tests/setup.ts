import { afterEach, vi } from 'vitest';

Object.defineProperty(window, 'scrollTo', {
  value: vi.fn(),
  writable: true,
});

afterEach(() => {
  window.localStorage.clear();
  vi.clearAllMocks();
});
