<script setup lang="ts">
interface TabItem {
  value: string;
  label: string;
}

defineProps<{
  tabs: TabItem[];
  modelValue: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();
</script>

<template>
  <div class="inline-flex rounded-md border border-brand-line bg-brand-surface p-1" role="tablist">
    <button
      v-for="tab in tabs"
      :key="tab.value"
      class="bb-focus h-8 rounded-sm px-3 text-13px font-600 transition-colors"
      :class="tab.value === modelValue ? 'bg-brand-ink text-brand-surface' : 'text-brand-muted hover:bg-brand-paper'"
      type="button"
      role="tab"
      :aria-selected="tab.value === modelValue"
      @click="emit('update:modelValue', tab.value)"
    >
      {{ tab.label }}
    </button>
  </div>
</template>
