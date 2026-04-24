<script setup lang="ts">
import { onMounted, ref } from 'vue';

import ContentPreviewCard from '@/features/content-preview/components/ContentPreviewCard.vue';
import { contentPreviewApi } from '@/shared/api/contentPreviewApi';
import type { ContentSummaryView } from '@/shared/api/types';

const projects = ref<ContentSummaryView[]>([]);

onMounted(async () => {
  projects.value = (await contentPreviewApi.listPublicContent({ type: 'project' })).items;
});
</script>

<template>
  <main class="mx-auto max-w-1180px px-4 py-10 sm:px-6 lg:px-8">
    <header class="mb-8 max-w-680px">
      <p class="m-0 text-13px font-700 text-brand-blue">Projects</p>
      <h1 class="m-0 mt-2 text-32px font-900">项目与作品</h1>
      <p class="m-0 mt-3 text-15px leading-7 text-brand-muted">项目页会承载背景、目标、技术栈、关联文章和阶段结果。</p>
    </header>
    <div class="grid gap-4 md:grid-cols-2">
      <ContentPreviewCard v-for="item in projects" :key="item.content_id" :item="item" />
    </div>
  </main>
</template>
