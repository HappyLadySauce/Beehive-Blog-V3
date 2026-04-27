<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, shallowRef, watch } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import type {
  ContentDetail,
  ContentRelation,
  ContentRelationType,
  ContentRevisionSummary,
  ContentStatus,
  ContentSummary,
  ContentTag,
  ContentType,
  ContentVisibility,
} from '@/features/studio'
import ActionTagButton from '@/shared/components/ActionTagButton.vue'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import FormField from '@/shared/components/FormField.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import ReadonlyField from '@/shared/components/ReadonlyField.vue'
import RichContentEditor from '@/shared/components/RichContentEditor.vue'
import SideDrawer from '@/shared/components/SideDrawer.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import StatusBadge from '@/shared/components/StatusBadge.vue'
import { useConfirm, useToast } from '@/shared/composables'

type ContentMode = 'view' | 'edit' | 'create'
type StudioTab = 'content' | 'tags'

const contentTypes: ContentType[] = ['article', 'note', 'project', 'experience', 'timeline_event', 'insight', 'portfolio', 'page']
const statuses: ContentStatus[] = ['draft', 'review', 'published', 'archived']
const visibilities: ContentVisibility[] = ['public', 'member', 'private']
const relationTypes: ContentRelationType[] = ['belongs_to', 'related_to', 'derived_from', 'references', 'part_of', 'depends_on', 'timeline_of']

const authStore = useAuthStore()
const { confirm } = useConfirm()
const { pushToast } = useToast()

const activeTab = shallowRef<StudioTab>('content')
const contents = shallowRef<ContentSummary[]>([])
const tags = shallowRef<ContentTag[]>([])
const relations = shallowRef<ContentRelation[]>([])
const revisions = shallowRef<ContentRevisionSummary[]>([])
const selectedContent = shallowRef<ContentDetail | null>(null)
const selectedTag = shallowRef<ContentTag | null>(null)
const contentMode = shallowRef<ContentMode>('view')
const isLoading = shallowRef(true)
const isTagsLoading = shallowRef(false)
const isDetailLoading = shallowRef(false)
const isMutating = shallowRef(false)
const errorMessage = shallowRef('')
const total = shallowRef(0)
let filterTimer: number | undefined

const filters = reactive({
  keyword: '',
  type: '',
  status: '',
  visibility: '',
})

const contentForm = reactive({
  type: 'article' as ContentType,
  title: '',
  slug: '',
  summary: '',
  body_markdown: '',
  body_json: '',
  cover_image_url: '',
  status: 'draft' as ContentStatus,
  visibility: 'private' as ContentVisibility,
  ai_access: 'denied' as 'allowed' | 'denied',
  comment_enabled: true,
  is_featured: false,
  sort_order: 0,
  change_summary: '',
  tag_ids: [] as string[],
})

const tagForm = reactive({
  name: '',
  slug: '',
  description: '',
  color: '',
})

const relationForm = reactive({
  to_content_id: '',
  relation_type: 'related_to' as ContentRelationType,
})

const drawerTitle = computed(() => {
  if (contentMode.value === 'create') {
    return 'New draft'
  }
  if (contentMode.value === 'edit') {
    return 'Edit content'
  }
  return 'Content details'
})

async function loadContents(): Promise<void> {
  isLoading.value = true
  errorMessage.value = ''
  try {
    const response = await studioApi.listContents(
      {
        keyword: filters.keyword.trim(),
        type: filters.type,
        status: filters.status,
        visibility: filters.visibility,
        page: 1,
        page_size: 50,
      },
      { accessToken: authStore.accessToken },
    )
    contents.value = response.items
    total.value = response.total
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to load content.'
    pushToast({ tone: 'danger', title: 'Content unavailable', message: errorMessage.value })
  } finally {
    isLoading.value = false
  }
}

