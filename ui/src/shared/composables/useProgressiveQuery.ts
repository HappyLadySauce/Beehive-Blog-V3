import {
  keepPreviousData,
  useQuery,
  useQueryClient,
  type QueryKey,
  type UseQueryOptions,
} from '@tanstack/vue-query'
import { computed, onBeforeUnmount, shallowRef, toValue, watch, type MaybeRefOrGetter, type ShallowRef } from 'vue'

const DEFAULT_DELAY_MS = 1000
const DEFAULT_MIN_VISIBLE_MS = 400

type ProgressiveQueryOptions<TQueryFnData, TError, TData, TQueryKey extends QueryKey> =
  Omit<UseQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'placeholderData'> & {
    placeholderData?: UseQueryOptions<TQueryFnData, TError, TData, TQueryKey>['placeholderData']
    delayMs?: number
    minVisibleMs?: number
    keepPreviousData?: boolean
  }

export function useProgressiveQuery<
  TQueryFnData = unknown,
  TError = Error,
  TData = TQueryFnData,
  TQueryKey extends QueryKey = QueryKey,
>(options: MaybeRefOrGetter<ProgressiveQueryOptions<TQueryFnData, TError, TData, TQueryKey>>) {
  const queryClient = useQueryClient()
  const hasResolvedOnce = shallowRef(false)
  const showBlockingLoading = shallowRef(false)
  const showRefreshingHint = shallowRef(false)
  const visibleSince = shallowRef(0)
  let revealTimer: number | undefined
  let hideTimer: number | undefined

  const normalizedOptions = computed(() => {
    const resolved = toValue(options)
    const {
      delayMs: _delayMs,
      minVisibleMs: _minVisibleMs,
      keepPreviousData: shouldKeepPreviousData,
      placeholderData,
      ...queryOptions
    } = resolved
    return {
      ...queryOptions,
      placeholderData: shouldKeepPreviousData === false
        ? placeholderData
        : (placeholderData ?? keepPreviousData),
    } satisfies UseQueryOptions<TQueryFnData, TError, TData, TQueryKey>
  })

  const query = useQuery(normalizedOptions)

  function clearTimers(): void {
    window.clearTimeout(revealTimer)
    window.clearTimeout(hideTimer)
    revealTimer = undefined
    hideTimer = undefined
  }

  function setVisible(mode: 'blocking' | 'refreshing' | null): void {
    showBlockingLoading.value = mode === 'blocking'
    showRefreshingHint.value = mode === 'refreshing'
    visibleSince.value = mode ? Date.now() : 0
  }

  function hideVisibleState(): void {
    const visible = showBlockingLoading.value || showRefreshingHint.value
    if (!visible) {
      setVisible(null)
      return
    }
    const minVisibleMs = toValue(options).minVisibleMs ?? DEFAULT_MIN_VISIBLE_MS
    const elapsed = Date.now() - visibleSince.value
    const remaining = Math.max(0, minVisibleMs - elapsed)
    if (remaining === 0) {
      setVisible(null)
      return
    }
    window.clearTimeout(hideTimer)
    hideTimer = window.setTimeout(() => {
      setVisible(null)
      hideTimer = undefined
    }, remaining)
  }

  watch(
    () => query.data.value,
    (value) => {
      if (value !== undefined) {
        hasResolvedOnce.value = true
      }
    },
    { immediate: true },
  )

  watch(
    () => query.isFetching.value,
    (isFetching) => {
      clearTimers()
      if (!isFetching) {
        hideVisibleState()
        return
      }

      const delayMs = toValue(options).delayMs ?? DEFAULT_DELAY_MS
      if (delayMs <= 0) {
        setVisible(hasResolvedOnce.value ? 'refreshing' : 'blocking')
        return
      }

      revealTimer = window.setTimeout(() => {
        revealTimer = undefined
        if (!query.isFetching.value) {
          return
        }
        setVisible(hasResolvedOnce.value ? 'refreshing' : 'blocking')
      }, delayMs)
    },
    { immediate: true },
  )

  onBeforeUnmount(() => {
    clearTimers()
  })

  async function invalidate() {
    return queryClient.invalidateQueries({ queryKey: normalizedOptions.value.queryKey })
  }

  return {
    ...query,
    hasResolvedOnce,
    isFetching: query.isFetching as ShallowRef<boolean>,
    showBlockingLoading,
    showRefreshingHint,
    invalidate,
  }
}
