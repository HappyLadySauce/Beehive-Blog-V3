<script setup lang="ts">
import { Inbox } from 'lucide-vue-next';

withDefaults(
  defineProps<{
    title: string;
    description?: string;
  }>(),
  {
    description: '',
  },
);
</script>

<template>
  <div class="grid place-items-center border border-dashed border-brand-line rounded-lg bg-brand-surface px-6 py-12 text-center">
    <div class="grid max-w-80 gap-3 justify-items-center">
      <Inbox class="h-8 w-8 text-brand-muted" aria-hidden="true" />
      <h2 class="m-0 text-18px font-700 text-brand-ink">{{ title }}</h2>
      <p v-if="description" class="m-0 text-14px leading-6 text-brand-muted">{{ description }}</p>
      <slot />
    </div>
  </div>
</template>
