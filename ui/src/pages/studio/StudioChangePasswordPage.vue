<script setup lang="ts">
import { reactive, shallowRef } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import BaseButton from '@/shared/components/BaseButton.vue'
import FormField from '@/shared/components/FormField.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import PasswordInput from '@/shared/components/PasswordInput.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'

const authStore = useAuthStore()
const form = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})
const isSaving = shallowRef(false)
const successMessage = shallowRef('')
const errorMessage = shallowRef('')

function resetForm(): void {
  form.currentPassword = ''
  form.newPassword = ''
  form.confirmPassword = ''
}

async function handleSubmit() {
  successMessage.value = ''
  errorMessage.value = ''

  if (form.currentPassword.length === 0 || form.newPassword.length < 12) {
    errorMessage.value = 'Enter the current password and a new password with at least 12 characters.'
    return
  }
  if (form.newPassword !== form.confirmPassword) {
    errorMessage.value = 'New password confirmation does not match.'
    return
  }

  isSaving.value = true
  try {
    await studioApi.changePassword(
      {
        old_password: form.currentPassword,
        new_password: form.newPassword,
      },
      { accessToken: authStore.accessToken },
    )
    resetForm()
    successMessage.value = 'Password changed.'
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to change password.'
  }
  finally {
    isSaving.value = false
  }
}
</script>

<template>
  <section class="password-page">
    <PageHeader
      eyebrow="Account"
      title="Change password"
      description="Rotate administrator credentials without leaving Studio."
    />

    <StatusAlert v-if="successMessage" tone="success" title="Password changed">
      {{ successMessage }}
    </StatusAlert>
    <StatusAlert v-if="errorMessage" tone="danger" title="Password change failed">
      {{ errorMessage }}
    </StatusAlert>

    <form class="password-page__form" novalidate @submit.prevent="handleSubmit">
      <FormField label="Current password" for-id="current-password">
        <PasswordInput id="current-password" v-model="form.currentPassword" autocomplete="current-password" />
      </FormField>

      <FormField label="New password" for-id="new-password">
        <PasswordInput id="new-password" v-model="form.newPassword" autocomplete="new-password" />
      </FormField>

      <FormField label="Confirm new password" for-id="confirm-password">
        <PasswordInput id="confirm-password" v-model="form.confirmPassword" autocomplete="new-password" />
      </FormField>

      <BaseButton type="submit" :busy="isSaving">Change password</BaseButton>
    </form>
  </section>
</template>

<style scoped>
.password-page {
  display: grid;
  gap: 24px;
}

.password-page__form {
  width: min(560px, 100%);
  display: grid;
  gap: 16px;
}
</style>
