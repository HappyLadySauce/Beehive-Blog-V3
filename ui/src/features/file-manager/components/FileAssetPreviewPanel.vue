<script setup lang="ts">
import { ExternalLink, FileImage, FileText } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '@/shared/components/BaseButton.vue'
import ReadonlyField from '@/shared/components/ReadonlyField.vue'

import type { FileAssetDetail } from '../types'

const props = defineProps<{
  asset: FileAssetDetail | null
  busyDelete?: boolean
}>()

const emit = defineEmits<{
  delete: []
}>()

const { t } = useI18n()
const isImage = computed(() => props.asset?.content_type.startsWith('image/') ?? false)
const isPdf = computed(() => props.asset?.content_type === 'application/pdf')

function formatFileSize(byteSize?: number): string {
  const value = byteSize ?? 0
  if (value < 1024) return `${value} B`
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`
  return `${(value / (1024 * 1024)).toFixed(1)} MB`
}
</script>

<template>
  <div v-if="asset" class="file-preview-panel">
    <div class="file-preview-panel__hero" :class="{ 'file-preview-panel__hero--empty': !isImage }">
      <img v-if="isImage" :src="asset.public_url" :alt="asset.file_name">
      <FileText v-else-if="isPdf" :size="36" aria-hidden="true" />
      <FileImage v-else :size="36" aria-hidden="true" />
    </div>

    <div class="file-preview-panel__grid">
      <ReadonlyField :label="t('files.fields.fileName')" :value="asset.file_name" />
      <ReadonlyField :label="t('files.fields.scope')" :value="asset.namespace" />
      <ReadonlyField :label="t('files.fields.status')" :value="asset.status" />
      <ReadonlyField :label="t('files.fields.visibility')" :value="asset.visibility" />
      <ReadonlyField :label="t('files.fields.contentType')" :value="asset.content_type" />
      <ReadonlyField :label="t('files.fields.byteSize')" :value="formatFileSize(asset.byte_size)" />
      <ReadonlyField :label="t('files.fields.ownerUserId')" :value="asset.owner_user_id" />
      <ReadonlyField :label="t('files.fields.objectKey')" :value="asset.object_key" />
      <ReadonlyField :label="t('files.fields.createdAt')" :value="String(asset.created_at)" />
      <ReadonlyField :label="t('files.fields.uploadedAt')" :value="asset.uploaded_at ? String(asset.uploaded_at) : t('common.none')" />
    </div>

    <div class="file-preview-panel__actions">
      <a v-if="asset.public_url" class="file-preview-panel__link" :href="asset.public_url" target="_blank" rel="noreferrer">
        <ExternalLink :size="16" aria-hidden="true" />
        {{ t('files.openAsset') }}
      </a>
      <BaseButton variant="danger" :busy="busyDelete" @click="emit('delete')">
        {{ t('files.deleteAsset') }}
      </BaseButton>
    </div>
  </div>
</template>

<style scoped>
.file-preview-panel {
  display: grid;
  gap: 16px;
}

.file-preview-panel__hero {
  aspect-ratio: 16 / 9;
  overflow: hidden;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  background: var(--bb-color-surface-elevated);
}

.file-preview-panel__hero img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: contain;
  background: var(--bb-color-surface);
}

.file-preview-panel__hero--empty {
  display: grid;
  place-items: center;
  color: var(--bb-color-muted);
}

.file-preview-panel__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.file-preview-panel__actions {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.file-preview-panel__link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--bb-color-primary);
  text-decoration: none;
}

@media (max-width: 720px) {
  .file-preview-panel__grid {
    grid-template-columns: 1fr;
  }

  .file-preview-panel__actions {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
