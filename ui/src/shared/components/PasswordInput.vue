<script setup lang="ts">
import { Eye, EyeOff } from 'lucide-vue-next'
import { shallowRef } from 'vue'

import BaseInput from './BaseInput.vue'

defineProps<{
  id?: string
  modelValue: string
  autocomplete?: string
  invalid?: boolean
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const isVisible = shallowRef(false)
</script>

<template>
  <div class="password-input">
    <BaseInput
      :id="id"
      :model-value="modelValue"
      :type="isVisible ? 'text' : 'password'"
      :autocomplete="autocomplete ?? 'current-password'"
      :invalid="invalid"
      @update:model-value="$emit('update:modelValue', $event)"
    />
    <button
      type="button"
      class="password-input__toggle"
      :aria-label="isVisible ? 'Hide password' : 'Show password'"
      @click="isVisible = !isVisible"
    >
      <EyeOff v-if="isVisible" :size="18" aria-hidden="true" />
      <Eye v-else :size="18" aria-hidden="true" />
    </button>
  </div>
</template>

<style scoped>
.password-input {
  position: relative;
}

.password-input :deep(.bb-input) {
  padding-right: 48px;
}

.password-input__toggle {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 36px;
  height: 36px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 0;
  border-radius: 6px;
  color: var(--bb-color-muted);
  background: transparent;
}

.password-input__toggle:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}
</style>
