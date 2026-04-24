<script setup lang="ts">
import { AlertCircle, CheckCircle2, Info, TriangleAlert } from 'lucide-vue-next';
import { computed } from 'vue';

const props = withDefaults(
  defineProps<{
    tone?: 'info' | 'success' | 'warning' | 'danger';
    title: string;
    description?: string;
  }>(),
  {
    tone: 'info',
    description: '',
  },
);

const icon = computed(() => {
  if (props.tone === 'success') {
    return CheckCircle2;
  }
  if (props.tone === 'warning') {
    return TriangleAlert;
  }
  if (props.tone === 'danger') {
    return AlertCircle;
  }
  return Info;
});
</script>

<template>
  <section
    class="grid grid-cols-[auto_1fr] gap-3 rounded-lg border p-4"
    :class="{
      'border-brand-blue/25 bg-brand-blue/8 text-brand-blue': tone === 'info',
      'border-brand-leaf/25 bg-brand-leaf/8 text-brand-leaf': tone === 'success',
      'border-brand-honey/25 bg-brand-honey/8 text-brand-honey': tone === 'warning',
      'border-red-500/25 bg-red-500/8 text-red-600': tone === 'danger',
    }"
  >
    <component :is="icon" class="mt-0.5 h-5 w-5" aria-hidden="true" />
    <div class="grid gap-1">
      <h2 class="m-0 text-14px font-800">{{ title }}</h2>
      <p v-if="description" class="m-0 text-13px leading-5 text-brand-muted">{{ description }}</p>
      <slot />
    </div>
  </section>
</template>
