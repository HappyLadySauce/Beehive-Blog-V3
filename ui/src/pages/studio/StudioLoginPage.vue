<script setup lang="ts">
import { reactive, shallowRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import SsoProviderButtons from '@/features/auth/components/SsoProviderButtons.vue'
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
const { t } = useI18n()
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
    errorMessage.value = t('studioLogin.validation.credentialsRequired')
    return
  }

  try {
    await authStore.login({
      login_identifier: form.loginIdentifier.trim(),
      password: form.password,
      client_type: 'web',
      user_agent: navigator.userAgent,
    })

    if (!authStore.isAdmin) {
      await authStore.logout()
      errorMessage.value = t('studioLogin.validation.adminOnly')
      return
    }

    await router.replace(resolveRedirectPath())
  }
  catch {
    errorMessage.value = authStore.errorMessage || t('studioLogin.fallback.signInFailed')
  }
}
</script>

<template>
  <form class="studio-login" novalidate @submit.prevent="handleSubmit">
    <PageHeader
      :eyebrow="t('studioLogin.eyebrow')"
      :title="t('studioLogin.title')"
      :description="t('studioLogin.description')"
    />

    <StatusAlert v-if="errorMessage" tone="danger" :title="t('studioLogin.blockedTitle')">
      {{ errorMessage }}
    </StatusAlert>

    <FormField :label="t('studioLogin.identifierLabel')" for-id="studio-login-identifier">
      <BaseInput
        id="studio-login-identifier"
        v-model="form.loginIdentifier"
        autocomplete="username"
        :placeholder="t('studioLogin.identifierPlaceholder')"
        required
        :invalid="Boolean(errorMessage)"
      />
    </FormField>

    <FormField :label="t('common.currentPassword')" for-id="studio-login-password">
      <PasswordInput
        id="studio-login-password"
        v-model="form.password"
        :invalid="Boolean(errorMessage)"
      />
    </FormField>

    <BaseButton type="submit" :busy="authStore.isLoading">{{ t('studioLogin.submit') }}</BaseButton>
    <SsoProviderButtons surface="studio" :return-to="resolveRedirectPath()" />
  </form>
</template>

<style scoped>
.studio-login {
  display: grid;
  gap: 16px;
}
</style>
