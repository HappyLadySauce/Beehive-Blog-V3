<script setup lang="ts">
defineProps<{
  label: string
  value?: string | number | null
}>()
</script>

<template>
  <div class="readonly-field">
    <span class="readonly-field__label">{{ label }}</span>
    <strong class="readonly-field__value">{{ value || 'None' }}</strong>
  </div>
</template>

<style scoped>
.readonly-field {
  display: grid;
  gap: 4px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 12px;
  background: var(--bb-color-subtle);
}

.readonly-field__label {
  color: var(--bb-color-muted);
  font-size: 0.78rem;
  font-weight: 760;
  text-transform: uppercase;
}

.readonly-field__value {
  color: var(--bb-color-text);
  overflow-wrap: anywhere;
}
</style>
