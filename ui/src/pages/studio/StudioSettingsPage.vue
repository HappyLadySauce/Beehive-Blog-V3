<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSelect, { type BaseSelectOption } from '@/shared/components/BaseSelect.vue'
import FormField from '@/shared/components/FormField.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import { useFileConfig } from '@/features/file-manager/useFileConfig'

const { t } = useI18n()
const { config, isSaving, errorMessage, loadConfig, saveConfig } = useFileConfig()

onMounted(() => {
  loadConfig()
})

const maxUploadBytesMB = computed({
  get: () => Math.round(config.value.max_upload_bytes / (1024 * 1024)),
  set: (mb: number) => {
    config.value.max_upload_bytes = mb * 1024 * 1024
  },
})

const presignTTLMinutes = computed({
  get: () => Math.round(config.value.presign_ttl_seconds / 60),
  set: (minutes: number) => {
    config.value.presign_ttl_seconds = minutes * 60
  },
})

const contentTypesText = computed({
  get: () => config.value.allowed_content_types.join(', '),
  set: (text: string) => {
    config.value.allowed_content_types = text
      .split(',')
      .map((s) => s.trim())
      .filter(Boolean)
  },
})

async function handleSave(): Promise<void> {
  try {
    await saveConfig({
      max_upload_bytes: config.value.max_upload_bytes,
      allowed_content_types: config.value.allowed_content_types,
      presign_ttl_seconds: config.value.presign_ttl_seconds,
    })
  } catch {
    // error displayed inline
  }
}
</script>

<template>
  <section class="settings-page">
    <StatusAlert v-if="errorMessage" tone="danger" title="Save failed">
      {{ errorMessage }}
    </StatusAlert>

    <form class="settings-page__form" :aria-label="t('settings.formLabel')" @submit.prevent="handleSave">
      <fieldset class="settings-page__section">
        <legend class="settings-page__section-title">{{ t('settings.fileConfig.title') }}</legend>
        <p class="settings-page__section-desc">{{ t('settings.fileConfig.description') }}</p>

        <FormField :label="t('settings.fileConfig.maxUploadBytes')" for-id="max-upload-bytes">
          <BaseInput
            id="max-upload-bytes"
            type="number"
            min="1"
            max="2147483648"
            :model-value="maxUploadBytesMB"
            @update:model-value="maxUploadBytesMB = Number($event)"
          />
          <span class="settings-page__hint">MB</span>
        </FormField>

        <FormField :label="t('settings.fileConfig.allowedContentTypes')" for-id="allowed-types">
          <BaseInput
            id="allowed-types"
            :model-value="contentTypesText"
            @update:model-value="contentTypesText = String($event)"
          />
          <span class="settings-page__hint">{{ t('settings.fileConfig.contentTypesHint') }}</span>
        </FormField>

        <FormField :label="t('settings.fileConfig.presignTTL')" for-id="presign-ttl">
          <BaseInput
            id="presign-ttl"
            type="number"
            min="1"
            max="86400"
            :model-value="presignTTLMinutes"
            @update:model-value="presignTTLMinutes = Number($event)"
          />
          <span class="settings-page__hint">{{ t('settings.fileConfig.ttlHint') }}</span>
        </FormField>
      </fieldset>

      <BaseButton type="submit" :busy="isSaving">{{ t('settings.fileConfig.save') }}</BaseButton>
    </form>
  </section>
</template>

<style scoped>
.settings-page {
  display: grid;
  gap: 16px;
}

.settings-page__form {
  display: grid;
  max-width: 640px;
  gap: 24px;
}

.settings-page__section {
  display: grid;
  gap: 12px;
  border: 1px solid var(--bb-color-line);
  border-radius: 10px;
  padding: 20px;
  margin: 0;
}

.settings-page__section-title {
  font-weight: 600;
  font-size: 1.05rem;
}

.settings-page__section-desc {
  margin: 0 0 8px;
  color: var(--bb-color-muted);
  font-size: 0.88rem;
}

.settings-page__hint {
  color: var(--bb-color-muted);
  font-size: 0.82rem;
}
</style>
