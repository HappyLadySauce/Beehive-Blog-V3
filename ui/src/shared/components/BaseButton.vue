<script setup lang="ts">
import { LoaderCircle } from 'lucide-vue-next';
import { computed } from 'vue';

type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'danger';
type ButtonSize = 'sm' | 'md';

const props = withDefaults(
  defineProps<{
    variant?: ButtonVariant;
    size?: ButtonSize;
    type?: 'button' | 'submit' | 'reset';
    disabled?: boolean;
    busy?: boolean;
  }>(),
  {
    variant: 'secondary',
    size: 'md',
    type: 'button',
    disabled: false,
    busy: false,
  },
);

const classes = computed(() => [
  'bb-focus inline-flex items-center justify-center gap-2 rounded-md border font-500 transition-colors disabled:cursor-not-allowed disabled:opacity-55',
  props.size === 'sm' ? 'h-8 px-3 text-13px' : 'h-10 px-4 text-14px',
  props.variant === 'primary' && 'border-brand-ink bg-brand-ink text-brand-surface hover:opacity-90',
  props.variant === 'secondary' && 'border-brand-line bg-brand-surface text-brand-ink hover:bg-brand-paper',
  props.variant === 'ghost' && 'border-transparent bg-transparent text-brand-muted hover:bg-brand-surface',
  props.variant === 'danger' && 'border-red-500/20 bg-red-500/10 text-red-600 hover:bg-red-500/15',
]);
</script>

<template>
  <button :type="type" :disabled="disabled || busy" :class="classes">
    <LoaderCircle v-if="busy" class="h-4 w-4 animate-spin" aria-hidden="true" />
    <slot />
  </button>
</template>
