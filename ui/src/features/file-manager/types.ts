export type FileCategoryKey = string
export type FileUploadVisibility = 'public' | 'private'
export type FileAssetStatus = 'pending' | 'uploaded' | 'deleted'

export interface FileCategory {
  category_key: string
  display_name: string
  description: string
  enabled: boolean
  is_default: boolean
  sort_order: number
  allowed_extensions: string[]
  created_at: number
  updated_at: number
}

export interface FileCategoryListResponse {
  items: FileCategory[]
}

export interface FileCategoryResponse {
  category: FileCategory
}

export interface FileCategoryCreateRequest {
  category_key: string
  display_name: string
  description?: string
  enabled?: boolean
  is_default?: boolean
  sort_order?: number
  allowed_extensions: string[]
}

export interface FileCategoryUpdateRequest {
  display_name: string
  description?: string
  enabled?: boolean
  sort_order?: number
}

export interface FileCategoryExtensionsUpdateRequest {
  allowed_extensions: string[]
}

export interface FileConfig {
  max_upload_bytes: number
  presign_ttl_seconds: number
}

export interface FileConfigResponse {
  config: FileConfig
}

export interface FileAsset {
  asset_id: string
  upload_id: string
  owner_user_id: string
  category_key: string
  visibility: FileUploadVisibility
  status: FileAssetStatus
  bucket: string
  object_key: string
  public_url: string
  file_name: string
  content_type: string
  byte_size: number
  created_at: number
  expires_at: number
  uploaded_at?: number
  deleted_at?: number
}

export type FileAssetSummary = FileAsset
export type FileAssetDetail = FileAsset

export interface FileUploadCreateRequest {
  category_key: string
  file_name: string
  content_type?: string
  byte_size: number
  visibility?: FileUploadVisibility
}

export interface FileUploadCreateResponse {
  asset: FileAsset
  upload_url: string
  headers: Record<string, string>
  expires_at: number
  max_bytes: number
}

export interface FileAssetResponse {
  asset: FileAsset
}

export interface FileAssetListParams {
  category_key?: string
  status?: string
  visibility?: string
  owner_user_id?: string
  keyword?: string
  page?: number
  page_size?: number
}

export interface FileAssetListResponse {
  items: FileAssetSummary[]
  total: number
  page: number
  page_size: number
}

export interface FileAssetDetailResponse {
  asset: FileAssetDetail
}

export interface FileAssetDeleteResponse {
  ok: boolean
}
