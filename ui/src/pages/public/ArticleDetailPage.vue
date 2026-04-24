<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import { contentPreviewApi } from '@/shared/api/contentPreviewApi';
import type { ContentSummaryView } from '@/shared/api/types';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import EmptyState from '@/shared/components/EmptyState.vue';

const route = useRoute();
const item = ref<ContentSummaryView | null>(null);
const slug = computed(() => String(route.params.slug));

onMounted(async () => {
  const response = await contentPreviewApi.listPublicContent();
  item.value = response.items.find((entry) => entry.slug === slug.value) ?? null;
});
</script>

<template>
  <main class="mx-auto grid max-w-1180px gap-6 px-4 py-10 sm:px-6 lg:grid-cols-[minmax(0,1fr)_280px] lg:px-8">
    <article v-if="item" class="bb-panel p-5 sm:p-8">
      <BaseBadge tone="leaf">Mock Article</BaseBadge>
      <h1 class="m-0 mt-4 text-34px font-900 leading-11">{{ item.title }}</h1>
      <p class="m-0 mt-4 text-16px leading-7 text-brand-muted">{{ item.summary }}</p>
      <div class="bb-reading mt-8 text-15px text-brand-ink">
        <p>这是文章详情页的首版渲染壳，用于验证阅读排版、侧栏目录和响应式布局。</p>
        <h2>内容边界</h2>
        <p>公开站只消费 gateway 暴露的 HTTP 契约，不直接访问后端 RPC 服务。</p>
        <h2>后续接入</h2>
        <p>当 content 服务稳定后，当前 mock adapter 会替换为真实 public content API。</p>
      </div>
    </article>
    <aside v-if="item" class="bb-panel h-max p-4 lg:sticky lg:top-24">
      <p class="m-0 text-13px font-800">目录</p>
      <nav class="mt-3 grid gap-2 text-13px text-brand-muted">
        <a href="#" class="bb-focus rounded-sm px-2 py-1 text-brand-leaf">内容边界</a>
        <a href="#" class="bb-focus rounded-sm px-2 py-1">后续接入</a>
      </nav>
    </aside>
    <EmptyState v-else title="文章不存在" description="当前 mock 数据中没有匹配的 slug。" />
  </main>
</template>
