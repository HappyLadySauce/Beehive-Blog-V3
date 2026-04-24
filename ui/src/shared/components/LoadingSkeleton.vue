<script setup lang="ts">
withDefaults(
  defineProps<{
    rows?: number;
  }>(),
  {
    rows: 3,
  },
);
</script>

<template>
  <div class="grid gap-3" aria-busy="true" aria-label="内容加载中">
    <span
      v-for="row in rows"
      :key="row"
      class="h-4 animate-pulse rounded-sm bg-brand-line"
      :class="row === rows ? 'w-2/3' : 'w-full'"
    />
  </div>
</template>
