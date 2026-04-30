import { computed, shallowRef } from 'vue'
import { defineStore } from 'pinia'

import { tokenStorage } from '@/shared/storage/tokenStorage'

import { authApi } from '../api/authApi'
import type { AuthLoginRequest, AuthRegisterRequest, AuthSessionSnapshot, AuthUserProfile } from '../types'

const accessTokenRefreshThresholdMs = 120_000
const proactiveRefreshLeadMs = 30_000

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
  const accessTokenExpiresAt = shallowRef(0)
  const sessionId = shallowRef('')
  const currentUser = shallowRef<AuthUserProfile | null>(null)
  const storedRefreshTokenValue = shallowRef(tokenStorage.getRefreshToken())
  const isLoading = shallowRef(false)
  const isRestoring = shallowRef(false)
  const errorMessage = shallowRef('')
  let refreshSessionPromise: Promise<boolean> | null = null
  let restoreSessionPromise: Promise<boolean> | null = null
  let proactiveRefreshTimer: ReturnType<typeof setTimeout> | null = null

  const isAuthenticated = computed(() => accessToken.value.length > 0 && currentUser.value !== null)
  const isAdmin = computed(() => normalizeAuthRole(currentUser.value?.role) === 'admin')
  const refreshToken = computed(() => storedRefreshTokenValue.value)

  function applySession(
    nextAccessToken: string,
    nextRefreshToken: string,
    nextSessionId: string,
    nextUser: AuthUserProfile | null,
    expiresInSeconds = 0,
  ): void {
    accessToken.value = nextAccessToken
    accessTokenExpiresAt.value = expiresInSeconds > 0 ? Date.now() + expiresInSeconds * 1000 : 0
    sessionId.value = nextSessionId
    currentUser.value = nextUser
    storedRefreshTokenValue.value = nextRefreshToken
    tokenStorage.setSnapshot({
      accessToken: nextAccessToken,
      refreshToken: nextRefreshToken,
      accessTokenExpiresAt: accessTokenExpiresAt.value,
      sessionId: nextSessionId,
      currentUser: nextUser,
    })
    scheduleProactiveRefresh()
  }

  function clearSession(): void {
    clearProactiveRefresh()
    accessToken.value = ''
    accessTokenExpiresAt.value = 0
    sessionId.value = ''
    currentUser.value = null
    storedRefreshTokenValue.value = null
    tokenStorage.clearSnapshot()
  }

  function clearProactiveRefresh(): void {
    if (proactiveRefreshTimer !== null) {
      clearTimeout(proactiveRefreshTimer)
      proactiveRefreshTimer = null
    }
  }

  function scheduleProactiveRefresh(): void {
    clearProactiveRefresh()
    if (!storedRefreshTokenValue.value || accessTokenExpiresAt.value <= 0) {
      return
    }
    const delay = accessTokenExpiresAt.value - Date.now() - accessTokenRefreshThresholdMs + proactiveRefreshLeadMs
    if (delay <= 0) {
      refreshSession()
      return
    }
    proactiveRefreshTimer = setTimeout(() => {
      proactiveRefreshTimer = null
      refreshSession()
    }, delay)
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
      applySession(response.access_token, response.refresh_token, response.session_id, response.user, response.expires_in)
    })
  }

  async function login(payload: AuthLoginRequest): Promise<void> {
    await runLoadingAction(async () => {
      const response = await authApi.login(payload)
      applySession(response.access_token, response.refresh_token, response.session_id, response.user, response.expires_in)
    })
  }

  function hydrateSession(snapshot: AuthSessionSnapshot): void {
    accessToken.value = snapshot.accessToken
    accessTokenExpiresAt.value = snapshot.accessTokenExpiresAt
    sessionId.value = snapshot.sessionId
    currentUser.value = snapshot.currentUser
    storedRefreshTokenValue.value = snapshot.refreshToken
  }

  function canReuseAccessToken(snapshot: AuthSessionSnapshot): boolean {
    return snapshot.accessToken.length > 0 && snapshot.accessTokenExpiresAt - Date.now() > accessTokenRefreshThresholdMs
  }

  function persistCurrentSnapshot(): void {
    if (!accessToken.value || !storedRefreshTokenValue.value || !sessionId.value || accessTokenExpiresAt.value <= 0) {
      return
    }

    tokenStorage.setSnapshot({
      accessToken: accessToken.value,
      refreshToken: storedRefreshTokenValue.value,
      accessTokenExpiresAt: accessTokenExpiresAt.value,
      sessionId: sessionId.value,
      currentUser: currentUser.value,
    })
  }

  async function loadCurrentUser(): Promise<boolean> {
    if (!accessToken.value) {
      return false
    }

    try {
      const response = await authApi.me({ accessToken: accessToken.value })
      currentUser.value = response.user
      persistCurrentSnapshot()
      return true
    }
    catch {
      return false
    }
  }

  async function runRefreshWithToken(
    refreshTokenValue: string,
    allowRecoveryRetry: boolean,
  ): Promise<boolean> {
    try {
      const response = await authApi.refresh({
        refresh_token: refreshTokenValue,
        user_agent: typeof navigator === 'undefined' ? undefined : navigator.userAgent,
      })
      applySession(
        response.access_token,
        response.refresh_token,
        response.session_id,
        response.user ?? currentUser.value,
        response.expires_in,
      )

      if (!response.user) {
        const profile = await authApi.me({ accessToken: response.access_token })
        currentUser.value = profile.user
      }

      return currentUser.value !== null
    }
    catch {
      const latestRefreshToken = tokenStorage.getRefreshToken()
      storedRefreshTokenValue.value = latestRefreshToken

      if (allowRecoveryRetry && latestRefreshToken && latestRefreshToken !== refreshTokenValue) {
        return runRefreshWithToken(latestRefreshToken, false)
      }

      clearSession()
      return false
    }
  }

  async function runRefreshSession(): Promise<boolean> {
    const storedRefreshToken = tokenStorage.getRefreshToken()
    storedRefreshTokenValue.value = storedRefreshToken
    if (!storedRefreshToken) {
      return false
    }

    return runRefreshWithToken(storedRefreshToken, true)
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
      const snapshot = tokenStorage.getSnapshot()
      if (snapshot && canReuseAccessToken(snapshot)) {
        hydrateSession(snapshot)
        scheduleProactiveRefresh()
        if (currentUser.value !== null || (await loadCurrentUser())) {
          return true
        }
      }

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
