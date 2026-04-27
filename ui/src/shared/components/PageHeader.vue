<script setup lang="ts">
defineProps<{
  eyebrow?: string
  title: string
  description?: string
}>()
</script>

<template>
  <header class="page-header">
    <div>
      <p v-if="eyebrow" class="page-header__eyebrow">{{ eyebrow }}</p>
      <h1>{{ title }}</h1>
      <p v-if="description">{{ description }}</p>
    </div>
    <div v-if="$slots.actions" class="page-header__actions">
      <slot name="actions" />
    </div>
  </header>
</template>

<style scoped>
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.page-header h1,
.page-header p {
  margin: 0;
}

.page-header__eyebrow {
  margin-bottom: 6px !important;
  color: var(--bb-color-primary);
  font-size: 0.78rem;
  font-weight: 800;
  text-transform: uppercase;
}

.page-header h1 {
  font-size: clamp(1.5rem, 4vw, 2.4rem);
  line-height: 1.08;
}

.page-header p {
  max-width: 680px;
  color: var(--bb-color-muted);
}

.page-header__actions {
  display: flex;
  gap: 8px;
}

@media (max-width: 640px) {
  .page-header {
    flex-direction: column;
  }
}
</style>
