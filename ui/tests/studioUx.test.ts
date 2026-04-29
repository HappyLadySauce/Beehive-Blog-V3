import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createMemoryHistory, createRouter } from 'vue-router'
import { beforeEach, describe, expect, it } from 'vitest'

import { useAuthStore } from '@/features/auth/stores/authStore'
import HomePage from '@/pages/public/HomePage.vue'
import StudioContentPage from '@/pages/studio/StudioContentPage.vue'
import StudioUsersPage from '@/pages/studio/StudioUsersPage.vue'
import { i18n, setLocale } from '@/shared/i18n'

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
  return { wrapper, router }
}

describe('studio UX flows', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    setLocale('en-US')
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
    expect(rows[0]?.find('[aria-label="Edit admin@beehive.local"]').exists()).toBe(true)
    expect(rows[0]?.find('[aria-label="Change password for admin@beehive.local"]').exists()).toBe(true)
    expect(rows[0]?.find('[aria-label="Delete admin@beehive.local"]').exists()).toBe(true)
    expect(rows[1]?.find('[aria-label="Edit editor@beehive.local"]').exists()).toBe(true)
    expect(rows[1]?.find('[aria-label="Reset password for editor@beehive.local"]').exists()).toBe(true)
    expect(rows[1]?.find('[aria-label="Delete editor@beehive.local"]').exists()).toBe(true)

    await rows[0]!.find('[aria-label="Edit admin@beehive.local"]').trigger('click')
    await flushPromises()

    expect(document.body.querySelector('[role="dialog"]')).not.toBeNull()
    expect(document.body.textContent).toContain('Edit user')
    expect(document.body.textContent).toContain('admin@beehive.local')
    expect(document.body.textContent).not.toContain('User ID')
    expect(document.body.textContent).not.toContain('Created')
    expect(document.body.querySelector('#edit-username')).not.toBeNull()
    expect(document.body.querySelector('#edit-email')).not.toBeNull()
  })

  it('does not expose removed content editor actions from the content list', async () => {
    const { wrapper } = await mountWithApp(StudioContentPage)
    await flushPromises()

    expect(wrapper.find('[aria-label="View v3 frontend integration notes"]').exists()).toBe(true)
    expect(wrapper.find('[aria-label="Edit v3 frontend integration notes"]').exists()).toBe(false)
    expect(wrapper.find('[aria-label="Archive v3 frontend integration notes"]').exists()).toBe(true)
    expect(wrapper.findAll('button').some((button) => button.text().includes('New draft'))).toBe(false)
  })
})
