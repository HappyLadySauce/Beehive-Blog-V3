<script setup lang="ts">
withDefaults(defineProps<{
  title: string
  description?: string
  align?: 'start' | 'center'
}>(), {
  align: 'start',
})
</script>

<template>
  <section class="empty-state" :class="`empty-state--${align}`">
    <div v-if="$slots.visual" class="empty-state__visual">
      <slot name="visual" />
    </div>
    <h2>{{ title }}</h2>
    <p v-if="description">{{ description }}</p>
    <div v-if="$slots.default" class="empty-state__actions">
      <slot />
    </div>
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
  box-shadow: var(--bb-shadow-soft);
}

.empty-state--center {
  justify-items: center;
  text-align: center;
}

.empty-state__visual,
.empty-state h2,
.empty-state p {
  margin: 0;
}

.empty-state__visual {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--bb-color-muted);
}

.empty-state h2 {
  color: var(--bb-color-text);
  font-size: 1rem;
}

.empty-state__actions {
  display: inline-flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 4px;
}

.empty-state--center .empty-state__actions {
  justify-content: center;
}
</style>
