<script setup lang="ts">
import { PackageOpen } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import EmptyState from './EmptyState.vue'
import InlineLoadingState from './InlineLoadingState.vue'
import PageLoadingState from './PageLoadingState.vue'

export interface DataTableColumn {
  key: string
  label: string
  width?: string
}

defineProps<{
  columns: DataTableColumn[]
  rows: Record<string, string | number | boolean | null | undefined>[]
  emptyText?: string
  emptyDescription?: string
  loading?: boolean
  loadingMode?: 'blocking' | 'refreshing'
  loadingTitle?: string
  loadingRows?: number
}>()

const { t } = useI18n()
</script>

<template>
  <div class="data-table" role="region" :aria-label="t('accessibility.dataTable')" tabindex="0">
    <PageLoadingState
      v-if="loading && (loadingMode ?? 'blocking') === 'blocking'"
      class="data-table__loading"
      :title="loadingTitle ?? t('common.loading')"
      :rows="loadingRows ?? 4"
    />
    <template v-else>
      <div v-if="loading && loadingMode === 'refreshing' && rows.length > 0" class="data-table__refreshing">
        <InlineLoadingState :title="loadingTitle ?? t('common.refreshing')" />
      </div>
      <table class="data-table__grid">
        <colgroup>
          <col v-for="column in columns" :key="column.key" :style="{ width: column.width }">
        </colgroup>
        <thead>
          <tr>
            <th v-for="column in columns" :key="column.key" scope="col">{{ column.label }}</th>
          </tr>
        </thead>
        <tbody v-if="rows.length > 0">
          <tr v-for="(row, rowIndex) in rows" :key="rowIndex">
            <td v-for="column in columns" :key="column.key">{{ row[column.key] }}</td>
          </tr>
        </tbody>
      </table>
      <div v-if="rows.length === 0" class="data-table__empty-panel">
        <EmptyState
          class="data-table__empty-state"
          align="center"
          :title="emptyText ?? t('accessibility.dataTableEmpty')"
          :description="emptyDescription"
        >
          <template #visual>
            <slot name="emptyVisual">
              <PackageOpen :size="44" aria-hidden="true" />
            </slot>
          </template>
          <slot name="emptyActions" />
        </EmptyState>
      </div>
    </template>
  </div>
</template>

<style scoped>
.data-table {
  overflow-x: auto;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-soft);
}

.data-table:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.data-table__loading {
  border: 0;
  border-radius: 0;
  box-shadow: none;
}

.data-table__refreshing {
  display: flex;
  justify-content: flex-end;
  padding: 10px 12px 0;
}

.data-table__grid {
  width: 100%;
  min-width: 560px;
  table-layout: fixed;
  border-collapse: collapse;
}

th,
td {
  border-bottom: 1px solid var(--bb-color-line);
  padding: 10px 12px;
  text-align: left;
}

th {
  color: var(--bb-color-muted);
  font-size: 0.8rem;
  text-transform: uppercase;
  background: var(--bb-color-subtle);
}

td {
  color: var(--bb-color-text);
  overflow-wrap: break-word;
}

.data-table tbody tr {
  transition: background-color 140ms ease;
}

.data-table tbody tr:hover {
  background: var(--bb-color-subtle);
}

.data-table__empty-panel {
  min-width: 560px;
  border-top: 1px solid var(--bb-color-line);
}

.data-table__empty-state {
  min-height: 180px;
  justify-items: center;
  text-align: center;
  border: 0;
  border-radius: 0;
  padding: 24px 16px;
  background: transparent;
  box-shadow: none;
}

.data-table__empty-state :deep(.empty-state__actions) {
  justify-content: center;
}

.data-table tbody tr:last-child td {
  border-bottom: 0;
}
</style>
