<script setup lang="ts">
import { Eye, Trash2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import type { FileAssetSummary } from '../types'

import IconActionButton from '@/shared/components/IconActionButton.vue'
import StatusBadge from '@/shared/components/StatusBadge.vue'

defineProps<{
  items: FileAssetSummary[]
}>()

const emit = defineEmits<{
  view: [asset: FileAssetSummary]
  delete: [asset: FileAssetSummary]
}>()

const { t, locale } = useI18n()

function formatFileSize(byteSize: number): string {
  if (byteSize < 1024) return `${byteSize} B`
  if (byteSize < 1024 * 1024) return `${(byteSize / 1024).toFixed(1)} KB`
  return `${(byteSize / (1024 * 1024)).toFixed(1)} MB`
}

function formatUnixTime(value?: number): string {
  if (!value) {
    return t('common.none')
  }
  return new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value * 1000))
}
</script>

<template>
  <table class="file-asset-table">
    <thead>
      <tr>
        <th scope="col">{{ t('files.columns.file') }}</th>
        <th scope="col">{{ t('files.columns.scope') }}</th>
        <th scope="col">{{ t('files.columns.type') }}</th>
        <th scope="col">{{ t('files.columns.size') }}</th>
        <th scope="col">{{ t('files.columns.status') }}</th>
        <th scope="col">{{ t('files.columns.owner') }}</th>
        <th scope="col">{{ t('files.columns.updated') }}</th>
        <th scope="col">{{ t('common.actions') }}</th>
      </tr>
    </thead>
    <tbody v-if="items.length > 0">
      <tr v-for="asset in items" :key="asset.asset_id">
        <td>
          <strong>{{ asset.file_name }}</strong>
          <span>{{ asset.object_key }}</span>
        </td>
        <td><StatusBadge :value="asset.scope" /></td>
        <td>{{ asset.content_type }}</td>
        <td>{{ formatFileSize(asset.byte_size) }}</td>
        <td><StatusBadge :value="asset.status" /></td>
        <td>{{ asset.owner_user_id }}</td>
        <td>{{ formatUnixTime(asset.deleted_at ?? asset.uploaded_at ?? asset.created_at) }}</td>
        <td>
          <div class="file-asset-table__actions">
            <IconActionButton :aria-label="t('files.actions.view', { name: asset.file_name })" :title="t('files.actions.view', { name: asset.file_name })" @click="emit('view', asset)">
              <Eye :size="17" aria-hidden="true" />
            </IconActionButton>
            <IconActionButton tone="danger" :aria-label="t('files.actions.delete', { name: asset.file_name })" :title="t('files.actions.delete', { name: asset.file_name })" @click="emit('delete', asset)">
              <Trash2 :size="17" aria-hidden="true" />
            </IconActionButton>
          </div>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<style scoped>
.file-asset-table {
  width: 100%;
  min-width: 980px;
  border-collapse: collapse;
}

.file-asset-table th,
.file-asset-table td {
  border-bottom: 1px solid var(--bb-color-line);
  padding: 10px 12px;
  text-align: left;
  vertical-align: middle;
}

.file-asset-table th {
  color: var(--bb-color-muted);
  font-size: 0.8rem;
  text-transform: uppercase;
  background: var(--bb-color-surface);
}

.file-asset-table tbody tr:nth-child(even) {
  background: var(--bb-color-subtle);
}

.file-asset-table tbody tr:hover {
  background: var(--bb-color-primary-soft);
}

.file-asset-table td:first-child {
  display: grid;
  gap: 3px;
}

.file-asset-table td:first-child span {
  color: var(--bb-color-muted);
}

.file-asset-table th:last-child,
.file-asset-table td:last-child {
  text-align: right;
}

.file-asset-table__actions {
  display: inline-flex;
  gap: 8px;
  justify-content: flex-end;
}
</style>
