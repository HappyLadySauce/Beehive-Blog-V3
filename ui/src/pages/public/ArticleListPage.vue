<script setup lang="ts">
import { onMounted, ref } from 'vue';

import ContentPreviewCard from '@/features/content-preview/components/ContentPreviewCard.vue';
import { contentPreviewApi } from '@/shared/api/contentPreviewApi';
import type { ContentSummaryView } from '@/shared/api/types';

const items = ref<ContentSummaryView[]>([]);

onMounted(async () => {
  items.value = (await contentPreviewApi.listPublicContent({ type: 'article' })).items;
});
</script>

<template>
  <main class="mx-auto max-w-1180px px-4 py-10 sm:px-6 lg:px-8">
    <header class="mb-8 max-w-680px">
      <p class="m-0 text-13px font-700 text-brand-leaf">Articles</p>
      <h1 class="m-0 mt-2 text-32px font-900">文章列表</h1>
      <p class="m-0 mt-3 text-15px leading-7 text-brand-muted">这里先展示契约形态的 mock 内容，等待 content 服务稳定后切换真实接口。</p>
    </header>
    <div class="grid gap-4 md:grid-cols-2">
      <ContentPreviewCard v-for="item in items" :key="item.content_id" :item="item" />
    </div>
  </main>
</template>
