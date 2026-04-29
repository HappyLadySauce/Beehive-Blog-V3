<script setup lang="ts">
import { nextTick, onBeforeUnmount, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = withDefaults(
  defineProps<{
    open: boolean
    title: string
    description?: string
    size?: 'sm' | 'md' | 'lg'
  }>(),
  {
    description: '',
    size: 'md',
  },
)

const emit = defineEmits<{
  close: []
}>()

const { t } = useI18n()
let previouslyFocused: HTMLElement | null = null

function closeDialog(): void {
  emit('close')
}

function onKeydown(event: KeyboardEvent): void {
  if (!props.open || event.key !== 'Escape') {
    return
  }
  event.preventDefault()
  closeDialog()
}

watch(
  () => props.open,
  async (open) => {
    if (!open) {
      previouslyFocused?.focus()
      previouslyFocused = null
      document.removeEventListener('keydown', onKeydown)
      return
    }
    previouslyFocused = document.activeElement instanceof HTMLElement ? document.activeElement : null
    document.addEventListener('keydown', onKeydown)
    await nextTick()
    document.querySelector<HTMLElement>('.modal-dialog__panel')?.focus()
  },
)

onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="modal-dialog" role="presentation" @mousedown.self="closeDialog">
      <section
        class="modal-dialog__panel"
        :class="`modal-dialog__panel--${size}`"
        role="dialog"
        aria-modal="true"
        :aria-label="title"
        tabindex="-1"
      >
        <header class="modal-dialog__header">
          <div class="modal-dialog__heading">
            <h2>{{ title }}</h2>
            <p v-if="description">{{ description }}</p>
          </div>
          <button class="modal-dialog__close" type="button" :aria-label="t('accessibility.closeDialog')" @click="closeDialog">×</button>
        </header>
        <div class="modal-dialog__body">
          <slot />
        </div>
        <footer v-if="$slots.footer" class="modal-dialog__footer">
          <slot name="footer" />
        </footer>
      </section>
    </div>
  </Teleport>
</template>

<style scoped>
.modal-dialog {
  position: fixed;
  z-index: 1200;
  inset: 0;
  display: grid;
  place-items: center;
  padding: 24px;
  background: color-mix(in srgb, var(--bb-color-text) 42%, transparent);
}

.modal-dialog__panel {
  width: min(100%, 620px);
  max-height: min(760px, calc(100vh - 48px));
  display: grid;
  grid-template-rows: auto minmax(0, 1fr) auto;
  overflow: hidden;
  border: 1px solid var(--bb-color-line);
  border-radius: 12px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-panel);
}

.modal-dialog__panel--sm {
  width: min(100%, 460px);
}

.modal-dialog__panel--lg {
  width: min(100%, 820px);
}

.modal-dialog__panel:focus {
  outline: none;
}

.modal-dialog__header {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 16px;
  align-items: start;
  border-bottom: 1px solid var(--bb-color-line);
  padding: 18px 20px;
}

.modal-dialog__heading {
  min-width: 0;
  display: grid;
  gap: 4px;
}

.modal-dialog__heading h2,
.modal-dialog__heading p {
  margin: 0;
}

.modal-dialog__heading h2 {
  color: var(--bb-color-text-strong);
  font-size: 1.1rem;
}

.modal-dialog__heading p {
  color: var(--bb-color-muted);
}

.modal-dialog__close {
  width: 36px;
  height: 36px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  color: var(--bb-color-muted);
  background: var(--bb-color-surface-elevated);
  font-size: 1.35rem;
  line-height: 1;
}

.modal-dialog__close:hover,
.modal-dialog__close:focus-visible {
  outline: none;
  color: var(--bb-color-text);
  border-color: var(--bb-color-primary);
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.modal-dialog__body {
  min-height: 0;
  overflow: auto;
  padding: 20px;
}

.modal-dialog__footer {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 10px;
  border-top: 1px solid var(--bb-color-line);
  padding: 14px 20px;
  background: var(--bb-color-subtle);
}

@media (max-width: 560px) {
  .modal-dialog {
    padding: 12px;
    place-items: end center;
  }

  .modal-dialog__panel {
    max-height: calc(100vh - 24px);
  }
}
</style>
