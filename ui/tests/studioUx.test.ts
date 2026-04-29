import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import type { VueWrapper } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { studioApi } from '@/features/studio'
import { useAuthStore } from '@/features/auth/stores/authStore'
import HomePage from '@/pages/public/HomePage.vue'
import StudioContentPage from '@/pages/studio/StudioContentPage.vue'
import StudioUsersPage from '@/pages/studio/StudioUsersPage.vue'
import { i18n, setLocale } from '@/shared/i18n'

function t(key: string, params?: Record<string, unknown>) {
  return String(i18n.global.t(key, params ?? {}))
}

const mountedWrappers: VueWrapper[] = []

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: HomePage },
      { path: '/login', component: { template: '<div />' } },
      { path: '/register', component: { template: '<div />' } },
      { path: '/account/profile', component: { template: '<div />' } },
      { path: '/studio', component: { template: '<div />' } },
      { path: '/studio/profile', component: { template: '<div />' } },
      { path: '/studio/change-password', component: { template: '<div />' } },
    ],
  })
}

async function mountWithApp(component: object, initialPath = '/') {
  const router = createTestRouter()
  router.push(initialPath)
  await router.isReady()
  const pinia = createPinia()
  setActivePinia(pinia)
  const authStore = useAuthStore()
  authStore.applySession('access', 'refresh', 'session', {
    user_id: 'user_mock_admin',
    username: 'admin',
    email: 'admin@beehive.local',
    nickname: 'Admin',
    avatar_url: '',
    role: 'admin',
    status: 'active',
  })
  const wrapper = mount(component, {
    attachTo: document.body,
    global: {
      plugins: [pinia, i18n, router],
    },
  })
  mountedWrappers.push(wrapper)
  return { wrapper, router }
}

async function waitForFilterDebounce() {
  await new Promise((resolve) => window.setTimeout(resolve, 350))
  await flushPromises()
}

