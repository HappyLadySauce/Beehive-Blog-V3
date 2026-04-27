<script setup lang="ts">
withDefaults(
  defineProps<{
    tone?: 'neutral' | 'success' | 'warning' | 'danger'
  }>(),
  {
    tone: 'neutral',
  },
)
</script>

<template>
  <span class="bb-badge" :class="`bb-badge--${tone}`">
    <slot />
  </span>
</template>

<style scoped>
.bb-badge {
  display: inline-flex;
  align-items: center;
  min-height: 24px;
  border: 1px solid var(--bb-color-line);
  border-radius: 999px;
  padding: 0 10px;
  color: var(--bb-color-muted);
  background: var(--bb-color-surface);
  font-size: 0.78rem;
  font-weight: 700;
}

.bb-badge--success {
  border-color: var(--bb-color-success);
  color: var(--bb-color-success);
  background: var(--bb-color-success-soft);
}

.bb-badge--warning {
  border-color: var(--bb-color-warning);
  color: var(--bb-color-warning);
  background: var(--bb-color-warning-soft);
}

.bb-badge--danger {
  border-color: var(--bb-color-danger);
  color: var(--bb-color-danger);
  background: var(--bb-color-danger-soft);
}
</style>
