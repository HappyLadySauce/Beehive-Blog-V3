<script setup lang="ts">
import SkeletonBlock from './SkeletonBlock.vue'

withDefaults(
  defineProps<{
    title?: string
    rows?: number
  }>(),
  {
    title: 'Loading',
    rows: 4,
  },
)
</script>

<template>
  <section class="page-loading" role="status" aria-live="polite" :aria-label="title">
    <SkeletonBlock width="38%" height="18px" />
    <div class="page-loading__rows">
      <SkeletonBlock v-for="row in rows" :key="row" height="46px" />
    </div>
  </section>
</template>

<style scoped>
.page-loading {
  display: grid;
  gap: 14px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 16px;
  background: var(--bb-color-surface);
}

.page-loading__rows {
  display: grid;
  gap: 10px;
}
</style>
