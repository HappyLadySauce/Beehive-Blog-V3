<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: boolean;
    label: string;
    description?: string;
    disabled?: boolean;
  }>(),
  {
    description: '',
    disabled: false,
  },
);

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
}>();
</script>

<template>
  <button
    type="button"
    class="bb-focus flex w-full items-center justify-between gap-4 rounded-md text-left disabled:cursor-not-allowed disabled:opacity-50"
    :disabled="disabled"
    :aria-pressed="modelValue"
    @click="emit('update:modelValue', !modelValue)"
  >
    <span class="grid gap-1">
      <span class="text-14px font-700 text-brand-ink">{{ label }}</span>
      <span v-if="description" class="text-13px leading-5 text-brand-muted">{{ description }}</span>
    </span>
    <span
      class="relative h-6 w-11 shrink-0 rounded-full border border-brand-line transition-colors"
      :class="modelValue ? 'bg-brand-blue' : 'bg-brand-paper'"
      aria-hidden="true"
    >
      <span
        class="absolute top-0.5 h-5 w-5 rounded-full bg-brand-surface shadow-sm transition-transform"
        :class="modelValue ? 'translate-x-5' : 'translate-x-0.5'"
      />
    </span>
  </button>
</template>
