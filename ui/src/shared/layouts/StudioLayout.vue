<script setup lang="ts">
import { Menu, PanelLeftClose, Search, X } from 'lucide-vue-next';
import { ref } from 'vue';
import { RouterLink, RouterView } from 'vue-router';

import { studioNavItems } from '@/features/navigation/navItems';
import { useAuthStore } from '@/features/auth/stores/authStore';
import BaseButton from '@/shared/components/BaseButton.vue';

const isMobileNavOpen = ref(false);
const authStore = useAuthStore();
</script>

<template>
  <div class="min-h-screen bg-#0f1011 text-#f7f8f8 lg:grid lg:grid-cols-[260px_1fr]">
    <aside class="hidden border-r border-white/8 bg-#0b0c0d lg:block">
      <div class="sticky top-0 grid h-screen grid-rows-[auto_1fr_auto] p-4">
        <RouterLink to="/studio" class="bb-focus flex items-center gap-3 rounded-md px-2 py-2">
          <span class="grid h-9 w-9 place-items-center rounded-md bg-#f7f8f8 text-#0f1011 font-800">B</span>
          <span>
            <span class="block text-15px font-800">Beehive Studio</span>
            <span class="text-12px text-white/48">内容生产工作台</span>
          </span>
        </RouterLink>

        <nav class="mt-8 grid content-start gap-1" aria-label="Studio 导航">
          <RouterLink
            v-for="item in studioNavItems"
            :key="item.to"
            :to="item.to"
            class="bb-focus flex h-10 items-center gap-3 rounded-md px-3 text-14px font-600 text-white/62 transition-colors hover:bg-white/6 hover:text-white"
            active-class="bg-white/8 text-white"
          >
            <component :is="item.icon" class="h-4 w-4" aria-hidden="true" />
            {{ item.label }}
          </RouterLink>
        </nav>

        <div class="rounded-lg border border-white/8 bg-white/4 p-3">
          <p class="m-0 text-12px text-white/45">当前身份</p>
          <p class="m-0 mt-1 truncate text-14px font-700">{{ authStore.currentUser?.nickname ?? 'Mock Creator' }}</p>
        </div>
      </div>
    </aside>

    <div class="min-w-0">
      <header class="sticky top-0 z-30 border-b border-white/8 bg-#101112/88 backdrop-blur">
        <div class="flex h-15 items-center justify-between gap-3 px-4 sm:px-6 lg:px-8">
          <div class="flex items-center gap-2">
            <BaseButton variant="ghost" size="sm" class="lg:hidden" aria-label="打开导航" @click="isMobileNavOpen = true">
              <Menu class="h-4 w-4" aria-hidden="true" />
            </BaseButton>
            <div class="hidden h-9 items-center gap-2 rounded-md border border-white/8 bg-white/4 px-3 text-13px text-white/50 sm:flex">
              <Search class="h-4 w-4" aria-hidden="true" />
              搜索内容、标签或草稿
            </div>
          </div>
          <RouterLink to="/">
            <BaseButton variant="ghost" size="sm">
              <PanelLeftClose class="h-4 w-4" aria-hidden="true" />
              返回公开站
            </BaseButton>
          </RouterLink>
        </div>
      </header>

      <RouterView />
    </div>

    <div v-if="isMobileNavOpen" class="fixed inset-0 z-40 bg-black/60 lg:hidden" @click="isMobileNavOpen = false">
      <aside class="h-full w-min-78 max-w-82 border-r border-white/10 bg-#0b0c0d p-4" @click.stop>
        <div class="flex items-center justify-between">
          <span class="text-15px font-800">Beehive Studio</span>
          <BaseButton variant="ghost" size="sm" aria-label="关闭导航" @click="isMobileNavOpen = false">
            <X class="h-4 w-4" aria-hidden="true" />
          </BaseButton>
        </div>
        <nav class="mt-6 grid gap-1" aria-label="移动端 Studio 导航">
          <RouterLink
            v-for="item in studioNavItems"
            :key="item.to"
            :to="item.to"
            class="bb-focus flex h-10 items-center gap-3 rounded-md px-3 text-14px font-600 text-white/68"
            active-class="bg-white/8 text-white"
            @click="isMobileNavOpen = false"
          >
            <component :is="item.icon" class="h-4 w-4" aria-hidden="true" />
            {{ item.label }}
          </RouterLink>
        </nav>
      </aside>
    </div>
  </div>
</template>
