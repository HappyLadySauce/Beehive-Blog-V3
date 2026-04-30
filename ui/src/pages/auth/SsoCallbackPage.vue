<script setup lang="ts">
import { computed, onMounted, shallowRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { authApi } from '@/features/auth/api/authApi'
import {
  clearStoredSsoFlow,
  getPendingSsoEmail,
  getStoredSsoFlow,
} from '@/features/auth/composables/useSsoFlow'
import { useAuthStore } from '@/features/auth/stores/authStore'
import type { AuthProvider, AuthSsoCallbackResponse } from '@/features/auth/types'
import PageHeader from '@/shared/components/PageHeader.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const status = shallowRef('Completing provider authorization...')
const errorMessage = shallowRef('')

const provider = computed(() => {
  const value = route.params.provider
  return (Array.isArray(value) ? value[0] : value) as AuthProvider
})

function getQueryValue(key: string): string {
  const value = route.query[key]
  return Array.isArray(value) ? value[0] ?? '' : value ?? ''
}

function notifyLogin(payload: AuthSsoCallbackResponse, returnTo: string): void {
  window.opener?.postMessage(
    {
      type: 'beehive:sso-login',
      payload,
      returnTo,
    },
    window.location.origin,
  )
}

function notifyEmail(code: string, state: string, redirectURI: string): void {
  window.opener?.postMessage(
    {
      type: 'beehive:sso-email',
      provider: provider.value,
      code,
      state,
      redirectURI,
    },
    window.location.origin,
  )
}

async function finishCallback(): Promise<void> {
  const stored = getStoredSsoFlow()
  const code = getQueryValue('code')
  const state = getQueryValue('state')

  if (!stored || stored.provider !== provider.value || stored.state !== state || !code) {
    throw new Error('SSO state is missing or does not match this callback.')
  }

  if (stored.surface === 'email') {
    if (!getPendingSsoEmail()) {
      throw new Error('Pending email update was not found.')
    }
    notifyEmail(code, state, stored.redirect_uri)
    status.value = 'Authorization confirmed. You can return to the profile page.'
    clearStoredSsoFlow()
    if (window.opener) {
      window.close()
    }
    return
  }

  const response = await authApi.finishSso({
    provider: provider.value,
    code,
    state,
    redirect_uri: stored.redirect_uri,
    client_type: 'web',
    user_agent: navigator.userAgent,
  })
  authStore.applySession(response.access_token, response.refresh_token, response.session_id, response.user, response.expires_in)
  notifyLogin(response, stored.return_to)
  clearStoredSsoFlow()
  if (window.opener) {
    window.close()
    return
  }
  await router.replace(stored.return_to || (authStore.isAdmin ? '/studio' : '/'))
}

onMounted(async () => {
  try {
    await finishCallback()
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to complete SSO callback.'
  }
})
</script>

<template>
  <section class="sso-callback">
    <PageHeader eyebrow="Identity" title="SSO callback" :description="status" />
    <StatusAlert v-if="errorMessage" tone="danger" title="SSO callback failed">
      {{ errorMessage }}
    </StatusAlert>
  </section>
</template>

<style scoped>
.sso-callback {
  display: grid;
  gap: 16px;
}
</style>
