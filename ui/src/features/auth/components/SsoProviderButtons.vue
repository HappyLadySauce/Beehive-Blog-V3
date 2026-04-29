<script setup lang="ts">
import { Github, MessageCircle, QrCode } from 'lucide-vue-next'
import { computed, onBeforeUnmount, onMounted, shallowRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/features/auth/stores/authStore'
import type { AuthProvider, AuthSsoCallbackResponse } from '@/features/auth/types'
import { useToast } from '@/shared/composables'

import SsoQrDialog from './SsoQrDialog.vue'
import { useSsoFlow, type SsoFlowSurface } from '../composables/useSsoFlow'

interface SsoLoginMessage {
  type: 'beehive:sso-login'
  payload: AuthSsoCallbackResponse
  returnTo: string
}

interface SsoEmailMessage {
  type: 'beehive:sso-email'
  provider: AuthProvider
  code: string
  state: string
  redirectURI: string
}

const props = withDefaults(
  defineProps<{
    surface?: SsoFlowSurface
    returnTo?: string
    email?: string
    accessToken?: string
  }>(),
  {
    surface: 'login',
    returnTo: '',
    email: '',
    accessToken: '',
  },
)

const emit = defineEmits<{
  emailAuthorized: [message: SsoEmailMessage]
}>()

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const { pushToast } = useToast()
const { start, errorMessage } = useSsoFlow(props.surface)
const activeProvider = shallowRef<AuthProvider | null>(null)
const activeAuthUrl = shallowRef('')
let authWindow: Window | null = null

const providers = [
  { provider: 'github', label: 'GitHub', icon: Github },
  { provider: 'qq', label: 'QQ', icon: MessageCircle },
  { provider: 'wechat', label: 'WeChat', icon: QrCode },
] as const

const isQrDialogOpen = computed(() => activeProvider.value === 'qq' || activeProvider.value === 'wechat')

function openAuthWindow(authURL = activeAuthUrl.value): void {
  if (!authURL) {
    return
  }
  const width = 520
  const height = 680
  const left = Math.max(0, window.screenX + (window.outerWidth - width) / 2)
  const top = Math.max(0, window.screenY + (window.outerHeight - height) / 2)
  authWindow = window.open(
    authURL,
    'beehive-sso',
    `popup=yes,width=${width},height=${height},left=${left},top=${top}`,
  )
  authWindow?.focus()
}

async function startProvider(provider: AuthProvider): Promise<void> {
  try {
    if (props.surface === 'email' && (!props.email || !props.accessToken)) {
      throw new Error(String(t('sso.emailSessionRequired')))
    }
    const response = await start(provider, {
      returnTo: props.returnTo || undefined,
      email: props.email || undefined,
      accessToken: props.accessToken || undefined,
    })
    activeProvider.value = provider
    activeAuthUrl.value = response.auth_url
    if (provider === 'github') {
      openAuthWindow(response.auth_url)
    }
  }
  catch {
    pushToast({
      tone: 'danger',
      title: String(t('sso.unavailableTitle')),
      message: errorMessage.value || String(t('sso.unavailableMessage')),
    })
  }
}

async function handleLoginMessage(message: SsoLoginMessage): Promise<void> {
  const payload = message.payload
  authStore.applySession(payload.access_token, payload.refresh_token, payload.session_id, payload.user)
  activeProvider.value = null
  activeAuthUrl.value = ''
  pushToast({ tone: 'success', title: String(t('sso.signedInTitle')), message: String(t('sso.welcomeBack', { email: payload.user.email })) })
  await router.push(message.returnTo || (authStore.isAdmin ? '/studio' : '/'))
}

function handleMessage(event: MessageEvent<SsoLoginMessage | SsoEmailMessage>): void {
  if (event.origin !== window.location.origin) {
    return
  }
  if (event.data?.type === 'beehive:sso-login') {
    void handleLoginMessage(event.data)
    return
  }
  if (event.data?.type === 'beehive:sso-email') {
    activeProvider.value = null
    activeAuthUrl.value = ''
    emit('emailAuthorized', event.data)
  }
}

onMounted(() => window.addEventListener('message', handleMessage))
onBeforeUnmount(() => {
  window.removeEventListener('message', handleMessage)
  authWindow = null
})
</script>

<template>
  <div class="sso-provider-buttons">
    <div class="sso-provider-buttons__divider">
      <span>{{ t('sso.divider') }}</span>
    </div>
    <div class="sso-provider-buttons__grid">
      <button
        v-for="item in providers"
        :key="item.provider"
        class="sso-provider-buttons__button"
        type="button"
        :aria-label="t('sso.continueWith', { provider: item.label })"
        @click="startProvider(item.provider)"
      >
        <component :is="item.icon" :size="17" aria-hidden="true" />
        {{ item.label }}
      </button>
    </div>

    <SsoQrDialog
      :open="isQrDialogOpen"
      :provider="activeProvider"
      :auth-url="activeAuthUrl"
      @close="activeProvider = null"
      @open-window="openAuthWindow()"
    />
  </div>
</template>

<style scoped>
.sso-provider-buttons {
  display: grid;
  gap: 12px;
}

.sso-provider-buttons__divider {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--bb-color-muted);
  font-size: 0.86rem;
}

.sso-provider-buttons__divider::before,
.sso-provider-buttons__divider::after {
  content: "";
  height: 1px;
  flex: 1;
  background: var(--bb-color-line);
}

.sso-provider-buttons__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.sso-provider-buttons__button {
  min-height: 42px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface-elevated);
  font-weight: 700;
}

.sso-provider-buttons__button:hover,
.sso-provider-buttons__button:focus-visible {
  outline: none;
  border-color: var(--bb-color-primary);
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

@media (max-width: 520px) {
  .sso-provider-buttons__grid {
    grid-template-columns: 1fr;
  }
}
</style>
