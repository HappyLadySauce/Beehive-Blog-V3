<script setup lang="ts">
withDefaults(
  defineProps<{
    tone?: 'neutral' | 'leaf' | 'honey' | 'blue' | 'violet';
  }>(),
  {
    tone: 'neutral',
  },
);
</script>

<template>
  <span
    class="inline-flex items-center rounded-sm border px-2 py-0.5 text-12px font-600"
    :class="{
      'border-brand-line bg-brand-paper text-brand-muted': tone === 'neutral',
      'border-brand-leaf/25 bg-brand-leaf/10 text-brand-leaf': tone === 'leaf',
      'border-brand-honey/25 bg-brand-honey/10 text-brand-honey': tone === 'honey',
      'border-brand-blue/25 bg-brand-blue/10 text-brand-blue': tone === 'blue',
      'border-brand-violet/25 bg-brand-violet/10 text-brand-violet': tone === 'violet',
    }"
  >
    <slot />
  </span>
</template>
