import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { useProgressiveQuery } from '@/shared/composables'

const QueryHarness = defineComponent({
  props: {
    queryKey: {
      type: Array,
      required: true,
    },
    queryFn: {
      type: Function,
      required: true,
    },
  },
  setup(props) {
    return useProgressiveQuery({
      queryKey: props.queryKey,
      queryFn: props.queryFn as () => Promise<unknown>,
    })
  },
  template: '<div />',
})

function mountHarness(queryKey: unknown[], queryFn: () => Promise<unknown>) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: 0,
        staleTime: 0,
      },
    },
  })
  return mount(QueryHarness, {
    props: { queryKey, queryFn },
    global: {
      plugins: [[VueQueryPlugin, { queryClient }]],
    },
  })
}

describe('useProgressiveQuery', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('does not show blocking loading before the delay threshold', async () => {
    let resolveQuery: ((value: { items: string[] }) => void) | null = null
    const wrapper = mountHarness(['progressive-query-fast'], () => new Promise((resolve) => {
      resolveQuery = resolve
    }))

    expect(wrapper.vm.showBlockingLoading).toBe(false)
    await vi.advanceTimersByTimeAsync(500)
    expect(wrapper.vm.showBlockingLoading).toBe(false)

    resolveQuery?.({ items: ['ready'] })
    await flushPromises()

    expect(wrapper.vm.showBlockingLoading).toBe(false)
    expect(wrapper.vm.hasResolvedOnce).toBe(true)
  })

  it('shows blocking loading after the delay when no data has resolved yet', async () => {
    const wrapper = mountHarness(['progressive-query-slow'], () => new Promise(() => {}))

    await vi.advanceTimersByTimeAsync(1000)

    expect(wrapper.vm.showBlockingLoading).toBe(true)
    expect(wrapper.vm.showRefreshingHint).toBe(false)
  })

  it('shows a refreshing hint instead of blocking once data already exists', async () => {
    let resolveInitial: ((value: { items: string[] }) => void) | null = null
    let resolveRefetch: ((value: { items: string[] }) => void) | null = null
    let calls = 0
    const wrapper = mountHarness(['progressive-query-refreshing'], () => {
      calls += 1
      if (calls === 1) {
        return new Promise((resolve) => {
          resolveInitial = resolve
        })
      }
      return new Promise((resolve) => {
        resolveRefetch = resolve
      })
    })

    resolveInitial?.({ items: ['first'] })
    await flushPromises()
    expect(wrapper.vm.hasResolvedOnce).toBe(true)

    void wrapper.vm.refetch()
    await vi.advanceTimersByTimeAsync(1000)

    expect(wrapper.vm.showBlockingLoading).toBe(false)
    expect(wrapper.vm.showRefreshingHint).toBe(true)

    resolveRefetch?.({ items: ['second'] })
    await flushPromises()
  })

  it('keeps the loading indicator visible for the minimum duration once shown', async () => {
    let resolveQuery: ((value: { items: string[] }) => void) | null = null
    const wrapper = mountHarness(['progressive-query-min-visible'], () => new Promise((resolve) => {
      resolveQuery = resolve
    }))

    await vi.advanceTimersByTimeAsync(1000)
    expect(wrapper.vm.showBlockingLoading).toBe(true)

    resolveQuery?.({ items: ['done'] })
    await flushPromises()
    expect(wrapper.vm.showBlockingLoading).toBe(true)

    await vi.advanceTimersByTimeAsync(399)
    expect(wrapper.vm.showBlockingLoading).toBe(true)

    await vi.advanceTimersByTimeAsync(1)
    expect(wrapper.vm.showBlockingLoading).toBe(false)
  })
})
