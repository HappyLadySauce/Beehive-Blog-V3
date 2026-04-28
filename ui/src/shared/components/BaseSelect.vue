<script setup lang="ts">
import { Check, ChevronDown } from 'lucide-vue-next'
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, shallowRef, useTemplateRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'

export interface BaseSelectOption {
  value: string
  label: string
  disabled?: boolean
}

const props = withDefaults(defineProps<{
  modelValue: string
  options: BaseSelectOption[]
  id?: string
  placeholder?: string
  disabled?: boolean
  invalid?: boolean
  ariaLabel?: string
}>(), {
  id: undefined,
  placeholder: '',
  disabled: false,
  invalid: false,
  ariaLabel: undefined,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { t } = useI18n()
const triggerRef = useTemplateRef<HTMLButtonElement>('trigger')
const panelRef = useTemplateRef<HTMLElement>('panel')
const isOpen = shallowRef(false)
const activeIndex = shallowRef(-1)
const listboxId = `bb-select-${Math.random().toString(36).slice(2)}`
const panelStyle = reactive({
  top: '0px',
  left: '0px',
  width: '0px',
  maxHeight: '280px',
})

const enabledOptions = computed(() => props.options.filter((option) => !option.disabled))
const selectedIndex = computed(() => props.options.findIndex((option) => option.value === props.modelValue))
const selectedOption = computed(() => props.options[selectedIndex.value])
const selectedLabel = computed(() => selectedOption.value?.label ?? props.placeholder)
const activeOptionId = computed(() => (activeIndex.value >= 0 ? `${listboxId}-option-${activeIndex.value}` : undefined))

function updatePanelPosition(): void {
  const trigger = triggerRef.value
  if (!trigger) {
    return
  }

  const rect = trigger.getBoundingClientRect()
  const viewportPadding = 12
  const width = Math.max(rect.width, 180)
  const left = Math.min(Math.max(viewportPadding, rect.left), window.innerWidth - width - viewportPadding)
  const belowSpace = window.innerHeight - rect.bottom - viewportPadding
  const aboveSpace = rect.top - viewportPadding
  const maxHeight = Math.max(180, Math.min(300, Math.max(belowSpace, aboveSpace)))
  const openAbove = belowSpace < 180 && aboveSpace > belowSpace

  panelStyle.left = `${left}px`
  panelStyle.width = `${width}px`
  panelStyle.maxHeight = `${maxHeight}px`
  panelStyle.top = openAbove ? `${Math.max(viewportPadding, rect.top - maxHeight - 6)}px` : `${rect.bottom + 6}px`
}

async function openSelect(): Promise<void> {
  if (props.disabled || isOpen.value) {
    return
  }
  isOpen.value = true
  activeIndex.value = selectedIndex.value >= 0 ? selectedIndex.value : firstEnabledIndex()
  await nextTick()
  updatePanelPosition()
}

function closeSelect(): void {
  isOpen.value = false
}

function toggleSelect(): void {
  if (isOpen.value) {
    closeSelect()
    return
  }
  void openSelect()
}

function firstEnabledIndex(): number {
  return props.options.findIndex((option) => !option.disabled)
}

function nextEnabledIndex(start: number, direction: 1 | -1): number {
  if (enabledOptions.value.length === 0) {
    return -1
  }
  let index = start
  for (let step = 0; step < props.options.length; step += 1) {
    index = (index + direction + props.options.length) % props.options.length
    if (!props.options[index]?.disabled) {
      return index
    }
  }
  return -1
}

function selectOption(index: number): void {
  const option = props.options[index]
  if (!option || option.disabled) {
    return
  }
  emit('update:modelValue', option.value)
  closeSelect()
  triggerRef.value?.focus()
}

function handleTriggerKeydown(event: KeyboardEvent): void {
  if (props.disabled) {
    return
  }

  if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
    event.preventDefault()
    if (!isOpen.value) {
      void openSelect()
      return
    }
    activeIndex.value = nextEnabledIndex(activeIndex.value, event.key === 'ArrowDown' ? 1 : -1)
    return
  }

  if (event.key === 'Home') {
    event.preventDefault()
    activeIndex.value = firstEnabledIndex()
    return
  }

  if (event.key === 'End') {
    event.preventDefault()
    activeIndex.value = nextEnabledIndex(0, -1)
    return
  }

  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault()
    if (!isOpen.value) {
      void openSelect()
      return
    }
    selectOption(activeIndex.value)
    return
  }

  if (event.key === 'Escape') {
    event.preventDefault()
    closeSelect()
  }
}