async function loadTags(): Promise<void> {
  isTagsLoading.value = true
  try {
    const response = await studioApi.listTags({ page: 1, page_size: 100 }, { accessToken: authStore.accessToken })
    tags.value = response.items
  } catch (error) {
    pushToast({ tone: 'danger', title: 'Tags unavailable', message: error instanceof Error ? error.message : 'Unable to load tags.' })
  } finally {
    isTagsLoading.value = false
  }
}

function scheduleLoadContents(): void {
  window.clearTimeout(filterTimer)
  filterTimer = window.setTimeout(() => {
    void loadContents()
  }, 300)
}

function openNewDraft(): void {
  contentMode.value = 'create'
  selectedContent.value = null
  resetContentForm()
}

async function openContent(content: ContentSummary, mode: ContentMode): Promise<void> {
  contentMode.value = mode
  isDetailLoading.value = true
  try {
    const response = await studioApi.getContent(content.content_id, { accessToken: authStore.accessToken })
    selectedContent.value = response.content
    hydrateContentForm(response.content)
    await Promise.all([loadRelations(response.content.content_id), loadRevisions(response.content.content_id)])
  } catch (error) {
    pushToast({ tone: 'danger', title: 'Content unavailable', message: error instanceof Error ? error.message : 'Unable to load content.' })
  } finally {
    isDetailLoading.value = false
  }
}

function closeContentDrawer(): void {
  selectedContent.value = null
  contentMode.value = 'view'
  relations.value = []
  revisions.value = []
}

async function saveContent(): Promise<void> {
  await runMutation(async () => {
    const payload = {
      type: contentForm.type,
      title: contentForm.title,
      slug: contentForm.slug,
      summary: contentForm.summary,
      body_markdown: contentForm.body_markdown,
      body_json: contentForm.body_json,
      cover_image_url: contentForm.cover_image_url,
      status: contentForm.status,
      visibility: contentForm.visibility,
      ai_access: contentForm.ai_access,
      source_type: 'manual',
      comment_enabled: contentForm.comment_enabled,
      is_featured: contentForm.is_featured,
      sort_order: Number(contentForm.sort_order) || 0,
      tag_ids: contentForm.tag_ids,
      change_summary: contentForm.change_summary,
    }
    const response = contentMode.value === 'create'
      ? await studioApi.createContent(payload, { accessToken: authStore.accessToken })
      : await studioApi.updateContent(selectedContent.value!.content_id, payload, { accessToken: authStore.accessToken })
    selectedContent.value = response.content
    contentMode.value = 'edit'
    hydrateContentForm(response.content)
    await loadContents()
    pushToast({ tone: 'success', title: 'Content saved', message: `${response.content.title} has been saved.` })
  })
}

async function archiveContent(content: ContentSummary): Promise<void> {
  const approved = await confirm({
    title: 'Archive content?',
    message: `${content.title} will be moved out of active publishing flows.`,
    confirmText: 'Archive',
    tone: 'danger',
  })
  if (!approved) {
    return
  }
  await runMutation(async () => {
    await studioApi.archiveContent(content.content_id, { accessToken: authStore.accessToken })
    await loadContents()
    pushToast({ tone: 'success', title: 'Content archived', message: `${content.title} has been archived.` })
  })
}

function editTag(tag: ContentTag | null): void {
  selectedTag.value = tag
  tagForm.name = tag?.name ?? ''
  tagForm.slug = tag?.slug ?? ''
  tagForm.description = tag?.description ?? ''
  tagForm.color = tag?.color ?? ''
}

async function saveTag(): Promise<void> {
  await runMutation(async () => {
    if (selectedTag.value) {
      await studioApi.updateTag(selectedTag.value.tag_id, tagForm, { accessToken: authStore.accessToken })
      pushToast({ tone: 'success', title: 'Tag updated' })
    } else {
      await studioApi.createTag(tagForm, { accessToken: authStore.accessToken })
      pushToast({ tone: 'success', title: 'Tag created' })
    }
    selectedTag.value = null
    tagForm.name = ''
    tagForm.slug = ''
    tagForm.description = ''
    tagForm.color = ''
    await loadTags()
  })
}

