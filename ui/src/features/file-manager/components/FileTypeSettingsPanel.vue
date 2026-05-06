<script setup lang="ts">
import { FolderCog } from 'lucide-vue-next'
import { computed, onMounted, reactive, shallowRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useAuthStore } from '@/features/auth/stores/authStore'
import {
  createStudioFileCategory,
  setStudioDefaultFileCategory,
  updateStudioFileCategory,
  updateStudioFileCategoryExtensions,
} from '@/features/file-manager/api'
import { DEFAULT_FILE_CATEGORY_EXTENSIONS, FILE_EXTENSION_OPTIONS } from '@/features/file-manager/constants'
import type { FileCategory } from '@/features/file-manager/types'
import { useFileCategories } from '@/features/file-manager/useFileCategories'
import { useFileConfig } from '@/features/file-manager/useFileConfig'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSelect, { type BaseSelectOption } from '@/shared/components/BaseSelect.vue'
import EmptyState from '@/shared/components/EmptyState.vue'
import FormField from '@/shared/components/FormField.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import { useToast } from '@/shared/composables'

type FileSizeUnit = 'B' | 'KB' | 'MB' | 'GB'

interface CategoryDraft {
  display_name: string
  description: string
  enabled: boolean
  sort_order: number
  allowed_extensions: string[]
}

const fileSizeUnits: BaseSelectOption[] = [
  { value: 'B', label: 'B' },
  { value: 'KB', label: 'KB' },
  { value: 'MB', label: 'MB' },
  { value: 'GB', label: 'GB' },
]

const { t } = useI18n()
const authStore = useAuthStore()
const { pushToast } = useToast()
const { config, isSaving, errorMessage, loadConfig, saveConfig } = useFileConfig()
const categoriesQuery = useFileCategories({ studio: true })

const categoryItems = computed(() => categoriesQuery.items.value)
const enabledCategoryItems = computed(() => categoriesQuery.enabledItems.value)
const categoryOptions = computed<BaseSelectOption[]>(() => (
  categoryItems.value.map((item) => ({ value: item.category_key, label: item.display_name }))
))
const categoryDrafts = reactive<Record<string, CategoryDraft>>({})
const categoryActionError = shallowRef('')
const activeCategoryKey = shallowRef('')
const categorySavingKey = shallowRef('')
const extensionSavingKey = shallowRef('')
const defaultSavingKey = shallowRef('')
const isCreatingCategory = shallowRef(false)
const maxUploadSizeValue = shallowRef('2')
const maxUploadSizeUnit = shallowRef<FileSizeUnit>('GB')
const createForm = reactive({
  category_key: '',
  display_name: '',
  description: '',
  enabled: true,
  is_default: false,
  sort_order: 100,
  allowed_extensions: [...DEFAULT_FILE_CATEGORY_EXTENSIONS],
})

const activeCategory = computed(() => (
  categoryItems.value.find((item) => item.category_key === activeCategoryKey.value) ?? null
))
const activeDraft = computed(() => (
  activeCategoryKey.value ? categoryDrafts[activeCategoryKey.value] ?? null : null
))
const categoriesErrorMessage = computed(() => {
  const error = categoriesQuery.error.value
  return error instanceof Error ? error.message : ''
})
const showCategoryLoading = computed(() => (
  categoriesQuery.showBlockingLoading.value && categoryItems.value.length === 0
))
const showCategoryEmpty = computed(() => (
  categoriesQuery.hasResolvedOnce.value
  && !categoriesErrorMessage.value
  && categoryItems.value.length === 0
))

onMounted(async () => {
  await loadConfig()
})

watch(
  () => config.value.max_upload_bytes,
  (value) => {
    const { amount, unit } = splitFileSize(value)
    maxUploadSizeValue.value = String(amount)
    maxUploadSizeUnit.value = unit
  },
  { immediate: true },
)

watch(
  categoryItems,
  (items) => {
    syncCategoryDrafts(items)
    if (!activeCategoryKey.value || !items.some((item) => item.category_key === activeCategoryKey.value)) {
      activeCategoryKey.value = items[0]?.category_key ?? ''
    }
  },
  { immediate: true },
)

