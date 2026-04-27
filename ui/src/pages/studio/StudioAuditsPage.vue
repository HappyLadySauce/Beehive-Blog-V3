<script setup lang="ts">
import { computed, onMounted, reactive, shallowRef } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import type { StudioAuditEvent } from '@/features/studio'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import FormField from '@/shared/components/FormField.vue'
import DataTable from '@/shared/components/DataTable.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import type { DataTableColumn } from '@/shared/components/DataTable.vue'

const authStore = useAuthStore()
const events = shallowRef<StudioAuditEvent[]>([])
const total = shallowRef(0)
const isLoading = shallowRef(true)
const errorMessage = shallowRef('')
const filters = reactive({
  eventType: '',
  result: '',
  userId: '',
})

const columns: DataTableColumn[] = [
  { key: 'createdAt', label: 'Time' },
  { key: 'eventType', label: 'Event' },
  { key: 'result', label: 'Result' },
  { key: 'userId', label: 'User' },
  { key: 'clientIp', label: 'IP' },
  { key: 'detail', label: 'Detail' },
]

const rows = computed(() =>
  events.value.map((event) => ({
    createdAt: formatUnixTime(event.created_at),
    eventType: event.event_type,
    result: event.result,
    userId: event.user_id || 'System',
    clientIp: event.client_ip || '-',
    detail: event.detail_json || '-',
  })),
)

function formatUnixTime(value?: number): string {
  if (!value) {
    return '-'
  }
  return new Intl.DateTimeFormat('en', {
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
    errorMessage.value = error instanceof Error ? error.message : 'Unable to load audit events.'
  }
  finally {
    isLoading.value = false
  }
}

onMounted(loadAudits)
</script>

<template>
  <section class="audits-page">
    <PageHeader
      eyebrow="Studio"
      title="Audit log"
      description="Inspect sensitive identity and account-management activity."
    />

    <form class="audits-page__filters" @submit.prevent="loadAudits">
      <FormField label="Event type" for-id="audit-event-type">
        <BaseInput id="audit-event-type" v-model="filters.eventType" placeholder="admin_update_user_status" />
      </FormField>
      <FormField label="User ID" for-id="audit-user-id">
        <BaseInput id="audit-user-id" v-model="filters.userId" placeholder="1" />
      </FormField>
      <label class="audits-page__select">
        <span>Result</span>
        <select v-model="filters.result">
          <option value="">All results</option>
          <option value="success">Success</option>
          <option value="failure">Failure</option>
        </select>
      </label>
      <BaseButton type="submit" :busy="isLoading">Apply</BaseButton>
    </form>

    <StatusAlert v-if="errorMessage" tone="danger" title="Audit log unavailable">
      {{ errorMessage }}
    </StatusAlert>
    <StatusAlert v-else-if="isLoading" tone="info" title="Loading audit events">
      Audit events are being loaded from gateway.
    </StatusAlert>

    <DataTable :columns="columns" :rows="rows" empty-text="No audit events found." />
    <p class="audits-page__count">{{ total }} total audit events</p>
  </section>
</template>

<style scoped>
.audits-page {
  display: grid;
  gap: 24px;
}

.audits-page__filters {
  display: grid;
  grid-template-columns: minmax(180px, 1fr) minmax(140px, 180px) minmax(140px, 180px) auto;
  align-items: end;
  gap: 12px;
}

.audits-page__select {
  display: grid;
  gap: 6px;
  color: var(--bb-color-muted);
  font-size: 0.92rem;
  font-weight: 650;
}

select {
  min-height: 44px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 0 10px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface);
}

select:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.audits-page__count {
  color: var(--bb-color-muted);
}

@media (max-width: 760px) {
  .audits-page__filters {
    grid-template-columns: 1fr;
  }
}
</style>