describe('studio UX flows', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    setLocale('en-US')
  })

  afterEach(() => {
    while (mountedWrappers.length > 0) {
      mountedWrappers.pop()?.unmount()
    }
    document.body.innerHTML = ''
    vi.restoreAllMocks()
  })

  it('does not render public homepage placeholder drafts', async () => {
    const { wrapper } = await mountWithApp(HomePage)

    expect(wrapper.text()).not.toContain('Featured drafts')
    expect(wrapper.text()).not.toContain('Ready for gateway integration')
    expect(wrapper.text()).not.toContain('Designing service boundaries for v3')
    expect(wrapper.text()).toContain('No public content is published yet')
  })

  it('shows icon actions and opens user edit in a modal dialog', async () => {
    const { wrapper } = await mountWithApp(StudioUsersPage)
    await flushPromises()

    const rows = wrapper.findAll('tbody tr')
    expect(rows[0]?.text()).not.toContain('Current user')
    expect(rows[0]?.find('[aria-label="Open profile"]').exists()).toBe(false)
    expect(rows[0]?.find(`[aria-label="${t('users.actions.editUser', { email: 'admin@beehive.local' })}"]`).exists()).toBe(true)
    expect(rows[0]?.find(`[aria-label="${t('users.actions.changePassword', { email: 'admin@beehive.local' })}"]`).exists()).toBe(true)
    expect(rows[0]?.find(`[aria-label="${t('users.actions.deleteUser', { email: 'admin@beehive.local' })}"]`).exists()).toBe(true)
    expect(rows[1]?.find(`[aria-label="${t('users.actions.editUser', { email: 'editor@beehive.local' })}"]`).exists()).toBe(true)
    expect(rows[1]?.find(`[aria-label="${t('users.actions.resetPassword', { email: 'editor@beehive.local' })}"]`).exists()).toBe(true)
    expect(rows[1]?.find(`[aria-label="${t('users.actions.deleteUser', { email: 'editor@beehive.local' })}"]`).exists()).toBe(true)

    await rows[0]!.find(`[aria-label="${t('users.actions.editUser', { email: 'admin@beehive.local' })}"]`).trigger('click')
    await flushPromises()

    expect(document.body.querySelector('[role="dialog"]')).not.toBeNull()
    expect(document.body.textContent).toContain(t('users.editDialog.title'))
    expect(document.body.textContent).toContain('admin@beehive.local')
    expect(document.body.textContent).not.toContain(t('users.viewDialog.userId'))
    expect(document.body.textContent).not.toContain(t('users.viewDialog.created'))
    expect(document.body.querySelector('#edit-username')).not.toBeNull()
    expect(document.body.querySelector('#edit-email')).not.toBeNull()
  })

  it('requests include_deleted and disables mutating actions for deleted users', async () => {
    const activeUser = {
      user_id: 'user_active',
      username: 'editor',
      email: 'editor@beehive.local',
      nickname: 'Editor',
      avatar_url: '',
      role: 'member',
      status: 'active',
      created_at: 1776781080,
      updated_at: 1777219320,
      last_login_at: 1777219320,
    }
    const deletedUser = {
      user_id: 'user_deleted',
      username: 'retired',
      email: 'retired@beehive.local',
      nickname: 'Retired',
      avatar_url: '',
      role: 'member',
      status: 'deleted',
      created_at: 1776781080,
      updated_at: 1777310000,
      deleted_at: 1777310000,
    }
    const listUsersSpy = vi.spyOn(studioApi, 'listUsers').mockImplementation(async (params) => {
      const items = params?.include_deleted ? [activeUser, deletedUser] : [activeUser]
      return {
        items,
        total: items.length,
        page: 1,
        page_size: 50,
      }
    })

    const { wrapper } = await mountWithApp(StudioUsersPage)
    await flushPromises()

    expect(listUsersSpy).toHaveBeenCalledWith(expect.objectContaining({ include_deleted: false }), expect.any(Object))
    expect(wrapper.text()).not.toContain('retired@beehive.local')

    await wrapper.get('.users-page__deleted-checkbox').setValue(true)
    await waitForFilterDebounce()

    expect(listUsersSpy.mock.calls.some(([params]) => params?.include_deleted === true)).toBe(true)
    expect(wrapper.text()).toContain('retired@beehive.local')
    const deletedAtPrefix = t('users.deletedAt', { value: '__VALUE__' }).replace('__VALUE__', '').trim()
    expect(wrapper.text()).toContain(deletedAtPrefix)

    const deletedRow = wrapper.findAll('tbody tr').find((row) => row.text().includes('retired@beehive.local'))
    expect(deletedRow).toBeDefined()
    expect(deletedRow!.find(`[aria-label="${t('users.actions.editDeletedUser', { email: 'retired@beehive.local' })}"]`).attributes('disabled')).toBeDefined()
    expect(deletedRow!.find(`[aria-label="${t('users.actions.resetDeletedUserPassword', { email: 'retired@beehive.local' })}"]`).attributes('disabled')).toBeDefined()
    expect(deletedRow!.find(`[aria-label="${t('users.actions.deleteDeletedUser', { email: 'retired@beehive.local' })}"]`).attributes('disabled')).toBeDefined()
  })

  it('does not expose removed content editor actions from the content list', async () => {
    const { wrapper } = await mountWithApp(StudioContentPage)
    await flushPromises()

    const title = 'v3 frontend integration notes'

    expect(wrapper.find(`[aria-label="${t('content.actions.view', { title })}"]`).exists()).toBe(true)
    expect(wrapper.find(`[aria-label="${t('content.actions.edit', { title })}"]`).exists()).toBe(false)
    expect(wrapper.find(`[aria-label="${t('content.actions.archive', { title })}"]`).exists()).toBe(true)
    expect(wrapper.findAll('button').some((button) => button.text().includes(t('content.newDraft')))).toBe(false)
  })

  it('renders localized content actions in Chinese', async () => {
    setLocale('zh-CN')
    const { wrapper } = await mountWithApp(StudioContentPage)
    await flushPromises()

    const title = 'v3 frontend integration notes'
    expect(wrapper.find(`[aria-label="${t('content.actions.view', { title })}"]`).exists()).toBe(true)
    expect(wrapper.find('[role="tablist"]').attributes('aria-label')).toBe(t('content.aria.workspace'))
  })
})
