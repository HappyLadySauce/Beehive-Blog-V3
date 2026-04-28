import { computed } from 'vue'
import { createI18n, useI18n } from 'vue-i18n'

import { messages } from './messages'

export type AppLocale = keyof typeof messages

export const availableLocales: AppLocale[] = ['zh-CN', 'en-US']

const localeStorageKey = 'beehive.ui.locale'
const fallbackLocale: AppLocale = 'zh-CN'

function isAppLocale(value: string | null | undefined): value is AppLocale {
  return availableLocales.includes(value as AppLocale)
}

function readStoredLocale(): AppLocale {
  if (typeof window === 'undefined') {
    return fallbackLocale
  }
  const stored = window.localStorage.getItem(localeStorageKey)
  return isAppLocale(stored) ? stored : fallbackLocale
}

function writeStoredLocale(locale: AppLocale): void {
  if (typeof window !== 'undefined') {
    window.localStorage.setItem(localeStorageKey, locale)
  }
}

function updateDocumentLocale(locale: AppLocale): void {
  if (typeof document !== 'undefined') {
    document.documentElement.lang = locale
  }
}

export const i18n = createI18n({
  legacy: false,
  locale: readStoredLocale(),
  fallbackLocale,
  messages,
})

export function setLocale(locale: AppLocale): void {
  i18n.global.locale.value = locale
  writeStoredLocale(locale)
  updateDocumentLocale(locale)
}

export function syncDocumentLocale(): void {
  updateDocumentLocale(i18n.global.locale.value as AppLocale)
}

export function useLocale() {
  const { locale } = useI18n()
  const currentLocale = computed(() => {
    const value = String(locale.value)
    return isAppLocale(value) ? value : fallbackLocale
  })

  return {
    locale: currentLocale,
    availableLocales,
    setLocale,
  }
}
