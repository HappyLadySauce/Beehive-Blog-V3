<script setup lang="ts">
defineProps<{
  label: string
  forId: string
  error?: string
  hint?: string
}>()
</script>

<template>
  <div class="form-field">
    <label class="form-field__label" :for="forId">{{ label }}</label>
    <slot :described-by="error ? `${forId}-error` : hint ? `${forId}-hint` : undefined" />
    <p v-if="error" :id="`${forId}-error`" class="form-field__error" role="alert">{{ error }}</p>
    <p v-else-if="hint" :id="`${forId}-hint`" class="form-field__hint">{{ hint }}</p>
  </div>
</template>

<style scoped>
.form-field {
  display: grid;
  gap: 7px;
}

.form-field__label {
  color: var(--bb-color-text);
  font-size: 0.925rem;
  font-weight: 650;
}

.form-field__hint,
.form-field__error {
  margin: 0;
  font-size: 0.86rem;
}

.form-field__hint {
  color: var(--bb-color-muted);
}

.form-field__error {
  color: var(--bb-color-danger);
}
</style>
