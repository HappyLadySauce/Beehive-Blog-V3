import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createMemoryHistory, createRouter } from 'vue-router'
import { beforeEach, describe, expect, it } from 'vitest'

import { useAuthStore } from '@/features/auth/stores/authStore'
import HomePage from '@/pages/public/HomePage.vue'
import StudioAuditsPage from '@/pages/studio/StudioAuditsPage.vue'
import ContentEditorPage from '@/pages/studio/ContentEditorPage.vue'
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
      { path: '/studio/content/new', component: { template: '<div />' } },
      { path: '/studio/content/:content_id/edit', component: { template: '<div />' } },
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

  it('routes new content actions to the full-screen editor', async () => {
    const { wrapper, router } = await mountWithApp(StudioContentPage)
    await flushPromises()

    expect(wrapper.find('[aria-label="View v3 frontend integration notes"]').exists()).toBe(true)
    expect(wrapper.find('[aria-label="Edit v3 frontend integration notes"]').exists()).toBe(true)
    expect(wrapper.find('[aria-label="Archive v3 frontend integration notes"]').exists()).toBe(true)

    await wrapper.findAll('button').find((button) => button.text().includes('New draft'))!.trigger('click')
    await flushPromises()

    expect(router.currentRoute.value.fullPath).toBe('/studio/content/new')
  })

  it('shows the tags tab as a searchable table layout', async () => {
    const { wrapper } = await mountWithApp(StudioContentPage)
    await flushPromises()

    await wrapper.findAll('button').find((button) => button.text().includes('Tags'))!.trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('New tag')
    expect(wrapper.text()).toContain('Name')
    expect(wrapper.text()).toContain('Slug')
    expect(wrapper.text()).toContain('Color')
    expect(wrapper.text()).toContain('Description')
    expect(wrapper.text()).toContain('Gateway')
  })

  it('filters tags from the search field', async () => {
    const { wrapper } = await mountWithApp(StudioContentPage)
    await flushPromises()

    await wrapper.findAll('button').find((button) => button.text().includes('Tags'))!.trigger('click')
    await flushPromises()

    await wrapper.get('#tag-search').setValue('identity')
    await new Promise((resolve) => window.setTimeout(resolve, 350))
    await flushPromises()

    expect(wrapper.text()).toContain('Identity')
    expect(wrapper.text()).not.toContain('Gateway')
  })

  it('opens the same tag drawer for creating and editing tags', async () => {
    const { wrapper } = await mountWithApp(StudioContentPage)
    await flushPromises()

    await wrapper.findAll('button').find((button) => button.text().includes('Tags'))!.trigger('click')
    await flushPromises()

    await wrapper.findAll('button').find((button) => button.text().includes('New tag'))!.trigger('click')
    await flushPromises()

    expect(document.body.textContent).toContain('Create tag')
    expect(document.body.querySelector('#tag-name')).not.toBeNull()

    document.body.querySelector<HTMLButtonElement>('.side-drawer__close')?.click()
    await flushPromises()

    await wrapper.find('[aria-label="Edit tag Gateway"]').trigger('click')
    await flushPromises()

    expect(document.body.textContent).toContain('Edit tag')
    expect((document.body.querySelector('#tag-name') as HTMLInputElement | null)?.value).toBe('Gateway')
  })

  it('shows the content empty state and routes its CTA to the editor', async () => {
    const { wrapper, router } = await mountWithApp(StudioContentPage)
    await flushPromises()

    await wrapper.get('#content-search').setValue('missing-content')
    await new Promise((resolve) => window.setTimeout(resolve, 350))
    await flushPromises()

    expect(wrapper.text()).toContain('No content yet.')
    expect(wrapper.text()).toContain('Create the first content item to start this workspace.')

    await wrapper.findAll('button').find((button) => button.text().includes('New draft'))!.trigger('click')
    await flushPromises()

    expect(router.currentRoute.value.fullPath).toBe('/studio/content/new')
  })

  it('filters audits automatically without an apply button', async () => {
    const { wrapper } = await mountWithApp(StudioAuditsPage)
    await flushPromises()

    expect(wrapper.text()).not.toContain('Apply')

    await wrapper.get('#audit-event-type').setValue('admin_update_user_status')
    await new Promise((resolve) => window.setTimeout(resolve, 350))
    await flushPromises()

    expect(wrapper.text()).toContain('admin_update_user_status')
    expect(wrapper.text()).not.toContain('login')
  })

  it('creates a draft from the full-screen content editor', async () => {
    const { wrapper, router } = await mountWithApp(ContentEditorPage, '/studio/content/new')
    await flushPromises()

    await wrapper.get('#editor-title').setValue('Editor integration draft')
    await wrapper.get('#editor-slug').setValue('editor-integration-draft')
    expect(wrapper.text()).toContain('Visual')
    expect(wrapper.text()).toContain('Markdown')
    await wrapper.findAll('button').find((button) => button.text().includes('Create draft'))!.trigger('click')
    await flushPromises()

    expect(router.currentRoute.value.fullPath).toMatch(/^\/studio\/content\/.+\/edit$/)
  })

  it('does not mark the editor dirty when switching to markdown without edits', async () => {
    const { wrapper } = await mountWithApp(ContentEditorPage, '/studio/content/new')
    await flushPromises()

    expect(wrapper.text()).toContain('Idle')

    await wrapper.findAll('button').find((button) => button.text().includes('Markdown'))!.trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('Idle')
    expect(wrapper.text()).not.toContain('Unsaved')
  })
})
