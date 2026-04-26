<script setup lang="ts">
import { Filter, Plus, RefreshCw } from 'lucide-vue-next';
import { computed, onMounted, reactive, ref } from 'vue';

import { GatewayHttpError } from '@/shared/api/httpClient';
import { contentPreviewApi } from '@/shared/api/contentPreviewApi';
import { appConfig } from '@/shared/config/env';
import { useAuthStore } from '@/features/auth/stores/authStore';
import type { ContentSummaryView, StudioContentListQuery } from '@/shared/api/types';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import BaseButton from '@/shared/components/BaseButton.vue';
import BaseInput from '@/shared/components/BaseInput.vue';
import BaseSelect, { type SelectOption } from '@/shared/components/BaseSelect.vue';
import BaseTextarea from '@/shared/components/BaseTextarea.vue';
import DataTable, { type TableColumn } from '@/shared/components/DataTable.vue';
import PageHeader from '@/shared/components/PageHeader.vue';
import EmptyState from '@/shared/components/EmptyState.vue';
import StatusAlert from '@/shared/components/StatusAlert.vue';
import { LoadingSkeleton } from '@/shared/components';

const items = ref<ContentSummaryView[]>([]);
const localDrafts = ref<ContentSummaryView[]>([]);
const isLoading = ref(false);
const errorMessage = ref('');
const emptyMessage = ref('');
const composerError = ref('');
const createNotice = ref('');
const isComposerOpen = ref(false);
const authStore = useAuthStore();
const isMockMode = appConfig.apiMode !== 'live';
const filters = reactive({
  keyword: '',
  type: '',
  status: '',
  visibility: '',
});
const draftForm = reactive({
  title: '',
  type: 'article',
  summary: '',
  visibility: 'private',
});
const itemCountLabel = computed(() => `${items.value.length} 条内容`);
const typeOptions: SelectOption[] = [
  { label: '全部类型', value: '' },
  { label: '文章', value: 'article' },
  { label: '笔记', value: 'note' },
  { label: '项目', value: 'project' },
  { label: '洞察', value: 'insight' },
];
const statusOptions: SelectOption[] = [
  { label: '全部状态', value: '' },
  { label: '草稿', value: 'draft' },
  { label: '审阅中', value: 'review' },
  { label: '已发布', value: 'published' },
  { label: '已归档', value: 'archived' },
];
const visibilityOptions: SelectOption[] = [
  { label: '全部可见性', value: '' },
  { label: '公开', value: 'public' },
  { label: '成员', value: 'member' },
  { label: '私密', value: 'private' },
];

function getErrorMessage(error: unknown): string {
  if (error instanceof GatewayHttpError && error.response) {
    const code = error.response.code;
    if (error.status === 401) {
      return '登录态失效，请先重新登录后再访问 Studio 内容管理。';
    }
    if (error.status === 403 || code === 120301) {
      return '当前账号无内容管理权限，请使用 admin 账号。';
    }
    return error.response.message ?? error.message;
  }
  if (error instanceof Error) {
    return error.message;
  }
  return '内容加载失败，请稍后重试。';
}

function shouldRefreshSession(error: unknown): boolean {
  return error instanceof GatewayHttpError && (error.status === 403 || error.response?.code === 120301);
}

const columns: TableColumn<ContentSummaryView>[] = [
  { key: 'title', label: '标题' },
  { key: 'type', label: '类型', width: '120px' },
  { key: 'status', label: '状态', width: '120px' },
];

function buildLocalSlug(title: string): string {
  return (
    title
      .trim()
      .toLowerCase()
      .replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-')
      .replace(/^-+|-+$/g, '')
      .slice(0, 48) || `local-draft-${Date.now()}`
  );
}

function matchesFilters(item: ContentSummaryView): boolean {
  const keyword = filters.keyword.trim().toLowerCase();
  if (keyword) {
    const haystack = `${item.title} ${item.summary} ${item.slug}`.toLowerCase();
    if (!haystack.includes(keyword)) {
      return false;
    }
  }
  if (filters.type && item.type !== filters.type) {
    return false;
  }
  if (filters.status && item.status !== filters.status) {
    return false;
  }
  if (filters.visibility && item.visibility !== filters.visibility) {
    return false;
  }
  return true;
}

