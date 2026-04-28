<script setup lang="ts">
import { computed } from 'vue'
import { EditorContent } from '@tiptap/vue-3'
import type { Editor } from '@tiptap/vue-3'

import ContentEditorBubbleMenu from './ContentEditorBubbleMenu.vue'
import type { ContentEditorSourceMode } from './useContentEditor'

const props = defineProps<{
  editor: Editor | null
  loading?: boolean
  sourceMode: ContentEditorSourceMode
  sourceContent: string
  wordCount: number
}>()

const emit = defineEmits<{
  sourceContent: [value: string]
}>()

const sourceLabel = computed(() => (props.sourceMode === 'html' ? 'HTML source editor' : 'Markdown source editor'))
</script>

<template>
  <section class="content-editor-canvas" aria-label="Writing canvas">
    <div v-if="loading || !editor" class="content-editor-canvas__loading">Loading editor...</div>
    <template v-else-if="sourceMode === 'visual'">
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
      <span>{{ wordCount }} words</span>
      <span>{{ sourceMode }}</span>
    </footer>
  </section>
</template>

<style scoped>
.content-editor-canvas {
  min-width: 0;
  min-height: 0;
  display: grid;
  grid-template-rows: minmax(0, 1fr) auto;
  overflow: hidden;
  background: var(--bb-color-subtle);
}

.content-editor-canvas__loading {
  display: grid;
  min-height: 420px;
  place-items: center;
  color: var(--bb-color-muted);
  font-weight: 700;
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
