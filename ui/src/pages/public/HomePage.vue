<script setup lang="ts">
import { ArrowRight, FileText, LayoutDashboard, LogIn, UserPlus } from 'lucide-vue-next'
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

import { useAuthStore } from '@/features/auth/stores/authStore'
import BaseBadge from '@/shared/components/BaseBadge.vue'
import EmptyState from '@/shared/components/EmptyState.vue'

const authStore = useAuthStore()
const isAdmin = computed(() => (authStore.currentUser?.role ?? '').toLowerCase().replace(/^role_/, '') === 'admin')
const primaryTarget = computed(() => {
  if (!authStore.isAuthenticated) {
    return { to: '/login', label: 'Login', icon: LogIn }
  }
  if (isAdmin.value) {
    return { to: '/studio', label: 'Open Studio', icon: LayoutDashboard }
  }
  return { to: '/account/profile', label: 'Open profile', icon: FileText }
})
</script>

<template>
  <section class="home-page">
    <section class="home-page__hero" aria-labelledby="home-title">
      <div class="home-page__hero-copy">
        <BaseBadge>Public Web</BaseBadge>
        <h1 id="home-title">Beehive Blog</h1>
        <p>
          A clean publishing and operations surface for authenticated content workflows.
        </p>
        <div class="home-page__actions">
          <RouterLink class="home-page__button home-page__button--primary" :to="primaryTarget.to">
            {{ primaryTarget.label }}
            <component :is="primaryTarget.icon" :size="17" aria-hidden="true" />
          </RouterLink>
          <RouterLink v-if="!authStore.isAuthenticated" class="home-page__button home-page__button--secondary" to="/register">
            Register
            <UserPlus :size="17" aria-hidden="true" />
          </RouterLink>
        </div>
      </div>
      <aside class="home-page__panel" aria-label="Current workspace status">
        <EmptyState
          title="No public content is published yet"
          description="Public article and project views will appear here after content is published through Studio."
        />
      </aside>
    </section>
  </section>
</template>

<style scoped>
.home-page {
  display: grid;
  gap: 24px;
}

.home-page__hero {
  min-height: 360px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(280px, 420px);
  gap: 24px;
  align-items: stretch;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 24px;
  background: linear-gradient(135deg, var(--bb-color-surface-elevated), var(--bb-color-primary-soft));
  box-shadow: var(--bb-shadow-panel);
}

.home-page__hero-copy {
  display: grid;
  align-content: center;
  justify-items: start;
  gap: 18px;
}

.home-page__hero-copy h1,
.home-page__hero-copy p {
  margin: 0;
}

.home-page__hero-copy h1 {
  color: var(--bb-color-text-strong);
  font-size: clamp(2.4rem, 7vw, 5rem);
  line-height: 0.96;
}

.home-page__hero-copy p {
  max-width: 620px;
  color: var(--bb-color-muted);
  font-size: 1.08rem;
}

.home-page__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.home-page__button {
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: 1px solid transparent;
  border-radius: 8px;
  padding: 0 16px;
  font-weight: 700;
  text-decoration: none;
  box-shadow: var(--bb-shadow-soft);
  transition: transform 160ms ease, background-color 160ms ease, border-color 160ms ease, color 160ms ease;
}

.home-page__button:hover {
  transform: translateY(-1px);
}

.home-page__button:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.home-page__button--primary {
  color: #fff;
  background: var(--bb-color-primary);
}

.home-page__button--primary:hover {
  background: var(--bb-color-primary-strong);
}

.home-page__button--secondary {
  color: var(--bb-color-text);
  border-color: var(--bb-color-line);
  background: var(--bb-color-surface-elevated);
}

.home-page__button--secondary:hover {
  border-color: var(--bb-color-primary);
}

.home-page__panel {
  min-height: 310px;
  display: grid;
  align-content: center;
}

@media (max-width: 840px) {
  .home-page__hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 520px) {
  .home-page__hero {
    min-height: 0;
    padding: 18px;
  }

  .home-page__actions,
  .home-page__button {
    width: 100%;
  }
}
</style>