async function deleteTag(tag: ContentTag): Promise<void> {
  const approved = await confirm({ title: 'Delete tag?', message: `${tag.name} will be removed if it is not in use.`, confirmText: 'Delete tag', tone: 'danger' })
  if (!approved) {
    return
  }
  await runMutation(async () => {
    await studioApi.deleteTag(tag.tag_id, { accessToken: authStore.accessToken })
    await loadTags()
    pushToast({ tone: 'success', title: 'Tag deleted' })
  })
}

async function createRelation(): Promise<void> {
  if (!selectedContent.value || relationForm.to_content_id.trim() === '') {
    return
  }
  await runMutation(async () => {
    await studioApi.createRelation(
      selectedContent.value!.content_id,
      { to_content_id: relationForm.to_content_id.trim(), relation_type: relationForm.relation_type },
      { accessToken: authStore.accessToken },
    )
    relationForm.to_content_id = ''
    await loadRelations(selectedContent.value!.content_id)
    pushToast({ tone: 'success', title: 'Relation created' })
  })
}

async function deleteRelation(relation: ContentRelation): Promise<void> {
  if (!selectedContent.value) {
    return
  }
  await runMutation(async () => {
    await studioApi.deleteRelation(selectedContent.value!.content_id, relation.relation_id, { accessToken: authStore.accessToken })
    await loadRelations(selectedContent.value!.content_id)
    pushToast({ tone: 'success', title: 'Relation deleted' })
  })
}

async function loadRelations(contentId: string): Promise<void> {
  const response = await studioApi.listRelations(contentId, { page: 1, page_size: 50 }, { accessToken: authStore.accessToken })
  relations.value = response.items
}

async function loadRevisions(contentId: string): Promise<void> {
  const response = await studioApi.listRevisions(contentId, { page: 1, page_size: 20 }, { accessToken: authStore.accessToken })
  revisions.value = response.items
}

async function runMutation(action: () => Promise<void>): Promise<void> {
  isMutating.value = true
  try {
    await action()
  } catch (error) {
    pushToast({ tone: 'danger', title: 'Operation failed', message: error instanceof Error ? error.message : 'Unable to update content.' })
  } finally {
    isMutating.value = false
  }
}

function resetContentForm(): void {
  Object.assign(contentForm, {
    type: 'article',
    title: '',
    slug: '',
    summary: '',
    body_markdown: '',
    body_json: '{"type":"doc","content":[{"type":"paragraph"}]}',
    cover_image_url: '',
    status: 'draft',
    visibility: 'private',
    ai_access: 'denied',
    comment_enabled: true,
    is_featured: false,
    sort_order: 0,
    change_summary: 'Initial draft',
    tag_ids: [],
  })
}

function hydrateContentForm(content: ContentDetail): void {
  Object.assign(contentForm, {
    type: content.type,
    title: content.title,
    slug: content.slug,
    summary: content.summary ?? '',
    body_markdown: content.body_markdown,
    body_json: content.body_json || '{"type":"doc","content":[{"type":"paragraph"}]}',
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
}

function formatUnixTime(value?: number): string {
  if (!value) {
    return 'None'
  }
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value * 1000))
}

watch(() => [filters.keyword, filters.type, filters.status, filters.visibility], scheduleLoadContents)

onMounted(() => {
  void loadContents()
  void loadTags()
})
onBeforeUnmount(() => window.clearTimeout(filterTimer))
</script>

