<script setup lang="ts">
withDefaults(
  defineProps<{
    tone?: 'neutral' | 'primary' | 'danger'
    disabled?: boolean
  }>(),
  {
    tone: 'neutral',
    disabled: false,
  },
)
</script>

<template>
  <button class="action-tag" :class="`action-tag--${tone}`" type="button" :disabled="disabled">
    <slot />
  </button>
</template>

<style scoped>
.action-tag {
  min-height: 30px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--bb-color-line);
  border-radius: 999px;
  padding: 0 12px;
  color: var(--bb-color-muted);
  background: var(--bb-color-surface);
  font-size: 0.82rem;
  font-weight: 700;
  transition: color 160ms ease, border-color 160ms ease, background-color 160ms ease, transform 160ms ease;
}

.action-tag:hover:not(:disabled) {
  transform: translateY(-1px);
  color: var(--bb-color-text);
  border-color: var(--bb-color-primary);
  background: var(--bb-color-primary-soft);
}

.action-tag:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.action-tag--primary {
  color: var(--bb-color-primary);
  border-color: var(--bb-color-primary);
  background: var(--bb-color-primary-soft);
}

.action-tag--danger {
  color: var(--bb-color-danger);
  border-color: var(--bb-color-danger);
  background: var(--bb-color-danger-soft);
}

.action-tag:disabled {
  opacity: 0.56;
}
</style>
