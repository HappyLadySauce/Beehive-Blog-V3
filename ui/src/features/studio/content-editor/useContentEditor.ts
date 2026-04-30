import CharacterCount from '@tiptap/extension-character-count'
import Highlight from '@tiptap/extension-highlight'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import TextAlign from '@tiptap/extension-text-align'
import Underline from '@tiptap/extension-underline'
import { useQueryClient } from '@tanstack/vue-query'
import StarterKit from '@tiptap/starter-kit'
import { type Editor, useEditor } from '@tiptap/vue-3'
import { marked } from 'marked'
import TurndownService from 'turndown'
import { computed, onBeforeUnmount, reactive, shallowRef, watch } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio/api/studioApi'
import { useProgressiveQuery } from '@/shared/composables'
import { i18n } from '@/shared/i18n'
import type {
  ContentAIAccess,
  ContentDetail,
  ContentStatus,
  ContentTag,
  ContentType,
  ContentVisibility,
  ContentWriteRequest,
} from '@/features/studio/types'
import { useAvatarUpload } from '@/features/file-manager/useAvatarUpload'

export type ContentEditorMode = 'create' | 'edit'
export type ContentEditorSaveState = 'idle' | 'dirty' | 'saving' | 'saved' | 'error'
export type ContentEditorSourceMode = 'visual' | 'html' | 'markdown'

export interface ContentEditorForm {
  type: ContentType
  title: string
  slug: string
  summary: string
  body_markdown: string
  body_json: string
  cover_image_url: string
  status: ContentStatus
  visibility: ContentVisibility
  ai_access: ContentAIAccess
  comment_enabled: boolean
  is_featured: boolean
  sort_order: number
  change_summary: string
  tag_ids: string[]
}

export const contentTypes: ContentType[] = ['article', 'note', 'project', 'experience', 'timeline_event', 'insight', 'portfolio', 'page']
export const contentStatuses: ContentStatus[] = ['draft', 'review', 'published', 'archived']
export const contentVisibilities: ContentVisibility[] = ['public', 'member', 'private']

const emptyDocument = '{"type":"doc","content":[{"type":"paragraph"}]}'
const turndown = new TurndownService({ headingStyle: 'atx', codeBlockStyle: 'fenced' })

export function createEmptyContentEditorForm(): ContentEditorForm {
  return {
    type: 'article',
    title: '',
    slug: '',
    summary: '',
    body_markdown: '',
    body_json: emptyDocument,
    cover_image_url: '',
    status: 'draft',
    visibility: 'private',
    ai_access: 'denied',
    comment_enabled: true,
    is_featured: false,
    sort_order: 0,
    change_summary: '',
    tag_ids: [],
  }
}