<template>
  <section class="content-page">
    <PageHeader
      eyebrow="Studio"
      title="Content"
      description="Manage content items, tags, relations, and revisions through the gateway content API."
    >
      <template #actions>
        <BaseButton @click="openNewDraft">New draft</BaseButton>
      </template>
    </PageHeader>

    <div class="content-page__tabs" role="tablist" aria-label="Content workspace">
      <button type="button" :class="{ active: activeTab === 'content' }" @click="activeTab = 'content'">Content</button>
      <button type="button" :class="{ active: activeTab === 'tags' }" @click="activeTab = 'tags'">Tags</button>
    </div>

    <template v-if="activeTab === 'content'">
      <div class="content-page__filters">
        <FormField label="Search" for-id="content-search">
          <BaseInput id="content-search" v-model="filters.keyword" placeholder="Title or slug" />
        </FormField>
        <label class="content-page__select">
          <span>Type</span>
          <select v-model="filters.type">
            <option value="">All types</option>
            <option v-for="type in contentTypes" :key="type" :value="type">{{ type }}</option>
          </select>
        </label>
        <label class="content-page__select">
          <span>Status</span>
          <select v-model="filters.status">
            <option value="">All statuses</option>
            <option v-for="status in statuses" :key="status" :value="status">{{ status }}</option>
          </select>
        </label>
        <label class="content-page__select">
          <span>Visibility</span>
          <select v-model="filters.visibility">
            <option value="">All visibility</option>
            <option v-for="visibility in visibilities" :key="visibility" :value="visibility">{{ visibility }}</option>
          </select>
        </label>
      </div>

      <StatusAlert v-if="errorMessage" tone="danger" title="Content unavailable">{{ errorMessage }}</StatusAlert>
      <PageLoadingState v-else-if="isLoading" title="Loading content" :rows="5" />

      <div v-else class="content-page__table" role="region" aria-label="Studio content" tabindex="0">
        <table>
          <thead>
            <tr>
              <th scope="col">Title</th>
              <th scope="col">Type</th>
              <th scope="col">Status</th>
              <th scope="col">Visibility</th>
              <th scope="col">Updated</th>
              <th scope="col">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="contents.length === 0">
              <td colspan="6">No content yet.</td>
            </tr>
            <tr v-for="content in contents" v-else :key="content.content_id">
              <td>
                <strong>{{ content.title }}</strong>
                <span>{{ content.slug }}</span>
              </td>
              <td><StatusBadge :value="content.type" /></td>
              <td><StatusBadge :value="content.status" /></td>
              <td><StatusBadge :value="content.visibility" /></td>
              <td>{{ formatUnixTime(content.updated_at) }}</td>
              <td>
                <div class="content-page__actions">
                  <ActionTagButton @click="openContent(content, 'view')">View</ActionTagButton>
                  <ActionTagButton tone="primary" @click="openContent(content, 'edit')">Edit</ActionTagButton>
                  <ActionTagButton tone="danger" :disabled="content.status === 'archived' || isMutating" @click="archiveContent(content)">Archive</ActionTagButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <p class="content-page__count">{{ total }} total content items</p>
    </template>

    <template v-else>
      <form class="content-page__tag-form" @submit.prevent="saveTag">
        <FormField label="Name" for-id="tag-name">
          <BaseInput id="tag-name" v-model="tagForm.name" />
        </FormField>
        <FormField label="Slug" for-id="tag-slug">
          <BaseInput id="tag-slug" v-model="tagForm.slug" />
        </FormField>
        <FormField label="Color" for-id="tag-color">
          <BaseInput id="tag-color" v-model="tagForm.color" placeholder="#0f8f83" />
        </FormField>
        <FormField label="Description" for-id="tag-description">
          <BaseInput id="tag-description" v-model="tagForm.description" />
        </FormField>
        <BaseButton type="submit" :busy="isMutating">{{ selectedTag ? 'Save tag' : 'Create tag' }}</BaseButton>
      </form>
      <PageLoadingState v-if="isTagsLoading" title="Loading tags" :rows="3" />
      <div v-else class="content-page__tag-list">
        <article v-for="tag in tags" :key="tag.tag_id" class="content-page__tag-card">
          <div>
            <strong>{{ tag.name }}</strong>
            <span>{{ tag.slug }}</span>
          </div>
          <div class="content-page__actions">
            <ActionTagButton @click="editTag(tag)">Edit</ActionTagButton>
            <ActionTagButton tone="danger" @click="deleteTag(tag)">Delete</ActionTagButton>
          </div>
        </article>
      </div>
    </template>

    <SideDrawer :open="contentMode === 'create' || selectedContent !== null" :title="drawerTitle" :description="selectedContent?.slug" size="lg" @close="closeContentDrawer">
      <PageLoadingState v-if="isDetailLoading" title="Loading content detail" :rows="4" />
      <div v-else class="content-page__drawer">
        <template v-if="contentMode === 'view' && selectedContent">
          <div class="content-page__detail-grid">
            <ReadonlyField label="Title" :value="selectedContent.title" />
            <ReadonlyField label="Slug" :value="selectedContent.slug" />
            <ReadonlyField label="Type" :value="selectedContent.type" />
            <ReadonlyField label="Status" :value="selectedContent.status" />
            <ReadonlyField label="Visibility" :value="selectedContent.visibility" />
            <ReadonlyField label="Updated" :value="formatUnixTime(selectedContent.updated_at)" />
          </div>
        </template>
        <template v-else>
          <div class="content-page__form-grid">
            <label class="content-page__select">
              <span>Type</span>
              <select v-model="contentForm.type">
                <option v-for="type in contentTypes" :key="type" :value="type">{{ type }}</option>
              </select>
            </label>
            <label class="content-page__select">
              <span>Status</span>
              <select v-model="contentForm.status">
                <option v-for="status in statuses" :key="status" :value="status">{{ status }}</option>
              </select>
            </label>
            <FormField label="Title" for-id="content-title">
              <BaseInput id="content-title" v-model="contentForm.title" />
            </FormField>
            <FormField label="Slug" for-id="content-slug">
              <BaseInput id="content-slug" v-model="contentForm.slug" />
            </FormField>
            <FormField label="Summary" for-id="content-summary">
              <BaseInput id="content-summary" v-model="contentForm.summary" />
            </FormField>
            <FormField label="Cover URL" for-id="content-cover">
              <BaseInput id="content-cover" v-model="contentForm.cover_image_url" />
            </FormField>
            <label class="content-page__select">
              <span>Visibility</span>
              <select v-model="contentForm.visibility">
                <option v-for="visibility in visibilities" :key="visibility" :value="visibility">{{ visibility }}</option>
              </select>
            </label>
            <label class="content-page__select">
              <span>AI access</span>
              <select v-model="contentForm.ai_access">
                <option value="denied">Denied</option>
                <option value="allowed">Allowed</option>
              </select>
            </label>
          </div>
          <div class="content-page__checks">
            <label><input v-model="contentForm.comment_enabled" type="checkbox" /> Comments</label>
            <label><input v-model="contentForm.is_featured" type="checkbox" /> Featured</label>
          </div>
          <div class="content-page__tag-picker">
            <span>Tags</span>
            <label v-for="tag in tags" :key="tag.tag_id">
              <input v-model="contentForm.tag_ids" type="checkbox" :value="tag.tag_id" />
              {{ tag.name }}
            </label>
          </div>
          <FormField label="Change summary" for-id="content-change-summary">
            <BaseInput id="content-change-summary" v-model="contentForm.change_summary" />
          </FormField>
          <RichContentEditor v-model="contentForm.body_json" v-model:plain-text="contentForm.body_markdown" />
        </template>

        <section v-if="selectedContent" class="content-page__subsection">
          <h3>Relations</h3>
          <form v-if="contentMode !== 'view'" class="content-page__relation-form" @submit.prevent="createRelation">
            <BaseInput v-model="relationForm.to_content_id" placeholder="Target content ID" />
            <select v-model="relationForm.relation_type">
              <option v-for="type in relationTypes" :key="type" :value="type">{{ type }}</option>
            </select>
            <BaseButton type="submit" variant="secondary" :busy="isMutating">Add</BaseButton>
          </form>
          <div class="content-page__mini-list">
            <article v-for="relation in relations" :key="relation.relation_id">
              <span>{{ relation.relation_type }} -> {{ relation.to_content_id }}</span>
              <ActionTagButton v-if="contentMode !== 'view'" tone="danger" @click="deleteRelation(relation)">Delete</ActionTagButton>
            </article>
            <p v-if="relations.length === 0">No relations.</p>
          </div>
        </section>

        <section v-if="selectedContent" class="content-page__subsection">
          <h3>Revisions</h3>
          <div class="content-page__mini-list">
            <article v-for="revision in revisions" :key="revision.revision_id">
              <span>#{{ revision.revision_no }} {{ revision.change_summary || 'No summary' }}</span>
              <small>{{ formatUnixTime(revision.created_at) }}</small>
            </article>
            <p v-if="revisions.length === 0">No revisions.</p>
          </div>
        </section>
      </div>
      <template #footer>
        <BaseButton v-if="contentMode !== 'view'" :busy="isMutating" @click="saveContent">Save content</BaseButton>
        <BaseButton variant="ghost" @click="closeContentDrawer">Close</BaseButton>
      </template>
    </SideDrawer>
  </section>
