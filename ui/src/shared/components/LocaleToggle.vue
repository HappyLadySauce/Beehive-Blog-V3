<script setup lang="ts">
import { Check, Globe2 } from 'lucide-vue-next'
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, shallowRef, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'

import { availableLocales, useLocale, type AppLocale } from '@/shared/i18n'

const { t } = useI18n()
const { locale, setLocale } = useLocale()
const triggerRef = useTemplateRef<HTMLButtonElement>('trigger')
const panelRef = useTemplateRef<HTMLElement>('panel')
const isOpen = shallowRef(false)
const panelStyle = reactive({
  top: '0px',
  left: '0px',
})

const options = computed(() => availableLocales.map((value) => ({
  value,
  label: value === 'zh-CN' ? t('locale.zhCN') : t('locale.enUS'),
})))
const currentLabel = computed(() => options.value.find((option) => option.value === locale.value)?.label ?? t('locale.label'))

function updatePanelPosition(): void {
  const trigger = triggerRef.value
  if (!trigger) {
    return
  }

  const rect = trigger.getBoundingClientRect()
  const width = 168
  const viewportPadding = 12
  const left = Math.min(Math.max(viewportPadding, rect.right - width), window.innerWidth - width - viewportPadding)
  panelStyle.left = `${left}px`
  panelStyle.top = `${rect.bottom + 8}px`
}

async function openPanel(): Promise<void> {
  isOpen.value = true
  await nextTick()
  updatePanelPosition()
  panelRef.value?.focus()
}

function closePanel(): void {
  isOpen.value = false
}

function togglePanel(): void {
  if (isOpen.value) {
    closePanel()
    return
  }
  void openPanel()
}

function selectLocale(value: AppLocale): void {
  setLocale(value)
  closePanel()
  triggerRef.value?.focus()
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
  closePanel()
}

function handleDocumentKeydown(event: KeyboardEvent): void {
  if (event.key === 'Escape') {
    closePanel()
  }
}

function handleWindowChange(): void {
  if (isOpen.value) {
    updatePanelPosition()
  }
}

onMounted(() => {
  document.addEventListener('pointerdown', handleDocumentPointerDown)
  document.addEventListener('keydown', handleDocumentKeydown)
  window.addEventListener('resize', handleWindowChange)
  window.addEventListener('scroll', handleWindowChange, true)
})

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', handleDocumentPointerDown)
  document.removeEventListener('keydown', handleDocumentKeydown)
  window.removeEventListener('resize', handleWindowChange)
  window.removeEventListener('scroll', handleWindowChange, true)
})
</script>

<template>
  <div class="locale-toggle">
    <button
      ref="trigger"
      class="locale-toggle__button"
      type="button"
      aria-haspopup="listbox"
      :aria-expanded="isOpen"
      :aria-label="t('locale.label')"
      :title="currentLabel"
      @click="togglePanel"
    >
      <Globe2 :size="18" aria-hidden="true" />
    </button>

    <Teleport to="body">
      <Transition name="locale-toggle-panel">
        <div
          v-if="isOpen"
          ref="panel"
          class="locale-toggle__panel"
          role="listbox"
          tabindex="-1"
          :aria-label="t('locale.label')"
          :style="panelStyle"
        >
          <button
            v-for="option in options"
            :key="option.value"
            class="locale-toggle__option"
            type="button"
            role="option"
            :aria-selected="option.value === locale"
            @click="selectLocale(option.value)"
          >
            <span>{{ option.label }}</span>
            <Check v-if="option.value === locale" :size="15" aria-hidden="true" />
          </button>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.locale-toggle {
  display: inline-flex;
}

.locale-toggle__button {
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

.locale-toggle__button:hover {
  transform: translateY(-1px);
  border-color: var(--bb-color-primary);
}

.locale-toggle__button:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.locale-toggle__panel {
  position: fixed;
  z-index: 2200;
  width: 168px;
  display: grid;
  gap: 4px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 6px;
  background: var(--bb-color-surface-elevated);
  box-shadow: var(--bb-shadow-panel);
}

.locale-toggle__panel:focus {
  outline: none;
}

.locale-toggle__option {
  min-height: 38px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  border: 0;
  border-radius: 6px;
  padding: 0 10px;
  color: var(--bb-color-text);
  background: transparent;
  text-align: left;
}

.locale-toggle__option:hover,
.locale-toggle__option:focus-visible {
  outline: none;
  background: var(--bb-color-subtle);
}

.locale-toggle__option[aria-selected='true'] {
  color: var(--bb-color-primary);
  font-weight: 750;
}

.locale-toggle-panel-enter-active,
.locale-toggle-panel-leave-active {
  transition: opacity 120ms ease, transform 120ms ease;
}

.locale-toggle-panel-enter-from,
.locale-toggle-panel-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
