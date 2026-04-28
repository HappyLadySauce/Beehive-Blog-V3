import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { createMemoryHistory, createRouter } from 'vue-router'
import { afterEach, describe, expect, it } from 'vitest'

import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import DataTable from '@/shared/components/DataTable.vue'
import FormField from '@/shared/components/FormField.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import UserAccountMenu from '@/shared/components/UserAccountMenu.vue'
import UserAvatar from '@/shared/components/UserAvatar.vue'

async function mountWithRouter(component: typeof UserAccountMenu, options: Parameters<typeof mount>[1]) {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/', component: { template: '<div />' } }],
  })
  router.push('/')
  await router.isReady()
  return mount(component, {
    ...options,
    attachTo: document.body,
    global: {
      ...(options?.global ?? {}),
      plugins: [router],
    },
  })
}

describe('shared components', () => {
  afterEach(() => {
    document.body.innerHTML = ''
  })

  it('disables busy buttons and exposes aria busy', () => {
    const wrapper = mount(BaseButton, {
      props: { busy: true },
      slots: { default: 'Save' },
    })

    expect(wrapper.get('button').attributes('disabled')).toBeDefined()
    expect(wrapper.get('button').attributes('aria-busy')).toBe('true')
  })

  it('emits input updates', async () => {
    const wrapper = mount(BaseInput, {
      props: { modelValue: '', id: 'email', type: 'email' },
    })

    await wrapper.get('input').setValue('admin@beehive.local')
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['admin@beehive.local'])
  })

  it('renders form field labels and errors', () => {
    const wrapper = mount(FormField, {
      props: { label: 'Email', forId: 'email', error: 'Required' },
      slots: { default: '<input id="email" />' },
    })

    expect(wrapper.get('label').text()).toBe('Email')
    expect(wrapper.get('[role="alert"]').text()).toBe('Required')
  })

  it('renders empty table state', () => {
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'title', label: 'Title' }],
        rows: [],
        emptyText: 'Nothing here',
      },
    })

    expect(wrapper.text()).toContain('Nothing here')
  })

  it('renders page loading skeleton rows', () => {
    const wrapper = mount(PageLoadingState, {
      props: { title: 'Loading users', rows: 3 },
    })

    expect(wrapper.attributes('aria-label')).toBe('Loading users')
    expect(wrapper.findAll('.skeleton-block')).toHaveLength(4)
  })

  it('renders avatar initials when no image is available', () => {
    const wrapper = mount(UserAvatar, {
      props: { name: 'Admin Editor' },
    })

    expect(wrapper.text()).toBe('AE')
  })

  it('emits logout from the account menu', async () => {
    const wrapper = await mountWithRouter(UserAccountMenu, {
      props: {
        user: {
          user_id: 'user_mock_admin',
          username: 'admin',
          email: 'admin@beehive.local',
          nickname: 'Admin',
          avatar_url: '',
          role: 'admin',
          status: 'active',
        },
      },
    })

    await wrapper.get('.account-menu__summary').trigger('click')
    await nextTick()
    document.body.querySelector<HTMLButtonElement>('.account-menu__item--danger')?.click()
    expect(wrapper.emitted('logout')).toHaveLength(1)
  })

  it('omits Studio shortcuts from the Studio account menu', async () => {
    const wrapper = await mountWithRouter(UserAccountMenu, {
      props: {
        surface: 'studio',
        user: {
          user_id: 'user_mock_admin',
          username: 'admin',
          email: 'admin@beehive.local',
          nickname: 'Admin',
          avatar_url: '',
          role: 'admin',
          status: 'active',
        },
      },
    })

    await wrapper.get('.account-menu__summary').trigger('click')
    await nextTick()

    expect(document.body.textContent).not.toContain('Studio')
    expect(document.body.textContent).not.toContain('Users')
    expect(document.body.textContent).not.toContain('Profile')
    expect(document.body.textContent).not.toContain('Change password')
    expect(document.body.textContent).toContain('Logout')
  })

  it('renders login and register actions in the public account menu', async () => {
    const wrapper = await mountWithRouter(UserAccountMenu, {
      props: {
        surface: 'public',
        user: null,
      },
    })

    await wrapper.get('.account-menu__summary').trigger('click')
    await nextTick()

    expect(document.body.textContent).toContain('Login')
    expect(document.body.textContent).toContain('Register')
  })
})
