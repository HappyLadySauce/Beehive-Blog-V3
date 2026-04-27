<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  value: string
}>()

const normalized = computed(() => props.value.toLowerCase())
const tone = computed(() => {
  switch (normalized.value) {
    case 'active':
    case 'published':
    case 'public':
    case 'allowed':
      return 'success'
    case 'pending':
    case 'review':
    case 'member':
    case 'private':
      return 'warning'
    case 'disabled':
    case 'locked':
    case 'deleted':
    case 'archived':
    case 'denied':
      return 'danger'
    default:
      return 'neutral'
  }
})
const label = computed(() => (props.value || 'unknown').replace(/_/g, ' '))
</script>

<template>
  <span class="status-badge" :class="`status-badge--${tone}`">{{ label }}</span>
</template>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  min-height: 26px;
  border: 1px solid var(--bb-color-line);
  border-radius: 999px;
  padding: 0 10px;
  color: var(--bb-color-muted);
  background: var(--bb-color-surface);
  font-size: 0.78rem;
  font-weight: 760;
  text-transform: capitalize;
}

.status-badge--success {
  border-color: var(--bb-color-success);
  color: var(--bb-color-success);
  background: var(--bb-color-success-soft);
}

.status-badge--warning {
  border-color: var(--bb-color-warning);
  color: var(--bb-color-warning);
  background: var(--bb-color-warning-soft);
}

.status-badge--danger {
  border-color: var(--bb-color-danger);
  color: var(--bb-color-danger);
  background: var(--bb-color-danger-soft);
}
</style>
