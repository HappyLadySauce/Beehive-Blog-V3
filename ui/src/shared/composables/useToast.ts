import { computed, readonly, shallowRef } from 'vue'

export type ToastTone = 'info' | 'success' | 'warning' | 'danger'

export interface ToastMessage {
  id: string
  tone: ToastTone
  title: string
  message: string
}

export interface PushToastInput {
  tone?: ToastTone
  title: string
  message?: string
  timeoutMs?: number
}

const toasts = shallowRef<ToastMessage[]>([])
const timers = new Map<string, number>()

function createToastId(): string {
  return `toast_${Date.now()}_${Math.random().toString(16).slice(2)}`
}

function dismissToast(id: string): void {
  const timer = timers.get(id)
  if (timer !== undefined) {
    window.clearTimeout(timer)
    timers.delete(id)
  }
  toasts.value = toasts.value.filter((toast) => toast.id !== id)
}

function pushToast(input: PushToastInput): string {
  const id = createToastId()
  const timeoutMs = input.timeoutMs ?? 4200
  const nextToast: ToastMessage = {
    id,
    tone: input.tone ?? 'info',
    title: input.title,
    message: input.message ?? '',
  }

  toasts.value = [...toasts.value, nextToast]
  if (timeoutMs > 0 && typeof window !== 'undefined') {
    timers.set(id, window.setTimeout(() => dismissToast(id), timeoutMs))
  }

  return id
}

const visibleToasts = computed(() => toasts.value)

export function useToast() {
  return {
    toasts: readonly(visibleToasts),
    pushToast,
    dismissToast,
  }
}
