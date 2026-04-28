<script setup lang="ts">
import { ExternalLink } from 'lucide-vue-next'
import { computed, nextTick, onBeforeUnmount, shallowRef, watch } from 'vue'

import type { AuthProvider } from '@/features/auth/types'
import BaseButton from '@/shared/components/BaseButton.vue'
import ModalDialog from '@/shared/components/ModalDialog.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'

const props = withDefaults(
  defineProps<{
    open: boolean
    provider: AuthProvider | null
    authUrl: string
  }>(),
  {
    provider: null,
    authUrl: '',
  },
)

const emit = defineEmits<{
  close: []
  openWindow: []
}>()

declare global {
  interface Window {
    WxLogin?: new (options: Record<string, string | boolean>) => unknown
  }
}

const wxContainerId = 'wechat-sso-qrcode'
const wxError = shallowRef('')
let wxScriptPromise: Promise<void> | null = null

const providerLabel = computed(() => {
  if (props.provider === 'wechat') {
    return 'WeChat'
  }
  if (props.provider === 'qq') {
    return 'QQ'
  }
  return 'SSO'
})

const dialogTitle = computed(() => `${providerLabel.value} sign in`)
const canRenderWechat = computed(() => props.provider === 'wechat' && props.authUrl.length > 0)
const canRenderIframe = computed(() => props.provider === 'qq' && props.authUrl.length > 0)

function loadWechatScript(): Promise<void> {
  if (window.WxLogin) {
    return Promise.resolve()
  }
  if (wxScriptPromise) {
    return wxScriptPromise
  }
  wxScriptPromise = new Promise((resolve, reject) => {
    const script = document.createElement('script')
    script.src = 'https://res.wx.qq.com/connect/zh_CN/htmledition/js/wxLogin.js'
    script.async = true
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('Unable to load WeChat QR script.'))
    document.head.appendChild(script)
  })
  return wxScriptPromise
}

function parseWechatURL(): Record<string, string> {
  const url = new URL(props.authUrl)
  return {
    appid: url.searchParams.get('appid') ?? '',
    scope: url.searchParams.get('scope') ?? 'snsapi_login',
    redirect_uri: url.searchParams.get('redirect_uri') ?? '',
    state: url.searchParams.get('state') ?? '',
  }
}

async function renderWechatQRCode(): Promise<void> {
  if (!canRenderWechat.value) {
    return
  }
  wxError.value = ''
  await nextTick()
  try {
    await loadWechatScript()
    const params = parseWechatURL()
    if (!params.appid || !params.redirect_uri || !params.state || !window.WxLogin) {
      throw new Error('WeChat authorization URL is incomplete.')
    }
    new window.WxLogin({
      self_redirect: false,
      id: wxContainerId,
      appid: params.appid,
      scope: params.scope,
      redirect_uri: params.redirect_uri,
      state: params.state,
      style: 'black',
      href: '',
    })
  }
  catch (error) {
    wxError.value = error instanceof Error ? error.message : 'WeChat QR code is unavailable.'
  }
}

watch(
  () => [props.open, props.provider, props.authUrl],
  () => {
    if (props.open && props.provider === 'wechat') {
      void renderWechatQRCode()
    }
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  wxError.value = ''
})
</script>

<template>
  <ModalDialog :open="open" :title="dialogTitle" description="Complete authorization in the secure provider window." size="sm" @close="emit('close')">
    <div class="sso-qr-dialog">
      <StatusAlert v-if="wxError" tone="warning" title="QR code unavailable">
        {{ wxError }}
      </StatusAlert>

      <div v-if="canRenderWechat" :id="wxContainerId" class="sso-qr-dialog__wechat" />

      <iframe
        v-else-if="canRenderIframe"
        class="sso-qr-dialog__frame"
        :src="authUrl"
        title="QQ authorization"
        sandbox="allow-forms allow-scripts allow-same-origin allow-popups"
      />

      <div v-else class="sso-qr-dialog__empty">
        Authorization is ready in a separate provider window.
      </div>
    </div>

    <template #footer>
      <BaseButton type="button" @click="emit('openWindow')">
        <ExternalLink :size="16" aria-hidden="true" />
        Open authorization window
      </BaseButton>
      <BaseButton variant="ghost" type="button" @click="emit('close')">Close</BaseButton>
    </template>
  </ModalDialog>
</template>

<style scoped>
.sso-qr-dialog {
  display: grid;
  gap: 14px;
}

.sso-qr-dialog__wechat {
  min-height: 390px;
  display: grid;
  place-items: center;
}

.sso-qr-dialog__frame {
  width: 100%;
  min-height: 420px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  background: var(--bb-color-surface-elevated);
}

.sso-qr-dialog__empty {
  min-height: 220px;
  display: grid;
  place-items: center;
  border: 1px dashed var(--bb-color-line);
  border-radius: 8px;
  color: var(--bb-color-muted);
  text-align: center;
}
</style>
