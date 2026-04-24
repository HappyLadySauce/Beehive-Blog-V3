<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string;
    description?: string;
    inverse?: boolean;
  }>(),
  {
    description: '',
    inverse: false,
  },
);
</script>

<template>
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div class="min-w-0">
      <h2 class="m-0 text-18px font-800" :class="inverse ? 'text-white' : 'text-brand-ink'">{{ title }}</h2>
      <p v-if="description" class="m-0 mt-1 text-13px leading-5" :class="inverse ? 'text-white/50' : 'text-brand-muted'">
        {{ description }}
      </p>
    </div>
    <div v-if="$slots.actions" class="flex shrink-0 flex-wrap gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>
