import type {
  FileAsset,
  FileCategory,
  FileCategoryCreateRequest,
  FileCategoryExtensionsUpdateRequest,
  FileCategoryUpdateRequest,
  FileUploadCreateRequest,
} from './types'

interface MockUploadState {
  asset: FileAsset
  publicUrl: string
}

const mockUploads = new Map<string, MockUploadState>()
const mockAssets = new Map<string, FileAsset>()
const mockCategories = new Map<string, FileCategory>()

bootstrapMockCategories()

export function createMockUpload(payload: FileUploadCreateRequest): FileAsset {
  const now = Math.floor(Date.now() / 1000)
  const uploadId = `mock_upload_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
  const asset: FileAsset = {
    asset_id: `mock_asset_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
    upload_id: uploadId,
    owner_user_id: 'user_mock_admin',
    category_key: payload.category_key,
    visibility: payload.visibility ?? 'public',
    status: 'pending',
    bucket: 'mock',
    object_key: `mock/${payload.category_key}/${payload.file_name}`,
    public_url: '',
    file_name: payload.file_name,
    content_type: payload.content_type || inferMockContentType(payload.file_name),
    byte_size: payload.byte_size,
    created_at: now,
    expires_at: now + 300,
  }
  mockUploads.set(uploadId, { asset, publicUrl: '' })
  mockAssets.set(asset.asset_id, asset)
  return asset
}

export function getMockUpload(uploadId: string): MockUploadState | undefined {
  return mockUploads.get(uploadId)
}

export function setMockUploadPublicUrl(uploadId: string, publicUrl: string): void {
  const state = mockUploads.get(uploadId)
  if (!state) {
    return
  }
  state.publicUrl = publicUrl
}

export function completeMockUpload(uploadId: string): FileAsset {
  const now = Math.floor(Date.now() / 1000)
  const state = mockUploads.get(uploadId)
  const asset = state?.asset ?? {
    asset_id: `mock_asset_${uploadId}`,
    upload_id: uploadId,
    owner_user_id: 'user_mock_admin',
    category_key: 'default',
    visibility: 'public',
    status: 'pending',
    bucket: 'mock',
    object_key: `mock/${uploadId}`,
    public_url: '',
    file_name: 'mock',
    content_type: 'image/png',
    byte_size: 0,
    created_at: now,
    expires_at: now + 300,
  } satisfies FileAsset

  const completed: FileAsset = {
    ...asset,
    status: 'uploaded',
    public_url: state?.publicUrl || asset.public_url || mockPublicUrl(uploadId),
    uploaded_at: now,
  }
  mockUploads.delete(uploadId)
  mockAssets.set(completed.asset_id, completed)
  return completed
}

export function listMockAssets(params: {
  category_key?: string
  status?: string
  visibility?: string
  owner_user_id?: string
  keyword?: string
  page?: number
  page_size?: number
} = {}): { items: FileAsset[]; total: number; page: number; page_size: number } {
  const keyword = params.keyword?.trim().toLowerCase() ?? ''
  const page = normalizePage(params.page)
  const pageSize = normalizePageSize(params.page_size)
  const items = Array.from(mockAssets.values())
    .filter((asset) => {
      const matchesCategory = !params.category_key || asset.category_key === params.category_key
      const matchesStatus = !params.status || asset.status === params.status
      const matchesVisibility = !params.visibility || asset.visibility === params.visibility
      const matchesOwner = !params.owner_user_id || asset.owner_user_id === params.owner_user_id
      const matchesKeyword = keyword === ''
        || asset.file_name.toLowerCase().includes(keyword)
        || asset.content_type.toLowerCase().includes(keyword)
        || asset.object_key.toLowerCase().includes(keyword)
      return matchesCategory && matchesStatus && matchesVisibility && matchesOwner && matchesKeyword
    })
    .sort((left, right) => {
      const leftUpdated = left.deleted_at ?? left.uploaded_at ?? left.created_at
      const rightUpdated = right.deleted_at ?? right.uploaded_at ?? right.created_at
      return rightUpdated - leftUpdated
    })
  const start = (page - 1) * pageSize
  return {
    items: items.slice(start, start + pageSize),
    total: items.length,
    page,
    page_size: pageSize,
  }
}

export function getMockAsset(assetId: string): FileAsset | undefined {
  return mockAssets.get(assetId)
}

export function deleteMockAsset(assetId: string): FileAsset | undefined {
  const asset = mockAssets.get(assetId)
  if (!asset) {
    return undefined
  }
  const deleted: FileAsset = {
    ...asset,
    status: 'deleted',
    deleted_at: Math.floor(Date.now() / 1000),
  }
  mockAssets.set(assetId, deleted)
  return deleted
}