export function useContentEditor(contentId?: string) {
  const authStore = useAuthStore()
  const queryClient = useQueryClient()
  const { uploadImage } = useAvatarUpload()
  const mode = shallowRef<ContentEditorMode>(contentId ? 'edit' : 'create')
  const sourceMode = shallowRef<ContentEditorSourceMode>('visual')
  const sourceContent = shallowRef('')
  const isSaving = shallowRef(false)
  const isSidebarOpen = shallowRef(true)
  const isFocusMode = shallowRef(false)
  const errorMessage = shallowRef('')
  const saveState = shallowRef<ContentEditorSaveState>('idle')
  const cleanSaveState = shallowRef<Extract<ContentEditorSaveState, 'idle' | 'saved'>>('idle')
  const cleanSnapshot = shallowRef('')
  const currentContentId = shallowRef(contentId ?? '')
  const tags = shallowRef<ContentTag[]>([])
  const hasInitialized = shallowRef(false)
  const canEditStatus = computed(() => mode.value !== 'create')
  const form = reactive(createEmptyContentEditorForm())
  let isHydrating = false
  let isSyncingSourceContent = false

  const tagsQuery = useProgressiveQuery({
    queryKey: ['editor-tags'],
    queryFn: () => studioApi.listTags({ page: 1, page_size: 100 }, { accessToken: authStore.accessToken }),
  })

  const contentQuery = useProgressiveQuery({
    queryKey: computed(() => ['editor-content', currentContentId.value]),
    queryFn: () => studioApi.getContent(currentContentId.value, { accessToken: authStore.accessToken }),
    enabled: computed(() => Boolean(currentContentId.value)),
  })

  const isLoading = computed(() => {
    if (mode.value === 'create') {
      return tagsQuery.showBlockingLoading.value && !hasInitialized.value
    }
    return (tagsQuery.showBlockingLoading.value || contentQuery.showBlockingLoading.value) && !hasInitialized.value
  })
  const isRefreshing = computed(() => tagsQuery.showRefreshingHint.value || contentQuery.showRefreshingHint.value)

  const editor = useEditor({
    extensions: [
      StarterKit,
      Link.configure({
        openOnClick: false,
        HTMLAttributes: { rel: 'noopener noreferrer', target: '_blank' },
      }),
      Image,
      Underline,
      Highlight,
      TextAlign.configure({ types: ['heading', 'paragraph'] }),
      Placeholder.configure({ placeholder: String(i18n.global.t('editor.placeholder')) }),
      CharacterCount,
    ],
    content: safeParseDocument(form.body_json),
    editorProps: {
      attributes: {
        class: 'content-editor-canvas__surface',
        'aria-label': String(i18n.global.t('editor.ariaLabel')),
      },
      handlePaste: (view, event) => {
        const files = getImageFiles(event.clipboardData)
        if (files.length === 0) {
          return false
        }
        event.preventDefault()
        void uploadAndInsertImages(editor.value, files, authStore.accessToken, uploadImage)
        return true
      },
      handleDrop: (view, event) => {
        const files = getImageFiles(event.dataTransfer)
        if (files.length === 0) {
          return false
        }
        event.preventDefault()
        void uploadAndInsertImages(editor.value, files, authStore.accessToken, uploadImage)
        return true
      },
    },
    onUpdate({ editor }) {
      if (sourceMode.value === 'visual') {
        syncEditorToForm(editor)
      }
    },
  })

  const wordCount = computed(() => {
    const count = editor.value?.storage.characterCount?.words?.()
    if (typeof count === 'number') {
      return count
    }
    return form.body_markdown.trim().split(/\s+/).filter(Boolean).length
  })

  watch(
    form,
    () => {
      if (isHydrating || saveState.value === 'saving') {
        return
      }
      refreshSaveState()
    },
    { deep: true },
  )

  watch(sourceContent, () => {
    if (isHydrating || isSyncingSourceContent || sourceMode.value === 'visual') {
      return
    }
    const nextMarkdown = sourceMode.value === 'markdown' ? sourceContent.value : htmlToMarkdown(sourceContent.value)
    if (nextMarkdown === form.body_markdown) {
      return
    }
    form.body_markdown = nextMarkdown
    refreshSaveState()
  })

  watch(
    () => tagsQuery.data.value,
    (value) => {
      if (!value) {
        return
      }
      tags.value = value.items
      if (mode.value === 'create' && !hasInitialized.value) {
        hydrate(createBlankContent())
        commitCleanSnapshot('idle')
        hasInitialized.value = true
      }
    },
    { immediate: true },
  )

  watch(
    () => contentQuery.data.value,
    (value) => {
      if (!value) {
        return
      }
      hydrate(value.content)
      commitCleanSnapshot('idle')
      hasInitialized.value = true
    },
    { immediate: true },
  )

  watch(
    () => [tagsQuery.error.value, contentQuery.error.value],
    ([tagsError, contentError]) => {
      const error = tagsError ?? contentError
      if (!error) {
        return
      }
      errorMessage.value = error instanceof Error ? error.message : String(i18n.global.t('editor.toast.unableToLoad'))
      saveState.value = 'error'
    },
    { immediate: true },
  )

  async function initialize(): Promise<void> {
    if (!tagsQuery.data.value) {
      await tagsQuery.refetch()
    }
    if (currentContentId.value && !contentQuery.data.value) {
      await contentQuery.refetch()
    }
  }

  async function save(): Promise<{ content: ContentDetail; wasCreated: boolean }> {
    syncSourceToVisual()
    const payload = buildPayload()
    const wasCreated = mode.value === 'create'
    isSaving.value = true
    saveState.value = 'saving'
    errorMessage.value = ''
    try {
      const response = wasCreated
        ? await studioApi.createContent(payload, { accessToken: authStore.accessToken })
        : await studioApi.updateContent(currentContentId.value, payload, { accessToken: authStore.accessToken })
      currentContentId.value = response.content.content_id
      mode.value = 'edit'
      hydrate(response.content)
      commitCleanSnapshot('saved')
      queryClient.setQueryData(['editor-content', response.content.content_id], response)
      await queryClient.invalidateQueries({ queryKey: ['studio-contents'] })
      return { content: response.content, wasCreated }
    } catch (error) {
      errorMessage.value = error instanceof Error ? error.message : String(i18n.global.t('editor.toast.unableToSave'))
      saveState.value = 'error'
      throw error
    } finally {
      isSaving.value = false
    }
  }

  function hydrate(content: ContentDetail): void {
    isHydrating = true
    Object.assign(form, {
      type: content.type,
      title: content.title,
      slug: content.slug,
      summary: content.summary ?? '',
      body_markdown: content.body_markdown,
      body_json: content.body_json || emptyDocument,
      cover_image_url: content.cover_image_url ?? '',
      status: content.status,
      visibility: content.visibility,
      ai_access: content.ai_access,
      comment_enabled: content.comment_enabled,
      is_featured: content.is_featured,
      sort_order: content.sort_order,
      change_summary: '',
      tag_ids: content.tags.map((tag) => tag.tag_id),
    })
    editor.value?.commands.setContent(safeParseDocument(form.body_json), false)
    sourceMode.value = 'visual'
    sourceContent.value = editor.value?.getHTML() ?? ''
    queueMicrotask(() => {
      isHydrating = false
    })
  }

  function setSourceMode(nextMode: ContentEditorSourceMode): void {
    if (sourceMode.value === nextMode) {
      return
    }
    if (sourceMode.value !== 'visual') {
      syncSourceToVisual()
    }
    sourceMode.value = nextMode
    isSyncingSourceContent = true
    if (nextMode === 'html') {
      sourceContent.value = editor.value?.getHTML() ?? markdownToHtml(form.body_markdown)
    } else if (nextMode === 'markdown') {
      sourceContent.value = editor.value ? htmlToMarkdown(editor.value.getHTML()) : form.body_markdown
    } else {
      sourceContent.value = ''
    }
    queueMicrotask(() => {
      isSyncingSourceContent = false
    })
  }

  function setSourceContent(value: string): void {
    sourceContent.value = value
  }

  function toggleSidebar(): void {
    isSidebarOpen.value = !isSidebarOpen.value
  }

  function toggleFocusMode(): void {
    isFocusMode.value = !isFocusMode.value
    if (isFocusMode.value) {
      isSidebarOpen.value = false
    }
  }

  function syncSourceToVisual(): void {
    if (!editor.value || sourceMode.value === 'visual') {
      return
    }
    const html = sourceMode.value === 'markdown' ? markdownToHtml(sourceContent.value) : sourceContent.value
    editor.value.commands.setContent(html, false)
    syncEditorToForm(editor.value)
  }

  function syncEditorToForm(editorInstance: Editor): void {
    form.body_json = JSON.stringify(editorInstance.getJSON())
    form.body_markdown = htmlToMarkdown(editorInstance.getHTML()) || editorInstance.getText()
  }

  function buildPayload(): ContentWriteRequest {
    const title = form.title.trim()
    if (!title) {
      throw new Error(String(i18n.global.t('editor.validation.titleRequired')))
    }
    const payload: ContentWriteRequest = {
      type: form.type,
      title,
      slug: form.slug.trim() || slugFromTitle(title),
      summary: form.summary.trim(),
      body_markdown: form.body_markdown,
      body_json: form.body_json,
      cover_image_url: form.cover_image_url.trim(),
      visibility: form.visibility,
      ai_access: form.ai_access,
      source_type: 'manual',
      comment_enabled: form.comment_enabled,
      is_featured: form.is_featured,
      sort_order: Number(form.sort_order) || 0,
      tag_ids: [...form.tag_ids],
      change_summary: form.change_summary.trim()
        || (mode.value === 'create'
          ? String(i18n.global.t('editor.changeSummaryDefaultCreate'))
          : String(i18n.global.t('editor.changeSummaryDefaultUpdate'))),
    }
    if (canEditStatus.value) {
      payload.status = form.status
    }
    return payload
  }

  onBeforeUnmount(() => {
    editor.value?.destroy()
  })

  return {
    mode,
    sourceMode,
    sourceContent,
    currentContentId,
    isLoading,
    isRefreshing,
    isSaving,
    isSidebarOpen,
    isFocusMode,
    errorMessage,
    saveState,
    wordCount,
    canEditStatus,
    tags,
    form,
    editor,
    initialize,
    save,
    setSourceMode,
    setSourceContent,
    toggleSidebar,
    toggleFocusMode,
  }

  function commitCleanSnapshot(nextState: Extract<ContentEditorSaveState, 'idle' | 'saved'>): void {
    cleanSnapshot.value = snapshotForm(form)
    cleanSaveState.value = nextState
    saveState.value = nextState
  }

  function refreshSaveState(): void {
    saveState.value = snapshotForm(form) === cleanSnapshot.value ? cleanSaveState.value : 'dirty'
  }
}

