<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    width?: string
    height?: string
    radius?: string
  }>(),
  {
    width: '100%',
    height: '16px',
    radius: '8px',
  },
)

const blockStyle = computed(() => ({
  width: props.width,
  height: props.height,
  borderRadius: props.radius,
}))
</script>

<template>
  <span class="skeleton-block" :style="blockStyle" aria-hidden="true" />
</template>

<style scoped>
.skeleton-block {
  display: block;
  overflow: hidden;
  background: linear-gradient(
    90deg,
    var(--bb-color-subtle) 0%,
    var(--bb-color-surface-elevated) 48%,
    var(--bb-color-subtle) 100%
  );
  background-size: 220% 100%;
  animation: skeleton-shimmer 1300ms ease-in-out infinite;
}

@keyframes skeleton-shimmer {
  0% {
    background-position: 120% 0;
  }

  100% {
    background-position: -120% 0;
  }
}
</style>
