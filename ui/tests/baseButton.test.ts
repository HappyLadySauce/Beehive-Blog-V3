import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';

import BaseButton from '@/shared/components/BaseButton.vue';

describe('BaseButton', () => {
  it('renders slot content and disabled state', () => {
    const wrapper = mount(BaseButton, {
      props: { busy: true },
      slots: { default: '保存' },
    });

    expect(wrapper.text()).toContain('保存');
    expect(wrapper.attributes('disabled')).toBeDefined();
  });
});
