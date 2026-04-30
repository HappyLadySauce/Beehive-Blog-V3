<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = withDefaults(defineProps<{
  title?: string
}>(), {
  title: undefined,
})

const { t } = useI18n()
const label = computed(() => props.title || t('common.refreshing'))
</script>

<template>
  <div class="inline-loading-state" role="status" aria-live="polite">
    <span class="inline-loading-state__dot" aria-hidden="true" />
    <span>{{ label }}</span>
  </div>
</template>

<style scoped>
.inline-loading-state {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--bb-color-muted);
  font-size: 0.88rem;
  font-weight: 700;
}

.inline-loading-state__dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: var(--bb-color-primary);
  animation: inline-loading-pulse 1s ease-in-out infinite;
}

@keyframes inline-loading-pulse {
  0%,
  100% {
    opacity: 0.35;
    transform: scale(0.9);
  }

  50% {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
