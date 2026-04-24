<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(
  defineProps<{
    modelValue: string;
    label: string;
    name?: string;
    id?: string;
    type?: string;
    placeholder?: string;
    autocomplete?: string;
    error?: string;
  }>(),
  {
    name: '',
    id: '',
    type: 'text',
    placeholder: '',
    autocomplete: 'off',
    error: '',
  },
);

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();

const inputId = computed(() => props.id || `bb-${props.name || props.label.replace(/\s+/g, '-').toLowerCase()}`);
</script>

<template>
  <label class="grid gap-2" :for="inputId">
    <span class="text-13px font-600 text-brand-muted">{{ label }}</span>
    <input
      :id="inputId"
      :name="name"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      class="bb-focus h-10 w-full border border-brand-line rounded-md bg-brand-surface px-3 text-14px text-brand-ink outline-none transition-colors placeholder:text-brand-muted"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <span v-if="error" class="text-12px text-red-600">{{ error }}</span>
  </label>
</template>
