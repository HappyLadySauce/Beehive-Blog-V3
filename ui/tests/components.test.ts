import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import DataTable from '@/shared/components/DataTable.vue'
import FormField from '@/shared/components/FormField.vue'
import UserAccountMenu from '@/shared/components/UserAccountMenu.vue'
import UserAvatar from '@/shared/components/UserAvatar.vue'

describe('shared components', () => {
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

  it('renders avatar initials when no image is available', () => {
    const wrapper = mount(UserAvatar, {
      props: { name: 'Admin Editor' },
    })

    expect(wrapper.text()).toBe('AE')
  })

  it('emits logout from the account menu', async () => {
    const wrapper = mount(UserAccountMenu, {
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
      global: {
        stubs: {
          RouterLink: {
            props: ['to'],
            template: '<a><slot /></a>',
          },
        },
      },
    })

    await wrapper.get('button.account-menu__item--danger').trigger('click')
    expect(wrapper.emitted('logout')).toHaveLength(1)
  })
})
