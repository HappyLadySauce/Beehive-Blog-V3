<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'

import { useAuthStore } from '@/features/auth/stores/authStore'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import ThemeToggle from '@/shared/components/ThemeToggle.vue'
import UserAccountMenu from '@/shared/components/UserAccountMenu.vue'

const route = useRoute()
const authStore = useAuthStore()
const isStudioForbidden = computed(() => {
  route.fullPath
  return window.history.state?.studio === 'forbidden' || window.history.state?.forbidden === true
})

onMounted(() => {
  if (!authStore.isAuthenticated && authStore.refreshToken) {
    void authStore.restoreSession()
  }
})

async function handleLogout(): Promise<void> {
  await authStore.logout()
}
</script>

<template>
  <div class="public-shell">
    <header class="public-shell__header">
      <RouterLink class="public-shell__brand" to="/">
        <span class="public-shell__brand-mark">B</span>
        <span>Beehive Blog</span>
      </RouterLink>
      <nav class="public-shell__nav" aria-label="Primary navigation">
        <RouterLink to="/">Home</RouterLink>
        <ThemeToggle />
        <UserAccountMenu :user="authStore.currentUser" surface="public" @logout="handleLogout" />
      </nav>
    </header>
    <StatusAlert v-if="isStudioForbidden" tone="warning" title="Studio access denied">
      Your account is authenticated but does not have the admin role required for Studio.
    </StatusAlert>
    <main class="public-shell__main">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.public-shell {
  width: min(1160px, calc(100% - 32px));
  display: grid;
  gap: 28px;
  margin: 0 auto;
  padding: 20px 0 48px;
}

.public-shell__header {
  min-height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 10px 12px;
  background: var(--bb-color-surface-glass);
  box-shadow: var(--bb-shadow-soft);
  backdrop-filter: blur(16px);
}

.public-shell__brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--bb-color-text);
  font-weight: 800;
  text-decoration: none;
}

.public-shell__brand-mark {
  width: 34px;
  height: 34px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  color: #fff;
  background: linear-gradient(135deg, var(--bb-color-primary), var(--bb-color-accent));
  box-shadow: var(--bb-shadow-soft);
}

.public-shell__brand:focus-visible,
.public-shell__nav a:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.public-shell__nav {
  display: flex;
  align-items: center;
  gap: 12px;
}

.public-shell__nav a {
  border-radius: 6px;
  padding: 8px;
  color: var(--bb-color-muted);
  text-decoration: none;
  transition: color 160ms ease, background-color 160ms ease;
}

.public-shell__nav a:hover {
  color: var(--bb-color-text);
  background: var(--bb-color-subtle);
}

.public-shell__nav a.router-link-active {
  color: var(--bb-color-primary);
}

.public-shell__main {
  min-width: 0;
}

@media (max-width: 640px) {
  .public-shell {
    width: min(100% - 20px, 1160px);
    padding-top: 10px;
  }

  .public-shell__header {
    align-items: flex-start;
    flex-direction: column;
  }

  .public-shell__nav {
    width: 100%;
    justify-content: space-between;
  }
}
</style>
