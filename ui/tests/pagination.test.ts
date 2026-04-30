import { mount } from '@vue/test-utils'
import { defineComponent, h, nextTick } from 'vue'
import { createMemoryHistory, createRouter, useRoute, useRouter } from 'vue-router'
import { beforeEach, describe, expect, it } from 'vitest'

import { createMockStudioApi } from '@/features/studio/api/mockStudioApi'
import { usePaginatedRouteState } from '@/shared/composables'

async function mountPaginationState(initialPath = '/?page=2&pageSize=50&keep=1', total = 120) {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/', component: { template: '<div />' } }],
  })
  router.push(initialPath)
  await router.isReady()

  const TestComponent = defineComponent({
    setup() {
      return usePaginatedRouteState({
        route: useRoute(),
        router: useRouter(),
        total,
      })
    },
    render() {
      return h('div')
    },
  })

  const wrapper = mount(TestComponent, {
    global: {
      plugins: [router],
    },
  })

  return { wrapper, router }
}

describe('pagination helpers', () => {
  beforeEach(() => {
    window.localStorage.clear()
  })

  it('reads pagination state from the route query and writes updates back', async () => {
    const { wrapper, router } = await mountPaginationState()

    expect((wrapper.vm as { page: number }).page).toBe(2)
    expect((wrapper.vm as { pageSize: number }).pageSize).toBe(50)

    await (wrapper.vm as { setPage: (value: number) => Promise<void> }).setPage(3)
    await nextTick()

    expect(router.currentRoute.value.query.page).toBe('3')
    expect(router.currentRoute.value.query.pageSize).toBe('50')
    expect(router.currentRoute.value.query.keep).toBe('1')
  })

  it('normalizes invalid query values and resets to page one when page size changes', async () => {
    const { wrapper, router } = await mountPaginationState('/?page=0&pageSize=13')

    expect((wrapper.vm as { page: number }).page).toBe(1)
    expect((wrapper.vm as { pageSize: number }).pageSize).toBe(10)

    await (wrapper.vm as { setPageSize: (value: number) => Promise<void> }).setPageSize(10)
    await nextTick()

    expect((wrapper.vm as { page: number }).page).toBe(1)
    expect(router.currentRoute.value.query.page).toBeUndefined()
    expect(router.currentRoute.value.query.pageSize).toBeUndefined()
  })

  it('slices mock studio lists with page and page_size while keeping total', async () => {
    const api = createMockStudioApi()

    const usersPage = await api.listUsers({ page: 2, page_size: 1 })
    expect(usersPage.total).toBe(12)
    expect(usersPage.items).toHaveLength(1)
    expect(usersPage.items[0]?.email).toBe('editor@beehive.local')

    const tagsPage = await api.listTags({ page: 2, page_size: 1 })
    expect(tagsPage.total).toBe(11)
    expect(tagsPage.items).toHaveLength(1)
    expect(tagsPage.items[0]?.name).toBe('Identity')
  })
})
