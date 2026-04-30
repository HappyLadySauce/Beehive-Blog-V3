import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { useAuthStore } from '@/features/auth/stores/authStore'
import StudioFilesPage from '@/pages/studio/StudioFilesPage.vue'
import { i18n, setLocale } from '@/shared/i18n'
import { useConfirm } from '@/shared/composables'

function imageFile(name = 'studio.png', size = 4, type = 'image/png'): File {
  return new File([new Uint8Array(size)], name, { type })
}

function attachFiles(input: HTMLInputElement, files: File[]): void {
  Object.defineProperty(input, 'files', {
    configurable: true,
    value: files,
  })
}

async function mountPage() {
  const pinia = createPinia()
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        refetchOnWindowFocus: false,
      },
    },
  })
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
  }, 900)

  return mount(StudioFilesPage, {
    attachTo: document.body,
    global: {
      plugins: [pinia, i18n, [VueQueryPlugin, { queryClient }]],
    },
  })
}

describe('studio files page', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    setLocale('en-US')
    Object.defineProperty(URL, 'createObjectURL', {
      configurable: true,
      value: vi.fn(() => 'blob:studio-file'),
    })
  })

  afterEach(() => {
    document.body.innerHTML = ''
    vi.restoreAllMocks()
  })

  it('uploads a file, opens the preview drawer, and removes the asset after confirmation', async () => {
    const wrapper = await mountPage()
    await flushPromises()

    const input = wrapper.get('.files-page__input').element as HTMLInputElement
    attachFiles(input, [imageFile()])
    await wrapper.get('.files-page__input').trigger('change')
    await flushPromises()

    expect(wrapper.text()).toContain('studio.png')
    expect(wrapper.find('[aria-label="View studio.png"]').exists()).toBe(true)

    await wrapper.get('[aria-label="View studio.png"]').trigger('click')
    await flushPromises()

    expect(document.body.textContent).toContain('File details')
    expect(document.body.textContent).toContain('studio.png')

    const { resolveConfirm } = useConfirm()
    await wrapper.get('[aria-label="Delete studio.png"]').trigger('click')
    resolveConfirm(true)
    await flushPromises()

    expect(wrapper.text()).not.toContain('studio.png')
  })
})
