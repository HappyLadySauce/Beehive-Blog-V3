import { readonly, shallowRef } from 'vue'

export interface ConfirmOptions {
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  tone?: 'default' | 'danger'
}

interface ConfirmState extends Required<ConfirmOptions> {
  id: string
}

const currentConfirm = shallowRef<ConfirmState | null>(null)
let resolver: ((confirmed: boolean) => void) | null = null

function resolveConfirm(confirmed: boolean): void {
  if (resolver) {
    resolver(confirmed)
    resolver = null
  }
  currentConfirm.value = null
}

function confirm(options: ConfirmOptions): Promise<boolean> {
  if (resolver) {
    resolveConfirm(false)
  }

  currentConfirm.value = {
    id: `confirm_${Date.now()}_${Math.random().toString(16).slice(2)}`,
    title: options.title,
    message: options.message,
    confirmText: options.confirmText ?? 'Confirm',
    cancelText: options.cancelText ?? 'Cancel',
    tone: options.tone ?? 'default',
  }

  return new Promise<boolean>((resolve) => {
    resolver = resolve
  })
}

export function useConfirm() {
  return {
    currentConfirm: readonly(currentConfirm),
    confirm,
    resolveConfirm,
  }
}
