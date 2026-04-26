<script setup lang="ts">
import { Search, X } from 'lucide-vue-next';
import { computed, onMounted, reactive, ref } from 'vue';

import ContentPreviewCard from '@/features/content-preview/components/ContentPreviewCard.vue';
import { contentPreviewApi } from '@/shared/api/contentPreviewApi';
import { appConfig } from '@/shared/config/env';
import { GatewayHttpError } from '@/shared/api/httpClient';
import { LoadingSkeleton } from '@/shared/components';
import type { ContentSummaryView, PublicContentQuery } from '@/shared/api/types';
import BaseButton from '@/shared/components/BaseButton.vue';
import BaseInput from '@/shared/components/BaseInput.vue';
import BaseSelect, { type SelectOption } from '@/shared/components/BaseSelect.vue';
import EmptyState from '@/shared/components/EmptyState.vue';
import PageHeader from '@/shared/components/PageHeader.vue';
import StatusAlert from '@/shared/components/StatusAlert.vue';

const items = ref<ContentSummaryView[]>([]);
const isLoading = ref(false);
const errorMessage = ref('');
const filters = reactive({
  keyword: '',
  type: 'article',
});

const modeLabel = computed(() => (appConfig.apiMode === 'live' ? 'live' : 'mock'));
const typeOptions: SelectOption[] = [
  { label: '文章', value: 'article' },
  { label: '笔记', value: 'note' },
  { label: '项目', value: 'project' },
  { label: '洞察', value: 'insight' },
];

function errorLabel(error: unknown): string {
  if (error instanceof GatewayHttpError) {
    return error.message;
  }
  return '加载文章列表失败，请稍后重试。';
}

async function loadArticles() {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    const query: PublicContentQuery = { type: filters.type };
    const keyword = filters.keyword.trim();
    if (keyword) {
      query.keyword = keyword;
    }
    items.value = (await contentPreviewApi.listPublicContent(query)).items;
  } catch (error) {
    errorMessage.value = errorLabel(error);
  } finally {
    isLoading.value = false;
  }
}

function resetFilters() {
  filters.keyword = '';
  filters.type = 'article';
  void loadArticles();
}

onMounted(() => {
  void loadArticles();
});
</script>

<template>
  <main class="mx-auto max-w-1180px px-4 py-10 sm:px-6 lg:px-8">
    <PageHeader
      class="mb-6"
      eyebrow="Articles"
      title="内容浏览"
      :description="`数据源：${modeLabel}。可按关键词与类型筛选公开内容。`"
    />

    <form class="mb-6 grid gap-3 rounded-lg border border-brand-line bg-brand-surface p-4 md:grid-cols-[minmax(0,1fr)_180px_auto]" @submit.prevent="loadArticles">
      <BaseInput v-model="filters.keyword" label="关键词" name="keyword" placeholder="搜索标题、摘要或 slug" />
      <BaseSelect v-model="filters.type" label="类型" name="type" :options="typeOptions" />
      <div class="flex items-end gap-2">
        <BaseButton type="submit" variant="primary" :busy="isLoading">
          <Search class="h-4 w-4" aria-hidden="true" />
          查询
        </BaseButton>
        <BaseButton type="button" variant="ghost" @click="resetFilters">
          <X class="h-4 w-4" aria-hidden="true" />
          重置
        </BaseButton>
      </div>
    </form>

    <LoadingSkeleton v-if="isLoading" :rows="3" />
    <StatusAlert v-else-if="errorMessage" tone="danger" title="内容列表加载失败" :description="errorMessage" />
    <EmptyState v-else-if="items.length === 0" title="没有匹配内容" description="可以换一个关键词或类型再试。" />
    <div v-else class="grid gap-4 md:grid-cols-2">
      <ContentPreviewCard v-for="item in items" :key="item.content_id" :item="item" />
    </div>
  </main>
</template>
