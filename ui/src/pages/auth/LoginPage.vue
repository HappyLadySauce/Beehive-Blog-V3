<script setup lang="ts">
import { ArrowRight, Home } from 'lucide-vue-next';
import { computed, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { useAuthStore } from '@/features/auth/stores/authStore';
import { appConfig } from '@/shared/config/env';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import BaseButton from '@/shared/components/BaseButton.vue';
import BaseInput from '@/shared/components/BaseInput.vue';
import StatusAlert from '@/shared/components/StatusAlert.vue';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const isLiveMode = computed(() => appConfig.apiMode === 'live');
const gatewayLabel = computed(() => {
  if (appConfig.apiMode !== 'live') {
    return 'mock adapter';
  }
  return appConfig.gatewayBaseUrl || 'Vite /api proxy';
});
const form = reactive({
  login_identifier: 'admin@beehive.local',
  password: 'Admin@123456',
});

function safeRedirectTarget(value: unknown): string {
  if (typeof value === 'string' && value.startsWith('/') && !value.startsWith('//')) {
    return value;
  }
  return '/studio';
}

async function submitLogin() {
  try {
    await authStore.login({
      ...form,
      client_type: 'web',
      device_id: 'browser',
      device_name: 'Beehive UI',
      user_agent: navigator.userAgent,
    });
    await router.push(safeRedirectTarget(route.query.redirect));
  } catch {
    // Store owns the user-facing error state.
    // 用户可见错误状态由 store 统一维护。
  }
}
</script>

<template>
  <main class="grid min-h-screen place-items-center bg-brand-paper px-4 py-8">
    <section class="grid w-full max-w-960px overflow-hidden rounded-lg border border-brand-line bg-brand-surface shadow-panel md:grid-cols-[0.9fr_1.1fr]">
      <aside class="bb-grid-bg hidden border-r border-brand-line p-8 md:grid">
        <div class="self-end">
          <p class="m-0 text-13px font-700 text-brand-leaf">{{ isLiveMode ? 'Auth Live Ready' : 'Auth Mock Ready' }}</p>
          <h1 class="m-0 mt-2 text-32px font-900 leading-10">{{ isLiveMode ? '当前正在连接真实 gateway。' : '先让前端独立跑起来，再切真实 gateway。' }}</h1>
          <p class="m-0 mt-4 text-14px leading-6 text-brand-muted">
            {{ isLiveMode ? 'Studio 仅管理员可进入，本地开发默认 admin 需要先执行 seed。' : '默认 mock 模式无需后端服务，admin 账号可进入 Studio。' }}
          </p>
        </div>
      </aside>
      <form class="grid gap-5 p-6 sm:p-8" @submit.prevent="submitLogin">
        <RouterLink to="/" class="bb-focus inline-flex w-max items-center gap-2 rounded-md text-13px font-700 text-brand-muted">
          <Home class="h-4 w-4" aria-hidden="true" />
          返回首页
        </RouterLink>
        <div>
          <h2 class="m-0 text-28px font-900">登录 Studio</h2>
          <p class="m-0 mt-2 text-14px leading-6 text-brand-muted">Studio 仅 admin 可进入，普通账号会回到公开站。</p>
          <div class="mt-3 flex flex-wrap gap-2">
            <BaseBadge :tone="isLiveMode ? 'leaf' : 'honey'">{{ isLiveMode ? 'live gateway' : 'mock mode' }}</BaseBadge>
            <BaseBadge tone="neutral">{{ gatewayLabel }}</BaseBadge>
          </div>
        </div>
        <StatusAlert
          v-if="isLiveMode"
          tone="info"
          title="当前是 live 联调模式"
          description="删库重刷后请先执行 sql/seed.ps1 写入默认 admin，再使用当前预填账号登录。"
        />
        <BaseInput v-model="form.login_identifier" label="登录标识" name="login_identifier" autocomplete="username" />
        <BaseInput v-model="form.password" label="密码" name="password" type="password" autocomplete="current-password" />
        <StatusAlert v-if="authStore.errorMessage" tone="danger" title="登录失败" :description="authStore.errorMessage" />
        <BaseButton type="submit" variant="primary" :busy="authStore.isLoading">
          登录
          <ArrowRight class="h-4 w-4" aria-hidden="true" />
        </BaseButton>
        <p class="m-0 text-13px text-brand-muted">
          需要公开站身份？
          <RouterLink to="/register" class="font-700 text-brand-blue">注册普通账号</RouterLink>
        </p>
      </form>
    </section>
  </main>
</template>
