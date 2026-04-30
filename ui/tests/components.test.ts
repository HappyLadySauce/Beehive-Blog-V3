import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { createMemoryHistory, createRouter } from 'vue-router'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'

import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import DataTable from '@/shared/components/DataTable.vue'
import EmptyState from '@/shared/components/EmptyState.vue'
import FormField from '@/shared/components/FormField.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import BaseSelect from '@/shared/components/BaseSelect.vue'
import LocaleToggle from '@/shared/components/LocaleToggle.vue'
import UserAccountMenu from '@/shared/components/UserAccountMenu.vue'
import UserAvatar from '@/shared/components/UserAvatar.vue'
import { i18n, setLocale } from '@/shared/i18n'

function t(key: string, params?: Record<string, unknown>) {
  return String(i18n.global.t(key, params ?? {}))
}

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
      plugins: [i18n, router],
    },
  })
}

describe('shared components', () => {
  beforeEach(() => {
    setLocale('en-US')
  })

  afterEach(() => {
    document.body.innerHTML = ''
    setLocale('en-US')
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
      global: {
        plugins: [i18n],
      },
    })

    expect(wrapper.text()).toContain('Nothing here')
  })

  it('renders the default empty table visual when no rows are present', () => {
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'title', label: 'Title' }],
        rows: [],
        emptyText: 'Nothing here',
      },
      global: {
        plugins: [i18n],
      },
    })

    expect(wrapper.find('.data-table__empty-state').exists()).toBe(true)
    expect(wrapper.find('svg').exists()).toBe(true)
    expect(wrapper.text()).toContain('Nothing here')
  })

  it('renders blocking table loading state', () => {
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'title', label: 'Title' }],
        rows: [],
        loading: true,
        loadingMode: 'blocking',
        loadingTitle: 'Loading table',
      },
      global: {
        plugins: [i18n],
      },
    })

    expect(wrapper.find('.page-loading').exists()).toBe(true)
    expect(wrapper.find('.page-loading').attributes('aria-label')).toBe('Loading table')
  })

  it('keeps table rows visible while showing a refreshing hint', () => {
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'title', label: 'Title' }],
        rows: [{ title: 'Admin' }],
        loading: true,
        loadingMode: 'refreshing',
      },
      global: {
        plugins: [i18n],
      },
    })

    expect(wrapper.text()).toContain('Admin')
    expect(wrapper.text()).toContain('Refreshing')
    expect(wrapper.find('.data-table__refreshing').exists()).toBe(true)
  })

  it('renders custom visuals in empty state slots', () => {
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'title', label: 'Title' }],
        rows: [],
        emptyText: 'Nothing here',
      },
      slots: {
        emptyVisual: '<div class="custom-empty-visual">custom visual</div>',
      },
      global: {
        plugins: [i18n],
      },
    })

    expect(wrapper.find('.custom-empty-visual').exists()).toBe(true)
  })

  it('renders empty state visual slots in the shared empty state component', () => {
    const wrapper = mount(EmptyState, {
      props: {
        title: 'No content yet',
        description: 'Create your first record.',
      },
      slots: {
        visual: '<div class="empty-state-visual">visual</div>',
      },
    })

    expect(wrapper.find('.empty-state-visual').exists()).toBe(true)
  })

  it('supports a centered empty state alignment mode', () => {
    const wrapper = mount(EmptyState, {
      props: {
        title: 'No content yet',
        description: 'Create your first record.',
        align: 'center',
      },
    })

    expect(wrapper.classes()).toContain('empty-state--center')
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

  it('selects options from the custom select listbox', async () => {
    const wrapper = mount(BaseSelect, {
      attachTo: document.body,
      props: {
        modelValue: '',
        options: [
          { value: '', label: 'All roles' },
          { value: 'member', label: 'Member' },
          { value: 'admin', label: 'Admin' },
        ],
        ariaLabel: 'Role',
      },
      global: {
        plugins: [i18n],
      },
    })

    await wrapper.get('button').trigger('click')
    document.body.querySelectorAll<HTMLElement>('[role="option"]')[2]?.click()

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['admin'])
  })

  it('persists locale changes from the locale toggle', async () => {
    setLocale('zh-CN')
    const wrapper = mount(LocaleToggle, {
      attachTo: document.body,
      global: {
        plugins: [i18n],
      },
    })

    await wrapper.get('button').trigger('click')
    document.body.querySelectorAll<HTMLElement>('[role="option"]')[1]?.click()

    expect(window.localStorage.getItem('beehive.ui.locale')).toBe('en-US')
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

    expect(document.body.textContent).not.toContain(t('account.studio'))
    expect(document.body.textContent).not.toContain(t('account.users'))
    expect(document.body.textContent).not.toContain(t('account.profile'))
    expect(document.body.textContent).not.toContain(t('users.passwordDialog.selfTitle'))
    expect(document.body.textContent).toContain(t('account.logout'))
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

    expect(document.body.textContent).toContain(t('account.login'))
    expect(document.body.textContent).toContain(t('account.register'))
  })
})
