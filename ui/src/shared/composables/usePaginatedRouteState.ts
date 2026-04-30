import { computed, shallowRef, toValue, watch, type MaybeRefOrGetter } from 'vue'
import type { RouteLocationNormalizedLoaded, Router } from 'vue-router'

const DEFAULT_PAGE = 1
const DEFAULT_PAGE_SIZE = 10
const DEFAULT_PAGE_SIZE_OPTIONS = [10, 20, 50, 100]

interface PaginatedRouteStateOptions {
  route: RouteLocationNormalizedLoaded
  router: Router
  pageParam?: string
  pageSizeParam?: string
  defaultPage?: number
  defaultPageSize?: number
  pageSizeOptions?: number[]
  total?: MaybeRefOrGetter<number>
}

export function usePaginatedRouteState(options: PaginatedRouteStateOptions) {
  const pageParam = options.pageParam ?? 'page'
  const pageSizeParam = options.pageSizeParam ?? 'pageSize'
  const pageSizeOptions = options.pageSizeOptions ?? DEFAULT_PAGE_SIZE_OPTIONS
  const defaultPage = options.defaultPage ?? DEFAULT_PAGE
  const defaultPageSize = sanitizePageSize(options.defaultPageSize ?? DEFAULT_PAGE_SIZE, pageSizeOptions)

  const page = shallowRef(readPositiveInt(options.route.query[pageParam], defaultPage))
  const pageSize = shallowRef(sanitizePageSize(readPositiveInt(options.route.query[pageSizeParam], defaultPageSize), pageSizeOptions))

  watch(
    () => [options.route.query[pageParam], options.route.query[pageSizeParam]],
    () => {
      page.value = readPositiveInt(options.route.query[pageParam], defaultPage)
      pageSize.value = sanitizePageSize(readPositiveInt(options.route.query[pageSizeParam], defaultPageSize), pageSizeOptions)
    },
  )

  watch(
    () => toValue(options.total),
    (total) => {
      const totalPages = Math.max(1, Math.ceil(Math.max(0, total ?? 0) / pageSize.value))
      if ((total ?? 0) > 0 && page.value > totalPages) {
        void setPage(totalPages)
      }
    },
  )

  const totalPages = computed(() => Math.max(1, Math.ceil(Math.max(0, toValue(options.total) ?? 0) / pageSize.value)))

  async function setPage(nextPage: number): Promise<void> {
    const normalizedPage = Math.min(Math.max(1, nextPage), totalPages.value)
    if (normalizedPage === page.value && String(options.route.query[pageParam] ?? '') === String(normalizedPage)) {
      return
    }
    page.value = normalizedPage
    await syncQuery()
  }

  async function setPageSize(nextPageSize: number): Promise<void> {
    const normalizedPageSize = sanitizePageSize(nextPageSize, pageSizeOptions)
    if (normalizedPageSize === pageSize.value && page.value === 1 && String(options.route.query[pageSizeParam] ?? '') === String(normalizedPageSize)) {
      return
    }
    pageSize.value = normalizedPageSize
    page.value = 1
    await syncQuery()
  }

  async function resetPage(): Promise<void> {
    if (page.value === 1 && String(options.route.query[pageParam] ?? '') !== '') {
      await syncQuery()
      return
    }
    if (page.value !== 1) {
      page.value = 1
      await syncQuery()
    }
  }

  async function syncQuery(extraQuery: Record<string, string | undefined> = {}): Promise<void> {
    const nextQuery: Record<string, string> = {}

    for (const [key, value] of Object.entries(options.route.query)) {
      if (value === undefined || value === null) {
        continue
      }
      if (key === pageParam || key === pageSizeParam) {
        continue
      }
      const normalized = Array.isArray(value) ? value[0] : value
      if (typeof normalized === 'string' && normalized !== '') {
        nextQuery[key] = normalized
      }
    }

    for (const [key, value] of Object.entries(extraQuery)) {
      if (value !== undefined && value !== '') {
        nextQuery[key] = value
      } else {
        delete nextQuery[key]
      }
    }

    if (page.value !== defaultPage) {
      nextQuery[pageParam] = String(page.value)
    }
    if (pageSize.value !== defaultPageSize) {
      nextQuery[pageSizeParam] = String(pageSize.value)
    }

    await options.router.replace({ query: nextQuery })
  }

  return {
    page,
    pageSize,
    pageSizeOptions,
    totalPages,
    setPage,
    setPageSize,
    resetPage,
    syncQuery,
  }
}

function readPositiveInt(value: unknown, fallback: number): number {
  const normalized = Array.isArray(value) ? value[0] : value
  const parsed = typeof normalized === 'string' ? Number(normalized) : Number.NaN
  if (!Number.isInteger(parsed) || parsed < 1) {
    return fallback
  }
  return parsed
}

function sanitizePageSize(value: number, options: number[]): number {
  if (!Number.isInteger(value) || value < 1) {
    return options[0] ?? DEFAULT_PAGE_SIZE
  }
  return options.includes(value) ? value : options[0] ?? DEFAULT_PAGE_SIZE
}
