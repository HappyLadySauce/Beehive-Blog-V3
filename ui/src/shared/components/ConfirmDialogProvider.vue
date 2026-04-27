<script setup lang="ts">
import { AlertTriangle } from 'lucide-vue-next'
import { computed } from 'vue'

import { useConfirm } from '@/shared/composables'

const { currentConfirm, resolveConfirm } = useConfirm()
const isOpen = computed(() => currentConfirm.value !== null)

function cancel(): void {
  resolveConfirm(false)
}

function accept(): void {
  resolveConfirm(true)
}
</script>

<template>
  <Teleport to="body">
    <Transition name="confirm-dialog">
      <div v-if="isOpen && currentConfirm" class="confirm-dialog" role="presentation" @keydown.esc="cancel">
        <button class="confirm-dialog__backdrop" type="button" aria-label="Cancel dialog" @click="cancel" />
        <section
          class="confirm-dialog__panel"
          role="dialog"
          aria-modal="true"
          :aria-labelledby="`${currentConfirm.id}-title`"
          :aria-describedby="`${currentConfirm.id}-message`"
        >
          <div class="confirm-dialog__icon" :class="{ 'confirm-dialog__icon--danger': currentConfirm.tone === 'danger' }">
            <AlertTriangle :size="20" aria-hidden="true" />
          </div>
          <div class="confirm-dialog__content">
            <h2 :id="`${currentConfirm.id}-title`">{{ currentConfirm.title }}</h2>
            <p :id="`${currentConfirm.id}-message`">{{ currentConfirm.message }}</p>
          </div>
          <div class="confirm-dialog__actions">
            <button class="confirm-dialog__button confirm-dialog__button--secondary" type="button" @click="cancel">
              {{ currentConfirm.cancelText }}
            </button>
            <button
              class="confirm-dialog__button"
              :class="{ 'confirm-dialog__button--danger': currentConfirm.tone === 'danger' }"
              type="button"
              autofocus
              @click="accept"
            >
              {{ currentConfirm.confirmText }}
            </button>
          </div>
        </section>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.confirm-dialog {
  position: fixed;
  z-index: 90;
  inset: 0;
  display: grid;
  place-items: center;
  padding: 20px;
}

.confirm-dialog__backdrop {
  position: absolute;
  inset: 0;
  border: 0;
  background: rgb(10 17 24 / 48%);
  backdrop-filter: blur(4px);
}

.confirm-dialog__panel {
  position: relative;
  width: min(440px, 100%);
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: 14px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 18px;
  background: var(--bb-color-surface-elevated);
  box-shadow: var(--bb-shadow-panel);
}

.confirm-dialog__icon {
  width: 38px;
  height: 38px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  color: var(--bb-color-warning);
  background: var(--bb-color-warning-soft);
}

.confirm-dialog__icon--danger {
  color: var(--bb-color-danger);
  background: var(--bb-color-danger-soft);
}

.confirm-dialog__content {
  display: grid;
  gap: 6px;
}

.confirm-dialog__content h2,
.confirm-dialog__content p {
  margin: 0;
}

.confirm-dialog__content h2 {
  font-size: 1rem;
}

.confirm-dialog__content p {
  color: var(--bb-color-muted);
}

.confirm-dialog__actions {
  grid-column: 1 / -1;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.confirm-dialog__button {
  min-height: 40px;
  border: 1px solid transparent;
  border-radius: 8px;
  padding: 0 14px;
  color: #fff;
  background: var(--bb-color-primary);
  font-weight: 700;
}

.confirm-dialog__button--secondary {
  color: var(--bb-color-text);
  border-color: var(--bb-color-line);
  background: var(--bb-color-surface);
}

.confirm-dialog__button--danger {
  background: var(--bb-color-danger);
}

.confirm-dialog__button:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.confirm-dialog-enter-active,
.confirm-dialog-leave-active {
  transition: opacity 140ms ease;
}

.confirm-dialog-enter-from,
.confirm-dialog-leave-to {
  opacity: 0;
}
</style>
