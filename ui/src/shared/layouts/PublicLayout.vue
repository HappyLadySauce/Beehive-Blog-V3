<script setup lang="ts">
import { LogIn, Search } from 'lucide-vue-next';
import { RouterLink, RouterView } from 'vue-router';

import { publicNavItems } from '@/features/navigation/navItems';
import BaseButton from '@/shared/components/BaseButton.vue';
</script>

<template>
  <div class="min-h-screen bg-brand-paper text-brand-ink">
    <header class="sticky top-0 z-30 border-b border-brand-line bg-brand-surface/92 backdrop-blur">
      <div class="mx-auto flex h-16 max-w-1180px items-center justify-between gap-4 px-4 sm:px-6 lg:px-8">
        <RouterLink to="/" class="bb-focus inline-flex items-center gap-3 rounded-md">
          <span class="grid h-9 w-9 place-items-center rounded-md bg-brand-ink text-brand-surface font-800">B</span>
          <span class="min-w-0">
            <span class="block text-15px font-800 leading-5">Beehive Blog</span>
            <span class="hidden text-12px text-brand-muted sm:block">Public Web</span>
          </span>
        </RouterLink>

        <nav class="hidden items-center gap-1 md:flex" aria-label="公开站导航">
          <RouterLink
            v-for="item in publicNavItems"
            :key="item.to"
            :to="item.to"
            class="bb-focus inline-flex h-9 items-center gap-2 rounded-md px-3 text-14px font-600 text-brand-muted transition-colors hover:bg-brand-paper hover:text-brand-ink"
            active-class="bg-brand-paper text-brand-ink"
          >
            <component :is="item.icon" class="h-4 w-4" aria-hidden="true" />
            {{ item.label }}
          </RouterLink>
        </nav>

        <div class="flex items-center gap-2">
          <BaseButton variant="ghost" size="sm" aria-label="搜索">
            <Search class="h-4 w-4" aria-hidden="true" />
          </BaseButton>
          <RouterLink to="/login">
            <BaseButton variant="secondary" size="sm">
              <LogIn class="h-4 w-4" aria-hidden="true" />
              <span class="hidden sm:inline">登录</span>
            </BaseButton>
          </RouterLink>
        </div>
      </div>
      <nav class="mx-auto flex max-w-1180px gap-1 overflow-x-auto px-4 pb-3 md:hidden" aria-label="移动端公开站导航">
        <RouterLink
          v-for="item in publicNavItems"
          :key="item.to"
          :to="item.to"
          class="bb-focus inline-flex h-8 shrink-0 items-center gap-2 rounded-md px-3 text-13px font-600 text-brand-muted"
          active-class="bg-brand-paper text-brand-ink"
        >
          <component :is="item.icon" class="h-4 w-4" aria-hidden="true" />
          {{ item.label }}
        </RouterLink>
      </nav>
    </header>

    <RouterView />
  </div>
</template>