function applyLoadedItems(loadedItems: ContentSummaryView[]) {
  const visibleLocalDrafts = localDrafts.value.filter(matchesFilters);
  items.value = [...visibleLocalDrafts, ...loadedItems];
}

async function load(hasRetriedSessionRefresh = false) {
  isLoading.value = true;
  errorMessage.value = '';
  emptyMessage.value = '';
  items.value = [];

  try {
    const query: StudioContentListQuery = {};
    const keyword = filters.keyword.trim();
    if (keyword) {
      query.keyword = keyword;
    }
    if (filters.type) {
      query.type = filters.type;
    }
    if (filters.status) {
      query.status = filters.status;
    }
    if (filters.visibility) {
      query.visibility = filters.visibility;
    }

    if (isMockMode) {
      applyLoadedItems((await contentPreviewApi.listPublicContent(query)).items);
      if (items.value.length === 0) {
        emptyMessage.value = '当前 mock 数据为空。';
      }
      return;
    }

    if (!authStore.accessToken) {
      errorMessage.value = '未检测到可用 token，请重新登录后再试。';
      return;
    }

    applyLoadedItems((await contentPreviewApi.listStudioContents(query, authStore.accessToken)).items);
    if (items.value.length === 0) {
      emptyMessage.value = '当前没有可展示的 studio 内容。';
    }
  } catch (error) {
    if (!hasRetriedSessionRefresh && !isMockMode && shouldRefreshSession(error) && (await authStore.refreshSession())) {
      await load(true);
      return;
    }
    errorMessage.value = getErrorMessage(error);
  } finally {
    isLoading.value = false;
  }
}

function clearFilters() {
  filters.keyword = '';
  filters.type = '';
  filters.status = '';
  filters.visibility = '';
  void load();
}

function submitFilters() {
  void load();
}

function openCreateComposer() {
  composerError.value = '';
  createNotice.value = '';
  isComposerOpen.value = true;
}

function resetDraftForm() {
  draftForm.title = '';
  draftForm.type = 'article';
  draftForm.summary = '';
  draftForm.visibility = 'private';
}

function createLocalDraft() {
  const title = draftForm.title.trim();
  const summary = draftForm.summary.trim();
  if (!title) {
    composerError.value = '请输入标题后再创建草稿。';
    return;
  }

  const now = Math.floor(Date.now() / 1000);
  const draft: ContentSummaryView = {
    content_id: `local_${now}_${localDrafts.value.length + 1}`,
    type: draftForm.type,
    title,
    slug: buildLocalSlug(title),
    summary: summary || '本地 UI 草稿，content 创建接口稳定后再接入真实保存。',
    cover_image_url: '',
    status: 'draft',
    visibility: draftForm.visibility,
    ai_access: 'none',
    published_at: 0,
    archived_at: 0,
    created_at: now,
    updated_at: now,
    tags: [],
  };

  localDrafts.value = [draft, ...localDrafts.value];
  composerError.value = '';
  createNotice.value = '本地草稿已加入当前列表，页面刷新后不会保留。';
  resetDraftForm();
  isComposerOpen.value = false;
  applyLoadedItems(items.value.filter((item) => !item.content_id.startsWith('local_')));
}

onMounted(() => {
  void load();
});
</script>