function handleDocumentPointerDown(event: PointerEvent): void {
  if (!isOpen.value) {
    return
  }
  const target = event.target
  if (!(target instanceof Node)) {
    return
  }
  if (triggerRef.value?.contains(target) || panelRef.value?.contains(target)) {
    return
  }
  closeSelect()
}

function handleWindowChange(): void {
  if (isOpen.value) {
    updatePanelPosition()
  }
}

watch(() => props.modelValue, () => {
  activeIndex.value = selectedIndex.value
})

onMounted(() => {
  document.addEventListener('pointerdown', handleDocumentPointerDown)
  window.addEventListener('resize', handleWindowChange)
  window.addEventListener('scroll', handleWindowChange, true)
})

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', handleDocumentPointerDown)
  window.removeEventListener('resize', handleWindowChange)
  window.removeEventListener('scroll', handleWindowChange, true)
})
</script>

<template>
  <button
    :id="id"
    ref="trigger"
    class="base-select"
    type="button"
    role="combobox"
    aria-haspopup="listbox"
    :aria-controls="listboxId"
    :aria-expanded="isOpen"
    :aria-activedescendant="activeOptionId"
    :aria-label="ariaLabel"
    :aria-invalid="invalid || undefined"
    :disabled="disabled"
    @click="toggleSelect"
    @keydown="handleTriggerKeydown"
  >
    <span class="base-select__value" :class="{ 'base-select__value--placeholder': !selectedOption }">
      {{ selectedLabel }}
    </span>
    <ChevronDown class="base-select__chevron" :class="{ 'base-select__chevron--open': isOpen }" :size="16" aria-hidden="true" />
  </button>

  <Teleport to="body">
    <Transition name="base-select-panel">
      <div v-if="isOpen" ref="panel" class="base-select__panel" :style="panelStyle">
        <ul :id="listboxId" class="base-select__list" role="listbox" :aria-label="ariaLabel">
          <li
            v-for="(option, index) in options"
            :id="`${listboxId}-option-${index}`"
            :key="option.value"
            class="base-select__option"
            :class="{
              'base-select__option--active': index === activeIndex,
              'base-select__option--selected': option.value === modelValue,
              'base-select__option--disabled': option.disabled,
            }"
            role="option"
            :aria-selected="option.value === modelValue"
            :aria-disabled="option.disabled || undefined"
            @mouseenter="activeIndex = index"
            @click="selectOption(index)"
          >
            <span>{{ option.label }}</span>
            <Check v-if="option.value === modelValue" :size="15" aria-hidden="true" />
          </li>
          <li v-if="options.length === 0" class="base-select__empty">{{ t('common.noOptions') }}</li>
        </ul>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.base-select {
  width: 100%;
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 0 12px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface);
  font: inherit;
  text-align: left;
  transition: border-color 160ms ease, box-shadow 160ms ease, background-color 160ms ease;
}

.base-select:hover:not(:disabled) {
  border-color: color-mix(in srgb, var(--bb-color-primary) 48%, var(--bb-color-line));
}

.base-select:focus-visible,
.base-select[aria-expanded='true'] {
  outline: none;
  border-color: var(--bb-color-primary);
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.base-select:disabled {
  cursor: not-allowed;
  color: var(--bb-color-muted);
  background: var(--bb-color-subtle);
}

.base-select[aria-invalid='true'] {
  border-color: var(--bb-color-danger);
}

.base-select__value {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.base-select__value--placeholder {
  color: var(--bb-color-muted);
}

.base-select__chevron {
  flex: 0 0 auto;
  color: var(--bb-color-muted);
  transition: transform 160ms ease;
}

.base-select__chevron--open {
  transform: rotate(180deg);
}

.base-select__panel {
  position: fixed;
  z-index: 2200;
  overflow: auto;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  background: var(--bb-color-surface-elevated);
  box-shadow: var(--bb-shadow-panel);
}

.base-select__list {
  display: grid;
  gap: 2px;
  margin: 0;
  padding: 6px;
  list-style: none;
}

.base-select__option,
.base-select__empty {
  min-height: 38px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  border-radius: 6px;
  padding: 0 10px;
  color: var(--bb-color-text);
}

.base-select__option {
  cursor: pointer;
}

.base-select__option--active {
  background: var(--bb-color-subtle);
}

.base-select__option--selected {
  color: var(--bb-color-primary);
  font-weight: 750;
}

.base-select__option--disabled {
  cursor: not-allowed;
  color: var(--bb-color-muted);
  opacity: 0.55;
}

.base-select__empty {
  color: var(--bb-color-muted);
}

.base-select-panel-enter-active,
.base-select-panel-leave-active {
  transition: opacity 120ms ease, transform 120ms ease;
}

.base-select-panel-enter-from,
.base-select-panel-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
