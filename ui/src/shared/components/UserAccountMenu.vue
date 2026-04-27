<script setup lang="ts">
import { ChevronDown, KeyRound, LayoutDashboard, LogIn, LogOut, ShieldCheck, User, UserPlus } from 'lucide-vue-next'
import { computed, useTemplateRef } from 'vue'
import { RouterLink } from 'vue-router'

import type { AuthUserProfile } from '@/features/auth/types'

import UserAvatar from './UserAvatar.vue'

const props = defineProps<{
  user: AuthUserProfile | null
  surface?: 'public' | 'studio'
}>()

const emit = defineEmits<{
  logout: []
}>()

const menuRef = useTemplateRef<HTMLDetailsElement>('accountMenu')
const isAuthenticated = computed(() => props.user !== null)
const isAdmin = computed(() => (props.user?.role ?? '').toLowerCase().replace(/^role_/, '') === 'admin')
const displayName = computed(() => props.user?.nickname || props.user?.username || 'Account')
const email = computed(() => props.user?.email || 'Sign in to continue')
const profilePath = computed(() => props.surface === 'studio' ? '/studio/profile' : '/account/profile')
const passwordPath = computed(() => props.surface === 'studio' ? '/studio/change-password' : '/account/change-password')

function closeMenu(): void {
  if (menuRef.value) {
    menuRef.value.open = false
  }
}

function handleLogout(): void {
  closeMenu()
  emit('logout')
}
</script>

<template>
  <details ref="accountMenu" class="account-menu">
    <summary class="account-menu__summary" aria-label="Open account menu">
      <UserAvatar :name="displayName" :src="user?.avatar_url" size="md" />
      <span class="account-menu__identity">
        <strong>{{ displayName }}</strong>
        <span>{{ email }}</span>
      </span>
      <ChevronDown :size="16" aria-hidden="true" />
    </summary>
    <div class="account-menu__panel" role="menu">
      <template v-if="isAuthenticated">
      <RouterLink class="account-menu__item" :to="profilePath" role="menuitem" @click="closeMenu">
        <User :size="16" aria-hidden="true" />
        Profile
      </RouterLink>
      <RouterLink class="account-menu__item" :to="passwordPath" role="menuitem" @click="closeMenu">
        <KeyRound :size="16" aria-hidden="true" />
        Change password
      </RouterLink>
      <RouterLink v-if="isAdmin" class="account-menu__item" to="/studio" role="menuitem" @click="closeMenu">
        <LayoutDashboard :size="16" aria-hidden="true" />
        Studio
      </RouterLink>
      <RouterLink v-if="isAdmin" class="account-menu__item" to="/studio/users" role="menuitem" @click="closeMenu">
        <ShieldCheck :size="16" aria-hidden="true" />
        Users
      </RouterLink>
      <button class="account-menu__item account-menu__item--danger" type="button" role="menuitem" @click="handleLogout">
        <LogOut :size="16" aria-hidden="true" />
        Logout
      </button>
      </template>
      <template v-else>
      <RouterLink class="account-menu__item" to="/login" role="menuitem" @click="closeMenu">
        <LogIn :size="16" aria-hidden="true" />
        Login
      </RouterLink>
      <RouterLink class="account-menu__item" to="/register" role="menuitem" @click="closeMenu">
        <UserPlus :size="16" aria-hidden="true" />
        Register
      </RouterLink>
      </template>
    </div>
  </details>
</template>

<style scoped>
.account-menu {
  position: relative;
}

.account-menu__summary {
  min-height: 44px;
  display: flex;
  align-items: center;
  gap: 10px;
  border-radius: 8px;
  padding: 4px 8px 4px 4px;
  list-style: none;
  color: var(--bb-color-text);
  cursor: pointer;
}

.account-menu__summary::-webkit-details-marker {
  display: none;
}

.account-menu__summary:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.account-menu__identity {
  min-width: 0;
  display: grid;
}

.account-menu__identity strong,
.account-menu__identity span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.account-menu__identity span {
  color: var(--bb-color-muted);
  font-size: 0.82rem;
}

.account-menu__panel {
  position: absolute;
  z-index: 20;
  top: calc(100% + 8px);
  right: 0;
  width: 220px;
  display: grid;
  gap: 4px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 6px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-panel);
}

.account-menu__item {
  min-height: 40px;
  display: flex;
  align-items: center;
  gap: 10px;
  border: 0;
  border-radius: 6px;
  padding: 0 10px;
  color: var(--bb-color-muted);
  background: transparent;
  text-align: left;
  text-decoration: none;
}

.account-menu__item:hover,
.account-menu__item:focus-visible,
.account-menu__item.router-link-exact-active {
  outline: none;
  color: var(--bb-color-text);
  background: var(--bb-color-subtle);
}

.account-menu__item--danger {
  color: var(--bb-color-danger);
}

@media (max-width: 640px) {
  .account-menu__identity {
    display: none;
  }
}
</style>