<template>
  <main class="grid gap-5 p-4 sm:p-6 lg:p-8">
    <PageHeader
      inverse
      :eyebrow="isMockMode ? 'Content Mock' : 'Studio Content'"
      title="内容中心"
      :description="
        isMockMode
          ? '展示管理台信息密度和表格布局，不接 content 真实服务。'
          : '展示 studio 内容列表，支持管理员鉴权与分页。'
      "
    >
      <template #actions>
        <BaseButton variant="secondary" title="创建本地 UI 草稿" @click="openCreateComposer">
          <Plus class="h-4 w-4" aria-hidden="true" />
          新建内容
        </BaseButton>
        <BaseButton variant="ghost" :busy="isLoading" @click="load">
          <RefreshCw class="h-4 w-4" aria-hidden="true" />
          刷新
        </BaseButton>
      </template>
    </PageHeader>

    <StatusAlert
      :tone="isMockMode ? 'warning' : 'success'"
      :title="isMockMode ? '当前使用 mock content adapter' : 'Studio 内容列表已接入真实 Gateway'"
      :description="
        isMockMode
          ? 'content 服务开发完成前，Studio 内容中心只验证 UI 结构、筛选和状态展示。'
          : '当前页面已接入 /api/v3/studio/content/items，支持鉴权与分页参数。'
      "
    />

    <StatusAlert v-if="createNotice" tone="success" title="草稿已创建" :description="createNotice" />

    <section v-if="isComposerOpen" class="grid gap-4 rounded-lg border border-white/10 bg-white/5 p-4" aria-label="新建内容草稿">
      <div>
        <h2 class="m-0 text-16px font-900 text-white">新建内容草稿</h2>
        <p class="m-0 mt-1 text-13px leading-5 text-white/48">当前只创建本地 UI 草稿，不调用 content 创建接口。</p>
      </div>
      <div class="grid gap-3 lg:grid-cols-[minmax(0,1fr)_150px_150px]">
        <BaseInput v-model="draftForm.title" label="标题" name="draft_title" placeholder="输入内容标题" />
        <BaseSelect v-model="draftForm.type" label="类型" name="draft_type" :options="typeOptions.filter((option) => option.value)" />
        <BaseSelect v-model="draftForm.visibility" label="可见性" name="draft_visibility" :options="visibilityOptions.filter((option) => option.value)" />
      </div>
      <BaseTextarea v-model="draftForm.summary" label="摘要" name="draft_summary" placeholder="补充摘要，留空时使用默认说明" />
      <StatusAlert v-if="composerError" tone="danger" title="无法创建草稿" :description="composerError" />
      <div class="flex flex-wrap gap-2">
        <BaseButton variant="primary" @click="createLocalDraft">保存本地草稿</BaseButton>
        <BaseButton variant="ghost" @click="isComposerOpen = false">取消</BaseButton>
      </div>
    </section>

    <form class="grid gap-3 rounded-lg border border-white/8 bg-white/5 p-4 lg:grid-cols-[minmax(0,1fr)_150px_150px_150px_auto]" @submit.prevent="submitFilters">
      <BaseInput v-model="filters.keyword" label="关键词" name="studio_keyword" placeholder="标题、摘要或 slug" />
      <BaseSelect v-model="filters.type" label="类型" name="studio_type" :options="typeOptions" />
      <BaseSelect v-model="filters.status" label="状态" name="studio_status" :options="statusOptions" />
      <BaseSelect v-model="filters.visibility" label="可见性" name="studio_visibility" :options="visibilityOptions" />
      <div class="flex items-end gap-2">
        <BaseButton type="submit" variant="secondary" :busy="isLoading">
          <Filter class="h-4 w-4" aria-hidden="true" />
          筛选
        </BaseButton>
        <BaseButton type="button" variant="ghost" @click="clearFilters">清空</BaseButton>
      </div>
    </form>

    <LoadingSkeleton v-if="isLoading" :rows="6" />
    <StatusAlert v-else-if="errorMessage" tone="danger" title="内容加载失败" :description="errorMessage" />

    <EmptyState v-else-if="emptyMessage" :title="emptyMessage" description="可以先发布一篇草稿内容，再刷新页面查看。" />
    <div v-else class="grid gap-3">
      <p class="m-0 text-13px text-white/45">{{ itemCountLabel }}</p>
      <DataTable :columns="columns" :rows="items" row-key="content_id" inverse>
        <template #cell-title="{ row }">
          <div class="min-w-0 max-w-560px">
            <p class="m-0 truncate text-14px font-800">{{ row.title }}</p>
            <p class="m-0 mt-1 line-clamp-2 text-13px leading-5 text-white/48">{{ row.summary }}</p>
          </div>
        </template>
        <template #cell-type="{ row }">
          <BaseBadge tone="neutral">{{ row.type }}</BaseBadge>
        </template>
        <template #cell-status="{ row }">
          <BaseBadge :tone="row.status === 'published' ? 'leaf' : 'honey'">{{ row.status }}</BaseBadge>
        </template>
      </DataTable>
    </div>
  </main>
</template>
