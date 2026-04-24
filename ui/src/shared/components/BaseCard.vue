<script setup lang="ts">
withDefaults(
  defineProps<{
    interactive?: boolean;
  }>(),
  {
    interactive: false,
  },
);
</script>

<template>
  <section
    class="bb-panel"
    :class="interactive ? 'transition-transform hover:-translate-y-0.5 hover:border-brand-blue/35' : ''"
  >
    <slot />
  </section>
</template>
