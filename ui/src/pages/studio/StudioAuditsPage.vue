<script setup lang="ts">
import { computed, onBeforeUnmount, reactive, shallowRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import type { StudioAuditEvent } from '@/features/studio'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSelect, { type BaseSelectOption } from '@/shared/components/BaseSelect.vue'
import DataTable from '@/shared/components/DataTable.vue'
import FormField from '@/shared/components/FormField.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import TablePagination from '@/shared/components/TablePagination.vue'
import type { DataTableColumn } from '@/shared/components/DataTable.vue'
import { usePaginatedRouteState, useProgressiveQuery } from '@/shared/composables'
import { useLocale } from '@/shared/i18n'

const authStore = useAuthStore()
const { t } = useI18n()
const { locale } = useLocale()
const route = useRoute()
const router = useRouter()
let filterTimer: number | undefined
const filters = reactive({
  eventType: readQueryString(route.query.eventType),
  result: readQueryString(route.query.result),
  userId: readQueryString(route.query.userId),
})
const appliedFilters = reactive({
  eventType: readQueryString(route.query.eventType),
  result: readQueryString(route.query.result),
  userId: readQueryString(route.query.userId),
})
const total = shallowRef(0)

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
const pagination = usePaginatedRouteState({
  route,
  router,
  total,
})

const auditsQuery = useProgressiveQuery({
  queryKey: computed(() => [
    'studio-audits',
    { ...appliedFilters },
    pagination.page.value,
    pagination.pageSize.value,
    authStore.currentUser?.user_id ?? 'anonymous',
  ]),
  queryFn: () => studioApi.listAudits(
    {
      event_type: appliedFilters.eventType.trim(),
      result: appliedFilters.result,
      user_id: appliedFilters.userId.trim(),
      page: pagination.page.value,
      page_size: pagination.pageSize.value,
    },
    { accessToken: authStore.accessToken },
  ),
})

const events = computed<StudioAuditEvent[]>(() => auditsQuery.data.value?.items ?? [])
const errorMessage = computed(() => {
  const error = auditsQuery.error.value
  return error instanceof Error ? error.message : ''
})

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

function scheduleLoadAudits(): void {
  window.clearTimeout(filterTimer)
  filterTimer = window.setTimeout(() => {
    appliedFilters.eventType = filters.eventType
    appliedFilters.userId = filters.userId
    appliedFilters.result = filters.result
    pagination.page.value = 1
    void pagination.syncQuery(buildAuditQuery())
  }, 300)
}

watch(() => auditsQuery.data.value?.total ?? 0, (value) => {
  total.value = value
}, { immediate: true })
watch(
  () => [route.query.eventType, route.query.userId, route.query.result],
  () => {
    filters.eventType = readQueryString(route.query.eventType)
    filters.userId = readQueryString(route.query.userId)
    filters.result = readQueryString(route.query.result)
    appliedFilters.eventType = filters.eventType
    appliedFilters.userId = filters.userId
    appliedFilters.result = filters.result
  },
)
watch(() => [filters.eventType, filters.userId, filters.result], scheduleLoadAudits)
onBeforeUnmount(() => window.clearTimeout(filterTimer))

function buildAuditQuery() {
  return {
    eventType: filters.eventType.trim() || undefined,
    userId: filters.userId.trim() || undefined,
    result: filters.result || undefined,
  }
}

function readQueryString(value: unknown): string {
  const normalized = Array.isArray(value) ? value[0] : value
  return typeof normalized === 'string' ? normalized : ''
}
</script>

<template>
  <section class="audits-page">
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

    <DataTable
      v-if="auditsQuery.hasResolvedOnce.value || auditsQuery.showBlockingLoading.value"
      :columns="columns"
      :rows="rows"
      :empty-text="t('audits.empty')"
      :loading="auditsQuery.showBlockingLoading.value || auditsQuery.showRefreshingHint.value"
      :loading-mode="auditsQuery.showBlockingLoading.value ? 'blocking' : 'refreshing'"
      :loading-title="auditsQuery.showBlockingLoading.value ? t('audits.loadingTitle') : t('common.refreshing')"
    />
    <div v-if="auditsQuery.data.value" class="studio-list-footer">
      <TablePagination
        :page="pagination.page.value"
        :page-size="pagination.pageSize.value"
        :total="total"
        :disabled="auditsQuery.isFetching.value"
        @update:page="pagination.setPage"
        @update:page-size="pagination.setPageSize"
      />
    </div>
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
