import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';

import BaseCheckbox from '@/shared/components/BaseCheckbox.vue';
import BaseSelect from '@/shared/components/BaseSelect.vue';
import BaseSwitch from '@/shared/components/BaseSwitch.vue';
import BaseTextarea from '@/shared/components/BaseTextarea.vue';
import DataTable, { type TableColumn } from '@/shared/components/DataTable.vue';
import LoadingSkeleton from '@/shared/components/LoadingSkeleton.vue';
import StatusAlert from '@/shared/components/StatusAlert.vue';

describe('shared UI components', () => {
  it('emits textarea changes', async () => {
    const wrapper = mount(BaseTextarea, {
      props: {
        modelValue: '',
        label: '摘要',
      },
    });

    await wrapper.get('textarea').setValue('新的摘要');

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['新的摘要']);
  });

  it('emits select changes', async () => {
    const wrapper = mount(BaseSelect, {
      props: {
        modelValue: 'draft',
        label: '状态',
        options: [
          { label: '草稿', value: 'draft' },
          { label: '发布', value: 'published' },
        ],
      },
    });

    await wrapper.get('select').setValue('published');

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['published']);
  });

  it('emits checkbox and switch toggles', async () => {
    const checkbox = mount(BaseCheckbox, {
      props: {
        modelValue: false,
        label: '允许评论',
      },
    });
    const switcher = mount(BaseSwitch, {
      props: {
        modelValue: false,
        label: '公开展示',
      },
    });

    await checkbox.get('input').setValue(true);
    await switcher.get('button').trigger('click');

    expect(checkbox.emitted('update:modelValue')?.[0]).toEqual([true]);
    expect(switcher.emitted('update:modelValue')?.[0]).toEqual([true]);
  });

  it('renders table rows and custom cells', () => {
    type Row = Record<string, unknown>;
    const columns: TableColumn<Row>[] = [
      { key: 'title', label: '标题' },
      { key: 'status', label: '状态' },
    ];
    const wrapper = mount(DataTable, {
      props: {
        columns,
        rows: [{ id: '1', title: '测试内容', status: 'published' }],
        rowKey: 'id',
      },
      slots: {
        'cell-status': '<span class="status">已发布</span>',
      },
    });

    expect(wrapper.text()).toContain('测试内容');
    expect(wrapper.find('.status').text()).toBe('已发布');
  });

  it('renders status alert and loading skeleton', () => {
    const alert = mount(StatusAlert, {
      props: {
        tone: 'success',
        title: '保存成功',
        description: '内容已写入本地状态。',
      },
    });
    const skeleton = mount(LoadingSkeleton, {
      props: {
        rows: 2,
      },
    });

    expect(alert.text()).toContain('保存成功');
    expect(skeleton.findAll('span')).toHaveLength(2);
  });
});
