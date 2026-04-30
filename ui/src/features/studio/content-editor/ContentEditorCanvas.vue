<script setup lang="ts">
import { computed } from 'vue'
import { EditorContent } from '@tiptap/vue-3'
import type { Editor } from '@tiptap/vue-3'
import { useI18n } from 'vue-i18n'

import ContentEditorBubbleMenu from './ContentEditorBubbleMenu.vue'
import type { ContentEditorSourceMode } from './useContentEditor'
import InlineLoadingState from '@/shared/components/InlineLoadingState.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'

const props = defineProps<{
  editor: Editor | null | undefined
  loading?: boolean
  refreshing?: boolean
  sourceMode: ContentEditorSourceMode
  sourceContent: string
  wordCount: number
}>()

const emit = defineEmits<{
  sourceContent: [value: string]
}>()

const { t } = useI18n()

const sourceLabel = computed(() => (props.sourceMode === 'html' ? t('editor.htmlSource') : t('editor.markdownSource')))
const sourceModeLabel = computed(() => t(`editor.modes.${props.sourceMode}`))
const wordCountLabel = computed(() => t('editor.wordCount', { count: props.wordCount }))
</script>

<template>
  <section class="content-editor-canvas" :aria-label="t('editor.writingCanvas')">
    <PageLoadingState v-if="loading || !editor" class="content-editor-canvas__loading" :title="t('editor.loadingEditor')" :rows="4" />
    <template v-else-if="sourceMode === 'visual'">
      <InlineLoadingState v-if="refreshing" class="content-editor-canvas__refreshing" />
      <ContentEditorBubbleMenu :editor="editor" />
      <EditorContent :editor="editor" />
    </template>
    <textarea
      v-else
      class="content-editor-canvas__source"
      :aria-label="sourceLabel"
      spellcheck="false"
      :value="sourceContent"
      @input="emit('sourceContent', ($event.target as HTMLTextAreaElement).value)"
    />
    <footer class="content-editor-canvas__footer">
      <span>{{ wordCountLabel }}</span>
      <span>{{ sourceModeLabel }}</span>
    </footer>
  </section>
</template>

<style scoped>
.content-editor-canvas {
  position: relative;
  min-width: 0;
  min-height: 0;
  display: grid;
  grid-template-rows: minmax(0, 1fr) auto;
  overflow: hidden;
  background: var(--bb-color-subtle);
}

.content-editor-canvas__loading {
  border: 0;
  border-radius: 0;
  box-shadow: none;
}

.content-editor-canvas__refreshing {
  position: absolute;
  top: 16px;
  right: 20px;
  z-index: 1;
  padding: 6px 10px;
  border: 1px solid var(--bb-color-line);
  border-radius: 999px;
  background: color-mix(in srgb, var(--bb-color-surface) 94%, transparent);
  box-shadow: var(--bb-shadow-soft);
}

:deep(.content-editor-canvas__surface) {
  width: min(840px, calc(100% - 48px));
  min-height: calc(100vh - 224px);
  margin: 28px auto;
  border: 1px solid var(--bb-color-line);
  border-radius: 10px;
  padding: 54px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-panel);
}

:deep(.content-editor-canvas__surface:focus) {
  outline: none;
}

:deep(.content-editor-canvas__surface p),
:deep(.content-editor-canvas__surface h1),
:deep(.content-editor-canvas__surface h2),
:deep(.content-editor-canvas__surface blockquote),
:deep(.content-editor-canvas__surface ul),
:deep(.content-editor-canvas__surface ol),
:deep(.content-editor-canvas__surface pre) {
  margin: 0 0 1rem;
}

:deep(.content-editor-canvas__surface h1) {
  font-size: 2.2rem;
  line-height: 1.2;
}

:deep(.content-editor-canvas__surface h2) {
  font-size: 1.55rem;
  line-height: 1.3;
}

:deep(.content-editor-canvas__surface blockquote) {
  border-left: 3px solid var(--bb-color-primary);
  padding-left: 14px;
  color: var(--bb-color-muted);
}

:deep(.content-editor-canvas__surface pre) {
  border-radius: 8px;
  padding: 14px;
  overflow: auto;
  color: var(--bb-color-surface);
  background: var(--bb-color-text-strong);
}

:deep(.content-editor-canvas__surface mark) {
  border-radius: 4px;
  padding: 0 2px;
  background: var(--bb-color-warning-soft);
}

:deep(.content-editor-canvas__surface p.is-editor-empty:first-child::before) {
  height: 0;
  float: left;
  color: var(--bb-color-muted);
  content: attr(data-placeholder);
  pointer-events: none;
}

.content-editor-canvas__source {
  width: min(940px, calc(100% - 48px));
  min-height: calc(100vh - 224px);
  margin: 28px auto;
  resize: none;
  border: 1px solid var(--bb-color-line);
  border-radius: 10px;
  padding: 28px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-panel);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.95rem;
  line-height: 1.7;
}

.content-editor-canvas__source:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus), var(--bb-shadow-panel);
}

.content-editor-canvas__footer {
  min-height: 36px;
  display: flex;
  justify-content: flex-end;
  gap: 14px;
  border-top: 1px solid var(--bb-color-line);
  padding: 8px 16px;
  color: var(--bb-color-muted);
  background: var(--bb-color-surface);
  font-size: 0.82rem;
  font-weight: 700;
  text-transform: capitalize;
}

@media (max-width: 760px) {
  :deep(.content-editor-canvas__surface),
  .content-editor-canvas__source {
    width: calc(100% - 20px);
    min-height: 360px;
    margin: 10px auto;
    padding: 22px;
  }
}
</style>
