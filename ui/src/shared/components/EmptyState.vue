<script setup lang="ts">
defineProps<{
  title: string
  description: string
}>()
</script>

<template>
  <section class="empty-state">
    <h2>{{ title }}</h2>
    <p>{{ description }}</p>
    <slot />
  </section>
</template>

<style scoped>
.empty-state {
  display: grid;
  justify-items: start;
  gap: 8px;
  border: 1px dashed var(--bb-color-line);
  border-radius: 8px;
  padding: 24px;
  color: var(--bb-color-muted);
  background: var(--bb-color-surface);
}

.empty-state h2,
.empty-state p {
  margin: 0;
}

.empty-state h2 {
  color: var(--bb-color-text);
  font-size: 1rem;
}
</style>
