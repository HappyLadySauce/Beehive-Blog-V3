<script setup lang="ts">
import { computed, shallowRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from './BaseButton.vue'
import BaseSelect, { type BaseSelectOption } from './BaseSelect.vue'

const props = withDefaults(defineProps<{
  page: number
  pageSize: number
  total: number
  pageSizeOptions?: number[]
  disabled?: boolean
}>(), {
  pageSizeOptions: () => [10, 20, 50, 100],
  disabled: false,
})

const emit = defineEmits<{
  'update:page': [value: number]
  'update:pageSize': [value: number]
}>()

const { t } = useI18n()
const pageSizeValue = shallowRef(String(props.pageSize))

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))
const currentPage = computed(() => Math.min(Math.max(1, props.page), totalPages.value))
const visiblePages = computed(() => {
  const maxVisible = 5
  const pages: number[] = []
  const halfWindow = Math.floor(maxVisible / 2)
  let start = Math.max(1, currentPage.value - halfWindow)
  let end = Math.min(totalPages.value, start + maxVisible - 1)
  start = Math.max(1, end - maxVisible + 1)
  for (let value = start; value <= end; value += 1) {
    pages.push(value)
  }
  return pages
})
const rangeStart = computed(() => (props.total === 0 ? 0 : (currentPage.value - 1) * props.pageSize + 1))
const rangeEnd = computed(() => Math.min(props.total, currentPage.value * props.pageSize))
const pageSizeSelectOptions = computed<BaseSelectOption[]>(() =>
  props.pageSizeOptions.map((option) => ({
    value: String(option),
    label: String(option),
  })),
)

watch(() => props.pageSize, (value) => {
  pageSizeValue.value = String(value)
}, { immediate: true })

function updatePage(page: number): void {
  const nextPage = Math.min(Math.max(1, page), totalPages.value)
  if (nextPage !== props.page) {
    emit('update:page', nextPage)
  }
}

function updatePageSize(value: string): void {
  pageSizeValue.value = value
  const parsed = Number(value)
  if (!Number.isFinite(parsed) || parsed === props.pageSize) {
    return
  }
  emit('update:pageSize', parsed)
}
</script>

<template>
  <div class="table-pagination" :aria-label="t('common.page')">
    <p class="table-pagination__range">{{ t('common.pageRange', { from: rangeStart, to: rangeEnd, total }) }}</p>

    <div class="table-pagination__controls">
      <label class="table-pagination__page-size">
        <span>{{ t('common.pageSize') }}</span>
        <BaseSelect
          :model-value="pageSizeValue"
          :options="pageSizeSelectOptions"
          :disabled="disabled"
          :aria-label="t('common.pageSize')"
          @update:model-value="updatePageSize"
        />
      </label>

      <div class="table-pagination__nav">
        <BaseButton
          variant="ghost"
          :disabled="disabled || currentPage <= 1"
          @click="updatePage(currentPage - 1)"
        >
          {{ t('common.previousPage') }}
        </BaseButton>

        <div class="table-pagination__pages">
          <BaseButton
            v-for="value in visiblePages"
            :key="value"
            :variant="value === currentPage ? 'primary' : 'ghost'"
            :disabled="disabled"
            @click="updatePage(value)"
          >
            {{ value }}
          </BaseButton>
        </div>

        <BaseButton
          variant="ghost"
          :disabled="disabled || currentPage >= totalPages"
          @click="updatePage(currentPage + 1)"
        >
          {{ t('common.nextPage') }}
        </BaseButton>
      </div>
    </div>
  </div>
</template>

<style scoped>
.table-pagination {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  color: var(--bb-color-muted);
}

.table-pagination__range {
  margin: 0;
}

.table-pagination__controls {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-end;
  gap: 16px;
}

.table-pagination__page-size {
  min-width: 168px;
  display: grid;
  gap: 6px;
  font-size: 0.88rem;
  font-weight: 700;
}

.table-pagination__nav,
.table-pagination__pages {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.table-pagination__pages :deep(.bb-button) {
  min-width: 44px;
  padding-inline: 12px;
}

@media (max-width: 760px) {
  .table-pagination,
  .table-pagination__controls {
    align-items: stretch;
  }

  .table-pagination__controls,
  .table-pagination__nav {
    width: 100%;
    justify-content: space-between;
  }

  .table-pagination__page-size {
    width: 100%;
  }
}
</style>