function normalizePage(page?: number): number {
  if (!Number.isInteger(page) || (page ?? 0) < 1) {
    return 1
  }
  return page as number
}

function normalizePageSize(pageSize?: number): number {
  if (!Number.isInteger(pageSize) || (pageSize ?? 0) < 1) {
    return 20
  }
  return Math.min(pageSize as number, 100)
}

function mockPublicUrl(uploadId: string): string {
  return `data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='1' height='1'%3E%3C/svg%3E#${encodeURIComponent(uploadId)}`
}

export function listMockCategories(includeDisabled = false): FileCategory[] {
  return Array.from(mockCategories.values())
    .filter((category) => includeDisabled || category.enabled)
    .sort((left, right) => {
      if (left.is_default !== right.is_default) {
        return left.is_default ? -1 : 1
      }
      if (left.sort_order !== right.sort_order) {
        return left.sort_order - right.sort_order
      }
      return left.display_name.localeCompare(right.display_name)
    })
}

export function createMockCategory(payload: FileCategoryCreateRequest): FileCategory {
  const now = Math.floor(Date.now() / 1000)
  const category: FileCategory = {
    category_key: payload.category_key,
    display_name: payload.display_name,
    description: payload.description ?? '',
    enabled: payload.enabled ?? true,
    is_default: payload.is_default ?? false,
    sort_order: payload.sort_order ?? 0,
    allowed_extensions: normalizeExtensions(payload.allowed_extensions),
    created_at: now,
    updated_at: now,
  }
  mockCategories.set(category.category_key, category)
  if (category.is_default) {
    return setMockDefaultCategory(category.category_key)
  }
  return category
}

export function updateMockCategory(categoryKey: string, payload: FileCategoryUpdateRequest): FileCategory {
  const existing = requireMockCategory(categoryKey)
  const updated: FileCategory = {
    ...existing,
    display_name: payload.display_name,
    description: payload.description ?? '',
    enabled: payload.enabled ?? existing.enabled,
    sort_order: payload.sort_order ?? existing.sort_order,
    updated_at: Math.floor(Date.now() / 1000),
  }
  mockCategories.set(categoryKey, updated)
  return updated
}

export function updateMockCategoryExtensions(categoryKey: string, payload: FileCategoryExtensionsUpdateRequest): FileCategory {
  const existing = requireMockCategory(categoryKey)
  const updated: FileCategory = {
    ...existing,
    allowed_extensions: normalizeExtensions(payload.allowed_extensions),
    updated_at: Math.floor(Date.now() / 1000),
  }
  mockCategories.set(categoryKey, updated)
  return updated
}

export function setMockDefaultCategory(categoryKey: string): FileCategory {
  const now = Math.floor(Date.now() / 1000)
  for (const [key, category] of mockCategories.entries()) {
    mockCategories.set(key, {
      ...category,
      is_default: key === categoryKey,
      updated_at: now,
    })
  }
  return requireMockCategory(categoryKey)
}

function requireMockCategory(categoryKey: string): FileCategory {
  const category = mockCategories.get(categoryKey)
  if (!category) {
    throw new Error('File category not found.')
  }
  return category
}

function bootstrapMockCategories(): void {
  if (mockCategories.size > 0) {
    return
  }
  const now = Math.floor(Date.now() / 1000)
  mockCategories.set('default', {
    category_key: 'default',
    display_name: '默认类型',
    description: 'System default file category.',
    enabled: true,
    is_default: true,
    sort_order: 0,
    allowed_extensions: ['.avif', '.jpeg', '.jpg', '.pdf', '.png', '.webp'],
    created_at: now,
    updated_at: now,
  })
}

function normalizeExtensions(extensions: string[]): string[] {
  return Array.from(new Set(extensions.map((item) => {
    const value = item.trim().toLowerCase()
    if (!value) {
      return ''
    }
    return value.startsWith('.') ? value : `.${value}`
  }).filter(Boolean))).sort()
}

function inferMockContentType(fileName: string): string {
  const ext = fileName.split('.').pop()?.trim().toLowerCase() ?? ''
  switch (ext) {
    case 'avif':
      return 'image/avif'
    case 'jpeg':
    case 'jpg':
      return 'image/jpeg'
    case 'pdf':
      return 'application/pdf'
    case 'png':
      return 'image/png'
    case 'webp':
      return 'image/webp'
    default:
      return 'application/octet-stream'
  }
}
