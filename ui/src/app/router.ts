import { createRouter, createWebHistory } from 'vue-router'

import { normalizeAuthRole, useAuthStore } from '@/features/auth/stores/authStore'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/shared/layouts/PublicLayout.vue'),
      children: [
        { path: '', name: 'public-home', component: () => import('@/pages/public/HomePage.vue') },
        { path: 'account/profile', name: 'account-profile', component: () => import('@/pages/studio/StudioProfilePage.vue'), meta: { requiresAuth: true } },
        { path: 'account/change-password', redirect: '/account/profile' },
      ],
    },
    {
      path: '/',
      component: () => import('@/shared/layouts/AuthLayout.vue'),
      children: [
        { path: 'login', name: 'auth-login', component: () => import('@/pages/auth/LoginPage.vue') },
        { path: 'register', name: 'auth-register', component: () => import('@/pages/auth/RegisterPage.vue') },
        { path: 'auth/sso/callback/:provider', name: 'auth-sso-callback', component: () => import('@/pages/auth/SsoCallbackPage.vue') },
        { path: 'studio/login', name: 'studio-login', component: () => import('@/pages/studio/StudioLoginPage.vue') },
      ],
    },
    {
      path: '/studio',
      component: () => import('@/shared/layouts/StudioLayout.vue'),
      meta: { requiresAuth: true, requiredRole: 'admin' },
      children: [
        { path: '', name: 'studio-dashboard', component: () => import('@/pages/studio/StudioDashboardPage.vue') },
        { path: 'content', name: 'studio-content', component: () => import('@/pages/studio/StudioContentPage.vue') },
        { path: 'content/new', redirect: '/studio/content' },
        { path: 'content/:content_id/edit', redirect: '/studio/content' },
        { path: 'users', name: 'studio-users', component: () => import('@/pages/studio/StudioUsersPage.vue') },
        { path: 'audits', name: 'studio-audits', component: () => import('@/pages/studio/StudioAuditsPage.vue') },
        { path: 'profile', redirect: '/account/profile' },
        { path: 'change-password', redirect: '/studio' },
        { path: 'settings', name: 'studio-settings', component: () => import('@/pages/studio/StudioSettingsPage.vue') },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach(async (to) => {
  if (!to.matched.some((route) => route.meta.requiresAuth === true)) {
    return true
  }

  const authStore = useAuthStore()
  const isReady = authStore.isAuthenticated || (await authStore.restoreSession())
  if (!isReady) {
    const isStudioRoute = to.fullPath.startsWith('/studio')
    return {
      name: isStudioRoute ? 'studio-login' : 'auth-login',
      query: { redirect: to.fullPath },
    }
  }

  const requiredRole = to.matched.find((route) => typeof route.meta.requiredRole === 'string')?.meta.requiredRole
  if (typeof requiredRole === 'string' && normalizeAuthRole(authStore.currentUser?.role) !== requiredRole) {
    return {
      path: '/',
      query: { studio: 'forbidden' },
      state: {
        studio: 'forbidden',
        forbidden: true,
        from: to.fullPath,
      },
    }
  }

  return true
})

router.afterEach((to, from) => {
  if (to.name !== 'public-home' || to.query.studio !== 'forbidden') {
    return
  }

  window.history.replaceState(
    {
      ...window.history.state,
      studio: 'forbidden',
      forbidden: true,
      from: window.history.state?.from ?? to.redirectedFrom?.fullPath ?? from.fullPath,
    },
    '',
    to.fullPath,
  )
})
