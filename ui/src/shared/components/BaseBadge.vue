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
  color: var(--bb-color-success);
}

.bb-badge--warning {
  color: var(--bb-color-warning);
}

.bb-badge--danger {
  color: var(--bb-color-danger);
}
</style>
