import { shallowRef } from 'vue'
import { useRoute } from 'vue-router'

import { authApi } from '../api/authApi'
import type { AuthProvider, AuthSsoStartResponse } from '../types'

const ssoStorageKey = 'beehive.sso.flow'
const pendingEmailStorageKey = 'beehive.sso.pending_email'

export type SsoFlowSurface = 'login' | 'register' | 'studio' | 'email'

export interface SsoFlowState {
  provider: AuthProvider
  state: string
  redirect_uri: string
  return_to: string
  surface: SsoFlowSurface
}

interface StartSsoOptions {
  returnTo?: string
  email?: string
  accessToken?: string
}

export function getStoredSsoFlow(): SsoFlowState | null {
  const raw = window.sessionStorage.getItem(ssoStorageKey)
  if (!raw) {
    return null
  }
  try {
    return JSON.parse(raw) as SsoFlowState
  }
  catch {
    return null
  }
}

export function clearStoredSsoFlow(): void {
  window.sessionStorage.removeItem(ssoStorageKey)
  window.sessionStorage.removeItem(pendingEmailStorageKey)
}

export function getPendingSsoEmail(): string {
  return window.sessionStorage.getItem(pendingEmailStorageKey) ?? ''
}

export function useSsoFlow(surface: SsoFlowSurface = 'login') {
  const route = useRoute()
  const isStarting = shallowRef(false)
  const errorMessage = shallowRef('')

  async function start(provider: AuthProvider, options: StartSsoOptions = {}): Promise<AuthSsoStartResponse> {
    isStarting.value = true
    errorMessage.value = ''
    try {
      const redirectURI = `${window.location.origin}/auth/sso/callback/${provider}`
      const state = crypto.randomUUID()
      const payload = {
        provider,
        redirect_uri: redirectURI,
        state,
      }
      const response = surface === 'email'
        ? await authApi.startEmailSso(payload, { accessToken: options.accessToken })
        : await authApi.startSso(payload)
      if (surface === 'email' && options.email) {
        window.sessionStorage.setItem(pendingEmailStorageKey, options.email)
      }
      window.sessionStorage.setItem(ssoStorageKey, JSON.stringify({
        provider,
        state: response.state,
        redirect_uri: redirectURI,
        return_to: options.returnTo ?? (typeof route.query.redirect === 'string' ? route.query.redirect : '/'),
        surface,
      } satisfies SsoFlowState))
      return response
    }
    catch (error) {
      errorMessage.value = error instanceof Error ? error.message : 'Unable to start SSO.'
      throw error
    }
    finally {
      isStarting.value = false
    }
  }

  return {
    isStarting,
    errorMessage,
    start,
  }
}
