<script setup lang="ts">
withDefaults(
  defineProps<{
    eyebrow?: string;
    title: string;
    description?: string;
    inverse?: boolean;
  }>(),
  {
    eyebrow: '',
    description: '',
    inverse: false,
  },
);
</script>

<template>
  <header class="flex flex-wrap items-end justify-between gap-4">
    <div class="min-w-0">
      <p v-if="eyebrow" class="m-0 text-13px font-800" :class="inverse ? 'text-white/52' : 'text-brand-leaf'">
        {{ eyebrow }}
      </p>
      <h1 class="m-0 mt-2 text-28px font-900 leading-9 sm:text-32px" :class="inverse ? 'text-white' : 'text-brand-ink'">
        {{ title }}
      </h1>
      <p v-if="description" class="m-0 mt-2 max-w-70ch text-14px leading-6" :class="inverse ? 'text-white/52' : 'text-brand-muted'">
        {{ description }}
      </p>
    </div>
    <div v-if="$slots.actions" class="flex shrink-0 flex-wrap gap-2">
      <slot name="actions" />
    </div>
  </header>
</template>