</template>

<style scoped>
.content-page,
.content-page__drawer,
.content-page__subsection {
  display: grid;
  gap: 20px;
}

.content-page__tabs {
  display: inline-flex;
  width: fit-content;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 4px;
  background: var(--bb-color-surface);
}

.content-page__tabs button {
  min-height: 36px;
  border: 0;
  border-radius: 6px;
  padding: 0 14px;
  color: var(--bb-color-muted);
  background: transparent;
  font-weight: 700;
}

.content-page__tabs button.active {
  color: var(--bb-color-primary);
  background: var(--bb-color-primary-soft);
}

.content-page__filters,
.content-page__form-grid {
  display: grid;
  grid-template-columns: minmax(180px, 1fr) repeat(3, minmax(140px, 170px));
  align-items: end;
  gap: 12px;
}

.content-page__select {
  display: grid;
  gap: 6px;
  color: var(--bb-color-muted);
  font-size: 0.92rem;
  font-weight: 650;
}

select {
  min-height: 44px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 0 10px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface);
}

select:focus-visible,
.content-page__table:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.content-page__table {
  overflow-x: auto;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-soft);
}

table {
  width: 100%;
  min-width: 920px;
  border-collapse: collapse;
}

th,
td {
  border-bottom: 1px solid var(--bb-color-line);
  padding: 12px;
  text-align: left;
  vertical-align: middle;
}

