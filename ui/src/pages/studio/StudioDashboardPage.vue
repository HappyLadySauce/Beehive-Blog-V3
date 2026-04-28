<script setup lang="ts">
import BaseBadge from '@/shared/components/BaseBadge.vue'
import EmptyState from '@/shared/components/EmptyState.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
</script>

<template>
  <section class="studio-page">
    <PageHeader
      eyebrow="Studio"
      title="Dashboard"
      description="Operational overview for publishing work and administrative activity."
    />

    <EmptyState
      title="Dashboard metrics are not configured yet"
      description="Live operational metrics will appear here after analytics and publishing signals are connected."
    />
    <BaseBadge tone="success">Admin only</BaseBadge>
  </section>
</template>

<style scoped>
.studio-page {
  display: grid;
  gap: 24px;
}

</style>
