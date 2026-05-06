import { shallowRef } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'

import { getStudioFileConfig, updateStudioFileConfig } from './api'
import type { FileConfig } from './types'

export function useFileConfig() {
  const authStore = useAuthStore()
  const config = shallowRef<FileConfig>({
    max_upload_bytes: 0,
    presign_ttl_seconds: 0,
  })
  const isSaving = shallowRef(false)
  const errorMessage = shallowRef('')

  async function loadConfig(): Promise<void> {
    const data = await getStudioFileConfig({ accessToken: authStore.accessToken })
    config.value = data.config
  }

  async function saveConfig(patch: Partial<FileConfig>): Promise<void> {
    isSaving.value = true
    errorMessage.value = ''
    try {
      const data = await updateStudioFileConfig(patch, { accessToken: authStore.accessToken })
      config.value = data.config
    } catch (error) {
      errorMessage.value = error instanceof Error ? error.message : 'Failed to save file config.'
      throw error
    } finally {
      isSaving.value = false
    }
  }

  return { config, isSaving, errorMessage, loadConfig, saveConfig }
}
