<script setup lang="ts" generic="TRow extends Record<string, unknown>">
export interface TableColumn<TRow> {
  key: string;
  label: string;
  align?: 'left' | 'right' | 'center';
  width?: string;
  render?: (row: TRow) => string;
}

defineProps<{
  columns: TableColumn<TRow>[];
  rows: TRow[];
  rowKey: keyof TRow;
  emptyText?: string;
  inverse?: boolean;
}>();
</script>

<template>
  <div class="overflow-hidden rounded-lg border" :class="inverse ? 'border-white/8 bg-white/5' : 'border-brand-line bg-brand-surface'">
    <div class="hidden overflow-x-auto md:block">
      <table class="w-full min-w-640px border-collapse text-left text-14px">
        <thead :class="inverse ? 'border-white/8 text-white/45' : 'border-brand-line text-brand-muted'" class="border-b">
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              class="px-4 py-3 text-12px font-800"
              :style="{ width: column.width }"
              :class="{
                'text-right': column.align === 'right',
                'text-center': column.align === 'center',
              }"
            >
              {{ column.label }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="rows.length === 0">
            <td :colspan="columns.length" class="px-4 py-8 text-center text-13px" :class="inverse ? 'text-white/45' : 'text-brand-muted'">
              {{ emptyText ?? '暂无数据' }}
            </td>
          </tr>
          <tr
            v-for="row in rows"
            :key="String(row[rowKey])"
            class="border-b last:border-b-0"
            :class="inverse ? 'border-white/6 text-white' : 'border-brand-line text-brand-ink'"
          >
            <td
              v-for="column in columns"
              :key="column.key"
              class="px-4 py-3 align-top"
              :class="{
                'text-right': column.align === 'right',
                'text-center': column.align === 'center',
              }"
            >
              <slot :name="`cell-${column.key}`" :row="row" :value="row[column.key]">
                {{ column.render ? column.render(row) : row[column.key] }}
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="grid md:hidden">
      <p v-if="rows.length === 0" class="m-0 px-4 py-8 text-center text-13px" :class="inverse ? 'text-white/45' : 'text-brand-muted'">
        {{ emptyText ?? '暂无数据' }}
      </p>
      <article
        v-for="row in rows"
        :key="String(row[rowKey])"
        class="grid gap-3 border-b p-4 last:border-b-0"
        :class="inverse ? 'border-white/6 text-white' : 'border-brand-line text-brand-ink'"
      >
        <div v-for="column in columns" :key="column.key" class="grid gap-1">
          <span class="text-11px font-800 uppercase" :class="inverse ? 'text-white/40' : 'text-brand-muted'">{{ column.label }}</span>
          <div class="min-w-0 text-14px">
            <slot :name="`cell-${column.key}`" :row="row" :value="row[column.key]">
              {{ column.render ? column.render(row) : row[column.key] }}
            </slot>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>
