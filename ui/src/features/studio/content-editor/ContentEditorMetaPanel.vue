<script setup lang="ts">
import FormField from '@/shared/components/FormField.vue'
import BaseInput from '@/shared/components/BaseInput.vue'

import type { ContentTag } from '../types'
import { contentStatuses, contentTypes, contentVisibilities, type ContentEditorForm } from './useContentEditor'

const form = defineModel<ContentEditorForm>({ required: true })

defineProps<{
  tags: ContentTag[]
}>()
</script>

<template>
  <aside class="content-editor-meta" aria-label="Content settings">
    <section class="content-editor-meta__section">
      <h2>Content settings</h2>
      <label class="content-editor-meta__select">
        <span>Type</span>
        <select v-model="form.type">
          <option v-for="type in contentTypes" :key="type" :value="type">{{ type }}</option>
        </select>
      </label>
      <label class="content-editor-meta__select">
        <span>Status</span>
        <select v-model="form.status">
          <option v-for="status in contentStatuses" :key="status" :value="status">{{ status }}</option>
        </select>
      </label>
      <label class="content-editor-meta__select">
        <span>Visibility</span>
        <select v-model="form.visibility">
          <option v-for="visibility in contentVisibilities" :key="visibility" :value="visibility">{{ visibility }}</option>
        </select>
      </label>
      <label class="content-editor-meta__select">
        <span>AI access</span>
        <select v-model="form.ai_access">
          <option value="denied">Denied</option>
          <option value="allowed">Allowed</option>
        </select>
      </label>
    </section>

    <section class="content-editor-meta__section">
      <h2>Metadata</h2>
      <FormField label="Slug" for-id="editor-slug">
        <BaseInput id="editor-slug" v-model="form.slug" />
      </FormField>
      <FormField label="Summary" for-id="editor-summary">
        <BaseInput id="editor-summary" v-model="form.summary" />
      </FormField>
      <FormField label="Cover URL" for-id="editor-cover">
        <BaseInput id="editor-cover" v-model="form.cover_image_url" />
      </FormField>
      <FormField label="Change summary" for-id="editor-change-summary">
        <BaseInput id="editor-change-summary" v-model="form.change_summary" />
      </FormField>
    </section>

    <section class="content-editor-meta__section">
      <h2>Options</h2>
      <label class="content-editor-meta__check">
        <input v-model="form.comment_enabled" type="checkbox" />
        Comments
      </label>
      <label class="content-editor-meta__check">
        <input v-model="form.is_featured" type="checkbox" />
        Featured
      </label>
      <FormField label="Sort order" for-id="editor-sort-order">
        <BaseInput id="editor-sort-order" v-model.number="form.sort_order" type="number" />
      </FormField>
    </section>

    <section class="content-editor-meta__section">
      <h2>Tags</h2>
      <div v-if="tags.length > 0" class="content-editor-meta__tags">
        <label v-for="tag in tags" :key="tag.tag_id">
          <input v-model="form.tag_ids" type="checkbox" :value="tag.tag_id" />
          {{ tag.name }}
        </label>
      </div>
      <p v-else class="content-editor-meta__empty">No tags are available.</p>
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

.content-editor-meta__select,
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

.content-editor-meta__select select {
  min-height: 40px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 0 10px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface-elevated);
}

.content-editor-meta__select select:focus-visible,
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
