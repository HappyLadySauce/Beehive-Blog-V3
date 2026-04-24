<script setup lang="ts">
import { computed } from 'vue';

import { appConfig } from '@/shared/config/env';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import BaseCard from '@/shared/components/BaseCard.vue';

const apiModeLabel = computed(() => (appConfig.apiMode === 'live' ? 'Live Gateway' : 'Mock Adapter'));
</script>

<template>
  <main class="grid gap-5 p-4 sm:p-6 lg:p-8">
    <header>
      <BaseBadge tone="honey">Settings</BaseBadge>
      <h1 class="m-0 mt-3 text-28px font-900 text-white">设置</h1>
      <p class="m-0 mt-2 text-14px leading-6 text-white/52">首版展示前端运行配置，后续扩展用户偏好和发布配置。</p>
    </header>

    <BaseCard class="grid gap-4 border-white/8 bg-white/5 p-5 text-white">
      <div class="grid gap-1">
        <span class="text-12px font-700 text-white/45">API 模式</span>
        <span class="text-15px font-800">{{ apiModeLabel }}</span>
      </div>
      <div class="grid gap-1">
        <span class="text-12px font-700 text-white/45">Gateway 地址</span>
        <span class="break-all text-15px font-800">{{ appConfig.gatewayBaseUrl }}</span>
      </div>
    </BaseCard>
  </main>
</template>
