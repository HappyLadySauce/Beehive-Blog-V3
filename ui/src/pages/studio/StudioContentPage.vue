<script setup lang="ts">
import { Archive, Eye, PackageOpen, Pencil, Tag, Trash2 } from 'lucide-vue-next'
import { computed, onBeforeUnmount, reactive, shallowRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

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
import EmptyState from '@/shared/components/EmptyState.vue'
import FormField from '@/shared/components/FormField.vue'
import IconActionButton from '@/shared/components/IconActionButton.vue'
import InlineLoadingState from '@/shared/components/InlineLoadingState.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import ReadonlyField from '@/shared/components/ReadonlyField.vue'
import SideDrawer from '@/shared/components/SideDrawer.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import StatusBadge from '@/shared/components/StatusBadge.vue'
import TablePagination from '@/shared/components/TablePagination.vue'
import { useConfirm, usePaginatedRouteState, useProgressiveQuery, useToast } from '@/shared/composables'
import { useLocale } from '@/shared/i18n'

type StudioTab = 'content' | 'tags'

const contentTypes: ContentType[] = ['article', 'note', 'project', 'experience', 'timeline_event', 'insight', 'portfolio', 'page']
const statuses: ContentStatus[] = ['draft', 'review', 'published', 'archived']
const visibilities: ContentVisibility[] = ['public', 'member', 'private']

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { t } = useI18n()
const { locale } = useLocale()
const { confirm } = useConfirm()
const { pushToast } = useToast()

const activeTab = shallowRef<StudioTab>(readTabQuery(route.query.tab))
const selectedContentId = shallowRef('')
const selectedTag = shallowRef<ContentTag | null>(null)
const isMutating = shallowRef(false)
const isTagDrawerOpen = shallowRef(false)
const errorMessage = shallowRef('')
let filterTimer: number | undefined
let tagFilterTimer: number | undefined

const filters = reactive({
  keyword: readQueryString(route.query.keyword),
  type: readQueryString(route.query.type),
  status: readQueryString(route.query.status),
  visibility: readQueryString(route.query.visibility),
})
const appliedFilters = reactive({
  keyword: readQueryString(route.query.keyword),
  type: readQueryString(route.query.type),
  status: readQueryString(route.query.status),
  visibility: readQueryString(route.query.visibility),
})
const tagFilters = reactive({
  keyword: readQueryString(route.query.tagKeyword),
})
const appliedTagFilters = reactive({
  keyword: readQueryString(route.query.tagKeyword),
})
const contentTotal = shallowRef(0)
const tagTotal = shallowRef(0)

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
const contentPagination = usePaginatedRouteState({
  route,
  router,
  total: contentTotal,
})
const tagPagination = usePaginatedRouteState({
  route,
  router,
  total: tagTotal,
  pageParam: 'tagPage',
  pageSizeParam: 'tagPageSize',
})

const contentQuery = useProgressiveQuery({
  queryKey: computed(() => ['studio-contents', { ...appliedFilters }, contentPagination.page.value, contentPagination.pageSize.value]),
  queryFn: () => studioApi.listContents(
    {
      keyword: appliedFilters.keyword.trim(),
      type: appliedFilters.type,
      status: appliedFilters.status,
      visibility: appliedFilters.visibility,
      page: contentPagination.page.value,
      page_size: contentPagination.pageSize.value,
    },
    { accessToken: authStore.accessToken },
  ),
})

const tagsQuery = useProgressiveQuery({
  queryKey: computed(() => ['studio-tags', { ...appliedTagFilters }, tagPagination.page.value, tagPagination.pageSize.value]),
  queryFn: () => studioApi.listTags(
    {
      keyword: appliedTagFilters.keyword.trim(),
      page: tagPagination.page.value,
      page_size: tagPagination.pageSize.value,
    },
    { accessToken: authStore.accessToken },
  ),
})

const detailQuery = useProgressiveQuery({
  queryKey: computed(() => ['studio-content-detail', selectedContentId.value]),
  queryFn: () => studioApi.getContent(selectedContentId.value, { accessToken: authStore.accessToken }),
  enabled: computed(() => Boolean(selectedContentId.value)),
})

const relationsQuery = useProgressiveQuery({
  queryKey: computed(() => ['studio-content-relations', selectedContentId.value]),
  queryFn: () => studioApi.listRelations(selectedContentId.value, { page: 1, page_size: 50 }, { accessToken: authStore.accessToken }),
  enabled: computed(() => Boolean(selectedContentId.value)),
})

const revisionsQuery = useProgressiveQuery({
  queryKey: computed(() => ['studio-content-revisions', selectedContentId.value]),
  queryFn: () => studioApi.listRevisions(selectedContentId.value, { page: 1, page_size: 20 }, { accessToken: authStore.accessToken }),
  enabled: computed(() => Boolean(selectedContentId.value)),
})

const contents = computed<ContentSummary[]>(() => contentQuery.data.value?.items ?? [])
const tags = computed<ContentTag[]>(() => tagsQuery.data.value?.items ?? [])
const selectedContent = computed<ContentDetail | null>(() => detailQuery.data.value?.content ?? null)
const relations = computed<ContentRelation[]>(() => relationsQuery.data.value?.items ?? [])
const revisions = computed<ContentRevisionSummary[]>(() => revisionsQuery.data.value?.items ?? [])
const contentErrorMessage = computed(() => {
  const error = contentQuery.error.value
  return error instanceof Error ? error.message : ''
})
const tagsErrorMessage = computed(() => {
  const error = tagsQuery.error.value
  return error instanceof Error ? error.message : ''
})
const hasContents = computed(() => contents.value.length > 0)
const hasTags = computed(() => tags.value.length > 0)
const detailRefreshing = computed(() =>
  detailQuery.showRefreshingHint.value
  || relationsQuery.showRefreshingHint.value
  || revisionsQuery.showRefreshingHint.value,
)

function scheduleLoadContents(): void {
  window.clearTimeout(filterTimer)
  filterTimer = window.setTimeout(() => {
    appliedFilters.keyword = filters.keyword
    appliedFilters.type = filters.type
    appliedFilters.status = filters.status
    appliedFilters.visibility = filters.visibility
    contentPagination.page.value = 1
    void contentPagination.syncQuery(buildContentQuery())
  }, 300)
}

function scheduleLoadTags(): void {
  window.clearTimeout(tagFilterTimer)
  tagFilterTimer = window.setTimeout(() => {
    appliedTagFilters.keyword = tagFilters.keyword
    tagPagination.page.value = 1
    void tagPagination.syncQuery(buildTagQuery())
  }, 300)
}

function openNewDraft(): void {
  void router.push('/studio/content/new')
}

function editContent(content: ContentSummary): void {
  void router.push(`/studio/content/${encodeURIComponent(content.content_id)}/edit`)
}

function openCreateTag(): void {
  selectedTag.value = null
  tagForm.name = ''
  tagForm.slug = ''
  tagForm.description = ''
  tagForm.color = ''
  isTagDrawerOpen.value = true
}

async function viewContent(content: ContentSummary): Promise<void> {
  selectedContentId.value = content.content_id
}

function closeContentDrawer(): void {
  selectedContentId.value = ''
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
    await contentQuery.refetch()
    pushToast({ tone: 'success', title: t('content.archivedTitle'), message: t('content.archivedMessage', { title: content.title }) })
  })
}

