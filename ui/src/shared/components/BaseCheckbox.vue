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
  <label class="flex items-start gap-3 rounded-md">
    <input
      type="checkbox"
      :checked="modelValue"
      :disabled="disabled"
      class="bb-focus mt-0.5 h-4 w-4 rounded-sm border-brand-line accent-brand-blue disabled:cursor-not-allowed disabled:opacity-50"
      @change="emit('update:modelValue', ($event.target as HTMLInputElement).checked)"
    />
    <span class="grid gap-1">
      <span class="text-14px font-700 text-brand-ink">{{ label }}</span>
      <span v-if="description" class="text-13px leading-5 text-brand-muted">{{ description }}</span>
    </span>
  </label>
</template>
