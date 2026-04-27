<script setup lang="ts">
import BaseButton from '@/shared/components/BaseButton.vue'
import DataTable from '@/shared/components/DataTable.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import type { DataTableColumn } from '@/shared/components/DataTable.vue'

const columns: DataTableColumn[] = [
  { key: 'title', label: 'Title' },
  { key: 'type', label: 'Type' },
  { key: 'status', label: 'Status' },
  { key: 'updatedAt', label: 'Updated' },
]

const rows = [
  { title: 'v3 frontend integration notes', type: 'Article', status: 'Draft', updatedAt: '2026-04-26' },
  { title: 'Auth refresh behavior', type: 'Article', status: 'Review', updatedAt: '2026-04-25' },
]
</script>

<template>
  <section class="content-page">
    <PageHeader
      eyebrow="Studio"
      title="Content"
      description="Content management placeholder while the content service remains under development."
    >
      <template #actions>
        <BaseButton>New draft</BaseButton>
      </template>
    </PageHeader>
    <DataTable :columns="columns" :rows="rows" empty-text="No content yet." />
  </section>
</template>

<style scoped>
.content-page {
  display: grid;
  gap: 24px;
}
</style>