async function handleSaveConfig(): Promise<void> {
  try {
    await saveConfig({
      max_upload_bytes: joinFileSize(Number(maxUploadSizeValue.value), maxUploadSizeUnit.value),
      presign_ttl_seconds: Math.max(60, Number(config.value.presign_ttl_seconds) || 300),
    })
    pushToast({
      tone: 'success',
      title: 'File config saved',
      message: 'System upload settings updated.',
    })
  } catch {
    // Inline alert handles the failure message.
    // 由内联提示承接失败信息。
  }
}

async function handleCreateCategory(): Promise<void> {
  categoryActionError.value = ''
  isCreatingCategory.value = true
  try {
    await createStudioFileCategory({
      category_key: createForm.category_key.trim(),
      display_name: createForm.display_name.trim(),
      description: createForm.description.trim(),
      enabled: createForm.enabled,
      is_default: createForm.is_default,
      sort_order: Number(createForm.sort_order) || 0,
      allowed_extensions: uniqueExtensions(createForm.allowed_extensions),
    }, { accessToken: authStore.accessToken })
    resetCreateForm()
    await categoriesQuery.refetch()
    pushToast({
      tone: 'success',
      title: 'Category created',
      message: 'The file category is now available.',
    })
  } catch (error) {
    categoryActionError.value = error instanceof Error ? error.message : 'Failed to create file category.'
  } finally {
    isCreatingCategory.value = false
  }
}

async function handleSaveCategory(categoryKey: string): Promise<void> {
  const draft = categoryDrafts[categoryKey]
  if (!draft) {
    return
  }
  categoryActionError.value = ''
  categorySavingKey.value = categoryKey
  try {
    await updateStudioFileCategory(categoryKey, {
      display_name: draft.display_name.trim(),
      description: draft.description.trim(),
      enabled: draft.enabled,
      sort_order: Number(draft.sort_order) || 0,
    }, { accessToken: authStore.accessToken })
    await categoriesQuery.refetch()
    pushToast({
      tone: 'success',
      title: 'Category updated',
      message: categoryKey,
    })
  } catch (error) {
    categoryActionError.value = error instanceof Error ? error.message : 'Failed to update file category.'
  } finally {
    categorySavingKey.value = ''
  }
}

async function handleSaveExtensions(categoryKey: string): Promise<void> {
  const draft = categoryDrafts[categoryKey]
  if (!draft) {
    return
  }
  categoryActionError.value = ''
  extensionSavingKey.value = categoryKey
  try {
    await updateStudioFileCategoryExtensions(categoryKey, {
      allowed_extensions: uniqueExtensions(draft.allowed_extensions),
    }, { accessToken: authStore.accessToken })
    await categoriesQuery.refetch()
    pushToast({
      tone: 'success',
      title: 'Extensions updated',
      message: categoryKey,
    })
  } catch (error) {
    categoryActionError.value = error instanceof Error ? error.message : 'Failed to update category extensions.'
  } finally {
    extensionSavingKey.value = ''
  }
}

async function handleSetDefaultCategory(categoryKey: string): Promise<void> {
  categoryActionError.value = ''
  defaultSavingKey.value = categoryKey
  try {
    await setStudioDefaultFileCategory(categoryKey, { accessToken: authStore.accessToken })
    await categoriesQuery.refetch()
    pushToast({
      tone: 'success',
      title: 'Default category updated',
      message: categoryKey,
    })
  } catch (error) {
    categoryActionError.value = error instanceof Error ? error.message : 'Failed to set default category.'
  } finally {
    defaultSavingKey.value = ''
  }
}

function toggleDraftExtension(categoryKey: string, extension: string, checked: boolean): void {
  const draft = categoryDrafts[categoryKey]
  if (!draft) {
    return
  }
  if (checked) {
    draft.allowed_extensions = uniqueExtensions([...draft.allowed_extensions, extension])
    return
  }
  draft.allowed_extensions = draft.allowed_extensions.filter((item) => item !== extension)
}

function toggleCreateExtension(extension: string, checked: boolean): void {
  if (checked) {
    createForm.allowed_extensions = uniqueExtensions([...createForm.allowed_extensions, extension])
    return
  }
  createForm.allowed_extensions = createForm.allowed_extensions.filter((item) => item !== extension)
}

