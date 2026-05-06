<script setup lang="ts">
import { computed, onMounted, reactive, shallowRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import {
  createStudioFileCategory,
  setStudioDefaultFileCategory,
  updateStudioFileCategory,
  updateStudioFileCategoryExtensions,
} from '@/features/file-manager/api'
import { DEFAULT_FILE_CATEGORY_EXTENSIONS, FILE_EXTENSION_OPTIONS } from '@/features/file-manager/constants'
import { useFileCategories } from '@/features/file-manager/useFileCategories'
import { useFileConfig } from '@/features/file-manager/useFileConfig'
import type { FileCategory } from '@/features/file-manager/types'
import { useAuthStore } from '@/features/auth/stores/authStore'
import { useToast } from '@/shared/composables'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSelect, { type BaseSelectOption } from '@/shared/components/BaseSelect.vue'
import FormField from '@/shared/components/FormField.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'

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

const activeCategory = computed(() => (
  categoryItems.value.find((item) => item.category_key === activeCategoryKey.value) ?? null
))

const activeDraft = computed(() => (
  activeCategoryKey.value ? categoryDrafts[activeCategoryKey.value] ?? null : null
))

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
    // error displayed inline
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
  <section class="settings-page">
    <StatusAlert v-if="errorMessage" tone="danger" :title="t('settings.fileConfig.saveFailedTitle')">
      {{ errorMessage }}
    </StatusAlert>
    <StatusAlert v-if="categoryActionError" tone="danger" :title="t('settings.fileCategories.actionFailedTitle')">
      {{ categoryActionError }}
    </StatusAlert>

    <form class="settings-page__form" :aria-label="t('settings.formLabel')" @submit.prevent="handleSaveConfig">
      <fieldset class="settings-page__section">
        <legend class="settings-page__section-title">{{ t('settings.fileConfig.title') }}</legend>
        <p class="settings-page__section-desc">{{ t('settings.fileConfig.description') }}</p>

        <div class="settings-page__inline-grid">
          <FormField :label="t('settings.fileConfig.maxUploadBytes')" for-id="max-upload-bytes">
            <BaseInput
              id="max-upload-bytes"
              type="number"
              min="1"
              :model-value="maxUploadSizeValue"
              @update:model-value="maxUploadSizeValue = String($event)"
            />
          </FormField>
          <FormField :label="t('settings.fileConfig.sizeUnit')" for-id="max-upload-unit">
            <BaseSelect id="max-upload-unit" v-model="maxUploadSizeUnit" :options="fileSizeUnits" />
          </FormField>
        </div>

        <FormField :label="t('settings.fileConfig.presignTTL')" for-id="presign-ttl">
          <BaseInput
            id="presign-ttl"
            type="number"
            min="60"
            :model-value="config.presign_ttl_seconds"
            @update:model-value="config.presign_ttl_seconds = Number($event)"
          />
          <span class="settings-page__hint">{{ t('settings.fileConfig.ttlHint') }}</span>
        </FormField>

        <BaseButton type="submit" :busy="isSaving">{{ t('settings.fileConfig.save') }}</BaseButton>
      </fieldset>
    </form>

    <section class="settings-page__section">
      <div class="settings-page__section-header">
        <div>
          <h2 class="settings-page__section-title">{{ t('settings.fileCategories.title') }}</h2>
          <p class="settings-page__section-desc">{{ t('settings.fileCategories.description') }}</p>
        </div>
      </div>

      <form class="settings-page__category-create" @submit.prevent="handleCreateCategory">
        <div class="settings-page__inline-grid">
          <FormField :label="t('settings.fileCategories.fields.categoryKey')" for-id="new-category-key">
            <BaseInput id="new-category-key" v-model="createForm.category_key" />
          </FormField>
          <FormField :label="t('settings.fileCategories.fields.displayName')" for-id="new-category-name">
            <BaseInput id="new-category-name" v-model="createForm.display_name" />
          </FormField>
        </div>

        <FormField :label="t('common.description')" for-id="new-category-description">
          <BaseInput id="new-category-description" v-model="createForm.description" />
        </FormField>

        <div class="settings-page__inline-grid">
          <FormField :label="t('settings.fileCategories.fields.sortOrder')" for-id="new-category-sort-order">
            <BaseInput id="new-category-sort-order" v-model.number="createForm.sort_order" type="number" />
          </FormField>
          <div class="settings-page__checks">
            <label class="settings-page__check">
              <input v-model="createForm.enabled" type="checkbox">
              {{ t('settings.fileCategories.fields.enabled') }}
            </label>
            <label class="settings-page__check">
              <input v-model="createForm.is_default" type="checkbox">
              {{ t('settings.fileCategories.fields.isDefault') }}
            </label>
          </div>
        </div>

        <div class="settings-page__extensions">
          <span class="settings-page__extensions-title">{{ t('settings.fileCategories.fields.allowedExtensions') }}</span>
          <div class="settings-page__extensions-grid">
            <label v-for="option in FILE_EXTENSION_OPTIONS" :key="option.value" class="settings-page__check">
              <input
                type="checkbox"
                :checked="createForm.allowed_extensions.includes(option.value)"
                @change="toggleCreateExtension(option.value, ($event.target as HTMLInputElement).checked)"
              >
              {{ option.label }}
            </label>
          </div>
        </div>

        <BaseButton type="submit" :busy="isCreatingCategory">{{ t('settings.fileCategories.create') }}</BaseButton>
      </form>

      <div class="settings-page__inline-grid settings-page__inline-grid--categories">
        <FormField :label="t('settings.fileCategories.activeCategory')" for-id="active-category">
          <BaseSelect id="active-category" v-model="activeCategoryKey" :options="categoryOptions" />
        </FormField>
      </div>

      <PageLoadingState v-if="categoriesQuery.showBlockingLoading.value && categoryItems.length === 0" :title="t('settings.fileCategories.loadingTitle')" :rows="4" />

      <div v-else-if="activeCategory && activeDraft" class="settings-page__category-editor">
        <div class="settings-page__category-meta">
          <div>
            <h3 class="settings-page__category-title">{{ activeCategory.display_name }}</h3>
            <p class="settings-page__category-key">{{ activeCategory.category_key }}</p>
          </div>
          <span v-if="activeCategory.is_default" class="settings-page__badge">{{ t('settings.fileCategories.defaultBadge') }}</span>
        </div>

        <div class="settings-page__inline-grid">
          <FormField :label="t('settings.fileCategories.fields.displayName')" for-id="edit-category-name">
            <BaseInput id="edit-category-name" v-model="activeDraft.display_name" />
          </FormField>
          <FormField :label="t('settings.fileCategories.fields.sortOrder')" for-id="edit-category-sort-order">
            <BaseInput id="edit-category-sort-order" v-model.number="activeDraft.sort_order" type="number" />
          </FormField>
        </div>

        <FormField :label="t('common.description')" for-id="edit-category-description">
          <BaseInput id="edit-category-description" v-model="activeDraft.description" />
        </FormField>

        <label class="settings-page__check">
          <input v-model="activeDraft.enabled" type="checkbox">
          {{ t('settings.fileCategories.fields.enabled') }}
        </label>

        <div class="settings-page__extensions">
          <span class="settings-page__extensions-title">{{ t('settings.fileCategories.fields.allowedExtensions') }}</span>
          <div class="settings-page__extensions-grid">
            <label v-for="option in FILE_EXTENSION_OPTIONS" :key="option.value" class="settings-page__check">
              <input
                type="checkbox"
                :checked="activeDraft.allowed_extensions.includes(option.value)"
                @change="toggleDraftExtension(activeCategory.category_key, option.value, ($event.target as HTMLInputElement).checked)"
              >
              {{ option.label }}
            </label>
          </div>
        </div>

        <div class="settings-page__actions">
          <BaseButton :busy="categorySavingKey === activeCategory.category_key" @click="handleSaveCategory(activeCategory.category_key)">
            {{ t('settings.fileCategories.saveCategory') }}
          </BaseButton>
          <BaseButton variant="secondary" :busy="extensionSavingKey === activeCategory.category_key" @click="handleSaveExtensions(activeCategory.category_key)">
            {{ t('settings.fileCategories.saveExtensions') }}
          </BaseButton>
          <BaseButton
            variant="ghost"
            :busy="defaultSavingKey === activeCategory.category_key"
            :disabled="activeCategory.is_default"
            @click="handleSetDefaultCategory(activeCategory.category_key)"
          >
            {{ t('settings.fileCategories.setDefault') }}
          </BaseButton>
        </div>
      </div>
    </section>
  </section>
</template>

<style scoped>
.settings-page {
  display: grid;
  gap: 16px;
}

.settings-page__form,
.settings-page__section,
.settings-page__category-create,
.settings-page__category-editor {
  display: grid;
  gap: 16px;
}

.settings-page__section {
  border: 1px solid var(--bb-color-line);
  border-radius: 12px;
  padding: 20px;
  margin: 0;
  background: var(--bb-color-surface);
}

.settings-page__section-header,
.settings-page__category-meta,
.settings-page__actions {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.settings-page__section-title,
.settings-page__category-title {
  margin: 0;
  font-weight: 700;
}

.settings-page__section-desc,
.settings-page__category-key,
.settings-page__hint {
  margin: 0;
  color: var(--bb-color-muted);
  font-size: 0.88rem;
}

.settings-page__inline-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.settings-page__inline-grid--categories {
  grid-template-columns: minmax(220px, 360px);
}

.settings-page__checks,
.settings-page__extensions-grid {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.settings-page__check {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--bb-color-text);
}

.settings-page__extensions {
  display: grid;
  gap: 10px;
}

.settings-page__extensions-title {
  font-weight: 700;
  color: var(--bb-color-text-strong);
}

.settings-page__badge {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  padding: 4px 10px;
  color: var(--bb-color-primary);
  background: var(--bb-color-primary-soft);
  font-size: 0.82rem;
  font-weight: 700;
}

@media (max-width: 900px) {
  .settings-page__inline-grid {
    grid-template-columns: 1fr;
  }

  .settings-page__actions {
    align-items: stretch;
  }
}
</style>
