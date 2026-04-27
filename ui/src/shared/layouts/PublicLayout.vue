<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'

import { useAuthStore } from '@/features/auth/stores/authStore'
import StatusAlert from '@/shared/components/StatusAlert.vue'
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
      <RouterLink class="public-shell__brand" to="/">Beehive Blog</RouterLink>
      <nav class="public-shell__nav" aria-label="Primary navigation">
        <RouterLink to="/">Home</RouterLink>
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
  width: min(1120px, calc(100% - 32px));
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
  border-bottom: 1px solid var(--bb-color-line);
}

.public-shell__brand {
  color: var(--bb-color-text);
  font-weight: 800;
  text-decoration: none;
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
}

.public-shell__nav a.router-link-active {
  color: var(--bb-color-primary);
}

.public-shell__main {
  min-width: 0;
}
</style>
