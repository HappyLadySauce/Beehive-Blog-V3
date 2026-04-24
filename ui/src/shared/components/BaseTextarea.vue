<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(
  defineProps<{
    modelValue: string;
    label: string;
    name?: string;
    id?: string;
    placeholder?: string;
    rows?: number;
    error?: string;
  }>(),
  {
    name: '',
    id: '',
    placeholder: '',
    rows: 4,
    error: '',
  },
);

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();

const textareaId = computed(() => props.id || `bb-${props.name || props.label.replace(/\s+/g, '-').toLowerCase()}`);
</script>

<template>
  <label class="grid gap-2" :for="textareaId">
    <span class="text-13px font-600 text-brand-muted">{{ label }}</span>
    <textarea
      :id="textareaId"
      :name="name"
      :rows="rows"
      :value="modelValue"
      :placeholder="placeholder"
      class="bb-focus min-h-24 w-full resize-y border border-brand-line rounded-md bg-brand-surface px-3 py-2 text-14px leading-6 text-brand-ink outline-none transition-colors placeholder:text-brand-muted"
      @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
    ></textarea>
    <span v-if="error" class="text-12px text-red-600">{{ error }}</span>
  </label>
</template>
