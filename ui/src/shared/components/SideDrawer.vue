<script setup lang="ts">
import { X } from 'lucide-vue-next'
import { nextTick, useTemplateRef, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    open: boolean
    title: string
    description?: string
    size?: 'md' | 'lg'
  }>(),
  {
    size: 'md',
  },
)

const emit = defineEmits<{
  close: []
}>()

const panelRef = useTemplateRef<HTMLElement>('panel')

watch(
  () => props.open,
  async (open) => {
    if (!open) {
      return
    }
    await nextTick()
    panelRef.value?.focus()
  },
)
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer-fade">
      <div v-if="open" class="side-drawer" role="presentation" @click.self="emit('close')">
        <aside
          ref="panel"
          class="side-drawer__panel"
          :class="`side-drawer__panel--${size}`"
          role="dialog"
          aria-modal="true"
          :aria-label="title"
          tabindex="-1"
          @keydown.esc="emit('close')"
        >
          <header class="side-drawer__header">
            <div>
              <h2>{{ title }}</h2>
              <p v-if="description">{{ description }}</p>
            </div>
            <button class="side-drawer__close" type="button" aria-label="Close drawer" @click="emit('close')">
              <X :size="18" aria-hidden="true" />
            </button>
          </header>
          <div class="side-drawer__body">
            <slot />
          </div>
          <footer v-if="$slots.footer" class="side-drawer__footer">
            <slot name="footer" />
          </footer>
        </aside>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.side-drawer {
  position: fixed;
  z-index: 900;
  inset: 0;
  display: flex;
  justify-content: flex-end;
  background: rgb(15 23 42 / 32%);
}

.side-drawer__panel {
  width: min(560px, 100vw);
  height: 100%;
  display: grid;
  grid-template-rows: auto 1fr auto;
  border-left: 1px solid var(--bb-color-line);
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-panel);
}

.side-drawer__panel--lg {
  width: min(880px, 100vw);
}

.side-drawer__panel:focus {
  outline: none;
}

.side-drawer__header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  border-bottom: 1px solid var(--bb-color-line);
  padding: 20px;
}

.side-drawer__header h2,
.side-drawer__header p {
  margin: 0;
}

.side-drawer__header h2 {
  font-size: 1.25rem;
}

.side-drawer__header p {
  color: var(--bb-color-muted);
}

.side-drawer__close {
  width: 36px;
  height: 36px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  color: var(--bb-color-muted);
  background: var(--bb-color-surface-elevated);
}

.side-drawer__close:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.side-drawer__body {
  overflow: auto;
  padding: 20px;
}

.side-drawer__footer {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 10px;
  border-top: 1px solid var(--bb-color-line);
  padding: 14px 20px;
  background: var(--bb-color-subtle);
}

.drawer-fade-enter-active,
.drawer-fade-leave-active {
  transition: opacity 160ms ease;
}

.drawer-fade-enter-from,
.drawer-fade-leave-to {
  opacity: 0;
}
</style>
