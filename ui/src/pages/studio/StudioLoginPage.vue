<script setup lang="ts">
import { reactive, shallowRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAuthStore } from '@/features/auth/stores/authStore'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import FormField from '@/shared/components/FormField.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import PasswordInput from '@/shared/components/PasswordInput.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const form = reactive({
  loginIdentifier: 'admin@beehive.local',
  password: '',
})
const errorMessage = shallowRef('')

function resolveRedirectPath(): string {
  const redirect = Array.isArray(route.query.redirect) ? route.query.redirect[0] : route.query.redirect
  if (typeof redirect !== 'string' || !redirect.startsWith('/studio') || redirect === '/studio/login') {
    return '/studio'
  }
  return redirect
}

async function handleSubmit() {
  errorMessage.value = ''

  if (form.loginIdentifier.trim().length === 0 || form.password.length < 8) {
    errorMessage.value = 'Enter an admin identifier and a password with at least 8 characters.'
    return
  }

  try {
    await authStore.login({
      login_identifier: form.loginIdentifier.trim(),
      password: form.password,
      client_type: 'studio',
      user_agent: navigator.userAgent,
    })

    if (!authStore.isAdmin) {
      await authStore.logout()
      errorMessage.value = 'This account is not allowed to access Studio.'
      return
    }

    await router.replace(resolveRedirectPath())
  }
  catch {
    errorMessage.value = authStore.errorMessage || 'Admin sign in failed.'
  }
}
</script>

<template>
  <form class="studio-login" novalidate @submit.prevent="handleSubmit">
    <PageHeader
      eyebrow="Studio"
      title="Admin sign in"
      description="Use an administrator account to access operational tools."
    />

    <StatusAlert v-if="errorMessage" tone="danger" title="Studio sign in blocked">
      {{ errorMessage }}
    </StatusAlert>

    <FormField label="Admin email or username" for-id="studio-login-identifier">
      <BaseInput
        id="studio-login-identifier"
        v-model="form.loginIdentifier"
        autocomplete="username"
        placeholder="admin@beehive.local"
        required
        :invalid="Boolean(errorMessage)"
      />
    </FormField>

    <FormField label="Password" for-id="studio-login-password">
      <PasswordInput
        id="studio-login-password"
        v-model="form.password"
        :invalid="Boolean(errorMessage)"
      />
    </FormField>

    <BaseButton type="submit" :busy="authStore.isLoading">Enter Studio</BaseButton>
  </form>
</template>

<style scoped>
.studio-login {
  display: grid;
  gap: 16px;
}
</style>
