<script setup lang="ts">
import { ArrowRight, CheckCircle2, Layers3, ShieldCheck, Sparkles } from 'lucide-vue-next'

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

const highlights = [
  {
    icon: Layers3,
    tone: 'neutral',
    label: 'Gateway first',
    title: 'One coherent surface',
    description: 'Public pages and Studio workflows share one Vue app while keeping the gateway as the only HTTP boundary.',
  },
  {
    icon: CheckCircle2,
    tone: 'success',
    label: 'Accessible',
    title: 'Interaction states included',
    description: 'Forms, tables, alerts, empty states, and actions are built with keyboard and feedback states first.',
  },
  {
    icon: ShieldCheck,
    tone: 'warning',
    label: 'Identity ready',
    title: 'Operational by default',
    description: 'Auth restoration and role checks are wired before sensitive Studio routes render.',
  },
] as const
</script>

<template>
  <section class="home-page">
    <section class="home-page__hero" aria-labelledby="home-title">
      <div class="home-page__hero-copy">
        <BaseBadge>Public Web</BaseBadge>
        <h1 id="home-title">Beehive Blog</h1>
        <p>
          A clean publishing and operations surface for technical notes, release context, and admin workflows.
        </p>
        <div class="home-page__actions">
          <BaseButton>
            Explore articles
            <ArrowRight :size="17" aria-hidden="true" />
          </BaseButton>
          <BaseButton variant="secondary">View projects</BaseButton>
        </div>
      </div>
      <aside class="home-page__visual" aria-label="Product preview">
        <div class="home-page__visual-toolbar">
          <span />
          <span />
          <span />
        </div>
        <div class="home-page__visual-card home-page__visual-card--primary">
          <Sparkles :size="20" aria-hidden="true" />
          <strong>Ready for gateway integration</strong>
          <span>Auth, Studio, and shared UI states are aligned.</span>
        </div>
        <div class="home-page__visual-grid">
          <span />
          <span />
          <span />
          <span />
        </div>
      </aside>
    </section>

    <section class="home-page__features" aria-label="Highlights">
      <article v-for="item in highlights" :key="item.title" class="home-page__feature">
        <component :is="item.icon" class="home-page__feature-icon" :size="22" aria-hidden="true" />
        <BaseBadge :tone="item.tone">{{ item.label }}</BaseBadge>
        <h2>{{ item.title }}</h2>
        <p>{{ item.description }}</p>
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

.home-page__hero {
  min-height: 360px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(320px, 440px);
  gap: 24px;
  align-items: stretch;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 24px;
  background: linear-gradient(135deg, var(--bb-color-surface-elevated), var(--bb-color-primary-soft));
  box-shadow: var(--bb-shadow-panel);
}

.home-page__hero-copy {
  display: grid;
  align-content: center;
  justify-items: start;
  gap: 18px;
}

.home-page__hero-copy h1,
.home-page__hero-copy p {
  margin: 0;
}

.home-page__hero-copy h1 {
  color: var(--bb-color-text-strong);
  font-size: clamp(2.4rem, 7vw, 5rem);
  line-height: 0.96;
}

.home-page__hero-copy p {
  max-width: 620px;
  color: var(--bb-color-muted);
  font-size: 1.08rem;
}

.home-page__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.home-page__visual {
  min-height: 310px;
  display: grid;
  align-content: start;
  gap: 18px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 18px;
  background: var(--bb-color-surface-glass);
  box-shadow: var(--bb-shadow-soft);
  backdrop-filter: blur(14px);
}

.home-page__visual-toolbar {
  display: flex;
  gap: 7px;
}

.home-page__visual-toolbar span {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: var(--bb-color-subtle-strong);
}

.home-page__visual-card {
  display: grid;
  gap: 6px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 18px;
  background: var(--bb-color-surface);
}

.home-page__visual-card--primary {
  color: var(--bb-color-text);
  background: linear-gradient(135deg, var(--bb-color-primary), var(--bb-color-accent));
}

.home-page__visual-card--primary,
.home-page__visual-card--primary span {
  color: #fff;
}

.home-page__visual-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.home-page__visual-grid span {
  min-height: 72px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  background: linear-gradient(135deg, var(--bb-color-surface), var(--bb-color-subtle));
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
  box-shadow: var(--bb-shadow-soft);
  transition: transform 160ms ease, border-color 160ms ease, box-shadow 160ms ease;
}

.home-page__feature:hover {
  transform: translateY(-2px);
  border-color: var(--bb-color-primary);
  box-shadow: var(--bb-shadow-panel);
}

.home-page__feature-icon {
  color: var(--bb-color-primary);
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
  .home-page__hero {
    grid-template-columns: 1fr;
  }

  .home-page__features {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 560px) {
  .home-page__hero {
    padding: 18px;
  }

  .home-page__visual {
    min-height: 240px;
  }
}
</style>
