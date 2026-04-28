<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useLocale, type AppLocale } from '@/shared/i18n'

import BaseSelect, { type BaseSelectOption } from './BaseSelect.vue'

const { t } = useI18n()
const { locale, setLocale } = useLocale()

const options = computed<BaseSelectOption[]>(() => [
  { value: 'zh-CN', label: t('locale.zhCN') },
  { value: 'en-US', label: t('locale.enUS') },
])

const selectedLocale = computed({
  get: () => locale.value,
  set: (value: string) => setLocale(value as AppLocale),
})
</script>

<template>
  <BaseSelect
    v-model="selectedLocale"
    class="locale-toggle"
    :options="options"
    :aria-label="t('locale.label')"
  />
</template>

<style scoped>
.locale-toggle {
  width: 112px;
  min-height: 42px;
  background: var(--bb-color-surface-elevated);
  box-shadow: var(--bb-shadow-soft);
}

@media (max-width: 640px) {
  .locale-toggle {
    width: 96px;
  }
}
</style>