function syncCategoryDrafts(items: FileCategory[]): void {
  const activeKeys = new Set(items.map((item) => item.category_key))
  for (const categoryKey of Object.keys(categoryDrafts)) {
    if (!activeKeys.has(categoryKey)) {
      delete categoryDrafts[categoryKey]
    }
  }
  for (const item of items) {
    categoryDrafts[item.category_key] = {
      display_name: item.display_name,
      description: item.description,
      enabled: item.enabled,
      sort_order: item.sort_order,
      allowed_extensions: [...item.allowed_extensions],
    }
  }
}

function resetCreateForm(): void {
  createForm.category_key = ''
  createForm.display_name = ''
  createForm.description = ''
  createForm.enabled = true
  createForm.is_default = false
  createForm.sort_order = 100
  createForm.allowed_extensions = [...DEFAULT_FILE_CATEGORY_EXTENSIONS]
}

function uniqueExtensions(extensions: readonly string[]): string[] {
  return [...new Set(extensions.map((item) => item.trim().toLowerCase()).filter(Boolean))].sort()
}

function splitFileSize(bytes: number): { amount: number; unit: FileSizeUnit } {
  const safeBytes = Math.max(1, bytes || 2 * 1024 * 1024 * 1024)
  if (safeBytes % (1024 * 1024 * 1024) === 0) {
    return { amount: safeBytes / (1024 * 1024 * 1024), unit: 'GB' }
  }
  if (safeBytes % (1024 * 1024) === 0) {
    return { amount: safeBytes / (1024 * 1024), unit: 'MB' }
  }
  if (safeBytes % 1024 === 0) {
    return { amount: safeBytes / 1024, unit: 'KB' }
  }
  return { amount: safeBytes, unit: 'B' }
}

function joinFileSize(amount: number, unit: FileSizeUnit): number {
  const safeAmount = Math.max(1, Number.isFinite(amount) ? amount : 1)
  const multipliers: Record<FileSizeUnit, number> = {
    B: 1,
    KB: 1024,
    MB: 1024 * 1024,
    GB: 1024 * 1024 * 1024,
  }
  return Math.round(safeAmount * multipliers[unit])
}
</script>

