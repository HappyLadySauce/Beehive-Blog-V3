<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3'
import {
  AlignCenter,
  AlignLeft,
  AlignRight,
  Bold,
  Code2,
  Heading1,
  Heading2,
  Highlighter,
  Italic,
  Link2,
  List,
  ListOrdered,
  Pilcrow,
  Quote,
  Redo2,
  Underline,
  Undo2,
} from 'lucide-vue-next'

import type { ContentEditorSourceMode } from './useContentEditor'

const props = defineProps<{
  editor: Editor | null
  sourceMode: ContentEditorSourceMode
}>()

const emit = defineEmits<{
  sourceMode: [mode: ContentEditorSourceMode]
}>()

const sourceModes: Array<{ value: ContentEditorSourceMode; label: string }> = [
  { value: 'visual', label: 'Visual' },
  { value: 'html', label: 'HTML' },
  { value: 'markdown', label: 'Markdown' },
]

function updateLink(): void {
  if (!props.editor) {
    return
  }
  const current = props.editor.getAttributes('link').href as string | undefined
  const next = window.prompt('Link URL', current ?? '')
  if (next === null) {
    return
  }
  if (next.trim() === '') {
    props.editor.chain().focus().unsetLink().run()
    return
  }
  props.editor.chain().focus().extendMarkRange('link').setLink({ href: next.trim() }).run()
}
</script>

<template>
  <div class="content-editor-toolbar" aria-label="Editor toolbar">
    <div class="content-editor-toolbar__modes" role="group" aria-label="Editor mode">
      <button
        v-for="mode in sourceModes"
        :key="mode.value"
        class="content-editor-toolbar__mode"
        type="button"
        :class="{ active: sourceMode === mode.value }"
        @click="emit('sourceMode', mode.value)"
      >
        {{ mode.label }}
      </button>
    </div>

    <div class="content-editor-toolbar__tools" role="group" aria-label="Format tools">
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('paragraph') }" aria-label="Paragraph" @click="editor?.chain().focus().setParagraph().run()">
        <Pilcrow :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('heading', { level: 1 }) }" aria-label="Heading 1" @click="editor?.chain().focus().toggleHeading({ level: 1 }).run()">
        <Heading1 :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('heading', { level: 2 }) }" aria-label="Heading 2" @click="editor?.chain().focus().toggleHeading({ level: 2 }).run()">
        <Heading2 :size="16" aria-hidden="true" />
      </button>
      <span class="content-editor-toolbar__divider" aria-hidden="true" />
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('bold') }" aria-label="Bold" @click="editor?.chain().focus().toggleBold().run()">
        <Bold :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('italic') }" aria-label="Italic" @click="editor?.chain().focus().toggleItalic().run()">
        <Italic :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('underline') }" aria-label="Underline" @click="editor?.chain().focus().toggleUnderline().run()">
        <Underline :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('highlight') }" aria-label="Highlight" @click="editor?.chain().focus().toggleHighlight().run()">
        <Highlighter :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('link') }" aria-label="Link" @click="updateLink">
        <Link2 :size="16" aria-hidden="true" />
      </button>
      <span class="content-editor-toolbar__divider" aria-hidden="true" />
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive({ textAlign: 'left' }) }" aria-label="Align left" @click="editor?.chain().focus().setTextAlign('left').run()">
        <AlignLeft :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive({ textAlign: 'center' }) }" aria-label="Align center" @click="editor?.chain().focus().setTextAlign('center').run()">
        <AlignCenter :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive({ textAlign: 'right' }) }" aria-label="Align right" @click="editor?.chain().focus().setTextAlign('right').run()">
        <AlignRight :size="16" aria-hidden="true" />
      </button>
      <span class="content-editor-toolbar__divider" aria-hidden="true" />
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('bulletList') }" aria-label="Bullet list" @click="editor?.chain().focus().toggleBulletList().run()">
        <List :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('orderedList') }" aria-label="Ordered list" @click="editor?.chain().focus().toggleOrderedList().run()">
        <ListOrdered :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('blockquote') }" aria-label="Quote" @click="editor?.chain().focus().toggleBlockquote().run()">
        <Quote :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual'" :class="{ active: editor?.isActive('codeBlock') }" aria-label="Code block" @click="editor?.chain().focus().toggleCodeBlock().run()">
        <Code2 :size="16" aria-hidden="true" />
      </button>
      <span class="content-editor-toolbar__divider" aria-hidden="true" />
      <button type="button" :disabled="!editor || sourceMode !== 'visual' || !editor.can().undo()" aria-label="Undo" @click="editor?.chain().focus().undo().run()">
        <Undo2 :size="16" aria-hidden="true" />
      </button>
      <button type="button" :disabled="!editor || sourceMode !== 'visual' || !editor.can().redo()" aria-label="Redo" @click="editor?.chain().focus().redo().run()">
        <Redo2 :size="16" aria-hidden="true" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.content-editor-toolbar {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid var(--bb-color-line);
  padding: 8px 16px;
  background: var(--bb-color-surface);
}

.content-editor-toolbar__modes,
.content-editor-toolbar__tools {
  display: inline-flex;
  align-items: center;
}

.content-editor-toolbar__modes {
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 2px;
  background: var(--bb-color-subtle);
}

.content-editor-toolbar__mode {
  min-height: 30px;
  border: 0;
  border-radius: 6px;
  padding: 0 10px;
  color: var(--bb-color-muted);
  background: transparent;
  font-size: 0.82rem;
  font-weight: 800;
}

.content-editor-toolbar__mode.active {
  color: var(--bb-color-text);
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-soft);
}

.content-editor-toolbar__tools {
  min-width: 0;
  flex-wrap: wrap;
  gap: 4px;
}

.content-editor-toolbar__tools button {
  width: 34px;
  height: 34px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid transparent;
  border-radius: 8px;
  color: var(--bb-color-muted);
  background: transparent;
}

.content-editor-toolbar__tools button:hover:not(:disabled),
.content-editor-toolbar__tools button.active {
  color: var(--bb-color-primary);
  border-color: var(--bb-color-primary);
  background: var(--bb-color-primary-soft);
}

.content-editor-toolbar__mode:focus-visible,
.content-editor-toolbar__tools button:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.content-editor-toolbar__tools button:disabled {
  opacity: 0.36;
}

.content-editor-toolbar__divider {
  width: 1px;
  height: 22px;
  margin: 0 4px;
  background: var(--bb-color-line);
}

@media (max-width: 980px) {
  .content-editor-toolbar {
    grid-template-columns: 1fr;
  }
}
</style>
