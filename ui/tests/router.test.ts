import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

import { router } from '@/app/router'
import { useAuthStore } from '@/features/auth/stores/authStore'

describe('router guards', () => {
  beforeEach(async () => {
    setActivePinia(createPinia())
    window.localStorage.clear()
    await router.push('/')
    await router.isReady()
  })

  it('redirects unauthenticated studio access to login', async () => {
    await router.push('/studio')
    expect(router.currentRoute.value.name).toBe('studio-login')
    expect(router.currentRoute.value.query.redirect).toBe('/studio')
  })

  it('keeps the public login route separate from studio login', async () => {
    await router.push('/login')
    expect(router.currentRoute.value.name).toBe('auth-login')

    await router.push('/studio/login')
    expect(router.currentRoute.value.name).toBe('studio-login')
  })

  it('allows admin access after login', async () => {
    const store = useAuthStore()
    await store.login({ login_identifier: 'admin@beehive.local', password: 'Admin@123456' })

    await router.push('/studio')
    expect(router.currentRoute.value.name).toBe('studio-dashboard')
  })

  it('allows admin users to access studio management routes', async () => {
    const store = useAuthStore()
    await store.login({ login_identifier: 'admin@beehive.local', password: 'Admin@123456' })

    await router.push('/studio/users')
    expect(router.currentRoute.value.name).toBe('studio-users')

    await router.push('/studio/audits')
    expect(router.currentRoute.value.name).toBe('studio-audits')
  })

  it('sends member users back to public home', async () => {
    const store = useAuthStore()
    await store.login({ login_identifier: 'member@beehive.local', password: 'Password123!' })

    await router.push('/studio')
    expect(router.currentRoute.value.name).toBe('public-home')
    expect(window.history.state.studio).toBe('forbidden')
    expect(window.history.state.from).toBe('/studio')
  })
})
