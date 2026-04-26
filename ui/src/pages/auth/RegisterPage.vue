<script setup lang="ts">
import { ArrowRight, Home } from 'lucide-vue-next';
import { computed, reactive } from 'vue';
import { useRouter } from 'vue-router';

import { useAuthStore } from '@/features/auth/stores/authStore';
import { appConfig } from '@/shared/config/env';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import BaseButton from '@/shared/components/BaseButton.vue';
import BaseInput from '@/shared/components/BaseInput.vue';
import StatusAlert from '@/shared/components/StatusAlert.vue';

const router = useRouter();
const authStore = useAuthStore();
const isLiveMode = computed(() => appConfig.apiMode === 'live');
const gatewayLabel = computed(() => {
  if (appConfig.apiMode !== 'live') {
    return 'mock adapter';
  }
  return appConfig.gatewayBaseUrl || 'Vite /api proxy';
});
const liveSeed = Date.now().toString().slice(-6);
const form = reactive({
  username: appConfig.apiMode === 'live' ? `member_${liveSeed}` : 'demo_member',
  email: appConfig.apiMode === 'live' ? `member-${liveSeed}@beehive.local` : 'member@beehive.local',
  password: 'Demo@123456',
  nickname: 'Member',
});

async function submitRegister() {
  try {
    await authStore.register(form);
    await router.push('/');
  } catch {
    // Store owns the user-facing error state.
    // 用户可见错误状态由 store 统一维护。
  }
}
</script>

<template>
  <main class="grid min-h-screen place-items-center bg-brand-paper px-4 py-8">
    <form class="grid w-full max-w-520px gap-5 rounded-lg border border-brand-line bg-brand-surface p-6 shadow-panel sm:p-8" @submit.prevent="submitRegister">
      <RouterLink to="/" class="bb-focus inline-flex w-max items-center gap-2 rounded-md text-13px font-700 text-brand-muted">
        <Home class="h-4 w-4" aria-hidden="true" />
        返回首页
      </RouterLink>
      <div>
        <h1 class="m-0 text-28px font-900">注册普通账号</h1>
        <p class="m-0 mt-2 text-14px leading-6 text-brand-muted">
          {{ isLiveMode ? '当前注册会写入本地 identity 数据库，角色固定为普通用户。' : 'mock 注册会创建普通用户，用于公开站身份流程。' }}
        </p>
        <div class="mt-3 flex flex-wrap gap-2">
          <BaseBadge :tone="isLiveMode ? 'leaf' : 'honey'">{{ isLiveMode ? 'live gateway' : 'mock mode' }}</BaseBadge>
          <BaseBadge tone="neutral">{{ gatewayLabel }}</BaseBadge>
        </div>
      </div>
      <StatusAlert
        v-if="isLiveMode"
        tone="info"
        title="普通用户提示"
        description="注册成功后会回到公开站。普通用户不能进入 Studio，Studio 请使用默认 admin 登录。"
      />
      <BaseInput v-model="form.username" label="用户名" name="username" autocomplete="username" />
      <BaseInput v-model="form.email" label="邮箱" name="email" type="email" autocomplete="email" />
      <BaseInput v-model="form.nickname" label="昵称" name="nickname" autocomplete="nickname" />
      <BaseInput v-model="form.password" label="密码" name="password" type="password" autocomplete="new-password" />
      <StatusAlert v-if="authStore.errorMessage" tone="danger" title="注册失败" :description="authStore.errorMessage" />
      <BaseButton type="submit" variant="primary" :busy="authStore.isLoading">
        注册普通账号
        <ArrowRight class="h-4 w-4" aria-hidden="true" />
      </BaseButton>
      <p class="m-0 text-13px text-brand-muted">
        已有账号？
        <RouterLink to="/login" class="font-700 text-brand-blue">去登录</RouterLink>
      </p>
    </form>
  </main>
</template>
