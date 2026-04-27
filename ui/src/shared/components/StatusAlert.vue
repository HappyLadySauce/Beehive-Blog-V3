<script setup lang="ts">
withDefaults(
  defineProps<{
    tone?: 'info' | 'success' | 'warning' | 'danger'
    title?: string
  }>(),
  {
    tone: 'info',
    title: '',
  },
)
</script>

<template>
  <section class="status-alert" :class="`status-alert--${tone}`" role="status" aria-live="polite">
    <strong v-if="title" class="status-alert__title">{{ title }}</strong>
    <div class="status-alert__body">
      <slot />
    </div>
  </section>
</template>

<style scoped>
.status-alert {
  display: grid;
  gap: 4px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 12px;
  background: var(--bb-color-subtle);
}

.status-alert--success {
  color: var(--bb-color-success);
}

.status-alert--warning {
  color: var(--bb-color-warning);
}

.status-alert--danger {
  color: var(--bb-color-danger);
}

.status-alert__title {
  font-size: 0.92rem;
}

.status-alert__body {
  color: var(--bb-color-text);
}
</style>
