<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  ariaLabel: string
  title?: string
  tone?: 'default' | 'primary' | 'danger'
  disabled?: boolean
}>(), {
  title: '',
  tone: 'default',
  disabled: false,
})

const buttonClass = computed(() => [
  'icon-action-button',
  `icon-action-button--${props.tone}`,
])
</script>

<template>
  <button
    :class="buttonClass"
    type="button"
    :disabled="disabled"
    :aria-label="ariaLabel"
    :title="title || ariaLabel"
  >
    <slot />
  </button>
</template>

<style scoped>
.icon-action-button {
  width: 32px;
  height: 32px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid transparent;
  border-radius: 8px;
  color: var(--bb-color-muted);
  background: transparent;
  text-decoration: none;
  transition: color 160ms ease, border-color 160ms ease, background-color 160ms ease, box-shadow 160ms ease, opacity 160ms ease;
}

.icon-action-button:hover:not(:disabled),
.icon-action-button:focus-visible {
  outline: none;
  color: var(--bb-color-text-strong);
  border-color: var(--bb-color-line);
  background: var(--bb-color-surface-elevated);
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.icon-action-button--primary {
  color: var(--bb-color-text);
}

.icon-action-button--danger {
  color: var(--bb-color-danger);
}

.icon-action-button:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
</style>
