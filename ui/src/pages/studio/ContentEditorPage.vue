<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ContentEditorCanvas from '@/features/studio/content-editor/ContentEditorCanvas.vue'
import ContentEditorHeader from '@/features/studio/content-editor/ContentEditorHeader.vue'
import ContentEditorMetaPanel from '@/features/studio/content-editor/ContentEditorMetaPanel.vue'
import ContentEditorToolbar from '@/features/studio/content-editor/ContentEditorToolbar.vue'
import { useContentEditor } from '@/features/studio/content-editor/useContentEditor'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import { useToast } from '@/shared/composables'

const route = useRoute()
const router = useRouter()
const { pushToast } = useToast()
const routeContentId = computed(() => {
  const value = route.params.content_id
  return typeof value === 'string' ? value : undefined
})

const editorState = useContentEditor(routeContentId.value)

const pageClasses = computed(() => ({
  'content-editor-page--focus': editorState.isFocusMode.value,
  'content-editor-page--sidebar-open': editorState.isSidebarOpen.value && !editorState.isFocusMode.value,
}))

onMounted(() => {
  void editorState.initialize()
})

async function saveContent(): Promise<void> {
  try {
    const result = await editorState.save()
    pushToast({ tone: 'success', title: result.wasCreated ? 'Draft created' : 'Content saved', message: result.content.title })
    if (result.wasCreated) {
      await router.replace(`/studio/content/${encodeURIComponent(result.content.content_id)}/edit`)
    }
  } catch (error) {
    pushToast({ tone: 'danger', title: 'Save failed', message: error instanceof Error ? error.message : 'Unable to save content.' })
  }
}

function backToContent(): void {
  void router.push('/studio/content')
}
</script>

<template>
  <main class="content-editor-page" :class="pageClasses" aria-label="Content editor">
    <ContentEditorHeader
      v-model:title="editorState.form.title"
      :mode="editorState.mode.value"
      :save-state="editorState.saveState.value"
      :saving="editorState.isSaving.value"
      :sidebar-open="editorState.isSidebarOpen.value"
      :focus-mode="editorState.isFocusMode.value"
      @back="backToContent"
      @save="saveContent"
      @toggle-sidebar="editorState.toggleSidebar"
      @toggle-focus="editorState.toggleFocusMode"
    />
    <StatusAlert v-if="editorState.errorMessage.value" class="content-editor-page__alert" tone="danger" title="Editor unavailable">
      {{ editorState.errorMessage.value }}
    </StatusAlert>
    <ContentEditorToolbar
      :editor="editorState.editor.value"
      :source-mode="editorState.sourceMode.value"
      @source-mode="editorState.setSourceMode"
    />
    <div class="content-editor-page__workspace">
      <ContentEditorCanvas
        :editor="editorState.editor.value"
        :loading="editorState.isLoading.value"
        :source-mode="editorState.sourceMode.value"
        :source-content="editorState.sourceContent.value"
        :word-count="editorState.wordCount.value"
        @source-content="editorState.setSourceContent"
      />
      <ContentEditorMetaPanel
        v-if="editorState.isSidebarOpen.value && !editorState.isFocusMode.value"
        v-model="editorState.form"
        :tags="editorState.tags.value"
      />
    </div>
  </main>
</template>

<style scoped>
.content-editor-page {
  position: fixed;
  z-index: 1000;
  inset: 0;
  display: grid;
  grid-template-rows: auto auto auto minmax(0, 1fr);
  min-width: 0;
  background: var(--bb-color-surface);
}

.content-editor-page__alert {
  border-radius: 0;
}

.content-editor-page__workspace {
  min-height: 0;
  display: grid;
  grid-template-columns: minmax(0, 1fr);
}

.content-editor-page--sidebar-open .content-editor-page__workspace {
  grid-template-columns: minmax(0, 1fr) minmax(300px, 360px);
}

.content-editor-page--focus :deep(.content-editor-toolbar) {
  display: none;
}

@media (max-width: 980px) {
  .content-editor-page {
    overflow: auto;
    grid-template-rows: auto auto auto auto;
  }

  .content-editor-page--sidebar-open .content-editor-page__workspace {
    grid-template-columns: 1fr;
  }
}
</style>
