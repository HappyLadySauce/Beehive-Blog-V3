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
  username: '',
  email: '',
  password: '',
})
const errorMessage = shallowRef('')

async function handleSubmit() {
  errorMessage.value = ''

  if (form.username.trim().length < 2 || !form.email.includes('@') || form.password.length < 8) {
    errorMessage.value = 'Enter a username, valid email, and password with at least 8 characters.'
    return
  }

  try {
    await authStore.register({
      username: form.username.trim(),
      email: form.email.trim(),
      password: form.password,
    })
    await router.push({ name: authStore.isAdmin ? 'studio-dashboard' : 'public-home' })
  }
  catch {
    errorMessage.value = authStore.errorMessage || 'Registration failed.'
  }
}
</script>

<template>
  <form class="auth-form" novalidate @submit.prevent="handleSubmit">
    <StatusAlert v-if="errorMessage" tone="danger" title="Registration blocked">
      {{ errorMessage }}
    </StatusAlert>

    <FormField label="Username" for-id="register-username">
      <BaseInput id="register-username" v-model="form.username" autocomplete="username" required />
    </FormField>

    <FormField label="Email" for-id="register-email">
      <BaseInput id="register-email" v-model="form.email" type="email" autocomplete="email" inputmode="email" required />
    </FormField>

    <FormField label="Password" for-id="register-password">
      <PasswordInput id="register-password" v-model="form.password" />
    </FormField>

    <BaseButton type="submit" :busy="authStore.isLoading">Create account</BaseButton>
  </form>
</template>

<style scoped>
.auth-form {
  display: grid;
  gap: 16px;
}
</style>