function editTag(tag: ContentTag): void {
  selectedTag.value = tag
  tagForm.name = tag.name
  tagForm.slug = tag.slug
  tagForm.description = tag.description ?? ''
  tagForm.color = tag.color ?? ''
  isTagDrawerOpen.value = true
}

function closeTagDrawer(): void {
  isTagDrawerOpen.value = false
  selectedTag.value = null
  tagForm.name = ''
  tagForm.slug = ''
  tagForm.description = ''
  tagForm.color = ''
}

async function saveTag(): Promise<void> {
  await runMutation(async () => {
    if (selectedTag.value) {
      await studioApi.updateTag(selectedTag.value.tag_id, tagForm, { accessToken: authStore.accessToken })
      pushToast({ tone: 'success', title: t('content.tags.updateTitle') })
    } else {
      await studioApi.createTag(tagForm, { accessToken: authStore.accessToken })
      pushToast({ tone: 'success', title: t('content.tags.createTitle') })
    }
    closeTagDrawer()
    await tagsQuery.refetch()
  })
}

async function deleteTag(tag: ContentTag): Promise<void> {
  const approved = await confirm({
    title: t('content.tags.deleteConfirmTitle'),
    message: t('content.tags.deleteConfirmMessage', { name: tag.name }),
    confirmText: t('content.tags.deleteConfirmAction'),
    tone: 'danger',
  })
  if (!approved) {
    return
  }
  await runMutation(async () => {
    await studioApi.deleteTag(tag.tag_id, { accessToken: authStore.accessToken })
    await tagsQuery.refetch()
    pushToast({ tone: 'success', title: t('content.tags.deleteTitle') })
  })
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

const tagDrawerTitle = computed(() => (selectedTag.value ? t('content.tags.editTitle') : t('content.tags.createDrawerTitle')))

watch(() => contentQuery.data.value?.total ?? 0, (value) => {
  contentTotal.value = value
}, { immediate: true })
watch(() => tagsQuery.data.value?.total ?? 0, (value) => {
  tagTotal.value = value
}, { immediate: true })
watch(
  () => [route.query.tab, route.query.keyword, route.query.type, route.query.status, route.query.visibility, route.query.tagKeyword],
  () => {
    activeTab.value = readTabQuery(route.query.tab)
    filters.keyword = readQueryString(route.query.keyword)
    filters.type = readQueryString(route.query.type)
    filters.status = readQueryString(route.query.status)
    filters.visibility = readQueryString(route.query.visibility)
    appliedFilters.keyword = filters.keyword
    appliedFilters.type = filters.type
    appliedFilters.status = filters.status
    appliedFilters.visibility = filters.visibility
    tagFilters.keyword = readQueryString(route.query.tagKeyword)
    appliedTagFilters.keyword = tagFilters.keyword
  },
)
watch(activeTab, (tab) => {
  const query = tab === 'content' ? buildContentQuery() : buildTagQuery()
  if (route.query.tab !== tab) {
    if (tab === 'content') {
      void contentPagination.syncQuery(query)
      return
    }
    void tagPagination.syncQuery(query)
  }
})
watch(() => [filters.keyword, filters.type, filters.status, filters.visibility], scheduleLoadContents)
watch(() => tagFilters.keyword, scheduleLoadTags)
onBeforeUnmount(() => {
  window.clearTimeout(filterTimer)
  window.clearTimeout(tagFilterTimer)
})

function buildContentQuery() {
  return {
    tab: activeTab.value,
    keyword: filters.keyword.trim() || undefined,
    type: filters.type || undefined,
    status: filters.status || undefined,
    visibility: filters.visibility || undefined,
  }
}

function buildTagQuery() {
  return {
    tab: activeTab.value,
    tagKeyword: tagFilters.keyword.trim() || undefined,
  }
}

function readQueryString(value: unknown): string {
  const normalized = Array.isArray(value) ? value[0] : value
  return typeof normalized === 'string' ? normalized : ''
}

function readTabQuery(value: unknown): StudioTab {
  const normalized = readQueryString(value)
  return normalized === 'tags' ? 'tags' : 'content'
}
</script>

<template>
  <section class="content-page">
    <div class="content-page__toolbar">
      <div class="content-page__tabs" role="tablist" aria-label="Content workspace">
        <button type="button" :class="{ active: activeTab === 'content' }" @click="activeTab = 'content'">{{ t('content.tabs.content') }}</button>
        <button type="button" :class="{ active: activeTab === 'tags' }" @click="activeTab = 'tags'">{{ t('content.tabs.tags') }}</button>
      </div>
      <BaseButton v-if="activeTab === 'content'" @click="openNewDraft">{{ t('content.newDraft') }}</BaseButton>
      <BaseButton v-else @click="openCreateTag">{{ t('content.tags.createAction') }}</BaseButton>
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

      <StatusAlert v-if="contentErrorMessage && !hasContents" tone="danger" :title="t('content.unavailableTitle')">{{ contentErrorMessage }}</StatusAlert>
      <PageLoadingState v-else-if="contentQuery.showBlockingLoading.value && !hasContents" :title="t('content.loadingTitle')" :rows="5" />

      <div v-else-if="contentQuery.hasResolvedOnce.value" class="content-page__table" role="region" aria-label="Studio content" tabindex="0">
        <div v-if="contentQuery.showRefreshingHint.value && hasContents" class="content-page__refreshing">
          <InlineLoadingState />
        </div>
        <table class="content-page__grid">
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
          <tbody v-if="contents.length > 0">
            <tr v-for="content in contents" :key="content.content_id">
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
        <div v-if="contents.length === 0" class="content-page__empty-panel">
          <EmptyState
            class="content-page__empty-state"
            align="center"
            :title="t('content.empty')"
            :description="t('content.emptyDescription')"
          >
            <template #visual>
              <PackageOpen :size="52" aria-hidden="true" />
            </template>
            <BaseButton class="content-page__empty-action" @click="openNewDraft">
              {{ t('content.emptyAction') }}
            </BaseButton>
          </EmptyState>
        </div>
      </div>
      <div v-if="contentQuery.data.value" class="studio-list-footer">
        <TablePagination
          :page="contentPagination.page.value"
          :page-size="contentPagination.pageSize.value"
          :total="contentTotal"
          :disabled="contentQuery.isFetching.value"
          @update:page="contentPagination.setPage"
          @update:page-size="contentPagination.setPageSize"
        />
      </div>
    </template>

    <template v-else>
      <div class="studio-list-shell">
        <div class="studio-list-filters content-page__tag-filters">
          <FormField :label="t('common.search')" for-id="tag-search">
            <BaseInput id="tag-search" v-model="tagFilters.keyword" :placeholder="t('content.tags.searchPlaceholder')" />
          </FormField>
        </div>
        <StatusAlert v-if="tagsErrorMessage && !hasTags" tone="danger" :title="t('content.tags.unavailableTitle')">{{ tagsErrorMessage }}</StatusAlert>
        <PageLoadingState v-else-if="tagsQuery.showBlockingLoading.value && !hasTags" :title="t('content.tags.loadingTitle')" :rows="3" />
        <div v-else-if="tagsQuery.hasResolvedOnce.value" class="studio-list-table content-page__tag-table" role="region" :aria-label="t('content.tabs.tags')" tabindex="0">
          <div v-if="tagsQuery.showRefreshingHint.value && hasTags" class="content-page__refreshing">
            <InlineLoadingState />
          </div>
          <table class="studio-list-grid content-page__tag-grid">
            <thead>
              <tr>
                <th scope="col">{{ t('content.columns.name') }}</th>
                <th scope="col">{{ t('content.columns.slug') }}</th>
                <th scope="col">{{ t('content.tags.fields.color') }}</th>
                <th scope="col">{{ t('content.tags.fields.description') }}</th>
                <th scope="col">{{ t('content.columns.updated') }}</th>
                <th scope="col">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody v-if="tags.length > 0">
              <tr v-for="tag in tags" :key="tag.tag_id">
                <td>
                  <strong>{{ tag.name }}</strong>
                </td>
                <td>{{ tag.slug }}</td>
                <td>
                  <span class="content-page__tag-color">
                    <span class="content-page__tag-swatch" :style="{ backgroundColor: tag.color || 'transparent' }" />
                    {{ tag.color || t('common.none') }}
                  </span>
                </td>
                <td>{{ tag.description || t('common.none') }}</td>
                <td>{{ formatUnixTime(tag.updated_at) }}</td>
                <td>
                  <div class="content-page__actions">
                    <IconActionButton :aria-label="t('content.actions.editTag', { name: tag.name })" :title="t('content.actions.editTag', { name: tag.name })" @click="editTag(tag)">
                      <Pencil :size="17" aria-hidden="true" />
                    </IconActionButton>
                    <IconActionButton tone="danger" :aria-label="t('content.actions.deleteTag', { name: tag.name })" :title="t('content.actions.deleteTag', { name: tag.name })" @click="deleteTag(tag)">
                      <Trash2 :size="17" aria-hidden="true" />
                    </IconActionButton>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <div v-if="tags.length === 0" class="studio-list-empty-panel">
            <EmptyState align="center" class="studio-list-empty-state" :title="t('content.tagsEmpty')" :description="t('content.tags.emptyDescription')">
              <template #visual>
                <Tag :size="52" aria-hidden="true" />
              </template>
              <BaseButton @click="openCreateTag">{{ t('content.tags.createAction') }}</BaseButton>
            </EmptyState>
          </div>
        </div>
        <div v-if="tagsQuery.data.value" class="studio-list-footer">
          <TablePagination
            :page="tagPagination.page.value"
            :page-size="tagPagination.pageSize.value"
            :total="tagTotal"
            :disabled="tagsQuery.isFetching.value"
            @update:page="tagPagination.setPage"
            @update:page-size="tagPagination.setPageSize"
          />
        </div>
      </div>
    </template>

    <SideDrawer :open="Boolean(selectedContentId)" :title="t('content.detailTitle')" :description="selectedContent?.slug" size="lg" @close="closeContentDrawer">
      <PageLoadingState v-if="detailQuery.showBlockingLoading.value && !selectedContent" :title="t('content.drawer.loadingTitle')" :rows="4" />
      <div v-else-if="selectedContent" class="content-page__drawer">
        <InlineLoadingState v-if="detailRefreshing" :title="t('common.refreshing')" />
        <div class="content-page__detail-grid">
          <ReadonlyField :label="t('content.columns.title')" :value="selectedContent.title" />
          <ReadonlyField :label="t('content.columns.slug')" :value="selectedContent.slug" />
          <ReadonlyField :label="t('content.columns.type')" :value="selectedContent.type" />
          <ReadonlyField :label="t('content.columns.status')" :value="selectedContent.status" />
          <ReadonlyField :label="t('content.columns.visibility')" :value="selectedContent.visibility" />
          <ReadonlyField :label="t('content.columns.updated')" :value="formatUnixTime(selectedContent.updated_at)" />
        </div>

        <section class="content-page__subsection">
          <h3>{{ t('content.drawer.relationsTitle') }}</h3>
          <div class="content-page__mini-list">
            <article v-for="relation in relations" :key="relation.relation_id">
              <span>{{ relation.relation_type }} -> {{ relation.to_content_id }}</span>
            </article>
            <p v-if="relations.length === 0">{{ t('content.drawer.relationsEmpty') }}</p>
          </div>
        </section>

        <section class="content-page__subsection">
          <h3>{{ t('content.drawer.revisionsTitle') }}</h3>
          <div class="content-page__mini-list">
            <article v-for="revision in revisions" :key="revision.revision_id">
              <span>#{{ revision.revision_no }} {{ revision.change_summary || t('content.drawer.noSummary') }}</span>
              <small>{{ formatUnixTime(revision.created_at) }}</small>
            </article>
            <p v-if="revisions.length === 0">{{ t('content.drawer.revisionsEmpty') }}</p>
          </div>
        </section>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="closeContentDrawer">{{ t('common.close') }}</BaseButton>
      </template>
    </SideDrawer>

    <SideDrawer
      :open="isTagDrawerOpen"
      :title="tagDrawerTitle"
      :description="selectedTag?.slug || t('content.tags.drawerDescription')"
      @close="closeTagDrawer"
    >
      <form class="content-page__tag-drawer" @submit.prevent="saveTag">
        <FormField :label="t('content.tags.fields.name')" for-id="tag-name">
          <BaseInput id="tag-name" v-model="tagForm.name" />
        </FormField>
        <FormField :label="t('content.tags.fields.slug')" for-id="tag-slug">
          <BaseInput id="tag-slug" v-model="tagForm.slug" />
        </FormField>
        <FormField :label="t('content.tags.fields.color')" for-id="tag-color">
          <BaseInput id="tag-color" v-model="tagForm.color" placeholder="#0f8f83" />
        </FormField>
        <FormField :label="t('content.tags.fields.description')" for-id="tag-description">
          <BaseInput id="tag-description" v-model="tagForm.description" />
        </FormField>
      </form>
      <template #footer>
        <BaseButton :busy="isMutating" @click="saveTag">{{ t('common.save') }}</BaseButton>
        <BaseButton variant="ghost" @click="closeTagDrawer">{{ t('common.close') }}</BaseButton>
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

.content-page__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
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

.content-page__refreshing {
  display: flex;
  justify-content: flex-end;
  padding: 12px 12px 0;
}

.content-page__grid {
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

.content-page__empty-panel {
  min-width: 920px;
  border-top: 1px solid var(--bb-color-line);
}

.content-page__empty-state {
  min-height: 232px;
  justify-items: center;
  text-align: center;
  border: 0;
  border-radius: 0;
  padding: 32px 24px;
  background: transparent;
  box-shadow: none;
}

.content-page__empty-action {
  margin-top: 4px;
}

.content-page__table tbody tr:last-child td {
  border-bottom: 0;
}

.content-page__table td:first-child {
  display: grid;
  gap: 3px;
}

.content-page__table td:first-child span,
.content-page__count {
  color: var(--bb-color-muted);
}

.content-page__empty-state :deep(.empty-state__visual) {
  color: var(--bb-color-muted);
}

.content-page__empty-state :deep(.empty-state__actions) {
  justify-content: center;
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

.content-page__tag-drawer,
.content-page__mini-list {
  display: grid;
  gap: 10px;
}

.content-page__mini-list article {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 12px;
  background: var(--bb-color-surface);
}

.content-page__tag-filters {
  grid-template-columns: minmax(220px, 420px);
}

.content-page__tag-table {
  min-width: 0;
}

.content-page__tag-grid {
  min-width: 920px;
}

.content-page__tag-table th:last-child,
.content-page__tag-table td:last-child {
  text-align: right;
}

.content-page__tag-table td:first-child strong {
  color: var(--bb-color-text-strong);
}

.content-page__tag-color {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.content-page__tag-swatch {
  width: 14px;
  height: 14px;
  border: 1px solid var(--bb-color-line);
  border-radius: 999px;
}

@media (max-width: 900px) {
  .content-page__filters,
  .content-page__tag-filters,
  .content-page__detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
