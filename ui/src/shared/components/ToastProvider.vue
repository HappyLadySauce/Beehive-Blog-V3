<script setup lang="ts">
import { AlertCircle, CheckCircle2, Info, TriangleAlert, X } from 'lucide-vue-next'
import { computed } from 'vue'

import { useToast } from '@/shared/composables'
import type { ToastTone } from '@/shared/composables'

const { toasts, dismissToast } = useToast()

const toneIcons: Record<ToastTone, typeof Info> = {
  info: Info,
  success: CheckCircle2,
  warning: TriangleAlert,
  danger: AlertCircle,
}

function iconFor(tone: ToastTone): typeof Info {
  return toneIcons[tone]
}

const hasToasts = computed(() => toasts.value.length > 0)
</script>

<template>
  <Teleport to="body">
    <section v-if="hasToasts" class="toast-region" aria-label="Notifications" aria-live="polite">
      <TransitionGroup name="toast-list" tag="div" class="toast-region__list">
        <article v-for="toast in toasts" :key="toast.id" class="toast" :class="`toast--${toast.tone}`">
          <component :is="iconFor(toast.tone)" class="toast__icon" :size="18" aria-hidden="true" />
          <div class="toast__content">
            <strong>{{ toast.title }}</strong>
            <p v-if="toast.message">{{ toast.message }}</p>
          </div>
          <button class="toast__close" type="button" :aria-label="`Dismiss ${toast.title}`" @click="dismissToast(toast.id)">
            <X :size="16" aria-hidden="true" />
          </button>
        </article>
      </TransitionGroup>
    </section>
  </Teleport>
</template>

<style scoped>
.toast-region {
  position: fixed;
  z-index: 80;
  right: 20px;
  bottom: 20px;
  width: min(420px, calc(100vw - 32px));
  pointer-events: none;
}

.toast-region__list {
  display: grid;
  gap: 10px;
}

.toast {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: 12px;
  align-items: start;
  border: 1px solid var(--bb-color-line);
  border-left-width: 4px;
  border-radius: 8px;
  padding: 12px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface-elevated);
  box-shadow: var(--bb-shadow-panel);
  pointer-events: auto;
}

.toast--info {
  border-left-color: var(--bb-color-info);
}

.toast--success {
  border-left-color: var(--bb-color-success);
}

.toast--warning {
  border-left-color: var(--bb-color-warning);
}

.toast--danger {
  border-left-color: var(--bb-color-danger);
}

.toast__icon {
  margin-top: 2px;
  color: var(--bb-color-primary);
}

.toast--success .toast__icon {
  color: var(--bb-color-success);
}

.toast--warning .toast__icon {
  color: var(--bb-color-warning);
}

.toast--danger .toast__icon {
  color: var(--bb-color-danger);
}

.toast__content {
  min-width: 0;
  display: grid;
  gap: 3px;
}

.toast__content strong,
.toast__content p {
  margin: 0;
}

.toast__content p {
  color: var(--bb-color-muted);
}

.toast__close {
  width: 30px;
  height: 30px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 0;
  border-radius: 6px;
  color: var(--bb-color-muted);
  background: transparent;
}

.toast__close:hover,
.toast__close:focus-visible {
  outline: none;
  color: var(--bb-color-text);
  background: var(--bb-color-subtle);
}

.toast-list-enter-active,
.toast-list-leave-active {
  transition: opacity 160ms ease, transform 160ms ease;
}

.toast-list-enter-from,
.toast-list-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
