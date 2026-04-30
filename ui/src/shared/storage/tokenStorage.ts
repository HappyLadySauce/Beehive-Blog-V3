import type { AuthSessionSnapshot } from '@/features/auth/types'

const refreshTokenKey = 'beehive.v3.auth.refresh_token'
const sessionSnapshotKey = 'beehive.v3.auth.session_snapshot'

export interface RefreshTokenStorage {
  getRefreshToken(): string | null
  setRefreshToken(refreshToken: string): void
  clearRefreshToken(): void
  readRefreshToken(): string | null
  writeRefreshToken(refreshToken: string): void
  getSnapshot(): AuthSessionSnapshot | null
  setSnapshot(snapshot: AuthSessionSnapshot): void
  clearSnapshot(): void
}

function resolveStorage(storage?: Storage): Storage | null {
  if (storage) {
    return storage
  }

  if (typeof window === 'undefined') {
    return null
  }

  return window.localStorage
}

export function createRefreshTokenStorage(storage?: Storage): RefreshTokenStorage {
  const target = resolveStorage(storage)
  let fallbackRefreshToken: string | null = null
  let fallbackSnapshot: AuthSessionSnapshot | null = null

  function readSnapshot(): AuthSessionSnapshot | null {
    const serialized = target?.getItem(sessionSnapshotKey)
    if (!serialized) {
      return fallbackSnapshot
    }

    try {
      const parsed = JSON.parse(serialized) as Partial<AuthSessionSnapshot>
      if (
        typeof parsed.accessToken !== 'string'
        || typeof parsed.refreshToken !== 'string'
        || typeof parsed.accessTokenExpiresAt !== 'number'
        || typeof parsed.sessionId !== 'string'
        || (parsed.currentUser !== null && typeof parsed.currentUser !== 'object')
      ) {
        return null
      }
      return {
        accessToken: parsed.accessToken,
        refreshToken: parsed.refreshToken,
        accessTokenExpiresAt: parsed.accessTokenExpiresAt,
        sessionId: parsed.sessionId,
        currentUser: parsed.currentUser ?? null,
      }
    }
    catch {
      return null
    }
  }

  function getRefreshToken(): string | null {
    return readSnapshot()?.refreshToken ?? target?.getItem(refreshTokenKey) ?? fallbackRefreshToken
  }

  function setRefreshToken(refreshToken: string): void {
    fallbackRefreshToken = refreshToken
    target?.setItem(refreshTokenKey, refreshToken)
    const snapshot = readSnapshot()
    if (snapshot) {
      setSnapshot({
        ...snapshot,
        refreshToken,
      })
    }
  }

  function clearRefreshToken(): void {
    fallbackRefreshToken = null
    target?.removeItem(refreshTokenKey)
    clearSnapshot()
  }

  function getSnapshot(): AuthSessionSnapshot | null {
    return readSnapshot()
  }

  function setSnapshot(snapshot: AuthSessionSnapshot): void {
    fallbackSnapshot = snapshot
    fallbackRefreshToken = snapshot.refreshToken
    target?.setItem(sessionSnapshotKey, JSON.stringify(snapshot))
    target?.setItem(refreshTokenKey, snapshot.refreshToken)
  }

  function clearSnapshot(): void {
    fallbackSnapshot = null
    fallbackRefreshToken = null
    target?.removeItem(sessionSnapshotKey)
    target?.removeItem(refreshTokenKey)
  }

  return {
    getRefreshToken,
    setRefreshToken,
    clearRefreshToken,
    readRefreshToken: getRefreshToken,
    writeRefreshToken: setRefreshToken,
    getSnapshot,
    setSnapshot,
    clearSnapshot,
  }
}

export const tokenStorage = createRefreshTokenStorage()
