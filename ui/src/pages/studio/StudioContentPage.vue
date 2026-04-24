<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { contentPreviewApi } from '@/shared/api/contentPreviewApi';
import type { ContentSummaryView } from '@/shared/api/types';
import BaseBadge from '@/shared/components/BaseBadge.vue';
import BaseButton from '@/shared/components/BaseButton.vue';
import DataTable, { type TableColumn } from '@/shared/components/DataTable.vue';
import PageHeader from '@/shared/components/PageHeader.vue';
import StatusAlert from '@/shared/components/StatusAlert.vue';

const items = ref<ContentSummaryView[]>([]);
const columns: TableColumn<ContentSummaryView>[] = [
  { key: 'title', label: '标题' },
  { key: 'type', label: '类型', width: '120px' },
  { key: 'status', label: '状态', width: '120px' },
];

onMounted(async () => {
  items.value = (await contentPreviewApi.listPublicContent()).items;
});
</script>

<template>
  <main class="grid gap-5 p-4 sm:p-6 lg:p-8">
    <PageHeader
      inverse
      eyebrow="Content Mock"
      title="内容中心"
      description="展示管理台信息密度和表格布局，不接 content 真实服务。"
    >
      <template #actions>
        <BaseButton variant="secondary">新建内容</BaseButton>
      </template>
    </PageHeader>

    <StatusAlert
      tone="info"
      title="当前使用 mock content adapter"
      description="content 服务开发完成前，Studio 内容中心只验证 UI 结构、筛选和状态展示。"
    />

    <DataTable :columns="columns" :rows="items" row-key="content_id" inverse>
      <template #cell-title="{ row }">
        <div class="min-w-0 max-w-560px">
          <p class="m-0 truncate text-14px font-800">{{ row.title }}</p>
          <p class="m-0 mt-1 line-clamp-2 text-13px leading-5 text-white/48">{{ row.summary }}</p>
        </div>
      </template>
      <template #cell-type="{ row }">
        <BaseBadge tone="neutral">{{ row.type }}</BaseBadge>
      </template>
      <template #cell-status="{ row }">
        <BaseBadge :tone="row.status === 'published' ? 'leaf' : 'honey'">{{ row.status }}</BaseBadge>
      </template>
    </DataTable>
  </main>
</template>
