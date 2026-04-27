<script setup lang="ts">
export interface DataTableColumn {
  key: string
  label: string
}

defineProps<{
  columns: DataTableColumn[]
  rows: Record<string, string | number | boolean | null | undefined>[]
  emptyText?: string
}>()
</script>

<template>
  <div class="data-table" role="region" aria-label="Data table" tabindex="0">
    <table>
      <thead>
        <tr>
          <th v-for="column in columns" :key="column.key" scope="col">{{ column.label }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="rows.length === 0">
          <td :colspan="columns.length">{{ emptyText ?? 'No records yet.' }}</td>
        </tr>
        <tr v-for="(row, rowIndex) in rows" v-else :key="rowIndex">
          <td v-for="column in columns" :key="column.key">{{ row[column.key] }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.data-table {
  overflow-x: auto;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  background: var(--bb-color-surface);
}

.data-table:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

table {
  width: 100%;
  min-width: 560px;
  border-collapse: collapse;
}

th,
td {
  border-bottom: 1px solid var(--bb-color-line);
  padding: 12px;
  text-align: left;
}

th {
  color: var(--bb-color-muted);
  font-size: 0.8rem;
  text-transform: uppercase;
}

tr:last-child td {
  border-bottom: 0;
}
</style>
