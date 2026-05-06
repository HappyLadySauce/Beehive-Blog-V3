import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createMemoryHistory, createRouter } from 'vue-router'
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

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/studio/files', component: StudioFilesPage },
      { path: '/studio/content/new', component: { template: '<div />' } },
    ],
  })
}

async function mountPage(initialPath = '/studio/files') {
  const router = createTestRouter()
  router.push(initialPath)
  await router.isReady()

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

  const wrapper = mount(StudioFilesPage, {
    attachTo: document.body,
    global: {
      plugins: [pinia, i18n, router, [VueQueryPlugin, { queryClient }]],
    },
  })

  return { wrapper, router }
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

  it('does not show unavailable alerts when the category and asset requests succeed', async () => {
    const { wrapper } = await mountPage()
    await flushPromises()

    expect(wrapper.text()).not.toContain('File categories unavailable')
    expect(wrapper.text()).not.toContain('File assets unavailable')
    expect(wrapper.text()).toContain('Upload file')
  })

  it('uploads a file, opens the preview drawer, and removes the asset after confirmation', async () => {
    const { wrapper } = await mountPage()
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

  it('shows the empty asset state after filtering to no matches', async () => {
    const { wrapper } = await mountPage()
    await flushPromises()

    await wrapper.get('#files-search').setValue('missing-file')
    await flushPromises()

    expect(wrapper.text()).toContain('No file assets yet.')
    expect(wrapper.text()).toContain('Upload the first file asset to see status and metadata here.')
  })

  it('restores the file types tab from the URL query', async () => {
    const { wrapper, router } = await mountPage('/studio/files?tab=types')
    await flushPromises()

    expect(router.currentRoute.value.query.tab).toBe('types')
    expect(wrapper.text()).toContain('Upload settings')
    expect(wrapper.text()).toContain('File types')
    expect(wrapper.text()).not.toContain('File assets unavailable')
  })
})
