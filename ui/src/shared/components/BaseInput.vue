<script setup lang="ts">
withDefaults(
  defineProps<{
    id?: string
    modelValue: string
    type?: string
    placeholder?: string
    autocomplete?: string
    inputmode?: 'none' | 'text' | 'decimal' | 'numeric' | 'tel' | 'search' | 'email' | 'url'
    disabled?: boolean
    required?: boolean
    invalid?: boolean
  }>(),
  {
    type: 'text',
    placeholder: '',
    autocomplete: 'off',
    inputmode: 'text',
    disabled: false,
    required: false,
    invalid: false,
  },
)

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<template>
  <input
    :id="id"
    class="bb-input"
    :class="{ 'bb-input--invalid': invalid }"
    :value="modelValue"
    :type="type"
    :placeholder="placeholder"
    :autocomplete="autocomplete"
    :inputmode="inputmode"
    :disabled="disabled"
    :required="required"
    :aria-invalid="invalid || undefined"
    @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
  />
</template>

<style scoped>
.bb-input {
  width: 100%;
  min-height: 44px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 0 12px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface);
}

.bb-input:focus-visible {
  outline: none;
  border-color: var(--bb-color-primary);
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.bb-input--invalid {
  border-color: var(--bb-color-danger);
}
</style>
