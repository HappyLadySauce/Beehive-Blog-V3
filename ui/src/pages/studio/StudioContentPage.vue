<script setup lang="ts">
import { Filter, RefreshCw } from 'lucide-vue-next';
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
import DataTable, { type TableColumn } from '@/shared/components/DataTable.vue';
import PageHeader from '@/shared/components/PageHeader.vue';
import EmptyState from '@/shared/components/EmptyState.vue';
import StatusAlert from '@/shared/components/StatusAlert.vue';
import { LoadingSkeleton } from '@/shared/components';

const items = ref<ContentSummaryView[]>([]);
const isLoading = ref(false);
const errorMessage = ref('');
const emptyMessage = ref('');
const authStore = useAuthStore();
const isMockMode = appConfig.apiMode !== 'live';
const filters = reactive({
  keyword: '',
  type: '',
  status: '',
  visibility: '',
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

const columns: TableColumn<ContentSummaryView>[] = [
  { key: 'title', label: '标题' },
  { key: 'type', label: '类型', width: '120px' },
  { key: 'status', label: '状态', width: '120px' },
];

async function load() {
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
      items.value = (await contentPreviewApi.listPublicContent(query)).items;
      if (items.value.length === 0) {
        emptyMessage.value = '当前 mock 数据为空。';
      }
      return;
    }

    if (!authStore.accessToken) {
      errorMessage.value = '未检测到可用 token，请重新登录后再试。';
      return;
    }

    items.value = (await contentPreviewApi.listStudioContents(query, authStore.accessToken)).items;
    if (items.value.length === 0) {
      emptyMessage.value = '当前没有可展示的 studio 内容。';
    }
  } catch (error) {
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
        <BaseButton variant="secondary" :disabled="isMockMode" title="content 创建接口稳定后启用">新建内容</BaseButton>
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

    <form class="grid gap-3 rounded-lg border border-white/8 bg-white/5 p-4 lg:grid-cols-[minmax(0,1fr)_150px_150px_150px_auto]" @submit.prevent="load">
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
