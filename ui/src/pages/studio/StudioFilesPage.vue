<script setup lang="ts">
import { Upload } from 'lucide-vue-next'
import { computed, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSelect, { type BaseSelectOption } from '@/shared/components/BaseSelect.vue'
import EmptyState from '@/shared/components/EmptyState.vue'
import FormField from '@/shared/components/FormField.vue'
import InlineLoadingState from '@/shared/components/InlineLoadingState.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import SideDrawer from '@/shared/components/SideDrawer.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import TablePagination from '@/shared/components/TablePagination.vue'

import FileAssetPreviewPanel from '@/features/file-manager/components/FileAssetPreviewPanel.vue'
import FileAssetTable from '@/features/file-manager/components/FileAssetTable.vue'
import { useFileManager } from '@/features/file-manager/useFileManager'

const { t } = useI18n()
const fileInput = useTemplateRef<HTMLInputElement>('fileInput')
const {
  filters,
  items,
  total,
  page,
  pageSize,
  selectedAssetId,
  selectedAsset,
  uploadNamespace,
  isUploading,
  uploadErrorMessage,
  isDeleting,
  listQuery,
  detailQuery,
  setPage,
  setPageSize,
  openAsset,
  closeAsset,
  uploadSelectedFile,
  removeAsset,
} = useFileManager()

const namespaceOptions = computed<BaseSelectOption[]>(() => [
  { value: '', label: t('contentType.all') },
  { value: 'avatar', label: t('files.scope.avatar') },
  { value: 'content_cover', label: t('files.scope.content_cover') },
  { value: 'content_image', label: t('files.scope.content_image') },
  { value: 'attachment', label: t('files.scope.attachment') },
])

const statusOptions = computed<BaseSelectOption[]>(() => [
  { value: '', label: t('contentStatus.all') },
  { value: 'pending', label: t('files.status.pending') },
  { value: 'uploaded', label: t('files.status.uploaded') },
  { value: 'deleted', label: t('files.status.deleted') },
])

const visibilityOptions = computed<BaseSelectOption[]>(() => [
  { value: '', label: t('visibility.all') },
  { value: 'public', label: t('visibility.public') },
  { value: 'private', label: t('visibility.private') },
])

const uploadNamespaceOptions = computed<BaseSelectOption[]>(() => [
  { value: 'content_image', label: t('files.scope.content_image') },
  { value: 'content_cover', label: t('files.scope.content_cover') },
  { value: 'avatar', label: t('files.scope.avatar') },
  { value: 'attachment', label: t('files.scope.attachment') },
])

const uploadAccept = computed(() => {
  const namespace = uploadNamespace.value
  if (namespace === 'attachment') {
    return 'image/png,image/jpeg,image/webp,image/avif,application/pdf'
  }
  return 'image/png,image/jpeg,image/webp,image/avif'
})

function openFilePicker(): void {
  fileInput.value?.click()
}

async function handleFileChange(event: Event): Promise<void> {
  const input = event.target
  if (!(input instanceof HTMLInputElement) || !input.files?.[0]) {
    return
  }
  try {
    await uploadSelectedFile(input.files[0])
  } finally {
    input.value = ''
  }
}

async function handleDeleteSelected(): Promise<void> {
  if (!selectedAsset.value) {
    return
  }
  await removeAsset(selectedAsset.value)
}
</script>

<template>
  <section class="files-page">
    <div class="files-page__toolbar">
      <div class="files-page__filters">
        <FormField :label="t('common.search')" for-id="files-search">
          <BaseInput id="files-search" v-model="filters.keyword" :placeholder="t('files.searchPlaceholder')" />
        </FormField>
        <FormField :label="t('files.fields.scope')" for-id="files-scope">
          <BaseSelect id="files-scope" v-model="filters.namespace" :options="namespaceOptions" />
        </FormField>
        <FormField :label="t('files.fields.status')" for-id="files-status">
          <BaseSelect id="files-status" v-model="filters.status" :options="statusOptions" />
        </FormField>
        <FormField :label="t('files.fields.visibility')" for-id="files-visibility">
          <BaseSelect id="files-visibility" v-model="filters.visibility" :options="visibilityOptions" />
        </FormField>
      </div>

      <div class="files-page__upload">
        <BaseSelect v-model="uploadNamespace" :options="uploadNamespaceOptions" :aria-label="t('files.uploadScope')" />
        <BaseButton :busy="isUploading" @click="openFilePicker">
          <Upload :size="16" aria-hidden="true" />
          {{ t('files.uploadAction') }}
        </BaseButton>
        <input ref="fileInput" class="files-page__input" type="file" :accept="uploadAccept" @change="handleFileChange">
      </div>
    </div>

    <StatusAlert v-if="uploadErrorMessage" tone="danger" :title="t('files.uploadFailedTitle')">
      {{ uploadErrorMessage }}
    </StatusAlert>
    <StatusAlert v-if="listQuery.error.value && items.length === 0" tone="danger" :title="t('files.unavailableTitle')">
      {{ listQuery.error.value instanceof Error ? listQuery.error.value.message : t('files.unavailableMessage') }}
    </StatusAlert>

    <PageLoadingState v-else-if="listQuery.showBlockingLoading.value && items.length === 0" :title="t('files.loadingTitle')" :rows="5" />

    <div v-else class="files-page__table-shell" role="region" :aria-label="t('files.regionLabel')" tabindex="0">
      <div v-if="listQuery.showRefreshingHint.value && items.length > 0" class="files-page__refreshing">
        <InlineLoadingState />
      </div>

      <FileAssetTable v-if="items.length > 0" :items="items" @view="openAsset" @delete="removeAsset" />

      <div v-else class="files-page__empty-panel">
        <EmptyState align="center" :title="t('files.empty')" :description="t('files.emptyDescription')">
          <template #actions>
            <BaseButton @click="openFilePicker">{{ t('files.uploadAction') }}</BaseButton>
          </template>
        </EmptyState>
      </div>
    </div>

    <div v-if="total > 0" class="files-page__footer">
      <TablePagination
        :page="page"
        :page-size="pageSize"
        :total="total"
        :disabled="listQuery.isFetching.value"
        @update:page="setPage"
        @update:page-size="setPageSize"
      />
    </div>

    <SideDrawer :open="Boolean(selectedAssetId)" :title="t('files.detailTitle')" :description="selectedAsset?.file_name" size="lg" @close="closeAsset">
      <PageLoadingState v-if="detailQuery.showBlockingLoading.value && !selectedAsset" :title="t('files.drawerLoadingTitle')" :rows="4" />
      <FileAssetPreviewPanel v-else :asset="selectedAsset" :busy-delete="isDeleting" @delete="handleDeleteSelected" />
      <template #footer>
        <BaseButton variant="ghost" @click="closeAsset">{{ t('common.close') }}</BaseButton>
      </template>
    </SideDrawer>
  </section>
</template>

<style scoped>
.files-page {
  display: grid;
  gap: 16px;
}

.files-page__toolbar {
  display: grid;
  gap: 12px;
}

.files-page__filters {
  display: grid;
  grid-template-columns: minmax(220px, 1fr) repeat(3, minmax(160px, 180px));
  gap: 12px;
  align-items: end;
}

.files-page__upload {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  align-items: center;
  flex-wrap: wrap;
}

.files-page__input {
  display: none;
}

.files-page__table-shell {
  overflow-x: auto;
  border: 1px solid var(--bb-color-line);
  border-radius: 10px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-soft);
}

.files-page__table-shell:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.files-page__refreshing {
  display: flex;
  justify-content: flex-end;
  padding: 12px 12px 0;
}

.files-page__empty-panel {
  min-width: 980px;
  padding: 24px;
}

.files-page__footer {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 900px) {
  .files-page__filters {
    grid-template-columns: 1fr;
  }

  .files-page__upload {
    justify-content: stretch;
  }
}
</style>
