import { computed } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { useProgressiveQuery } from '@/shared/composables'

import { listFileCategories, listStudioFileCategories } from './api'
import { DEFAULT_FILE_CATEGORY_KEY, DEFAULT_FILE_CATEGORY_EXTENSIONS } from './constants'
import type { FileCategory } from './types'

export function useFileCategories(options: { studio?: boolean } = {}) {
  const authStore = useAuthStore()
  const query = useProgressiveQuery({
    queryKey: computed(() => ['file-categories', options.studio === true ? 'studio' : 'public']),
    queryFn: () => (
      options.studio === true
        ? listStudioFileCategories({ accessToken: authStore.accessToken })
        : listFileCategories({ accessToken: authStore.accessToken })
    ),
    delayMs: 250,
  })

  const items = computed<FileCategory[]>(() => query.data.value?.items ?? [])
  const enabledItems = computed(() => items.value.filter((item) => item.enabled))
  const defaultCategory = computed(() => (
    items.value.find((item) => item.is_default)
    ?? items.value.find((item) => item.category_key === DEFAULT_FILE_CATEGORY_KEY)
    ?? items.value[0]
    ?? null
  ))

  function findCategory(categoryKey: string): FileCategory | null {
    return items.value.find((item) => item.category_key === categoryKey) ?? null
  }

  function resolveAllowedExtensions(categoryKey: string): string[] {
    return findCategory(categoryKey)?.allowed_extensions
      ?? defaultCategory.value?.allowed_extensions
      ?? [...DEFAULT_FILE_CATEGORY_EXTENSIONS]
  }

  return {
    ...query,
    items,
    enabledItems,
    defaultCategory,
    findCategory,
    resolveAllowedExtensions,
  }
}
