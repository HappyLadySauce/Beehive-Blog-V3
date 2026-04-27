const refreshTokenKey = 'beehive.v3.auth.refresh_token'

export interface RefreshTokenStorage {
  getRefreshToken(): string | null
  setRefreshToken(refreshToken: string): void
  clearRefreshToken(): void
  readRefreshToken(): string | null
  writeRefreshToken(refreshToken: string): void
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

  function getRefreshToken(): string | null {
    return target?.getItem(refreshTokenKey) ?? fallbackRefreshToken
  }

  function setRefreshToken(refreshToken: string): void {
    fallbackRefreshToken = refreshToken
    target?.setItem(refreshTokenKey, refreshToken)
  }

  function clearRefreshToken(): void {
    fallbackRefreshToken = null
    target?.removeItem(refreshTokenKey)
  }

  return {
    getRefreshToken,
    setRefreshToken,
    clearRefreshToken,
    readRefreshToken: getRefreshToken,
    writeRefreshToken: setRefreshToken,
  }
}

export const tokenStorage = createRefreshTokenStorage()
