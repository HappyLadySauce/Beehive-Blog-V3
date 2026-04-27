<script setup lang="ts">
import BaseBadge from '@/shared/components/BaseBadge.vue'
import DataTable from '@/shared/components/DataTable.vue'
import EmptyState from '@/shared/components/EmptyState.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import type { DataTableColumn } from '@/shared/components/DataTable.vue'

const columns: DataTableColumn[] = [
  { key: 'title', label: 'Item' },
  { key: 'owner', label: 'Owner' },
  { key: 'status', label: 'Status' },
]

const rows = [
  { title: 'Review gateway release notes', owner: 'Admin', status: 'Review' },
  { title: 'Prepare identity flow article', owner: 'Editor', status: 'Draft' },
]
</script>

<template>
  <section class="studio-page">
    <PageHeader
      eyebrow="Studio"
      title="Dashboard"
      description="Operational overview for publishing work, backed by local mock data until live content is ready."
    />

    <section class="studio-page__metrics" aria-label="Dashboard metrics">
      <article><span>Drafts</span><strong>12</strong></article>
      <article><span>Ready</span><strong>4</strong></article>
      <article><span>Blocked</span><strong>1</strong></article>
    </section>

    <DataTable :columns="columns" :rows="rows" />
    <EmptyState title="Publishing automation is pending" description="Content APIs can be wired here when the service is ready." />
    <BaseBadge tone="success">Admin only</BaseBadge>
  </section>
</template>

<style scoped>
.studio-page {
  display: grid;
  gap: 24px;
}

.studio-page__metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.studio-page__metrics article {
  display: grid;
  gap: 4px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 16px;
  background: var(--bb-color-surface);
}

.studio-page__metrics span {
  color: var(--bb-color-muted);
  font-weight: 700;
}

.studio-page__metrics strong {
  font-size: 2rem;
}

@media (max-width: 760px) {
  .studio-page__metrics {
    grid-template-columns: 1fr;
  }
}
</style>
