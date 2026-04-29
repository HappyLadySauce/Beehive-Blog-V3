<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSelect, { type BaseSelectOption } from '@/shared/components/BaseSelect.vue'
import FormField from '@/shared/components/FormField.vue'
import ImageUploader from '@/shared/components/ImageUploader.vue'
import ReadonlyField from '@/shared/components/ReadonlyField.vue'

import type { ContentTag } from '../types'
import { contentStatuses, contentTypes, contentVisibilities, type ContentEditorForm } from './useContentEditor'

const form = defineModel<ContentEditorForm>({ required: true })
const { t } = useI18n()

defineProps<{
  canEditStatus: boolean
  tags: ContentTag[]
}>()

const typeOptions = computed<BaseSelectOption[]>(() => contentTypes.map((type) => ({ value: type, label: t(`contentType.${type}`) })))
const statusOptions = computed<BaseSelectOption[]>(() => contentStatuses.map((status) => ({ value: status, label: t(`contentStatus.${status}`) })))
const visibilityOptions = computed<BaseSelectOption[]>(() => contentVisibilities.map((visibility) => ({ value: visibility, label: t(`visibility.${visibility}`) })))
const aiAccessOptions = computed<BaseSelectOption[]>(() => [
  { value: 'denied', label: t('aiAccess.denied') },
  { value: 'allowed', label: t('aiAccess.allowed') },
])
</script>

<template>
  <aside class="content-editor-meta" :aria-label="t('editor.settings')">
    <section class="content-editor-meta__section">
      <h2>{{ t('editor.settings') }}</h2>
      <FormField :label="t('editor.type')" for-id="editor-type">
        <BaseSelect id="editor-type" v-model="form.type" :options="typeOptions" :aria-label="t('editor.type')" />
      </FormField>
      <FormField v-if="canEditStatus" :label="t('editor.status')" for-id="editor-status">
        <BaseSelect id="editor-status" v-model="form.status" :options="statusOptions" :aria-label="t('editor.status')" />
      </FormField>
      <ReadonlyField v-else :label="t('editor.status')" :value="t('contentStatus.draft')" />
      <FormField :label="t('editor.visibility')" for-id="editor-visibility">
        <BaseSelect id="editor-visibility" v-model="form.visibility" :options="visibilityOptions" :aria-label="t('editor.visibility')" />
      </FormField>
      <FormField :label="t('editor.aiAccess')" for-id="editor-ai-access">
        <BaseSelect id="editor-ai-access" v-model="form.ai_access" :options="aiAccessOptions" :aria-label="t('editor.aiAccess')" />
      </FormField>
    </section>

    <section class="content-editor-meta__section">
      <h2>{{ t('editor.metadata') }}</h2>
      <FormField :label="t('editor.slug')" for-id="editor-slug">
        <BaseInput id="editor-slug" v-model="form.slug" />
      </FormField>
      <FormField :label="t('editor.summary')" for-id="editor-summary">
        <BaseInput id="editor-summary" v-model="form.summary" />
      </FormField>
      <div class="content-editor-meta__field">
        <span>{{ t('editor.coverImage') }}</span>
        <ImageUploader v-model="form.cover_image_url" scope="content_cover" :label="t('editor.uploadCover')" />
      </div>
      <FormField :label="t('editor.changeSummary')" for-id="editor-change-summary">
        <BaseInput id="editor-change-summary" v-model="form.change_summary" />
      </FormField>
    </section>

    <section class="content-editor-meta__section">
      <h2>{{ t('editor.options') }}</h2>
      <label class="content-editor-meta__check">
        <input v-model="form.comment_enabled" type="checkbox" />
        {{ t('editor.comments') }}
      </label>
      <label class="content-editor-meta__check">
        <input v-model="form.is_featured" type="checkbox" />
        {{ t('editor.featured') }}
      </label>
      <FormField :label="t('editor.sortOrder')" for-id="editor-sort-order">
        <BaseInput id="editor-sort-order" v-model.number="form.sort_order" type="number" />
      </FormField>
    </section>

    <section class="content-editor-meta__section">
      <h2>{{ t('editor.tags') }}</h2>
      <div v-if="tags.length > 0" class="content-editor-meta__tags">
        <label v-for="tag in tags" :key="tag.tag_id">
          <input v-model="form.tag_ids" type="checkbox" :value="tag.tag_id" />
          {{ tag.name }}
        </label>
      </div>
      <p v-else class="content-editor-meta__empty">{{ t('editor.noTags') }}</p>
    </section>
  </aside>
</template>

<style scoped>
.content-editor-meta {
  overflow: auto;
  border-left: 1px solid var(--bb-color-line);
  background: var(--bb-color-surface);
}

.content-editor-meta__section {
  display: grid;
  gap: 12px;
  border-bottom: 1px solid var(--bb-color-line);
  padding: 16px;
}

.content-editor-meta__section h2,
.content-editor-meta__empty {
  margin: 0;
}

.content-editor-meta__section h2 {
  color: var(--bb-color-text-strong);
  font-size: 0.95rem;
}

.content-editor-meta__field,
.content-editor-meta__check,
.content-editor-meta__tags label {
  display: grid;
  gap: 6px;
  color: var(--bb-color-muted);
  font-size: 0.9rem;
  font-weight: 700;
}

.content-editor-meta__check,
.content-editor-meta__tags label {
  grid-template-columns: auto 1fr;
  align-items: center;
}

.content-editor-meta__check input:focus-visible,
.content-editor-meta__tags input:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.content-editor-meta__tags {
  display: grid;
  gap: 8px;
}

.content-editor-meta__empty {
  color: var(--bb-color-muted);
}

@media (max-width: 980px) {
  .content-editor-meta {
    border-top: 1px solid var(--bb-color-line);
    border-left: 0;
  }
}
</style>
