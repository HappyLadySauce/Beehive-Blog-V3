import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'

function installMatchMedia(matches: boolean): void {
  Object.defineProperty(window, 'matchMedia', {
    configurable: true,
    writable: true,
    value: vi.fn().mockImplementation((query: string) => ({
      matches,
      media: query,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    })),
  })
}

describe('theme infrastructure', () => {
  beforeEach(() => {
    vi.resetModules()
    document.documentElement.removeAttribute('data-theme')
    document.documentElement.style.colorScheme = ''
  })

  it('uses system preference by default and persists manual theme changes', async () => {
    installMatchMedia(true)
    const { useTheme } = await import('@/shared/composables/useTheme')

    const { theme, resolvedTheme, setTheme } = useTheme()

    expect(theme.value).toBe('system')
    expect(resolvedTheme.value).toBe('dark')
    expect(document.documentElement.dataset.theme).toBe('dark')

    setTheme('light')

    expect(theme.value).toBe('light')
    expect(resolvedTheme.value).toBe('light')
    expect(window.localStorage.getItem('beehive.ui.theme')).toBe('light')
    expect(document.documentElement.dataset.theme).toBe('light')
  })

  it('toggles between resolved light and dark themes from the button', async () => {
    installMatchMedia(false)
    const { default: ThemeToggle } = await import('@/shared/components/ThemeToggle.vue')
    const wrapper = mount(ThemeToggle)

    expect(wrapper.get('button').attributes('aria-label')).toBe('Switch to dark theme')

    await wrapper.get('button').trigger('click')

    expect(document.documentElement.dataset.theme).toBe('dark')
    expect(wrapper.get('button').attributes('aria-label')).toBe('Switch to light theme')
  })
})

describe('toast infrastructure', () => {
  beforeEach(() => {
    vi.resetModules()
  })

  it('pushes and dismisses visible toasts by tone', async () => {
    const { useToast } = await import('@/shared/composables/useToast')
    const { pushToast, dismissToast, toasts } = useToast()

    const id = pushToast({
      tone: 'success',
      title: 'Saved',
      message: 'Your changes are available.',
      timeoutMs: 0,
    })

    expect(toasts.value).toEqual([
      expect.objectContaining({
        id,
        tone: 'success',
        title: 'Saved',
      }),
    ])

    dismissToast(id)
    expect(toasts.value).toHaveLength(0)
  })

  it('renders toast provider content with live region semantics', async () => {
    const { useToast } = await import('@/shared/composables/useToast')
    const { default: ToastProvider } = await import('@/shared/components/ToastProvider.vue')
    const wrapper = mount(ToastProvider, { attachTo: document.body })

    useToast().pushToast({ tone: 'warning', title: 'Review needed', timeoutMs: 0 })
    await nextTick()

    expect(document.body.querySelector('[aria-live="polite"]')?.textContent).toContain('Review needed')

    wrapper.unmount()
  })
})

describe('confirm infrastructure', () => {
  beforeEach(() => {
    vi.resetModules()
    document.body.innerHTML = ''
  })

  it('resolves true when the provider confirm button is clicked', async () => {
    const { useConfirm } = await import('@/shared/composables/useConfirm')
    const { default: ConfirmDialogProvider } = await import('@/shared/components/ConfirmDialogProvider.vue')
    const wrapper = mount(ConfirmDialogProvider, { attachTo: document.body })
    const { confirm } = useConfirm()

    const result = confirm({
      title: 'Reset password?',
      message: 'This action changes user credentials.',
      confirmText: 'Reset password',
      tone: 'danger',
    })
    await nextTick()

    expect(document.body.textContent).toContain('Reset password?')
    const buttons = [...document.body.querySelectorAll('button')]
    buttons.find((button) => button.textContent?.includes('Reset password'))?.click()

    await expect(result).resolves.toBe(true)
    wrapper.unmount()
  })

  it('resolves false when the provider cancel button is clicked', async () => {
    const { useConfirm } = await import('@/shared/composables/useConfirm')
    const { default: ConfirmDialogProvider } = await import('@/shared/components/ConfirmDialogProvider.vue')
    const wrapper = mount(ConfirmDialogProvider, { attachTo: document.body })
    const { confirm } = useConfirm()

    const result = confirm({
      title: 'Discard changes?',
      message: 'Unsaved edits will be lost.',
    })
    await nextTick()

    const buttons = [...document.body.querySelectorAll('button')]
    buttons.find((button) => button.textContent?.includes('Cancel'))?.click()

    await expect(result).resolves.toBe(false)
    wrapper.unmount()
  })
})
