<script setup lang="ts">
import { ArrowRight, Home } from 'lucide-vue-next';
import { reactive } from 'vue';
import { useRouter } from 'vue-router';

import { useAuthStore } from '@/features/auth/stores/authStore';
import BaseButton from '@/shared/components/BaseButton.vue';
import BaseInput from '@/shared/components/BaseInput.vue';

const router = useRouter();
const authStore = useAuthStore();
const form = reactive({
  username: 'demo_creator',
  email: 'creator@beehive.local',
  password: 'Demo@123456',
  nickname: 'Creator',
});

async function submitRegister() {
  await authStore.register(form);
  await router.push('/studio');
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
        <h1 class="m-0 text-28px font-900">创建账号</h1>
        <p class="m-0 mt-2 text-14px leading-6 text-brand-muted">首版默认 mock 注册，可通过环境变量切换 live gateway。</p>
      </div>
      <BaseInput v-model="form.username" label="用户名" name="username" autocomplete="username" />
      <BaseInput v-model="form.email" label="邮箱" name="email" type="email" autocomplete="email" />
      <BaseInput v-model="form.nickname" label="昵称" name="nickname" autocomplete="nickname" />
      <BaseInput v-model="form.password" label="密码" name="password" type="password" autocomplete="new-password" />
      <p v-if="authStore.errorMessage" class="m-0 rounded-md bg-red-500/10 px-3 py-2 text-13px text-red-600">
        {{ authStore.errorMessage }}
      </p>
      <BaseButton type="submit" variant="primary" :busy="authStore.isLoading">
        注册并进入 Studio
        <ArrowRight class="h-4 w-4" aria-hidden="true" />
      </BaseButton>
      <p class="m-0 text-13px text-brand-muted">
        已有账号？
        <RouterLink to="/login" class="font-700 text-brand-blue">去登录</RouterLink>
      </p>
    </form>
  </main>
</template>
