<script setup lang="ts">
import { EditorContent, useEditor } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { Bold, Heading1, Heading2, Italic, List, ListOrdered, Pilcrow, Quote } from 'lucide-vue-next'
import { onBeforeUnmount, watch } from 'vue'

const model = defineModel<string>({ default: '' })
const plainText = defineModel<string>('plainText', { default: '' })

const editor = useEditor({
  extensions: [StarterKit],
  content: safeParse(model.value),
  editorProps: {
    attributes: {
      class: 'rich-editor__surface',
      'aria-label': 'Content body editor',
    },
  },
  onUpdate({ editor }) {
    model.value = JSON.stringify(editor.getJSON())
    plainText.value = editor.getText()
  },
})

watch(
  () => model.value,
  (value) => {
    if (!editor.value) {
      return
    }
    const next = safeParse(value)
    if (JSON.stringify(editor.value.getJSON()) !== JSON.stringify(next)) {
      editor.value.commands.setContent(next, false)
    }
  },
)

onBeforeUnmount(() => {
  editor.value?.destroy()
})

function safeParse(value: string): Record<string, unknown> {
  if (!value.trim()) {
    return { type: 'doc', content: [{ type: 'paragraph' }] }
  }
  try {
    return JSON.parse(value) as Record<string, unknown>
  } catch {
    return { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: value }] }] }
  }
}
</script>

<template>
  <div class="rich-editor">
    <div v-if="editor" class="rich-editor__toolbar" aria-label="Editor toolbar">
      <button type="button" :class="{ active: editor.isActive('paragraph') }" aria-label="Paragraph" @click="editor.chain().focus().setParagraph().run()">
        <Pilcrow :size="16" aria-hidden="true" />
      </button>
      <button type="button" :class="{ active: editor.isActive('heading', { level: 1 }) }" aria-label="Heading 1" @click="editor.chain().focus().toggleHeading({ level: 1 }).run()">
        <Heading1 :size="16" aria-hidden="true" />
      </button>
      <button type="button" :class="{ active: editor.isActive('heading', { level: 2 }) }" aria-label="Heading 2" @click="editor.chain().focus().toggleHeading({ level: 2 }).run()">
        <Heading2 :size="16" aria-hidden="true" />
      </button>
      <button type="button" :class="{ active: editor.isActive('bold') }" aria-label="Bold" @click="editor.chain().focus().toggleBold().run()">
        <Bold :size="16" aria-hidden="true" />
      </button>
      <button type="button" :class="{ active: editor.isActive('italic') }" aria-label="Italic" @click="editor.chain().focus().toggleItalic().run()">
        <Italic :size="16" aria-hidden="true" />
      </button>
      <button type="button" :class="{ active: editor.isActive('bulletList') }" aria-label="Bullet list" @click="editor.chain().focus().toggleBulletList().run()">
        <List :size="16" aria-hidden="true" />
      </button>
      <button type="button" :class="{ active: editor.isActive('orderedList') }" aria-label="Ordered list" @click="editor.chain().focus().toggleOrderedList().run()">
        <ListOrdered :size="16" aria-hidden="true" />
      </button>
      <button type="button" :class="{ active: editor.isActive('blockquote') }" aria-label="Quote" @click="editor.chain().focus().toggleBlockquote().run()">
        <Quote :size="16" aria-hidden="true" />
      </button>
    </div>
    <EditorContent :editor="editor" />
  </div>
</template>

<style scoped>
.rich-editor {
  overflow: hidden;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  background: var(--bb-color-surface);
}

.rich-editor__toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  border-bottom: 1px solid var(--bb-color-line);
  padding: 8px;
  background: var(--bb-color-subtle);
}

.rich-editor__toolbar button {
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

.rich-editor__toolbar button:hover,
.rich-editor__toolbar button.active {
  color: var(--bb-color-primary);
  border-color: var(--bb-color-primary);
  background: var(--bb-color-primary-soft);
}

.rich-editor__toolbar button:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

:deep(.rich-editor__surface) {
  min-height: 240px;
  padding: 16px;
  color: var(--bb-color-text);
}

:deep(.rich-editor__surface:focus) {
  outline: none;
}

:deep(.rich-editor__surface p),
:deep(.rich-editor__surface h1),
:deep(.rich-editor__surface h2) {
  margin: 0 0 0.8rem;
}
</style>
