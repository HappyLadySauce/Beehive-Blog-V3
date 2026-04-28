<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3'
import { Bold, Highlighter, Italic, Link2, Underline } from 'lucide-vue-next'
import { computed, onBeforeUnmount, onMounted, shallowRef } from 'vue'

const props = defineProps<{
  editor: Editor
}>()

const visible = shallowRef(false)
const position = shallowRef({ left: 0, top: 0 })

const bubbleStyle = computed(() => ({
  left: `${position.value.left}px`,
  top: `${position.value.top}px`,
}))

function updatePosition(): void {
  const { state, view } = props.editor
  const { from, to, empty } = state.selection
  if (empty || !props.editor.isFocused) {
    visible.value = false
    return
  }
  const start = view.coordsAtPos(from)
  const end = view.coordsAtPos(to)
  position.value = {
    left: Math.max(12, (start.left + end.right) / 2),
    top: Math.max(12, start.top - 52),
  }
  visible.value = true
}

function hideLater(): void {
  window.setTimeout(() => {
    if (!props.editor.isFocused) {
      visible.value = false
    }
  }, 120)
}

function updateLink(): void {
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

onMounted(() => {
  props.editor.on('selectionUpdate', updatePosition)
  props.editor.on('focus', updatePosition)
  props.editor.on('blur', hideLater)
})

onBeforeUnmount(() => {
  props.editor.off('selectionUpdate', updatePosition)
  props.editor.off('focus', updatePosition)
  props.editor.off('blur', hideLater)
})
</script>

<template>
  <div v-show="visible" class="content-editor-bubble" :style="bubbleStyle" role="toolbar" aria-label="Selection tools">
    <button type="button" :class="{ active: editor.isActive('bold') }" aria-label="Bold" @click="editor.chain().focus().toggleBold().run()">
      <Bold :size="15" aria-hidden="true" />
    </button>
    <button type="button" :class="{ active: editor.isActive('italic') }" aria-label="Italic" @click="editor.chain().focus().toggleItalic().run()">
      <Italic :size="15" aria-hidden="true" />
    </button>
    <button type="button" :class="{ active: editor.isActive('underline') }" aria-label="Underline" @click="editor.chain().focus().toggleUnderline().run()">
      <Underline :size="15" aria-hidden="true" />
    </button>
    <button type="button" :class="{ active: editor.isActive('highlight') }" aria-label="Highlight" @click="editor.chain().focus().toggleHighlight().run()">
      <Highlighter :size="15" aria-hidden="true" />
    </button>
    <button type="button" :class="{ active: editor.isActive('link') }" aria-label="Link" @click="updateLink">
      <Link2 :size="15" aria-hidden="true" />
    </button>
  </div>
</template>

<style scoped>
.content-editor-bubble {
  position: fixed;
  z-index: 1300;
  display: inline-flex;
  gap: 2px;
  border: 1px solid var(--bb-color-line);
  border-radius: 10px;
  padding: 4px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-panel);
  transform: translateX(-50%);
}

.content-editor-bubble button {
  width: 30px;
  height: 30px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid transparent;
  border-radius: 7px;
  color: var(--bb-color-muted);
  background: transparent;
}

.content-editor-bubble button:hover,
.content-editor-bubble button:focus-visible,
.content-editor-bubble button.active {
  outline: none;
  color: var(--bb-color-primary);
  border-color: var(--bb-color-primary);
  background: var(--bb-color-primary-soft);
}
</style>
