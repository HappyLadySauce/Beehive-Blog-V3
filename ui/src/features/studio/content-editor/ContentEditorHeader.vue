<script setup lang="ts">
import { ArrowLeft, CheckCircle2, Columns3, EyeOff, Loader2, Save } from 'lucide-vue-next'

import BaseButton from '@/shared/components/BaseButton.vue'

import type { ContentEditorMode, ContentEditorSaveState } from './useContentEditor'

const title = defineModel<string>('title', { required: true })

const props = defineProps<{
  mode: ContentEditorMode
  saveState: ContentEditorSaveState
  saving: boolean
  sidebarOpen: boolean
  focusMode: boolean
}>()

const emit = defineEmits<{
  back: []
  save: []
  toggleSidebar: []
  toggleFocus: []
}>()

function saveLabel(): string {
  if (props.saving) {
    return 'Saving'
  }
  return props.mode === 'create' ? 'Create draft' : 'Save content'
}
</script>

<template>
  <header class="content-editor-header">
    <button class="content-editor-header__icon" type="button" aria-label="Back to content list" @click="emit('back')">
      <ArrowLeft :size="18" aria-hidden="true" />
    </button>
    <div class="content-editor-header__title">
      <input
        id="editor-title"
        v-model="title"
        class="content-editor-header__input"
        type="text"
        placeholder="Untitled"
        autocomplete="off"
        aria-label="Content title"
      />
      <span class="content-editor-header__status" :class="`content-editor-header__status--${saveState}`">
        <Loader2 v-if="saveState === 'saving'" :size="14" aria-hidden="true" />
        <CheckCircle2 v-else-if="saveState === 'saved'" :size="14" aria-hidden="true" />
        {{ saveState }}
      </span>
    </div>
    <div class="content-editor-header__actions">
      <button
        class="content-editor-header__icon"
        type="button"
        :aria-label="focusMode ? 'Exit focus mode' : 'Enter focus mode'"
        :title="focusMode ? 'Exit focus mode' : 'Focus mode'"
        :class="{ active: focusMode }"
        @click="emit('toggleFocus')"
      >
        <EyeOff :size="18" aria-hidden="true" />
      </button>
      <button
        class="content-editor-header__icon"
        type="button"
        :aria-label="sidebarOpen ? 'Hide settings panel' : 'Show settings panel'"
        :title="sidebarOpen ? 'Hide settings panel' : 'Show settings panel'"
        :class="{ active: sidebarOpen }"
        @click="emit('toggleSidebar')"
      >
        <Columns3 :size="18" aria-hidden="true" />
      </button>
      <BaseButton :busy="saving" @click="emit('save')">
        <Save :size="16" aria-hidden="true" />
        {{ saveLabel() }}
      </BaseButton>
    </div>
  </header>
</template>

<style scoped>
.content-editor-header {
  min-height: 68px;
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: 14px;
  border-bottom: 1px solid var(--bb-color-line);
  padding: 10px 18px;
  background: var(--bb-color-surface);
}

.content-editor-header__icon {
  width: 40px;
  height: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  color: var(--bb-color-muted);
  background: var(--bb-color-surface-elevated);
}

.content-editor-header__icon:hover,
.content-editor-header__icon:focus-visible,
.content-editor-header__icon.active {
  outline: none;
  color: var(--bb-color-text);
  border-color: var(--bb-color-primary);
  background: var(--bb-color-primary-soft);
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.content-editor-header__title {
  min-width: 0;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
}

.content-editor-header__input {
  min-width: 0;
  border: 0;
  color: var(--bb-color-text-strong);
  background: transparent;
  font-size: clamp(1.15rem, 1.8vw, 1.55rem);
  font-weight: 820;
}

.content-editor-header__input:focus-visible {
  outline: none;
}

.content-editor-header__status {
  min-height: 26px;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  border: 1px solid var(--bb-color-line);
  border-radius: 999px;
  padding: 0 10px;
  color: var(--bb-color-muted);
  font-size: 0.78rem;
  font-weight: 800;
  text-transform: capitalize;
}

.content-editor-header__status--saved {
  color: var(--bb-color-success);
  border-color: var(--bb-color-success);
  background: var(--bb-color-success-soft);
}

.content-editor-header__status--dirty {
  color: var(--bb-color-warning);
  border-color: var(--bb-color-warning);
  background: var(--bb-color-warning-soft);
}

.content-editor-header__status--error {
  color: var(--bb-color-danger);
  border-color: var(--bb-color-danger);
  background: var(--bb-color-danger-soft);
}

.content-editor-header__actions {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

@media (max-width: 780px) {
  .content-editor-header {
    grid-template-columns: auto 1fr;
  }

  .content-editor-header__actions {
    grid-column: 1 / -1;
    justify-content: flex-end;
  }
}

@media (max-width: 560px) {
  .content-editor-header__title {
    grid-template-columns: 1fr;
  }

  .content-editor-header__actions :deep(.bb-button) {
    flex: 1;
  }
}
</style>
