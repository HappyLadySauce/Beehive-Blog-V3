<script setup lang="ts">
import { FileClock, GitBranch, ShieldCheck, Sparkles } from 'lucide-vue-next';

import BaseBadge from '@/shared/components/BaseBadge.vue';
import BaseCard from '@/shared/components/BaseCard.vue';
</script>

<template>
  <main class="grid gap-5 p-4 sm:p-6 lg:p-8">
    <header class="flex flex-wrap items-end justify-between gap-3">
      <div>
        <BaseBadge tone="violet">Studio</BaseBadge>
        <h1 class="m-0 mt-3 text-28px font-900 text-white">仪表盘</h1>
        <p class="m-0 mt-2 text-14px leading-6 text-white/52">高密度工作台骨架，后续接入内容、审阅和索引状态。</p>
      </div>
    </header>

    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <BaseCard class="border-white/8 bg-white/5 p-5 text-white">
        <FileClock class="h-5 w-5 text-brand-honey" aria-hidden="true" />
        <p class="m-0 mt-4 text-24px font-900">12</p>
        <p class="m-0 text-13px text-white/48">最近编辑</p>
      </BaseCard>
      <BaseCard class="border-white/8 bg-white/5 p-5 text-white">
        <ShieldCheck class="h-5 w-5 text-brand-leaf" aria-hidden="true" />
        <p class="m-0 mt-4 text-24px font-900">3</p>
        <p class="m-0 text-13px text-white/48">待审草稿</p>
      </BaseCard>
      <BaseCard class="border-white/8 bg-white/5 p-5 text-white">
        <GitBranch class="h-5 w-5 text-brand-blue" aria-hidden="true" />
        <p class="m-0 mt-4 text-24px font-900">8</p>
        <p class="m-0 text-13px text-white/48">内容关系</p>
      </BaseCard>
      <BaseCard class="border-white/8 bg-white/5 p-5 text-white">
        <Sparkles class="h-5 w-5 text-brand-violet" aria-hidden="true" />
        <p class="m-0 mt-4 text-24px font-900">Mock</p>
        <p class="m-0 text-13px text-white/48">AI 草稿队列</p>
      </BaseCard>
    </section>

    <section class="grid gap-4 xl:grid-cols-[1.2fr_0.8fr]">
      <BaseCard class="border-white/8 bg-white/5 p-5 text-white">
        <h2 class="m-0 text-18px font-800">工作流</h2>
        <div class="mt-5 grid gap-3">
          <div class="grid grid-cols-[120px_1fr] gap-3 rounded-md border border-white/8 bg-black/12 p-3 text-13px">
            <span class="text-white/45">笔记</span>
            <span>整理为草稿，等待人工审阅</span>
          </div>
          <div class="grid grid-cols-[120px_1fr] gap-3 rounded-md border border-white/8 bg-black/12 p-3 text-13px">
            <span class="text-white/45">文章</span>
            <span>补充关联项目、标签与发布状态</span>
          </div>
          <div class="grid grid-cols-[120px_1fr] gap-3 rounded-md border border-white/8 bg-black/12 p-3 text-13px">
            <span class="text-white/45">公开站</span>
            <span>只展示 public/member 可见内容</span>
          </div>
        </div>
      </BaseCard>
      <BaseCard class="border-white/8 bg-white/5 p-5 text-white">
        <h2 class="m-0 text-18px font-800">联调状态</h2>
        <p class="m-0 mt-4 text-14px leading-6 text-white/56">Auth 可切换 mock/live，content 保持 mock，避免依赖未完成服务。</p>
      </BaseCard>
    </section>
  </main>
</template>
