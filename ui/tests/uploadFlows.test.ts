import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { defineComponent, nextTick, shallowRef } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import StudioProfilePage from '@/pages/studio/StudioProfilePage.vue'
import ImageUploader from '@/shared/components/ImageUploader.vue'
import { i18n, setLocale } from '@/shared/i18n'

function t(key: string, params?: Record<string, unknown>) {
  return String(i18n.global.t(key, params ?? {}))
}

function imageFile(name = 'image.png', size = 4, type = 'image/png'): File {
  return new File([new Uint8Array(size)], name, { type })
}

function attachFiles(input: HTMLInputElement, files: File[]): void {
  Object.defineProperty(input, 'files', {
    configurable: true,
    value: files,
  })
}

function installSession() {
  const authStore = useAuthStore()
  authStore.applySession('access-token', 'refresh-token', 'session-id', {
    user_id: 'user_mock_admin',
    username: 'admin',
    email: 'admin@beehive.local',
    nickname: 'Admin',
    avatar_url: '',
    role: 'admin',
    status: 'active',
  }, 900)
  return authStore
}

async function mountStudioProfilePage() {
  const pinia = createPinia()
  setActivePinia(pinia)
  installSession()

  return mount(StudioProfilePage, {
    attachTo: document.body,
    global: {
      plugins: [pinia, i18n],
      stubs: {
        SsoProviderButtons: {
          template: '<div data-testid="sso-provider-buttons" />',
        },
      },
    },
  })
}

const ImageUploaderHost = defineComponent({
  components: {
    ImageUploader,
  },
  setup() {
    const coverUrl = shallowRef('')
    return {
      coverUrl,
    }
  },
  template: `
    <form>
      <ImageUploader v-model="coverUrl" scope="content_cover" />
      <output data-testid="cover-url">{{ coverUrl }}</output>
    </form>
  `,
})

async function mountImageUploaderHost() {
  const pinia = createPinia()
  setActivePinia(pinia)
  installSession()

  return mount(ImageUploaderHost, {
    attachTo: document.body,
    global: {
      plugins: [pinia, i18n],
    },
  })
}

describe('upload flows', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    setLocale('en-US')
    vi.restoreAllMocks()
  })

  afterEach(() => {
    document.body.innerHTML = ''
    vi.restoreAllMocks()
  })

  it('uploads an avatar on the profile page and persists it on save', async () => {
    Object.defineProperty(URL, 'createObjectURL', {
      configurable: true,
      value: vi.fn(() => 'blob:profile-avatar'),
    })
    const updateProfileSpy = vi.spyOn(studioApi, 'updateProfile').mockResolvedValue({
      user: {
        user_id: 'user_mock_admin',
        username: 'admin',
        email: 'admin@beehive.local',
        nickname: 'Admin',
        avatar_url: 'blob:profile-avatar',
        role: 'admin',
        status: 'active',
      },
    })

    const wrapper = await mountStudioProfilePage()
    const input = wrapper.get('.avatar-uploader__input').element as HTMLInputElement
    attachFiles(input, [imageFile('avatar.png')])

    await wrapper.get('.avatar-uploader__input').trigger('change')
    await flushPromises()
    await nextTick()

    expect(wrapper.get('.user-avatar img').attributes('src')).toBe('blob:profile-avatar')

    await wrapper.get('form').trigger('submit.prevent')
    await flushPromises()

    expect(updateProfileSpy).toHaveBeenCalledTimes(1)
    expect(updateProfileSpy).toHaveBeenCalledWith(
      {
        nickname: 'Admin',
        avatar_url: 'blob:profile-avatar',
      },
      { accessToken: 'access-token' },
    )
    expect(wrapper.text()).toContain(t('profile.status.profileSaved'))
  })

  it('shows localized upload validation errors on the profile page and keeps the previous avatar value', async () => {
    setLocale('zh-CN')
    const updateProfileSpy = vi.spyOn(studioApi, 'updateProfile').mockResolvedValue({
      user: {
        user_id: 'user_mock_admin',
        username: 'admin',
        email: 'admin@beehive.local',
        nickname: 'Admin',
        avatar_url: '',
        role: 'admin',
        status: 'active',
      },
    })
    const wrapper = await mountStudioProfilePage()
    const input = wrapper.get('.avatar-uploader__input').element as HTMLInputElement
    attachFiles(input, [imageFile('avatar.txt', 4, 'text/plain')])

    await wrapper.get('.avatar-uploader__input').trigger('change')
    await flushPromises()
    await nextTick()

    expect(wrapper.text()).toContain(t('uploads.fileTypeUnsupported'))
    expect(wrapper.find('.user-avatar img').exists()).toBe(false)

    await wrapper.get('form').trigger('submit.prevent')
    await flushPromises()

    expect(updateProfileSpy).toHaveBeenCalledWith(
      {
        nickname: 'Admin',
        avatar_url: '',
      },
      { accessToken: 'access-token' },
    )
  })

  it('uploads and removes a content cover image through the image uploader host', async () => {
    Object.defineProperty(URL, 'createObjectURL', {
      configurable: true,
      value: vi.fn(() => 'blob:content-cover'),
    })
    const wrapper = await mountImageUploaderHost()
    const input = wrapper.get('.image-uploader__input').element as HTMLInputElement
    attachFiles(input, [imageFile('cover.png')])

    await wrapper.get('.image-uploader__input').trigger('change')
    await flushPromises()
    await nextTick()

    expect(wrapper.get('[data-testid="cover-url"]').text()).toBe('blob:content-cover')
    expect(wrapper.get('.image-uploader__preview img').attributes('src')).toBe('blob:content-cover')

    await wrapper.get('button.bb-button--ghost').trigger('click')
    await flushPromises()

    expect(wrapper.get('[data-testid="cover-url"]').text()).toBe('')
    expect(wrapper.find('.image-uploader__preview img').exists()).toBe(false)
  })

  it('shows a size error when a content cover exceeds the allowed limit', async () => {
    const wrapper = await mountImageUploaderHost()
    const input = wrapper.get('.image-uploader__input').element as HTMLInputElement
    attachFiles(input, [imageFile('cover.png', 5 * 1024 * 1024 + 1)])

    await wrapper.get('.image-uploader__input').trigger('change')
    await flushPromises()
    await nextTick()

    expect(wrapper.text()).toContain(t('uploads.fileTooLarge'))
    expect(wrapper.get('[data-testid="cover-url"]').text()).toBe('')
  })

  it('renders localized profile and uploader copy when locale changes', async () => {
    setLocale('zh-CN')
    const wrapper = await mountStudioProfilePage()

    expect(wrapper.text()).toContain(t('profile.title'))
    expect(wrapper.text()).toContain(t('profile.description'))
    expect(wrapper.text()).toContain(t('uploads.avatarHint'))
    expect(wrapper.text()).toContain(t('profile.saveProfile'))
  })
})
