<script setup lang="ts">
import { computed } from 'vue';

export interface SelectOption {
  label: string;
  value: string;
}

const props = withDefaults(
  defineProps<{
    modelValue: string;
    label: string;
    options: SelectOption[];
    name?: string;
    id?: string;
    error?: string;
  }>(),
  {
    name: '',
    id: '',
    error: '',
  },
);

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();

const selectId = computed(() => props.id || `bb-${props.name || props.label.replace(/\s+/g, '-').toLowerCase()}`);
</script>

<template>
  <label class="grid gap-2" :for="selectId">
    <span class="text-13px font-600 text-brand-muted">{{ label }}</span>
    <select
      :id="selectId"
      :name="name"
      :value="modelValue"
      class="bb-focus h-10 w-full border border-brand-line rounded-md bg-brand-surface px-3 text-14px text-brand-ink outline-none transition-colors"
      @change="emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
    >
      <option v-for="option in options" :key="option.value" :value="option.value">
        {{ option.label }}
      </option>
    </select>
    <span v-if="error" class="text-12px text-red-600">{{ error }}</span>
  </label>
</template>
