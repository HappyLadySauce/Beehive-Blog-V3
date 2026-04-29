<script setup lang="ts">
import { Archive, Eye, Pencil, Trash2 } from 'lucide-vue-next'
import { computed, onBeforeUnmount, onMounted, reactive, shallowRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import type {
  ContentDetail,
  ContentRelation,
  ContentRevisionSummary,
  ContentStatus,
  ContentSummary,
  ContentTag,
  ContentType,
  ContentVisibility,
} from '@/features/studio'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSelect, { type BaseSelectOption } from '@/shared/components/BaseSelect.vue'
import FormField from '@/shared/components/FormField.vue'
import IconActionButton from '@/shared/components/IconActionButton.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import ReadonlyField from '@/shared/components/ReadonlyField.vue'
import SideDrawer from '@/shared/components/SideDrawer.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import StatusBadge from '@/shared/components/StatusBadge.vue'
import { useConfirm, useToast } from '@/shared/composables'
import { useLocale } from '@/shared/i18n'

type StudioTab = 'content' | 'tags'

const contentTypes: ContentType[] = ['article', 'note', 'project', 'experience', 'timeline_event', 'insight', 'portfolio', 'page']
const statuses: ContentStatus[] = ['draft', 'review', 'published', 'archived']
const visibilities: ContentVisibility[] = ['public', 'member', 'private']

const router = useRouter()
const authStore = useAuthStore()
const { t } = useI18n()
const { locale } = useLocale()
const { confirm } = useConfirm()
const { pushToast } = useToast()

const activeTab = shallowRef<StudioTab>('content')
const contents = shallowRef<ContentSummary[]>([])
const tags = shallowRef<ContentTag[]>([])
const relations = shallowRef<ContentRelation[]>([])
const revisions = shallowRef<ContentRevisionSummary[]>([])
const selectedContent = shallowRef<ContentDetail | null>(null)
const selectedTag = shallowRef<ContentTag | null>(null)
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

const typeOptions = computed<BaseSelectOption[]>(() => [
  { value: '', label: t('contentType.all') },
  ...contentTypes.map((type) => ({ value: type, label: t(`contentType.${type}`) })),
])
const statusOptions = computed<BaseSelectOption[]>(() => [
  { value: '', label: t('contentStatus.all') },
  ...statuses.map((status) => ({ value: status, label: t(`contentStatus.${status}`) })),
])
const visibilityOptions = computed<BaseSelectOption[]>(() => [
  { value: '', label: t('visibility.all') },
  ...visibilities.map((visibility) => ({ value: visibility, label: t(`visibility.${visibility}`) })),
])

const tagForm = reactive({
  name: '',
  slug: '',
  description: '',
  color: '',
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
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : t('content.unavailableTitle')
    pushToast({ tone: 'danger', title: t('content.unavailableTitle'), message: errorMessage.value })
  }
  finally {
    isLoading.value = false
  }
}

async function loadTags(): Promise<void> {
  isTagsLoading.value = true
  try {
    const response = await studioApi.listTags({ page: 1, page_size: 100 }, { accessToken: authStore.accessToken })
    tags.value = response.items
  }
  catch (error) {
    pushToast({ tone: 'danger', title: 'Tags unavailable', message: error instanceof Error ? error.message : 'Unable to load tags.' })
  }
  finally {
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
  void router.push('/studio/content/new')
}

function editContent(content: ContentSummary): void {
  void router.push(`/studio/content/${encodeURIComponent(content.content_id)}/edit`)
}

async function viewContent(content: ContentSummary): Promise<void> {
  isDetailLoading.value = true
  try {
    const response = await studioApi.getContent(content.content_id, { accessToken: authStore.accessToken })
    selectedContent.value = response.content
    await Promise.all([loadRelations(response.content.content_id), loadRevisions(response.content.content_id)])
  }
  catch (error) {
    pushToast({ tone: 'danger', title: t('content.unavailableTitle'), message: error instanceof Error ? error.message : t('content.unavailableTitle') })
  }
  finally {
    isDetailLoading.value = false
  }
}

function closeContentDrawer(): void {
  selectedContent.value = null
  relations.value = []
  revisions.value = []
}

async function archiveContent(content: ContentSummary): Promise<void> {
  const approved = await confirm({
    title: t('content.archiveConfirmTitle'),
    message: t('content.archiveConfirmMessage', { title: content.title }),
    confirmText: t('content.archiveConfirmAction'),
    tone: 'danger',
  })
  if (!approved) {
    return
  }
  await runMutation(async () => {
    await studioApi.archiveContent(content.content_id, { accessToken: authStore.accessToken })
    await loadContents()
    pushToast({ tone: 'success', title: t('content.archivedTitle'), message: t('content.archivedMessage', { title: content.title }) })
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
  }
  catch (error) {
    pushToast({ tone: 'danger', title: 'Operation failed', message: error instanceof Error ? error.message : 'Unable to update content.' })
  }
  finally {
    isMutating.value = false
  }
}

function formatUnixTime(value?: number): string {
  if (!value) {
    return t('common.none')
  }
  return new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value * 1000))
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
      :eyebrow="t('content.eyebrow')"
      :title="t('content.title')"
      :description="t('content.description')"
    >
      <template #actions>
        <BaseButton @click="openNewDraft">{{ t('content.newDraft') }}</BaseButton>
      </template>
    </PageHeader>

    <div class="content-page__tabs" role="tablist" aria-label="Content workspace">
      <button type="button" :class="{ active: activeTab === 'content' }" @click="activeTab = 'content'">{{ t('content.tabs.content') }}</button>
      <button type="button" :class="{ active: activeTab === 'tags' }" @click="activeTab = 'tags'">{{ t('content.tabs.tags') }}</button>
    </div>

    <template v-if="activeTab === 'content'">
      <div class="content-page__filters">
        <FormField :label="t('common.search')" for-id="content-search">
          <BaseInput id="content-search" v-model="filters.keyword" :placeholder="t('content.searchPlaceholder')" />
        </FormField>
        <FormField :label="t('content.columns.type')" for-id="content-type-filter">
          <BaseSelect id="content-type-filter" v-model="filters.type" :options="typeOptions" :aria-label="t('content.columns.type')" />
        </FormField>
        <FormField :label="t('content.columns.status')" for-id="content-status-filter">
          <BaseSelect id="content-status-filter" v-model="filters.status" :options="statusOptions" :aria-label="t('content.columns.status')" />
        </FormField>
        <FormField :label="t('content.columns.visibility')" for-id="content-visibility-filter">
          <BaseSelect id="content-visibility-filter" v-model="filters.visibility" :options="visibilityOptions" :aria-label="t('content.columns.visibility')" />
        </FormField>
      </div>

      <StatusAlert v-if="errorMessage" tone="danger" :title="t('content.unavailableTitle')">{{ errorMessage }}</StatusAlert>
      <PageLoadingState v-else-if="isLoading" :title="t('content.loadingTitle')" :rows="5" />

      <div v-else class="content-page__table" role="region" aria-label="Studio content" tabindex="0">
        <table>
          <thead>
            <tr>
              <th scope="col">{{ t('content.columns.title') }}</th>
              <th scope="col">{{ t('content.columns.type') }}</th>
              <th scope="col">{{ t('content.columns.status') }}</th>
              <th scope="col">{{ t('content.columns.visibility') }}</th>
              <th scope="col">{{ t('content.columns.updated') }}</th>
              <th scope="col">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="contents.length === 0">
              <td colspan="6">{{ t('content.empty') }}</td>
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
                  <IconActionButton :aria-label="t('content.actions.view', { title: content.title })" :title="t('content.actions.view', { title: content.title })" @click="viewContent(content)">
                    <Eye :size="17" aria-hidden="true" />
                  </IconActionButton>
                  <IconActionButton tone="primary" :aria-label="t('content.actions.edit', { title: content.title })" :title="t('content.actions.edit', { title: content.title })" @click="editContent(content)">
                    <Pencil :size="17" aria-hidden="true" />
                  </IconActionButton>
                  <IconActionButton
                    tone="danger"
                    :disabled="content.status === 'archived' || isMutating"
                    :aria-label="t('content.actions.archive', { title: content.title })"
                    :title="t('content.actions.archive', { title: content.title })"
                    @click="archiveContent(content)"
                  >
                    <Archive :size="17" aria-hidden="true" />
                  </IconActionButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <p class="content-page__count">{{ t('content.count', { count: total }) }}</p>
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
        <BaseButton type="submit" :busy="isMutating">{{ selectedTag ? t('common.save') : t('common.save') }}</BaseButton>
      </form>
      <PageLoadingState v-if="isTagsLoading" title="Loading tags" :rows="3" />
      <div v-else class="content-page__tag-list">
        <article v-for="tag in tags" :key="tag.tag_id" class="content-page__tag-card">
          <div>
            <strong>{{ tag.name }}</strong>
            <span>{{ tag.slug }}</span>
          </div>
          <div class="content-page__actions">
            <IconActionButton :aria-label="t('content.actions.editTag', { name: tag.name })" :title="t('content.actions.editTag', { name: tag.name })" @click="editTag(tag)">
              <Pencil :size="17" aria-hidden="true" />
            </IconActionButton>
            <IconActionButton tone="danger" :aria-label="t('content.actions.deleteTag', { name: tag.name })" :title="t('content.actions.deleteTag', { name: tag.name })" @click="deleteTag(tag)">
              <Trash2 :size="17" aria-hidden="true" />
            </IconActionButton>
          </div>
        </article>
      </div>
    </template>

    <SideDrawer :open="selectedContent !== null" :title="t('content.detailTitle')" :description="selectedContent?.slug" size="lg" @close="closeContentDrawer">
      <PageLoadingState v-if="isDetailLoading" title="Loading content detail" :rows="4" />
      <div v-else-if="selectedContent" class="content-page__drawer">
        <div class="content-page__detail-grid">
          <ReadonlyField :label="t('content.columns.title')" :value="selectedContent.title" />
          <ReadonlyField :label="t('content.columns.slug')" :value="selectedContent.slug" />
          <ReadonlyField :label="t('content.columns.type')" :value="selectedContent.type" />
          <ReadonlyField :label="t('content.columns.status')" :value="selectedContent.status" />
          <ReadonlyField :label="t('content.columns.visibility')" :value="selectedContent.visibility" />
          <ReadonlyField :label="t('content.columns.updated')" :value="formatUnixTime(selectedContent.updated_at)" />
        </div>

        <section class="content-page__subsection">
          <h3>Relations</h3>
          <div class="content-page__mini-list">
            <article v-for="relation in relations" :key="relation.relation_id">
              <span>{{ relation.relation_type }} -> {{ relation.to_content_id }}</span>
            </article>
            <p v-if="relations.length === 0">No relations.</p>
          </div>
        </section>

        <section class="content-page__subsection">
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
        <BaseButton variant="ghost" @click="closeContentDrawer">{{ t('common.close') }}</BaseButton>
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

.content-page__filters {
  display: grid;
  grid-template-columns: minmax(180px, 1fr) repeat(3, minmax(140px, 170px));
  align-items: end;
  gap: 12px;
}

.content-page__table:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.content-page__table {
  overflow-x: auto;
  border: 1px solid var(--bb-color-line);
  border-radius: 10px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-soft);
}

.content-page__table table {
  width: 100%;
  min-width: 920px;
  border-collapse: collapse;
}

.content-page__table th,
.content-page__table td {
  border-bottom: 1px solid var(--bb-color-line);
  padding: 12px;
  text-align: left;
  vertical-align: middle;
}

.content-page__table th {
  color: var(--bb-color-muted);
  font-size: 0.8rem;
  text-transform: uppercase;
  background: var(--bb-color-surface);
}

.content-page__table tbody tr:nth-child(even) {
  background: var(--bb-color-subtle);
}

.content-page__table tbody tr:hover {
  background: var(--bb-color-primary-soft);
}

.content-page__table tbody tr:last-child td {
  border-bottom: 0;
}

.content-page__table td:first-child {
  display: grid;
  gap: 3px;
}

.content-page__table td:first-child span,
.content-page__count,
.content-page__tag-card span {
  color: var(--bb-color-muted);
}

.content-page__actions {
  display: inline-flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 8px;
  align-items: center;
}

.content-page__table td:last-child,
.content-page__table th:last-child {
  text-align: right;
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
.content-page__mini-list {
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

@media (max-width: 900px) {
  .content-page__filters,
  .content-page__tag-form,
  .content-page__detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
