<script setup lang="ts">
import { Upload } from 'lucide-vue-next'
import { computed, useTemplateRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import FileAssetPreviewPanel from '@/features/file-manager/components/FileAssetPreviewPanel.vue'
import FileAssetTable from '@/features/file-manager/components/FileAssetTable.vue'
import { DEFAULT_FILE_CATEGORY_KEY, DEFAULT_FILE_MAX_UPLOAD_BYTES, buildAcceptAttribute } from '@/features/file-manager/constants'
import { useFileCategories } from '@/features/file-manager/useFileCategories'
import { useFileConfig } from '@/features/file-manager/useFileConfig'
import { useFileManager } from '@/features/file-manager/useFileManager'
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

const { t } = useI18n()
const fileInput = useTemplateRef<HTMLInputElement>('fileInput')
const categoriesQuery = useFileCategories({ studio: true })
const { config, loadConfig } = useFileConfig()
const {
  filters,
  items,
  total,
  page,
  pageSize,
  selectedAssetId,
  selectedAsset,
  uploadSelection,
  isUploading,
  uploadErrorMessage,
  isDeleting,
  listQuery,
  detailQuery,
  setPage,
  setPageSize,
  openAsset,
  closeAsset,
  updateUploadSelection,
  uploadSelectedFile,
  removeAsset,
} = useFileManager()

void loadConfig()

const categoryOptions = computed<BaseSelectOption[]>(() => [
  { value: '', label: t('files.filters.allCategories') },
  ...categoriesQuery.items.value.map((category) => ({
    value: category.category_key,
    label: category.display_name,
  })),
])

const uploadCategoryOptions = computed<BaseSelectOption[]>(() => (
  categoriesQuery.enabledItems.value.map((category) => ({
    value: category.category_key,
    label: category.display_name,
  }))
))

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

const uploadAccept = computed(() => buildAcceptAttribute(categoriesQuery.resolveAllowedExtensions(uploadSelection.categoryKey)))

watch(
  () => [categoriesQuery.enabledItems.value, categoriesQuery.defaultCategory.value] as const,
  () => {
    const selected = uploadSelection.categoryKey
    const exists = categoriesQuery.enabledItems.value.some((item) => item.category_key === selected)
    const nextCategoryKey = exists
      ? selected
      : categoriesQuery.defaultCategory.value?.category_key ?? DEFAULT_FILE_CATEGORY_KEY
    updateUploadSelection({
      categoryKey: nextCategoryKey,
      allowedExtensions: categoriesQuery.resolveAllowedExtensions(nextCategoryKey),
    })
  },
  { immediate: true },
)

watch(
  () => config.value.max_upload_bytes,
  (value) => {
    updateUploadSelection({
      maxUploadBytes: value > 0 ? value : DEFAULT_FILE_MAX_UPLOAD_BYTES,
    })
  },
  { immediate: true },
)

watch(
  () => uploadSelection.categoryKey,
  (value) => {
    updateUploadSelection({
      categoryKey: value,
      allowedExtensions: categoriesQuery.resolveAllowedExtensions(value),
    })
  },
)

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
        <FormField :label="t('files.fields.category')" for-id="files-category">
          <BaseSelect id="files-category" v-model="filters.category_key" :options="categoryOptions" />
        </FormField>
        <FormField :label="t('files.fields.status')" for-id="files-status">
          <BaseSelect id="files-status" v-model="filters.status" :options="statusOptions" />
        </FormField>
        <FormField :label="t('files.fields.visibility')" for-id="files-visibility">
          <BaseSelect id="files-visibility" v-model="filters.visibility" :options="visibilityOptions" />
        </FormField>
      </div>

      <div class="files-page__upload">
        <BaseSelect
          v-model="uploadSelection.categoryKey"
          :options="uploadCategoryOptions"
          :aria-label="t('files.uploadCategory')"
          :disabled="uploadCategoryOptions.length === 0"
        />
        <BaseButton :busy="isUploading" :disabled="uploadCategoryOptions.length === 0" @click="openFilePicker">
          <Upload :size="16" aria-hidden="true" />
          {{ t('files.uploadAction') }}
        </BaseButton>
        <input ref="fileInput" class="files-page__input" type="file" :accept="uploadAccept" @change="handleFileChange">
      </div>
    </div>

    <StatusAlert v-if="uploadErrorMessage" tone="danger" :title="t('files.uploadFailedTitle')">
      {{ uploadErrorMessage }}
    </StatusAlert>
    <StatusAlert v-if="categoriesQuery.error" tone="danger" :title="t('files.categoriesUnavailableTitle')">
      {{ categoriesQuery.error instanceof Error ? categoriesQuery.error.message : t('files.categoriesUnavailableMessage') }}
    </StatusAlert>
    <StatusAlert v-if="listQuery.error && items.length === 0" tone="danger" :title="t('files.unavailableTitle')">
      {{ listQuery.error instanceof Error ? listQuery.error.message : t('files.unavailableMessage') }}
    </StatusAlert>

    <PageLoadingState
      v-else-if="(listQuery.showBlockingLoading || categoriesQuery.showBlockingLoading) && items.length === 0"
      :title="t('files.loadingTitle')"
      :rows="5"
    />

    <div v-else class="files-page__table-shell" role="region" :aria-label="t('files.regionLabel')" tabindex="0">
      <div v-if="(listQuery.showRefreshingHint || categoriesQuery.showRefreshingHint) && items.length > 0" class="files-page__refreshing">
        <InlineLoadingState />
      </div>

      <FileAssetTable v-if="items.length > 0" :items="items" @view="openAsset" @delete="removeAsset" />

      <div v-else class="files-page__empty-panel">
        <EmptyState align="center" :title="t('files.empty')" :description="t('files.emptyDescription')">
          <template #actions>
            <BaseButton :disabled="uploadCategoryOptions.length === 0" @click="openFilePicker">{{ t('files.uploadAction') }}</BaseButton>
          </template>
        </EmptyState>
      </div>
    </div>

    <div v-if="total > 0" class="files-page__footer">
      <TablePagination
        :page="page"
        :page-size="pageSize"
        :total="total"
        :disabled="listQuery.isFetching"
        @update:page="setPage"
        @update:page-size="setPageSize"
      />
    </div>

    <SideDrawer :open="Boolean(selectedAssetId)" :title="t('files.detailTitle')" :description="selectedAsset?.file_name" size="lg" @close="closeAsset">
      <PageLoadingState v-if="detailQuery.showBlockingLoading && !selectedAsset" :title="t('files.drawerLoadingTitle')" :rows="4" />
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
