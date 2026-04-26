<script setup lang="ts">
import { ArrowUpRight } from 'lucide-vue-next';
import { computed } from 'vue';

import type { ContentSummaryView } from '@/shared/api/types';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import BaseCard from '@/shared/components/BaseCard.vue';

const props = defineProps<{
  item: ContentSummaryView;
}>();

const detailPath = computed(() => (props.item.type === 'article' ? `/articles/${props.item.slug}` : '/projects'));
const typeLabel = computed(() => {
  const labels: Record<string, string> = {
    article: '文章',
    project: '项目',
    note: '笔记',
    insight: '洞察',
  };
  return labels[props.item.type] ?? props.item.type;
});
const publishedLabel = computed(() => {
  if (!props.item.published_at) {
    return '未发布';
  }
  return new Intl.DateTimeFormat('zh-CN', { month: 'short', day: 'numeric' }).format(new Date(props.item.published_at * 1000));
});
const readingMinutes = computed(() => Math.max(1, Math.ceil((props.item.summary.length + props.item.title.length) / 120)));
</script>

<template>
  <RouterLink :to="detailPath" class="group block h-full">
    <BaseCard interactive class="grid h-full gap-4 p-5">
      <div class="flex items-start justify-between gap-4">
        <BaseBadge :tone="item.type === 'project' ? 'blue' : 'leaf'">{{ typeLabel }}</BaseBadge>
        <ArrowUpRight class="h-4 w-4 shrink-0 text-brand-muted transition-colors group-hover:text-brand-blue" aria-hidden="true" />
      </div>
      <div class="grid gap-2">
        <h3 class="m-0 text-20px font-800 leading-7 text-brand-ink">{{ item.title }}</h3>
        <p class="m-0 text-14px leading-6 text-brand-muted">{{ item.summary }}</p>
      </div>
      <div class="flex flex-wrap items-center gap-2 text-12px text-brand-muted">
        <span>{{ publishedLabel }}</span>
        <span aria-hidden="true">/</span>
        <span>{{ readingMinutes }} 分钟阅读</span>
        <span aria-hidden="true">/</span>
        <span>{{ item.visibility }}</span>
      </div>
      <div class="mt-auto flex flex-wrap gap-2">
        <BaseBadge v-for="tag in item.tags" :key="tag.tag_id" tone="neutral">{{ tag.name }}</BaseBadge>
      </div>
    </BaseCard>
  </RouterLink>
</template>