<template>
  <section class="file-types-panel">
    <StatusAlert v-if="errorMessage" tone="danger" :title="t('files.types.config.saveFailedTitle')">
      {{ errorMessage }}
    </StatusAlert>
    <StatusAlert v-if="categoryActionError" tone="danger" :title="t('files.types.categories.actionFailedTitle')">
      {{ categoryActionError }}
    </StatusAlert>
    <StatusAlert v-if="categoriesErrorMessage" tone="danger" :title="t('files.types.categories.unavailableTitle')">
      {{ categoriesErrorMessage || t('files.types.categories.unavailableMessage') }}
    </StatusAlert>

    <section class="file-types-panel__section">
      <div class="file-types-panel__heading">
        <h2>{{ t('files.types.config.title') }}</h2>
        <p>{{ t('files.types.config.description') }}</p>
      </div>

      <form class="file-types-panel__form" :aria-label="t('files.types.config.formLabel')" @submit.prevent="handleSaveConfig">
        <div class="file-types-panel__inline-grid">
          <FormField :label="t('files.types.config.maxUploadBytes')" for-id="files-max-upload-bytes">
            <BaseInput
              id="files-max-upload-bytes"
              type="number"
              min="1"
              :model-value="maxUploadSizeValue"
              @update:model-value="maxUploadSizeValue = String($event)"
            />
          </FormField>
          <FormField :label="t('files.types.config.sizeUnit')" for-id="files-max-upload-unit">
            <BaseSelect id="files-max-upload-unit" v-model="maxUploadSizeUnit" :options="fileSizeUnits" />
          </FormField>
        </div>

        <FormField :label="t('files.types.config.presignTTL')" for-id="files-presign-ttl">
          <BaseInput
            id="files-presign-ttl"
            type="number"
            min="60"
            :model-value="config.presign_ttl_seconds"
            @update:model-value="config.presign_ttl_seconds = Number($event)"
          />
          <span class="file-types-panel__hint">{{ t('files.types.config.ttlHint') }}</span>
        </FormField>

        <div class="file-types-panel__actions">
          <BaseButton type="submit" :busy="isSaving">{{ t('files.types.config.save') }}</BaseButton>
        </div>
      </form>
    </section>

    <section class="file-types-panel__section">
      <div class="file-types-panel__heading">
        <h2>{{ t('files.types.categories.title') }}</h2>
        <p>{{ t('files.types.categories.description') }}</p>
      </div>

      <PageLoadingState v-if="showCategoryLoading" :title="t('files.types.categories.loadingTitle')" :rows="4" />

      <div v-else class="file-types-panel__categories">
        <div v-if="showCategoryEmpty" class="studio-list-empty-panel file-types-panel__empty-panel">
          <EmptyState
            class="studio-list-empty-state"
            align="center"
            :title="t('files.types.categories.empty')"
            :description="t('files.types.categories.emptyDescription')"
          >
            <template #visual>
              <FolderCog :size="52" aria-hidden="true" />
            </template>
          </EmptyState>
        </div>

        <form class="file-types-panel__form file-types-panel__category-create" @submit.prevent="handleCreateCategory">
          <div class="file-types-panel__inline-grid">
            <FormField :label="t('files.types.categories.fields.categoryKey')" for-id="new-category-key">
              <BaseInput id="new-category-key" v-model="createForm.category_key" />
            </FormField>
            <FormField :label="t('files.types.categories.fields.displayName')" for-id="new-category-name">
              <BaseInput id="new-category-name" v-model="createForm.display_name" />
            </FormField>
          </div>

          <FormField :label="t('common.description')" for-id="new-category-description">
            <BaseInput id="new-category-description" v-model="createForm.description" />
          </FormField>

          <div class="file-types-panel__inline-grid">
            <FormField :label="t('files.types.categories.fields.sortOrder')" for-id="new-category-sort-order">
              <BaseInput id="new-category-sort-order" v-model.number="createForm.sort_order" type="number" />
            </FormField>
            <div class="file-types-panel__checks">
              <label class="file-types-panel__check">
                <input v-model="createForm.enabled" type="checkbox">
                {{ t('files.types.categories.fields.enabled') }}
              </label>
              <label class="file-types-panel__check">
                <input v-model="createForm.is_default" type="checkbox">
                {{ t('files.types.categories.fields.isDefault') }}
              </label>
            </div>
          </div>

          <div class="file-types-panel__extensions">
            <span class="file-types-panel__extensions-title">{{ t('files.types.categories.fields.allowedExtensions') }}</span>
            <div class="file-types-panel__extensions-grid">
              <label v-for="option in FILE_EXTENSION_OPTIONS" :key="option.value" class="file-types-panel__check">
                <input
                  type="checkbox"
                  :checked="createForm.allowed_extensions.includes(option.value)"
                  @change="toggleCreateExtension(option.value, ($event.target as HTMLInputElement).checked)"
                >
                {{ option.label }}
              </label>
            </div>
          </div>

          <div class="file-types-panel__actions">
            <BaseButton type="submit" :busy="isCreatingCategory">{{ t('files.types.categories.create') }}</BaseButton>
          </div>
        </form>

        <div v-if="categoryItems.length > 0" class="file-types-panel__editor-shell">
          <div class="file-types-panel__inline-grid file-types-panel__inline-grid--categories">
            <FormField :label="t('files.types.categories.activeCategory')" for-id="active-category">
              <BaseSelect id="active-category" v-model="activeCategoryKey" :options="categoryOptions" />
            </FormField>
          </div>

          <div v-if="activeCategory && activeDraft" class="file-types-panel__editor">
            <div class="file-types-panel__meta">
              <div>
                <h3>{{ activeCategory.display_name }}</h3>
                <p>{{ activeCategory.category_key }}</p>
              </div>
              <span v-if="activeCategory.is_default" class="file-types-panel__badge">{{ t('files.types.categories.defaultBadge') }}</span>
            </div>

            <div class="file-types-panel__inline-grid">
              <FormField :label="t('files.types.categories.fields.displayName')" for-id="edit-category-name">
                <BaseInput id="edit-category-name" v-model="activeDraft.display_name" />
              </FormField>
              <FormField :label="t('files.types.categories.fields.sortOrder')" for-id="edit-category-sort-order">
                <BaseInput id="edit-category-sort-order" v-model.number="activeDraft.sort_order" type="number" />
              </FormField>
            </div>

            <FormField :label="t('common.description')" for-id="edit-category-description">
              <BaseInput id="edit-category-description" v-model="activeDraft.description" />
            </FormField>

            <label class="file-types-panel__check">
              <input v-model="activeDraft.enabled" type="checkbox">
              {{ t('files.types.categories.fields.enabled') }}
            </label>

            <div class="file-types-panel__extensions">
              <span class="file-types-panel__extensions-title">{{ t('files.types.categories.fields.allowedExtensions') }}</span>
              <div class="file-types-panel__extensions-grid">
                <label v-for="option in FILE_EXTENSION_OPTIONS" :key="option.value" class="file-types-panel__check">
                  <input
                    type="checkbox"
                    :checked="activeDraft.allowed_extensions.includes(option.value)"
                    @change="toggleDraftExtension(activeCategory.category_key, option.value, ($event.target as HTMLInputElement).checked)"
                  >
                  {{ option.label }}
                </label>
              </div>
            </div>

            <div class="file-types-panel__actions">
              <BaseButton :busy="categorySavingKey === activeCategory.category_key" @click="handleSaveCategory(activeCategory.category_key)">
                {{ t('files.types.categories.saveCategory') }}
              </BaseButton>
              <BaseButton variant="secondary" :busy="extensionSavingKey === activeCategory.category_key" @click="handleSaveExtensions(activeCategory.category_key)">
                {{ t('files.types.categories.saveExtensions') }}
              </BaseButton>
              <BaseButton
                variant="ghost"
                :busy="defaultSavingKey === activeCategory.category_key"
                :disabled="activeCategory.is_default || !enabledCategoryItems.length"
                @click="handleSetDefaultCategory(activeCategory.category_key)"
              >
                {{ t('files.types.categories.setDefault') }}
              </BaseButton>
            </div>
          </div>
        </div>
      </div>
    </section>
  </section>
