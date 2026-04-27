<script setup lang="ts">
import { reactive, shallowRef } from 'vue'
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/features/auth/stores/authStore'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import FormField from '@/shared/components/FormField.vue'
import PasswordInput from '@/shared/components/PasswordInput.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'

const router = useRouter()
const authStore = useAuthStore()
const form = reactive({
  loginIdentifier: '',
  password: '',
})
const errorMessage = shallowRef('')

async function handleSubmit() {
  errorMessage.value = ''

  if (form.loginIdentifier.trim().length === 0 || form.password.length < 8) {
    errorMessage.value = 'Enter an identifier and a password with at least 8 characters.'
    return
  }

  try {
    await authStore.login({
      login_identifier: form.loginIdentifier.trim(),
      password: form.password,
      client_type: 'web',
      user_agent: navigator.userAgent,
    })
    await router.push({ name: authStore.isAdmin ? 'studio-dashboard' : 'public-home' })
  }
  catch {
    errorMessage.value = authStore.errorMessage || 'Sign in failed.'
  }
}
</script>

<template>
  <form class="auth-form" novalidate @submit.prevent="handleSubmit">
    <StatusAlert v-if="errorMessage" tone="danger" title="Sign in blocked">
      {{ errorMessage }}
    </StatusAlert>

    <FormField label="Email or username" for-id="login-identifier">
      <BaseInput
        id="login-identifier"
        v-model="form.loginIdentifier"
        autocomplete="username"
        placeholder="admin@example.com"
        required
        :invalid="Boolean(errorMessage)"
      />
    </FormField>

    <FormField label="Password" for-id="login-password">
      <PasswordInput
        id="login-password"
        v-model="form.password"
        :invalid="Boolean(errorMessage)"
      />
    </FormField>

    <BaseButton type="submit" :busy="authStore.isLoading">Sign in</BaseButton>
  </form>
</template>

<style scoped>
.auth-form {
  display: grid;
  gap: 16px;
}
</style>
