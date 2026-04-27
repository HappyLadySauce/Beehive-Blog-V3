import { computed, shallowRef } from 'vue'
import { defineStore } from 'pinia'

import { tokenStorage } from '@/shared/storage/tokenStorage'

import { authApi } from '../api/authApi'
import type { AuthLoginRequest, AuthRegisterRequest, AuthUserProfile } from '../types'

function normalizeError(error: unknown): string {
  if (error instanceof Error) {
    return error.message
  }
  return 'Unexpected authentication error'
}

export function normalizeAuthRole(role: string | undefined): string {
  return (role ?? '').toLowerCase().replace(/^role_/, '')
}

export const useAuthStore = defineStore('auth', () => {
  const accessToken = shallowRef('')
  const sessionId = shallowRef('')
  const currentUser = shallowRef<AuthUserProfile | null>(null)
  const storedRefreshTokenValue = shallowRef(tokenStorage.getRefreshToken())
  const isLoading = shallowRef(false)
  const isRestoring = shallowRef(false)
  const errorMessage = shallowRef('')
  let refreshSessionPromise: Promise<boolean> | null = null
  let restoreSessionPromise: Promise<boolean> | null = null

  const isAuthenticated = computed(() => accessToken.value.length > 0 && currentUser.value !== null)
  const isAdmin = computed(() => normalizeAuthRole(currentUser.value?.role) === 'admin')
  const refreshToken = computed(() => storedRefreshTokenValue.value)

  function applySession(
    nextAccessToken: string,
    nextRefreshToken: string,
    nextSessionId: string,
    nextUser: AuthUserProfile | null,
  ): void {
    accessToken.value = nextAccessToken
    sessionId.value = nextSessionId
    currentUser.value = nextUser
    storedRefreshTokenValue.value = nextRefreshToken
    tokenStorage.setRefreshToken(nextRefreshToken)
  }

  function clearSession(): void {
    accessToken.value = ''
    sessionId.value = ''
    currentUser.value = null
    storedRefreshTokenValue.value = null
    tokenStorage.clearRefreshToken()
  }

  function setCurrentUser(user: AuthUserProfile): void {
    currentUser.value = user
  }

  async function runLoadingAction(action: () => Promise<void>): Promise<void> {
    isLoading.value = true
    errorMessage.value = ''

    try {
      await action()
    }
    catch (error) {
      errorMessage.value = normalizeError(error)
      throw error
    }
    finally {
      isLoading.value = false
    }
  }

  async function register(payload: AuthRegisterRequest): Promise<void> {
    await runLoadingAction(async () => {
      const response = await authApi.register(payload)
      applySession(response.access_token, response.refresh_token, response.session_id, response.user)
    })
  }

  async function login(payload: AuthLoginRequest): Promise<void> {
    await runLoadingAction(async () => {
      const response = await authApi.login(payload)
      applySession(response.access_token, response.refresh_token, response.session_id, response.user)
    })
  }

  async function loadCurrentUser(): Promise<boolean> {
    if (!accessToken.value) {
      return false
    }

    try {
      const response = await authApi.me({ accessToken: accessToken.value })
      currentUser.value = response.user
      return true
    }
    catch {
      return false
    }
  }

  async function runRefreshSession(): Promise<boolean> {
    const storedRefreshToken = tokenStorage.getRefreshToken()
    storedRefreshTokenValue.value = storedRefreshToken
    if (!storedRefreshToken) {
      return false
    }

    try {
      const response = await authApi.refresh({
        refresh_token: storedRefreshToken,
        user_agent: typeof navigator === 'undefined' ? undefined : navigator.userAgent,
      })
      applySession(
        response.access_token,
        response.refresh_token,
        response.session_id,
        response.user ?? currentUser.value,
      )

      if (!response.user) {
        const profile = await authApi.me({ accessToken: response.access_token })
        currentUser.value = profile.user
      }

      return currentUser.value !== null
    }
    catch {
      clearSession()
      return false
    }
  }

  function refreshSession(): Promise<boolean> {
    if (refreshSessionPromise) {
      return refreshSessionPromise
    }

    refreshSessionPromise = runRefreshSession().finally(() => {
      refreshSessionPromise = null
    })
    return refreshSessionPromise
  }

  async function runRestoreSession(): Promise<boolean> {
    if (isAuthenticated.value) {
      return true
    }

    isRestoring.value = true
    try {
      if (accessToken.value && (await loadCurrentUser())) {
        return true
      }
      return await refreshSession()
    }
    finally {
      isRestoring.value = false
    }
  }

  function restoreSession(): Promise<boolean> {
    if (isAuthenticated.value) {
      return Promise.resolve(true)
    }
    if (restoreSessionPromise) {
      return restoreSessionPromise
    }

    restoreSessionPromise = runRestoreSession().finally(() => {
      restoreSessionPromise = null
    })
    return restoreSessionPromise
  }

  async function logout(): Promise<void> {
    try {
      const storedRefreshToken = tokenStorage.getRefreshToken()
      if (accessToken.value) {
        await authApi.logout(storedRefreshToken ? { refresh_token: storedRefreshToken } : {}, {
          accessToken: accessToken.value,
        })
      }
    }
    finally {
      clearSession()
    }
  }

  return {
    accessToken,
    sessionId,
    currentUser,
    isLoading,
    isRestoring,
    errorMessage,
    isAuthenticated,
    isAdmin,
    refreshToken,
    applySession,
    clearSession,
    setCurrentUser,
    register,
    login,
    loadCurrentUser,
    restoreSession,
    refreshSession,
    logout,
  }
})
