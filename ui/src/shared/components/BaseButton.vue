<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    type?: 'button' | 'submit' | 'reset'
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
    busy?: boolean
    disabled?: boolean
  }>(),
  {
    type: 'button',
    variant: 'primary',
    busy: false,
    disabled: false,
  },
)

const buttonClass = computed(() => ['bb-button', `bb-button--${props.variant}`])
</script>

<template>
  <button :type="type" :class="buttonClass" :disabled="disabled || busy" :aria-busy="busy">
    <span v-if="busy" class="bb-button__spinner" aria-hidden="true" />
    <slot />
  </button>
</template>

<style scoped>
.bb-button {
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: 1px solid transparent;
  border-radius: 8px;
  padding: 0 16px;
  font-weight: 650;
  transition: background-color 160ms ease, border-color 160ms ease, color 160ms ease;
}

.bb-button:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.bb-button--primary {
  color: #fff;
  background: var(--bb-color-primary);
}

.bb-button--primary:hover:not(:disabled) {
  background: var(--bb-color-primary-strong);
}

.bb-button--secondary {
  color: var(--bb-color-text);
  border-color: var(--bb-color-line);
  background: var(--bb-color-surface);
}

.bb-button--ghost {
  color: var(--bb-color-muted);
  background: transparent;
}

.bb-button--danger {
  color: #fff;
  background: var(--bb-color-danger);
}

.bb-button:disabled {
  opacity: 0.62;
}

.bb-button__spinner {
  width: 16px;
  height: 16px;
  border: 2px solid currentColor;
  border-right-color: transparent;
  border-radius: 999px;
  animation: bb-spin 700ms linear infinite;
}

@keyframes bb-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
