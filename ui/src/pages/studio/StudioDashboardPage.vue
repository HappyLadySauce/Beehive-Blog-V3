<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import BaseBadge from '@/shared/components/BaseBadge.vue'
import EmptyState from '@/shared/components/EmptyState.vue'
import PageHeader from '@/shared/components/PageHeader.vue'

const { t } = useI18n()
</script>

<template>
  <section class="studio-page">
    <PageHeader
      :eyebrow="t('dashboard.eyebrow')"
      :title="t('dashboard.title')"
      :description="t('dashboard.description')"
    />

    <EmptyState
      :title="t('dashboard.emptyTitle')"
      :description="t('dashboard.emptyDescription')"
    />
    <BaseBadge tone="success">{{ t('dashboard.adminOnly') }}</BaseBadge>
  </section>
</template>

<style scoped>
.studio-page {
  display: grid;
  gap: 24px;
}

</style>
