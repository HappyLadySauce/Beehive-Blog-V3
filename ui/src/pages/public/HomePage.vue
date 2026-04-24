<script setup lang="ts">
import { ArrowRight, BrainCircuit, Layers3, Sparkles } from 'lucide-vue-next';
import { onMounted, ref } from 'vue';

import ContentPreviewCard from '@/features/content-preview/components/ContentPreviewCard.vue';
import { contentPreviewApi } from '@/shared/api/contentPreviewApi';
import type { ContentSummaryView } from '@/shared/api/types';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import BaseButton from '@/shared/components/BaseButton.vue';
import BaseCard from '@/shared/components/BaseCard.vue';

const latestItems = ref<ContentSummaryView[]>([]);

onMounted(async () => {
  latestItems.value = (await contentPreviewApi.listPublicContent({ page_size: 3 })).items;
});
</script>

<template>
  <main>
    <section class="bb-grid-bg border-b border-brand-line bg-brand-surface">
      <div class="mx-auto grid max-w-1180px gap-8 px-4 py-10 sm:px-6 md:grid-cols-[1.05fr_0.95fr] md:py-14 lg:px-8">
        <div class="grid content-center gap-6">
          <BaseBadge tone="honey">Public Web + Studio + AI</BaseBadge>
          <div class="grid gap-4">
            <h1 class="m-0 max-w-12em text-36px font-900 leading-11 text-brand-ink md:text-48px md:leading-14">
              把创作、沉淀与公开表达放进同一个知识系统
            </h1>
            <p class="m-0 max-w-58ch text-16px leading-7 text-brand-muted md:text-18px">
              Beehive Blog v3 首屏直接呈现内容入口、项目线索与知识流转状态，避免空泛营销页。
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
          <BaseCard class="p-4">
            <div class="flex items-center gap-3">
              <Sparkles class="h-5 w-5 text-brand-honey" aria-hidden="true" />
              <div>
                <p class="m-0 text-13px text-brand-muted">今日焦点</p>
                <p class="m-0 text-18px font-800">内容沉淀链路</p>
              </div>
            </div>
          </BaseCard>
          <div class="grid gap-3 sm:grid-cols-2">
            <BaseCard class="p-4">
              <Layers3 class="h-5 w-5 text-brand-blue" aria-hidden="true" />
              <p class="m-0 mt-4 text-24px font-900">3</p>
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
        </div>
        <RouterLink to="/articles" class="bb-focus rounded-md text-14px font-700 text-brand-blue">查看全部</RouterLink>
      </div>
      <div class="grid gap-4 md:grid-cols-3">
        <ContentPreviewCard v-for="item in latestItems" :key="item.content_id" :item="item" />
      </div>
    </section>
  </main>
</template>
