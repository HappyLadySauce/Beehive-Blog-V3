<script setup lang="ts">
import { PackageOpen, Upload } from 'lucide-vue-next'
import { computed, shallowRef, useTemplateRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import FileAssetPreviewPanel from '@/features/file-manager/components/FileAssetPreviewPanel.vue'
import FileAssetTable from '@/features/file-manager/components/FileAssetTable.vue'
import FileTypeSettingsPanel from '@/features/file-manager/components/FileTypeSettingsPanel.vue'
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

type FilesTab = 'assets' | 'types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
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

const activeTab = shallowRef<FilesTab>(readTabQuery(route.query.tab))

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
const hasAssets = computed(() => items.value.length > 0)
const categoriesErrorMessage = computed(() => {
  const error = categoriesQuery.error.value
  return error instanceof Error ? error.message : ''
})
const assetErrorMessage = computed(() => {
  const error = listQuery.error.value
  return error instanceof Error ? error.message : ''
})
const showAssetBlockingLoading = computed(() => (
  (listQuery.showBlockingLoading.value || categoriesQuery.showBlockingLoading.value)
  && !hasAssets.value
  && !assetErrorMessage.value
))
const showAssetRefreshing = computed(() => (
  (listQuery.showRefreshingHint.value || categoriesQuery.showRefreshingHint.value) && hasAssets.value
))

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

watch(
  () => route.query.tab,
  (value) => {
    activeTab.value = readTabQuery(value)
  },
)

watch(activeTab, async (tab) => {
  if (readTabQuery(route.query.tab) === tab) {
    return
  }
  await router.replace({
    query: {
      ...route.query,
      tab: tab === 'assets' ? undefined : tab,
    },
  })
})

function readTabQuery(value: unknown): FilesTab {
  const normalized = Array.isArray(value) ? value[0] : value
  return normalized === 'types' ? 'types' : 'assets'
}

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
    <div class="studio-workspace-toolbar">
      <div class="studio-workspace-tabs" role="tablist" :aria-label="t('files.tabs.workspaceLabel')">
        <button type="button" :class="{ active: activeTab === 'assets' }" @click="activeTab = 'assets'">
          {{ t('files.tabs.assets') }}
        </button>
        <button type="button" :class="{ active: activeTab === 'types' }" @click="activeTab = 'types'">
          {{ t('files.tabs.types') }}
        </button>
      </div>

      <div v-if="activeTab === 'assets'" class="files-page__upload">
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

    <template v-if="activeTab === 'assets'">
      <div class="studio-list-shell">
        <div class="studio-list-filters files-page__filters">
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

        <StatusAlert v-if="uploadErrorMessage" tone="danger" :title="t('files.uploadFailedTitle')">
          {{ uploadErrorMessage }}
        </StatusAlert>
        <StatusAlert v-if="categoriesErrorMessage" tone="danger" :title="t('files.categoriesUnavailableTitle')">
          {{ categoriesErrorMessage || t('files.categoriesUnavailableMessage') }}
        </StatusAlert>
        <StatusAlert v-if="assetErrorMessage && !hasAssets" tone="danger" :title="t('files.unavailableTitle')">
          {{ assetErrorMessage || t('files.unavailableMessage') }}
        </StatusAlert>
        <StatusAlert v-else-if="assetErrorMessage && hasAssets" tone="danger" :title="t('files.unavailableTitle')">
          {{ assetErrorMessage || t('files.unavailableMessage') }}
        </StatusAlert>

        <PageLoadingState
          v-if="showAssetBlockingLoading"
          :title="t('files.loadingTitle')"
          :rows="5"
        />

        <div
          v-else-if="listQuery.hasResolvedOnce.value"
          class="studio-list-table files-page__table-shell"
          role="region"
          :aria-label="t('files.regionLabel')"
          tabindex="0"
        >
          <div v-if="showAssetRefreshing" class="files-page__refreshing">
            <InlineLoadingState />
          </div>

          <FileAssetTable v-if="hasAssets" :items="items" @view="openAsset" @delete="removeAsset" />

          <div v-else class="studio-list-empty-panel files-page__empty-panel">
            <EmptyState
              class="studio-list-empty-state"
              align="center"
              :title="t('files.empty')"
              :description="t('files.emptyDescription')"
            >
              <template #visual>
                <PackageOpen :size="52" aria-hidden="true" />
              </template>
              <BaseButton :disabled="uploadCategoryOptions.length === 0" @click="openFilePicker">
                {{ t('files.uploadAction') }}
              </BaseButton>
            </EmptyState>
          </div>
        </div>

        <div v-if="listQuery.data.value" class="studio-list-footer">
          <TablePagination
            :page="page"
            :page-size="pageSize"
            :total="total"
            :disabled="listQuery.isFetching.value"
            @update:page="setPage"
            @update:page-size="setPageSize"
          />
        </div>
      </div>

      <SideDrawer :open="Boolean(selectedAssetId)" :title="t('files.detailTitle')" :description="selectedAsset?.file_name" size="lg" @close="closeAsset">
        <PageLoadingState v-if="detailQuery.showBlockingLoading.value && !selectedAsset" :title="t('files.drawerLoadingTitle')" :rows="4" />
        <FileAssetPreviewPanel v-else :asset="selectedAsset" :busy-delete="isDeleting" @delete="handleDeleteSelected" />
        <template #footer>
          <BaseButton variant="ghost" @click="closeAsset">{{ t('common.close') }}</BaseButton>
        </template>
      </SideDrawer>
    </template>

    <FileTypeSettingsPanel v-else />
  </section>
</template>

<style scoped>
.files-page {
  display: grid;
  gap: 16px;
}

.files-page__filters {
  grid-template-columns: minmax(220px, 1fr) repeat(3, minmax(160px, 180px));
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
}

.files-page__refreshing {
  display: flex;
  justify-content: flex-end;
  padding: 12px 12px 0;
}

.files-page__empty-panel {
  min-width: 980px;
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
