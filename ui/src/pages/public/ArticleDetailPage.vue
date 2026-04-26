<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { contentPreviewApi } from '@/shared/api/contentPreviewApi';
import type { ContentDetailView } from '@/shared/api/types';
import { GatewayHttpError } from '@/shared/api/httpClient';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import EmptyState from '@/shared/components/EmptyState.vue';
import { LoadingSkeleton } from '@/shared/components';

const route = useRoute();
const item = ref<ContentDetailView | null>(null);
const slug = computed(() => String(route.params.slug));
const isLoading = ref(false);
const notFound = ref(false);
const errorMessage = ref('');

const statusText = computed(() => {
  if (!item.value) {
    return '';
  }
  return `${item.value.type.toUpperCase()} / ${item.value.status} / ${item.value.visibility}`;
});
const bodyBlocks = computed(() => {
  if (!item.value?.body_markdown) {
    return [];
  }
  return item.value.body_markdown
    .split(/\n{2,}/)
    .map((block) => block.trim())
    .filter((block) => block.replace(/^#\s+/, '') !== item.value?.title)
    .filter(Boolean);
});

function isContentNotFound(error: unknown): boolean {
  return (
    (error instanceof GatewayHttpError && (error.status === 404 || error.response?.code === 120501)) ||
    (error instanceof Error && error.message.includes('not found'))
  );
}

function errorLabel(error: unknown): string {
  if (error instanceof GatewayHttpError && error.response) {
    return error.response.message ?? error.message;
  }
  if (error instanceof Error) {
    return error.message;
  }
  return '加载文章内容失败，请稍后重试。';
}

function formatUnixTime(value?: number): string {
  if (!value) {
    return '未更新';
  }
  return new Intl.DateTimeFormat('zh-CN', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value * 1000));
}

async function load() {
  isLoading.value = true;
  notFound.value = false;
  errorMessage.value = '';
  item.value = null;

  try {
    const response = await contentPreviewApi.getPublicContentBySlug(slug.value);
    item.value = response.content;
  } catch (error) {
    if (isContentNotFound(error)) {
      notFound.value = true;
      return;
    }
    errorMessage.value = errorLabel(error);
  } finally {
    isLoading.value = false;
  }
}

onMounted(load);
watch(() => route.params.slug, () => {
  void load();
});
</script>

<template>
  <main class="mx-auto grid max-w-1180px gap-6 px-4 py-10 sm:px-6 lg:grid-cols-[minmax(0,1fr)_280px] lg:px-8">
    <LoadingSkeleton v-if="isLoading" class="lg:col-span-2" :rows="8" />
    <article v-else-if="item" class="bb-panel p-5 sm:p-8">
      <BaseBadge tone="leaf">{{ statusText }}</BaseBadge>
      <h1 class="m-0 mt-4 text-34px font-900 leading-11">{{ item.title }}</h1>
      <p class="m-0 mt-3 text-14px text-brand-muted">{{ item.summary }}</p>
      <div class="mt-4 flex flex-wrap gap-3 text-12px text-brand-muted">
        <span>Published：{{ formatUnixTime(item.published_at) }}</span>
        <span>Updated：{{ formatUnixTime(item.updated_at) }}</span>
        <span>Owner：{{ item.owner_user_id }}</span>
      </div>
      <div class="bb-reading mt-8 grid gap-4 text-15px text-brand-ink">
        <template v-for="block in bodyBlocks" :key="block">
          <h2 v-if="block.startsWith('## ')" class="m-0 mt-2 text-24px font-900 leading-8">
            {{ block.replace(/^##\s+/, '') }}
          </h2>
          <h2 v-else-if="block.startsWith('# ')" class="m-0 mt-2 text-26px font-900 leading-9">
            {{ block.replace(/^#\s+/, '') }}
          </h2>
          <pre v-else-if="block.startsWith('```')" class="m-0 overflow-x-auto rounded-lg border border-brand-line bg-brand-paper p-4 text-13px leading-6"><code>{{ block.replace(/^```[a-zA-Z]*\n?/, '').replace(/```$/, '') }}</code></pre>
          <p v-else class="m-0 whitespace-pre-line text-15px leading-8">
            {{ block }}
          </p>
        </template>
      </div>
    </article>
    <aside v-if="item" class="bb-panel h-max p-5">
      <h2 class="m-0 text-15px font-900">内容信息</h2>
      <dl class="mt-4 grid gap-3 text-13px">
        <div class="grid gap-1">
          <dt class="text-brand-muted">类型</dt>
          <dd class="m-0 font-700">{{ item.type }}</dd>
        </div>
        <div class="grid gap-1">
          <dt class="text-brand-muted">可见性</dt>
          <dd class="m-0 font-700">{{ item.visibility }}</dd>
        </div>
        <div class="grid gap-1">
          <dt class="text-brand-muted">AI 访问</dt>
          <dd class="m-0 font-700">{{ item.ai_access }}</dd>
        </div>
      </dl>
    </aside>
    <article v-else-if="errorMessage" class="bb-panel p-5 text-brand-ink sm:p-8">
      <p class="m-0 text-13px text-red-600">{{ errorMessage }}</p>
    </article>
    <EmptyState v-else-if="notFound" title="文章未找到" description="该 slug 不存在或该内容未发布公开，无法在公共端访问。" />
    <EmptyState v-else title="加载中断" description="内容信息尚未返回。" />
  </main>
</template>
