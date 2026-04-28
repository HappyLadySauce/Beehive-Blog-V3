<script setup lang="ts">
import { FileText, LayoutDashboard, ScrollText, Settings, Users } from 'lucide-vue-next'
import { RouterLink, RouterView, useRouter } from 'vue-router'

import { useAuthStore } from '@/features/auth/stores/authStore'
import ThemeToggle from '@/shared/components/ThemeToggle.vue'
import UserAccountMenu from '@/shared/components/UserAccountMenu.vue'

const authStore = useAuthStore()
const router = useRouter()
const navItems = [
  { label: 'Dashboard', to: '/studio', icon: LayoutDashboard },
  { label: 'Content', to: '/studio/content', icon: FileText },
  { label: 'Users', to: '/studio/users', icon: Users },
  { label: 'Audits', to: '/studio/audits', icon: ScrollText },
  { label: 'Settings', to: '/studio/settings', icon: Settings },
] as const

async function handleLogout() {
  await authStore.logout()
  await router.push({ name: 'studio-login' })
}
</script>

<template>
  <div class="studio-shell">
    <aside class="studio-shell__sidebar">
      <RouterLink class="studio-shell__brand" to="/">Beehive Studio</RouterLink>
      <nav class="studio-shell__nav" aria-label="Studio navigation">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          v-slot="{ href, navigate, isExactActive }"
          :to="item.to"
          custom
        >
          <a
            :href="href"
            class="studio-shell__nav-link"
            :class="{ 'studio-shell__nav-link--active': isExactActive }"
            :aria-current="isExactActive ? 'page' : undefined"
            @click="navigate"
          >
            <component :is="item.icon" :size="17" aria-hidden="true" />
            {{ item.label }}
          </a>
        </RouterLink>
      </nav>
    </aside>
    <div class="studio-shell__workspace">
      <header class="studio-shell__topbar">
        <div class="studio-shell__topbar-title">Admin workspace</div>
        <div class="studio-shell__topbar-actions">
          <ThemeToggle />
          <UserAccountMenu :user="authStore.currentUser" surface="studio" @logout="handleLogout" />
        </div>
      </header>
      <main class="studio-shell__main">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<style scoped>
.studio-shell {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr);
  background: var(--bb-gradient-page);
}

.studio-shell__sidebar {
  display: grid;
  align-content: start;
  gap: 24px;
  border-right: 1px solid var(--bb-color-line);
  padding: 20px;
  background: var(--bb-color-surface-glass);
  backdrop-filter: blur(16px);
}

.studio-shell__brand {
  min-height: 42px;
  display: inline-flex;
  align-items: center;
  color: var(--bb-color-text);
  font-weight: 800;
  text-decoration: none;
}

.studio-shell__nav {
  display: grid;
  gap: 6px;
}

.studio-shell__nav-link {
  min-height: 42px;
  display: flex;
  align-items: center;
  gap: 10px;
  border-radius: 6px;
  padding: 10px;
  color: var(--bb-color-muted);
  text-decoration: none;
  transition: color 160ms ease, background-color 160ms ease, transform 160ms ease;
}

.studio-shell__nav-link:hover {
  transform: translateX(2px);
  color: var(--bb-color-text);
  background: var(--bb-color-subtle);
}

.studio-shell__nav-link:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.studio-shell__nav-link--active {
  color: var(--bb-color-primary);
  background: var(--bb-color-primary-soft);
}

.studio-shell__workspace {
  min-width: 0;
}

.studio-shell__topbar {
  position: sticky;
  z-index: 10;
  top: 0;
  min-height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  border-bottom: 1px solid var(--bb-color-line);
  padding: 0 24px;
  background: var(--bb-color-surface-glass);
  backdrop-filter: blur(16px);
}

.studio-shell__topbar-title {
  color: var(--bb-color-muted);
  font-weight: 800;
}

.studio-shell__topbar-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.studio-shell__main {
  padding: 24px;
}

@media (max-width: 820px) {
  .studio-shell {
    grid-template-columns: 1fr;
  }

  .studio-shell__sidebar {
    position: static;
    border-right: 0;
    border-bottom: 1px solid var(--bb-color-line);
  }

  .studio-shell__nav {
    grid-template-columns: repeat(5, minmax(0, 1fr));
    overflow-x: auto;
  }

  .studio-shell__nav-link {
    justify-content: center;
    white-space: nowrap;
  }

  .studio-shell__topbar {
    padding: 0 16px;
  }
}
</style>
