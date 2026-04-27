import { computed, readonly, shallowRef } from 'vue'

export type ThemePreference = 'system' | 'light' | 'dark'
export type ResolvedTheme = 'light' | 'dark'

const storageKey = 'beehive.ui.theme'
const theme = shallowRef<ThemePreference>('system')
const systemTheme = shallowRef<ResolvedTheme>('light')
let isInitialized = false
let mediaQuery: MediaQueryList | null = null

function readStoredTheme(): ThemePreference {
  if (typeof window === 'undefined') {
    return 'system'
  }

  const stored = window.localStorage.getItem(storageKey)
  return stored === 'light' || stored === 'dark' || stored === 'system' ? stored : 'system'
}

function readSystemTheme(): ResolvedTheme {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') {
    return 'light'
  }
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

function applyTheme(nextTheme: ResolvedTheme): void {
  if (typeof document === 'undefined') {
    return
  }
  document.documentElement.dataset.theme = nextTheme
  document.documentElement.style.colorScheme = nextTheme
}

function initializeTheme(): void {
  if (isInitialized) {
    return
  }

  theme.value = readStoredTheme()
  systemTheme.value = readSystemTheme()
  applyTheme(resolvedTheme.value)

  if (typeof window !== 'undefined' && typeof window.matchMedia === 'function') {
    mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    mediaQuery.addEventListener('change', (event) => {
      systemTheme.value = event.matches ? 'dark' : 'light'
      if (theme.value === 'system') {
        applyTheme(systemTheme.value)
      }
    })
  }

  isInitialized = true
}

const resolvedTheme = computed<ResolvedTheme>(() => (theme.value === 'system' ? systemTheme.value : theme.value))

function setTheme(nextTheme: ThemePreference): void {
  initializeTheme()
  theme.value = nextTheme

  if (typeof window !== 'undefined') {
    window.localStorage.setItem(storageKey, nextTheme)
  }
  applyTheme(resolvedTheme.value)
}

function toggleTheme(): void {
  initializeTheme()
  setTheme(resolvedTheme.value === 'dark' ? 'light' : 'dark')
}

export function useTheme() {
  initializeTheme()

  return {
    theme: readonly(theme),
    resolvedTheme,
    setTheme,
    toggleTheme,
  }
}
