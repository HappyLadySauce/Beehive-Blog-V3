import { shallowRef } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { requestJson } from '@/shared/api/httpClient'

export interface FileConfigData {
  max_upload_bytes: number
  allowed_content_types: string[]
  presign_ttl_seconds: number
}

export function useFileConfig() {
  const authStore = useAuthStore()
  const config = shallowRef<FileConfigData>({
    max_upload_bytes: 0,
    allowed_content_types: [],
    presign_ttl_seconds: 0,
  })
  const isSaving = shallowRef(false)
  const errorMessage = shallowRef('')

  async function loadConfig(): Promise<void> {
    const data = await requestJson<{ config: FileConfigData }>(
      '/api/v3/studio/file/config',
      { method: 'GET', accessToken: authStore.accessToken },
    )
    config.value = data.config
  }

  async function saveConfig(patch: Partial<FileConfigData>): Promise<void> {
    isSaving.value = true
    errorMessage.value = ''
    try {
      const data = await requestJson<{ config: FileConfigData }>(
        '/api/v3/studio/file/config',
        {
          method: 'PUT',
          body: JSON.stringify(patch),
          accessToken: authStore.accessToken,
        },
      )
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
