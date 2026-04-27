<script setup lang="ts">
import { Moon, Sun } from 'lucide-vue-next'
import { computed } from 'vue'

import { useTheme } from '@/shared/composables'

const { resolvedTheme, toggleTheme } = useTheme()
const label = computed(() => (resolvedTheme.value === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'))
</script>

<template>
  <button class="theme-toggle" type="button" :aria-label="label" :title="label" @click="toggleTheme">
    <Moon v-if="resolvedTheme === 'light'" :size="18" aria-hidden="true" />
    <Sun v-else :size="18" aria-hidden="true" />
  </button>
</template>

<style scoped>
.theme-toggle {
  width: 42px;
  height: 42px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface-elevated);
  box-shadow: var(--bb-shadow-soft);
  transition: transform 160ms ease, border-color 160ms ease, background-color 160ms ease;
}

.theme-toggle:hover {
  transform: translateY(-1px);
  border-color: var(--bb-color-primary);
}

.theme-toggle:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}
</style>
