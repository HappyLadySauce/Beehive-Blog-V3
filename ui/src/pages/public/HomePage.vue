<script setup lang="ts">
import BaseBadge from '@/shared/components/BaseBadge.vue'
import BaseButton from '@/shared/components/BaseButton.vue'
import DataTable from '@/shared/components/DataTable.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import type { DataTableColumn } from '@/shared/components/DataTable.vue'

const columns: DataTableColumn[] = [
  { key: 'title', label: 'Title' },
  { key: 'status', label: 'Status' },
  { key: 'updatedAt', label: 'Updated' },
]

const rows = [
  { title: 'Designing service boundaries for v3', status: 'Draft', updatedAt: '2026-04-25' },
  { title: 'Gateway-first frontend integration', status: 'Review', updatedAt: '2026-04-24' },
  { title: 'Operational notes for identity flows', status: 'Ready', updatedAt: '2026-04-22' },
]
</script>

<template>
  <section class="home-page">
    <PageHeader
      eyebrow="Public Web"
      title="Beehive Blog"
      description="A production-minded publishing surface for technical notes, project journals, and release context."
    >
      <template #actions>
        <BaseButton>Explore articles</BaseButton>
        <BaseButton variant="secondary">View projects</BaseButton>
      </template>
    </PageHeader>

    <section class="home-page__features" aria-label="Highlights">
      <article class="home-page__feature">
        <BaseBadge>Gateway first</BaseBadge>
        <h2>One frontend entry</h2>
        <p>Public pages and Studio workflows share a single Vue app and gateway boundary.</p>
      </article>
      <article class="home-page__feature">
        <BaseBadge tone="success">Accessible</BaseBadge>
        <h2>Reusable components</h2>
        <p>Forms, tables, alerts, empty states, and actions use clear contracts and keyboard-friendly states.</p>
      </article>
      <article class="home-page__feature">
        <BaseBadge tone="warning">Mock ready</BaseBadge>
        <h2>Prepared for auth</h2>
        <p>Studio routes restore the session before checking the required admin role.</p>
      </article>
    </section>

    <section class="home-page__section" aria-labelledby="featured-heading">
      <div>
        <h2 id="featured-heading">Featured drafts</h2>
        <p>Representative content shown with the shared table primitive.</p>
      </div>
      <DataTable :columns="columns" :rows="rows" />
    </section>
  </section>
</template>

<style scoped>
.home-page,
.home-page__section {
  display: grid;
  gap: 24px;
}

.home-page__features {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.home-page__feature {
  display: grid;
  gap: 10px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 18px;
  background: var(--bb-color-surface);
}

.home-page__feature h2,
.home-page__feature p,
.home-page__section h2,
.home-page__section p {
  margin: 0;
}

.home-page__feature p,
.home-page__section p {
  color: var(--bb-color-muted);
}

@media (max-width: 860px) {
  .home-page__features {
    grid-template-columns: 1fr;
  }
}
</style>
