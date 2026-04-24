<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(
  defineProps<{
    label: string;
    variant?: 'ghost' | 'secondary' | 'danger';
    size?: 'sm' | 'md';
    disabled?: boolean;
  }>(),
  {
    variant: 'ghost',
    size: 'md',
    disabled: false,
  },
);

const classes = computed(() => [
  'bb-focus inline-flex items-center justify-center rounded-md border transition-colors disabled:cursor-not-allowed disabled:opacity-55',
  props.size === 'sm' ? 'h-8 w-8' : 'h-10 w-10',
  props.variant === 'ghost' && 'border-transparent bg-transparent text-brand-muted hover:bg-brand-surface hover:text-brand-ink',
  props.variant === 'secondary' && 'border-brand-line bg-brand-surface text-brand-ink hover:bg-brand-paper',
  props.variant === 'danger' && 'border-red-500/20 bg-red-500/10 text-red-600 hover:bg-red-500/15',
]);
</script>

<template>
  <button type="button" :aria-label="label" :title="label" :disabled="disabled" :class="classes">
    <slot />
  </button>
</template>
