<script setup lang="ts">
import { ArrowRight, BrainCircuit, Layers3, RefreshCw, Sparkles } from 'lucide-vue-next';
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import ContentPreviewCard from '@/features/content-preview/components/ContentPreviewCard.vue';
import { contentPreviewApi } from '@/shared/api/contentPreviewApi';
import { appConfig } from '@/shared/config/env';
import { GatewayHttpError } from '@/shared/api/httpClient';
import type { ContentSummaryView } from '@/shared/api/types';
import EmptyState from '@/shared/components/EmptyState.vue';
import { LoadingSkeleton } from '@/shared/components';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import BaseButton from '@/shared/components/BaseButton.vue';
import BaseCard from '@/shared/components/BaseCard.vue';
import StatusAlert from '@/shared/components/StatusAlert.vue';

const latestItems = ref<ContentSummaryView[]>([]);
const route = useRoute();
const isLoading = ref(false);
const errorMessage = ref('');
const modeLabel = computed(() => (appConfig.apiMode === 'live' ? 'live' : 'mock'));
const featuredItem = computed(() => latestItems.value[0] ?? null);
const previewItems = computed(() => latestItems.value.slice(1, 4));
const studioForbidden = computed(() => route.query.studio === 'forbidden');

function errorLabel(error: unknown): string {
  if (error instanceof GatewayHttpError && error.response) {
    return error.response.message ?? error.message;
  }
  if (error instanceof Error) {
    return error.message;
  }
  return '首页内容加载失败，请稍后重试。';
}

async function loadLatest() {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    latestItems.value = (await contentPreviewApi.listPublicContent({ page_size: 4 })).items;
  } catch (error) {
    errorMessage.value = errorLabel(error);
  } finally {
    isLoading.value = false;
  }
}

onMounted(() => {
  void loadLatest();
});
</script>

<template>
  <main>
    <div v-if="studioForbidden" class="mx-auto max-w-1180px px-4 pt-4 sm:px-6 lg:px-8">
      <StatusAlert tone="warning" title="无法进入 Studio" description="Studio 仅管理员可访问。普通账号可以继续浏览公开内容。" />
    </div>
    <section class="bb-grid-bg border-b border-brand-line bg-brand-surface">
      <div class="mx-auto grid max-w-1180px gap-8 px-4 py-10 sm:px-6 md:grid-cols-[0.95fr_1.05fr] md:py-14 lg:px-8">
        <div class="grid content-center gap-6">
          <BaseBadge tone="honey">Beehive Blog v3</BaseBadge>
          <div class="grid gap-4">
            <h1 class="m-0 max-w-12em text-36px font-900 leading-11 text-brand-ink md:text-48px md:leading-14">
              用真实内容驱动公开表达与创作工作台
            </h1>
            <p class="m-0 max-w-58ch text-16px leading-7 text-brand-muted md:text-18px">
              首页优先展示精选文章、项目线索与最新动态，Studio 负责生产、审阅和发布。
            </p>
          </div>
          <div class="flex flex-wrap gap-3">
            <RouterLink to="/articles">
              <BaseButton variant="primary">
                阅读文章
                <ArrowRight class="h-4 w-4" aria-hidden="true" />
              </BaseButton>
            </RouterLink>
            <RouterLink to="/studio">
              <BaseButton variant="secondary">
                进入 Studio
              </BaseButton>
            </RouterLink>
          </div>
        </div>

        <div class="grid gap-3">
          <BaseCard v-if="featuredItem" interactive class="p-5">
            <RouterLink :to="featuredItem.type === 'article' ? `/articles/${featuredItem.slug}` : '/projects'" class="grid gap-4 text-brand-ink">
              <div class="flex items-center justify-between gap-3">
                <BaseBadge :tone="featuredItem.type === 'project' ? 'blue' : 'leaf'">精选内容</BaseBadge>
                <Sparkles class="h-5 w-5 text-brand-honey" aria-hidden="true" />
              </div>
              <div>
                <h2 class="m-0 text-24px font-900 leading-8">{{ featuredItem.title }}</h2>
                <p class="m-0 mt-3 line-clamp-3 text-14px leading-6 text-brand-muted">{{ featuredItem.summary }}</p>
              </div>
              <div class="flex flex-wrap gap-2">
                <BaseBadge v-for="tag in featuredItem.tags" :key="tag.tag_id" tone="neutral">{{ tag.name }}</BaseBadge>
              </div>
            </RouterLink>
          </BaseCard>
          <LoadingSkeleton v-else-if="isLoading" :rows="4" />
          <StatusAlert v-else-if="errorMessage" tone="warning" title="公共内容暂不可用" :description="errorMessage" />
          <div class="grid gap-3 sm:grid-cols-2">
            <BaseCard class="p-4">
              <Layers3 class="h-5 w-5 text-brand-blue" aria-hidden="true" />
              <p class="m-0 mt-4 text-24px font-900">{{ latestItems.length }}</p>
              <p class="m-0 text-13px text-brand-muted">首批内容类型预览</p>
            </BaseCard>
            <BaseCard class="p-4">
              <BrainCircuit class="h-5 w-5 text-brand-leaf" aria-hidden="true" />
              <p class="m-0 mt-4 text-24px font-900">AI</p>
              <p class="m-0 text-13px text-brand-muted">可审阅协作边界</p>
            </BaseCard>
          </div>
        </div>
      </div>
    </section>

    <section class="mx-auto max-w-1180px px-4 py-10 sm:px-6 lg:px-8">
      <div class="mb-5 flex flex-wrap items-end justify-between gap-3">
        <div>
          <p class="m-0 text-13px font-700 text-brand-leaf">Latest</p>
          <h2 class="m-0 mt-1 text-26px font-900">最新内容</h2>
          <p class="m-0 mt-2 text-12px text-brand-muted">数据源：{{ modeLabel }}</p>
        </div>
        <div class="flex flex-wrap gap-2">
          <BaseButton variant="ghost" size="sm" :busy="isLoading" @click="loadLatest">
            <RefreshCw class="h-4 w-4" aria-hidden="true" />
            刷新
          </BaseButton>
          <RouterLink to="/articles" class="bb-focus rounded-md text-14px font-700 text-brand-blue">查看全部</RouterLink>
        </div>
      </div>
      <LoadingSkeleton v-if="isLoading" :rows="3" />
      <StatusAlert v-else-if="errorMessage" tone="danger" title="最新内容加载失败" :description="errorMessage" />
      <EmptyState v-else-if="latestItems.length === 0" title="暂无内容" description="当前没有可公开展示的内容。" />
      <div v-else class="grid gap-4 md:grid-cols-3">
        <ContentPreviewCard v-for="item in previewItems.length > 0 ? previewItems : latestItems" :key="item.content_id" :item="item" />
      </div>
    </section>
  </main>
</template>
