<script setup lang="ts">
import { ChevronDown, KeyRound, LayoutDashboard, LogIn, LogOut, ShieldCheck, User, UserPlus } from 'lucide-vue-next'
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, shallowRef, useTemplateRef, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import type { AuthUserProfile } from '@/features/auth/types'

import UserAvatar from './UserAvatar.vue'

const props = defineProps<{
  user: AuthUserProfile | null
  surface?: 'public' | 'studio'
}>()

const emit = defineEmits<{
  logout: []
  changePassword: []
}>()

const route = useRoute()
const triggerRef = useTemplateRef<HTMLButtonElement>('menuTrigger')
const panelRef = useTemplateRef<HTMLElement>('menuPanel')
const isOpen = shallowRef(false)
const panelStyle = reactive({
  top: '0px',
  left: '0px',
})

const isAuthenticated = computed(() => props.user !== null)
const isAdmin = computed(() => (props.user?.role ?? '').toLowerCase().replace(/^role_/, '') === 'admin')
const showAdminLinks = computed(() => isAdmin.value && props.surface !== 'studio')
const showProfileLink = computed(() => props.surface !== 'studio')
const displayName = computed(() => props.user?.nickname || props.user?.username || 'Account')
const email = computed(() => props.user?.email || 'Sign in to continue')

function updatePanelPosition(): void {
  const trigger = triggerRef.value
  if (!trigger) {
    return
  }

  const rect = trigger.getBoundingClientRect()
  const width = Math.min(220, Math.max(180, window.innerWidth - 24))
  const left = Math.min(Math.max(12, rect.right - width), window.innerWidth - width - 12)
  panelStyle.top = `${rect.bottom + 8}px`
  panelStyle.left = `${left}px`
}

async function openMenu(): Promise<void> {
  isOpen.value = true
  await nextTick()
  updatePanelPosition()
  panelRef.value?.focus()
}

function closeMenu(): void {
  isOpen.value = false
}

function toggleMenu(): void {
  if (isOpen.value) {
    closeMenu()
    return
  }
  void openMenu()
}

function handleLogout(): void {
  closeMenu()
  emit('logout')
}

function handleChangePassword(): void {
  closeMenu()
  emit('changePassword')
}

function handleDocumentPointerDown(event: PointerEvent): void {
  if (!isOpen.value) {
    return
  }
  const target = event.target
  if (!(target instanceof Node)) {
    return
  }
  if (triggerRef.value?.contains(target) || panelRef.value?.contains(target)) {
    return
  }
  closeMenu()
}

function handleDocumentKeydown(event: KeyboardEvent): void {
  if (event.key === 'Escape') {
    closeMenu()
  }
}

function handleWindowChange(): void {
  if (isOpen.value) {
    updatePanelPosition()
  }
}

watch(() => route.fullPath, closeMenu)

onMounted(() => {
  document.addEventListener('pointerdown', handleDocumentPointerDown)
  document.addEventListener('keydown', handleDocumentKeydown)
  window.addEventListener('resize', handleWindowChange)
  window.addEventListener('scroll', handleWindowChange, true)
})

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', handleDocumentPointerDown)
  document.removeEventListener('keydown', handleDocumentKeydown)
  window.removeEventListener('resize', handleWindowChange)
  window.removeEventListener('scroll', handleWindowChange, true)
})
</script>

<template>
  <div class="account-menu">
    <button
      ref="menuTrigger"
      class="account-menu__summary"
      type="button"
      aria-haspopup="menu"
      :aria-expanded="isOpen"
      aria-label="Open account menu"
      @click="toggleMenu"
    >
      <UserAvatar :name="displayName" :src="user?.avatar_url" size="md" />
      <span class="account-menu__identity">
        <strong>{{ displayName }}</strong>
        <span>{{ email }}</span>
      </span>
      <ChevronDown :size="16" aria-hidden="true" />
    </button>

    <Teleport to="body">
      <Transition name="account-menu-fade">
        <div
          v-if="isOpen"
          ref="menuPanel"
          class="account-menu__panel"
          role="menu"
          tabindex="-1"
          :style="panelStyle"
        >
          <template v-if="isAuthenticated">
            <RouterLink v-if="showProfileLink" class="account-menu__item" to="/account/profile" role="menuitem" @click="closeMenu">
              <User :size="16" aria-hidden="true" />
              Profile
            </RouterLink>
            <button class="account-menu__item" type="button" role="menuitem" @click="handleChangePassword">
              <KeyRound :size="16" aria-hidden="true" />
              Change password
            </button>
            <RouterLink v-if="showAdminLinks" class="account-menu__item" to="/studio" role="menuitem" @click="closeMenu">
              <LayoutDashboard :size="16" aria-hidden="true" />
              Studio
            </RouterLink>
            <RouterLink v-if="showAdminLinks" class="account-menu__item" to="/studio/users" role="menuitem" @click="closeMenu">
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
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.account-menu {
  position: relative;
  z-index: 120;
}

.account-menu__summary {
  min-height: 44px;
  display: flex;
  align-items: center;
  gap: 10px;
  border: 0;
  border-radius: 8px;
  padding: 4px 8px 4px 4px;
  color: var(--bb-color-text);
  background: transparent;
  cursor: pointer;
}

.account-menu__summary:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.account-menu__identity {
  min-width: 0;
  display: grid;
  text-align: left;
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
  position: fixed;
  z-index: 1200;
  width: min(220px, calc(100vw - 24px));
  display: grid;
  gap: 4px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 6px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-panel);
}

.account-menu__panel:focus {
  outline: none;
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

.account-menu-fade-enter-active,
.account-menu-fade-leave-active {
  transition: opacity 120ms ease, transform 120ms ease;
}

.account-menu-fade-enter-from,
.account-menu-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

@media (max-width: 640px) {
  .account-menu__identity {
    display: none;
  }
}
</style>
