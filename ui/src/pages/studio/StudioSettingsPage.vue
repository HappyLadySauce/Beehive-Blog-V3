<script setup lang="ts">
import { Settings2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import BaseButton from '@/shared/components/BaseButton.vue'
import EmptyState from '@/shared/components/EmptyState.vue'

const { t } = useI18n()
const router = useRouter()

async function openFileTypesWorkspace(): Promise<void> {
  await router.push({
    path: '/studio/files',
    query: { tab: 'types' },
  })
}
</script>

<template>
  <section class="settings-page">
    <EmptyState
      class="settings-page__empty"
      align="center"
      :title="t('settings.fileMovedTitle')"
      :description="t('settings.fileMovedMessage')"
    >
      <template #visual>
        <Settings2 :size="52" aria-hidden="true" />
      </template>
      <BaseButton @click="openFileTypesWorkspace">{{ t('settings.openFileTypesAction') }}</BaseButton>
    </EmptyState>
  </section>
</template>

<style scoped>
.settings-page {
  display: grid;
  gap: 16px;
}

.settings-page__empty {
  min-height: 240px;
}
</style>
