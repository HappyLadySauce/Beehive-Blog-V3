<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, shallowRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import type { StudioAuditEvent } from '@/features/studio'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSelect, { type BaseSelectOption } from '@/shared/components/BaseSelect.vue'
import FormField from '@/shared/components/FormField.vue'
import DataTable from '@/shared/components/DataTable.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import type { DataTableColumn } from '@/shared/components/DataTable.vue'
import { useLocale } from '@/shared/i18n'

const authStore = useAuthStore()
const { t } = useI18n()
const { locale } = useLocale()
const events = shallowRef<StudioAuditEvent[]>([])
const total = shallowRef(0)
const isLoading = shallowRef(true)
const errorMessage = shallowRef('')
let filterTimer: number | undefined
const filters = reactive({
  eventType: '',
  result: '',
  userId: '',
})

const resultOptions = computed<BaseSelectOption[]>(() => [
  { value: '', label: t('audits.allResults') },
  { value: 'success', label: t('audits.success') },
  { value: 'failure', label: t('audits.failure') },
])
const columns = computed<DataTableColumn[]>(() => [
  { key: 'createdAt', label: t('audits.columns.time') },
  { key: 'eventType', label: t('audits.columns.event') },
  { key: 'result', label: t('audits.columns.result') },
  { key: 'userId', label: t('audits.columns.user') },
  { key: 'clientIp', label: t('audits.columns.ip') },
  { key: 'detail', label: t('audits.columns.detail') },
])

const rows = computed(() =>
  events.value.map((event) => ({
    createdAt: formatUnixTime(event.created_at),
    eventType: event.event_type,
    result: event.result,
    userId: event.user_id || t('common.system'),
    clientIp: event.client_ip || t('common.none'),
    detail: event.detail_json || t('common.none'),
  })),
)

function formatUnixTime(value?: number): string {
  if (!value) {
    return t('common.none')
  }
  return new Intl.DateTimeFormat(locale.value, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value * 1000))
}

async function loadAudits(): Promise<void> {
  isLoading.value = true
  errorMessage.value = ''
  try {
    const response = await studioApi.listAudits(
      {
        event_type: filters.eventType.trim(),
        result: filters.result,
        user_id: filters.userId.trim(),
        page: 1,
        page_size: 50,
      },
      { accessToken: authStore.accessToken },
    )
    events.value = response.items
    total.value = response.total
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : t('audits.unavailableTitle')
  }
  finally {
    isLoading.value = false
  }
}

function scheduleLoadAudits(): void {
  window.clearTimeout(filterTimer)
  filterTimer = window.setTimeout(() => {
    void loadAudits()
  }, 300)
}

watch(() => [filters.eventType, filters.userId, filters.result], scheduleLoadAudits)

onMounted(loadAudits)
onBeforeUnmount(() => window.clearTimeout(filterTimer))
</script>

<template>
  <section class="audits-page">
    <PageHeader
      :eyebrow="t('audits.eyebrow')"
      :title="t('audits.title')"
      :description="t('audits.description')"
    />

    <div class="studio-list-filters audits-page__filters">
      <FormField :label="t('audits.eventType')" for-id="audit-event-type">
        <BaseInput id="audit-event-type" v-model="filters.eventType" :placeholder="t('audits.placeholders.eventType')" />
      </FormField>
      <FormField :label="t('audits.userId')" for-id="audit-user-id">
        <BaseInput id="audit-user-id" v-model="filters.userId" :placeholder="t('audits.placeholders.userId')" />
      </FormField>
      <FormField :label="t('audits.result')" for-id="audit-result">
        <BaseSelect id="audit-result" v-model="filters.result" :options="resultOptions" :aria-label="t('audits.result')" />
      </FormField>
    </div>

    <StatusAlert v-if="errorMessage" tone="danger" :title="t('audits.unavailableTitle')">
      {{ errorMessage }}
    </StatusAlert>
    <StatusAlert v-else-if="isLoading" tone="info" :title="t('audits.loadingTitle')">
      {{ t('audits.loadingMessage') }}
    </StatusAlert>

    <DataTable :columns="columns" :rows="rows" :empty-text="t('audits.empty')" />
    <p class="studio-list-count">{{ t('audits.count', { count: total }) }}</p>
  </section>
</template>

<style scoped>
.audits-page {
  display: grid;
  gap: 24px;
}

.audits-page__filters {
  grid-template-columns: minmax(220px, 1fr) repeat(2, minmax(160px, 180px));
}

@media (max-width: 760px) {
  .audits-page__filters {
    grid-template-columns: 1fr;
  }
}
</style>
