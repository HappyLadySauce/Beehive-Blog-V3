<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { appConfig } from '@/shared/config/env'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import FormField from '@/shared/components/FormField.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'

const { t } = useI18n()
</script>

<template>
  <section class="settings-page">
    <StatusAlert v-if="appConfig.apiMode === 'mock'" tone="info" :title="t('settings.mockModeTitle')">
      {{ t('settings.mockModeMessage') }}
    </StatusAlert>
    <form class="settings-page__form" :aria-label="t('settings.formLabel')">
      <FormField :label="t('settings.siteTitle')" for-id="site-title">
        <BaseInput id="site-title" model-value="Beehive Blog" />
      </FormField>
      <FormField :label="t('settings.defaultAuthor')" for-id="default-author">
        <BaseInput id="default-author" model-value="Admin Editor" />
      </FormField>
      <BaseButton type="button">{{ t('settings.save') }}</BaseButton>
    </form>
  </section>
</template>

<style scoped>
.settings-page {
  display: grid;
  gap: 16px;
}

.settings-page__form {
  display: grid;
  max-width: 560px;
  gap: 16px;
}
</style>
