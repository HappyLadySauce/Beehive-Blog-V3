<script setup lang="ts">
interface TabItem {
  value: string;
  label: string;
}

const props = defineProps<{
  tabs: TabItem[];
  modelValue: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();

function selectByOffset(offset: number) {
  const currentIndex = props.tabs.findIndex((tab) => tab.value === props.modelValue);
  const nextIndex = (currentIndex + offset + props.tabs.length) % props.tabs.length;
  const next = props.tabs[nextIndex];
  if (next) {
    emit('update:modelValue', next.value);
  }
}
</script>

<template>
  <div class="inline-flex rounded-md border border-brand-line bg-brand-surface p-1" role="tablist" @keydown.left.prevent="selectByOffset(-1)" @keydown.right.prevent="selectByOffset(1)">
    <button
      v-for="tab in tabs"
      :key="tab.value"
      class="bb-focus h-8 rounded-sm px-3 text-13px font-600 transition-colors"
      :class="tab.value === modelValue ? 'bg-brand-ink text-brand-surface' : 'text-brand-muted hover:bg-brand-paper'"
      type="button"
      role="tab"
      :aria-selected="tab.value === modelValue"
      :tabindex="tab.value === modelValue ? 0 : -1"
      @click="emit('update:modelValue', tab.value)"
    >
      {{ tab.label }}
    </button>
  </div>
</template>
