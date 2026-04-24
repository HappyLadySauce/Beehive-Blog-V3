<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { contentPreviewApi } from '@/shared/api/contentPreviewApi';
import type { ContentSummaryView } from '@/shared/api/types';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import BaseButton from '@/shared/components/BaseButton.vue';
import BaseCard from '@/shared/components/BaseCard.vue';

const items = ref<ContentSummaryView[]>([]);

onMounted(async () => {
  items.value = (await contentPreviewApi.listPublicContent()).items;
});
</script>

<template>
  <main class="grid gap-5 p-4 sm:p-6 lg:p-8">
    <header class="flex flex-wrap items-end justify-between gap-3">
      <div>
        <BaseBadge tone="blue">Content Mock</BaseBadge>
        <h1 class="m-0 mt-3 text-28px font-900 text-white">内容中心</h1>
        <p class="m-0 mt-2 text-14px leading-6 text-white/52">展示管理台信息密度和表格布局，不接 content 真实服务。</p>
      </div>
      <BaseButton variant="secondary">新建内容</BaseButton>
    </header>

    <BaseCard class="overflow-hidden border-white/8 bg-white/5 text-white">
      <div class="grid grid-cols-[1.4fr_120px_120px] gap-3 border-b border-white/8 px-4 py-3 text-12px font-700 text-white/45 max-md:hidden">
        <span>标题</span>
        <span>类型</span>
        <span>状态</span>
      </div>
      <div v-for="item in items" :key="item.content_id" class="grid gap-2 border-b border-white/6 px-4 py-4 last:border-b-0 md:grid-cols-[1.4fr_120px_120px] md:items-center">
        <div class="min-w-0">
          <p class="m-0 truncate text-14px font-800">{{ item.title }}</p>
          <p class="m-0 mt-1 line-clamp-2 text-13px leading-5 text-white/48">{{ item.summary }}</p>
        </div>
        <BaseBadge tone="neutral">{{ item.type }}</BaseBadge>
        <BaseBadge :tone="item.status === 'published' ? 'leaf' : 'honey'">{{ item.status }}</BaseBadge>
      </div>
    </BaseCard>
  </main>
</template>