th {
  color: var(--bb-color-muted);
  font-size: 0.8rem;
  text-transform: uppercase;
  background: var(--bb-color-subtle);
}

td:first-child {
  display: grid;
  gap: 3px;
}

td:first-child span,
.content-page__count,
.content-page__tag-card span {
  color: var(--bb-color-muted);
}

.content-page__actions,
.content-page__checks,
.content-page__relation-form {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.content-page__detail-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.content-page__tag-form {
  display: grid;
  grid-template-columns: repeat(4, minmax(140px, 1fr)) auto;
  align-items: end;
  gap: 12px;
}

.content-page__tag-list,
.content-page__mini-list,
.content-page__tag-picker {
  display: grid;
  gap: 10px;
}

.content-page__tag-card,
.content-page__mini-list article {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 12px;
  background: var(--bb-color-surface);
}

.content-page__tag-card div:first-child {
  display: grid;
}

.content-page__tag-picker {
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 12px;
  background: var(--bb-color-subtle);
}

.content-page__tag-picker span,
.content-page__subsection h3 {
  margin: 0;
  color: var(--bb-color-text);
  font-weight: 760;
}

@media (max-width: 900px) {
  .content-page__filters,
  .content-page__form-grid,
  .content-page__tag-form,
  .content-page__detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
