import type { FileAsset, FileUploadCreateRequest } from './types'

interface MockUploadState {
  asset: FileAsset
  publicUrl: string
}

const mockUploads = new Map<string, MockUploadState>()
const mockAssets = new Map<string, FileAsset>()

export function createMockUpload(payload: FileUploadCreateRequest): FileAsset {
  const now = Math.floor(Date.now() / 1000)
  const uploadId = `mock_upload_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
  const asset: FileAsset = {
    asset_id: `mock_asset_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
    upload_id: uploadId,
    owner_user_id: 'user_mock_admin',
    namespace: payload.namespace,
    visibility: payload.visibility ?? 'public',
    status: 'pending',
    bucket: 'mock',
    object_key: `mock/${payload.namespace}/${payload.file_name}`,
    public_url: '',
    file_name: payload.file_name,
    content_type: payload.content_type,
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
    namespace: 'avatar',
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
  namespace?: string
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
      const matchesNamespace = !params.namespace || asset.namespace === params.namespace
      const matchesStatus = !params.status || asset.status === params.status
      const matchesVisibility = !params.visibility || asset.visibility === params.visibility
      const matchesOwner = !params.owner_user_id || asset.owner_user_id === params.owner_user_id
      const matchesKeyword = keyword === ''
        || asset.file_name.toLowerCase().includes(keyword)
        || asset.content_type.toLowerCase().includes(keyword)
        || asset.object_key.toLowerCase().includes(keyword)
      return matchesNamespace && matchesStatus && matchesVisibility && matchesOwner && matchesKeyword
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
