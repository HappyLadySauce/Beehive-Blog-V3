const refreshTokenKey = 'beehive.auth.refreshToken';

function canUseStorage(): boolean {
  return typeof window !== 'undefined' && typeof window.localStorage !== 'undefined';
}

export const tokenStorage = {
  readRefreshToken(): string | null {
    if (!canUseStorage()) {
      return null;
    }
    return window.localStorage.getItem(refreshTokenKey);
  },
  writeRefreshToken(refreshToken: string): void {
    if (!canUseStorage()) {
      return;
    }
    window.localStorage.setItem(refreshTokenKey, refreshToken);
  },
  clearRefreshToken(): void {
    if (!canUseStorage()) {
      return;
    }
    window.localStorage.removeItem(refreshTokenKey);
  },
};