function createBlankContent(): ContentDetail {
  const now = Math.floor(Date.now() / 1000)
  return {
    content_id: '',
    type: 'article',
    title: '',
    slug: '',
    summary: '',
    body_markdown: '',
    body_json: emptyDocument,
    cover_image_url: '',
    status: 'draft',
    visibility: 'private',
    ai_access: 'denied',
    owner_user_id: '',
    author_user_id: '',
    source_type: 'manual',
    current_revision_id: '',
    comment_enabled: true,
    is_featured: false,
    sort_order: 0,
    created_at: now,
    updated_at: now,
    tags: [],
  }
}

function safeParseDocument(value: string): Record<string, unknown> {
  if (!value.trim()) {
    return { type: 'doc', content: [{ type: 'paragraph' }] }
  }
  try {
    return JSON.parse(value) as Record<string, unknown>
  } catch {
    return { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: value }] }] }
  }
}

function markdownToHtml(markdown: string): string {
  return marked.parse(markdown, { async: false }) as string
}

function htmlToMarkdown(html: string): string {
  return turndown.turndown(html).trim()
}

function slugFromTitle(title: string): string {
  return title
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-')
    .replace(/^-+|-+$/g, '')
    || `content-${Date.now()}`
}

function snapshotForm(form: ContentEditorForm): string {
  return JSON.stringify({
    type: form.type,
    title: form.title,
    slug: form.slug,
    summary: form.summary,
    body_markdown: form.body_markdown,
    body_json: form.body_json,
    cover_image_url: form.cover_image_url,
    status: form.status,
    visibility: form.visibility,
    ai_access: form.ai_access,
    comment_enabled: form.comment_enabled,
    is_featured: form.is_featured,
    sort_order: form.sort_order,
    change_summary: form.change_summary,
    tag_ids: form.tag_ids,
  })
}

function getImageFiles(dataTransfer: DataTransfer | null): File[] {
  if (!dataTransfer) return []
  return Array.from(dataTransfer.files).filter((file) => file.type.startsWith('image/'))
}

async function uploadAndInsertImages(
  editor: Editor | null,
  files: File[],
  accessToken: string | undefined,
  uploader: (file: File, accessToken: string | undefined, scope: 'content_image') => Promise<string>,
): Promise<void> {
  if (!editor || editor.isDestroyed) return

  for (const file of files) {
    const blobUrl = URL.createObjectURL(file)
    editor.chain().focus().setImage({ src: blobUrl, alt: file.name }).run()
    try {
      const remoteUrl = await uploader(file, accessToken, 'content_image')
      const html = editor.getHTML()
      if (html.includes(blobUrl)) {
        editor.commands.setContent(html.replaceAll(blobUrl, remoteUrl), false)
      }
    } finally {
      window.setTimeout(() => URL.revokeObjectURL(blobUrl), 1_500)
    }
  }
}
