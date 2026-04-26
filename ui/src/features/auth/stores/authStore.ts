import { defineStore } from 'pinia';

import { authApi } from '@/shared/api/authApi';
import type { AuthLoginRequest, AuthRegisterRequest, AuthUserProfile } from '@/shared/api/types';
import { tokenStorage } from '@/shared/storage/tokenStorage';

interface AuthState {
  accessToken: string;
  currentUser: AuthUserProfile | null;
  isLoading: boolean;
  errorMessage: string;
}

function normalizeError(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }
  return 'Unexpected authentication error';
}

export function normalizeAuthRole(role: string | undefined): string {
  return (role ?? '').toLowerCase().replace(/^role_/, '');
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    accessToken: '',
    currentUser: null,
    isLoading: false,
    errorMessage: '',
  }),
  getters: {
    isAuthenticated: (state) => state.accessToken.length > 0 && state.currentUser !== null,
    isAdmin: (state) => normalizeAuthRole(state.currentUser?.role) === 'admin',
    refreshToken: () => tokenStorage.readRefreshToken(),
  },
  actions: {
    applySession(accessToken: string, refreshToken: string, user: AuthUserProfile) {
      this.accessToken = accessToken;
      this.currentUser = user;
      tokenStorage.writeRefreshToken(refreshToken);
    },
    clearSession() {
      this.accessToken = '';
      this.currentUser = null;
      tokenStorage.clearRefreshToken();
    },
    async register(payload: AuthRegisterRequest) {
      this.isLoading = true;
      this.errorMessage = '';
      try {
        const response = await authApi.register(payload);
        this.applySession(response.access_token, response.refresh_token, response.user);
      } catch (error) {
        this.errorMessage = normalizeError(error);
        throw error;
      } finally {
        this.isLoading = false;
      }
    },
    async login(payload: AuthLoginRequest) {
      this.isLoading = true;
      this.errorMessage = '';
      try {
        const response = await authApi.login(payload);
        this.applySession(response.access_token, response.refresh_token, response.user);
      } catch (error) {
        this.errorMessage = normalizeError(error);
        throw error;
      } finally {
        this.isLoading = false;
      }
    },
    async loadCurrentUser() {
      if (!this.accessToken) {
        return;
      }
      const response = await authApi.me({ accessToken: this.accessToken });
      this.currentUser = response.user;
    },
    async restoreSession(): Promise<boolean> {
      if (this.isAuthenticated) {
        return true;
      }
      return this.refreshSession();
    },
    async refreshSession(): Promise<boolean> {
      const refreshToken = tokenStorage.readRefreshToken();
      if (!refreshToken) {
        return false;
      }
      try {
        const response = await authApi.refresh({ refresh_token: refreshToken });
        const profile = await authApi.me({ accessToken: response.access_token });
        this.applySession(response.access_token, response.refresh_token, profile.user);
        return true;
      } catch {
        this.clearSession();
        return false;
      }
    },
    async logout() {
      try {
        const refreshToken = tokenStorage.readRefreshToken();
        await authApi.logout(refreshToken ? { refresh_token: refreshToken } : {}, { accessToken: this.accessToken });
      } finally {
        this.clearSession();
      }
    },
  },
});