</template>

<style scoped>
.file-types-panel,
.file-types-panel__categories,
.file-types-panel__editor,
.file-types-panel__editor-shell,
.file-types-panel__form {
  display: grid;
  gap: 16px;
}

.file-types-panel__section {
  display: grid;
  gap: 16px;
  border: 1px solid var(--bb-color-line);
  border-radius: 12px;
  padding: 20px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-soft);
}

.file-types-panel__heading,
.file-types-panel__meta {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.file-types-panel__heading h2,
.file-types-panel__meta h3 {
  margin: 0;
  font-weight: 700;
}

.file-types-panel__heading p,
.file-types-panel__meta p,
.file-types-panel__hint {
  margin: 0;
  color: var(--bb-color-muted);
  font-size: 0.88rem;
}

.file-types-panel__inline-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.file-types-panel__inline-grid--categories {
  grid-template-columns: minmax(220px, 360px);
}

.file-types-panel__checks,
.file-types-panel__extensions-grid,
.file-types-panel__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
}

.file-types-panel__actions {
  justify-content: flex-start;
}

.file-types-panel__check {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--bb-color-text);
}

.file-types-panel__extensions {
  display: grid;
  gap: 10px;
}

.file-types-panel__extensions-title {
  font-weight: 700;
  color: var(--bb-color-text-strong);
}

.file-types-panel__badge {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  padding: 4px 10px;
  color: var(--bb-color-primary);
  background: var(--bb-color-primary-soft);
  font-size: 0.82rem;
  font-weight: 700;
}

.file-types-panel__editor-shell {
  border-top: 1px solid var(--bb-color-line);
  padding-top: 16px;
}

.file-types-panel__empty-panel {
  border-top: 0;
}

@media (max-width: 900px) {
  .file-types-panel__inline-grid {
    grid-template-columns: 1fr;
  }

  .file-types-panel__actions {
    align-items: stretch;
  }
}
</style>
