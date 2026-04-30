export type FileUploadNamespace = string
export type FileUploadVisibility = 'public' | 'private'
export type FileAssetStatus = 'pending' | 'uploaded' | 'deleted'

export interface FileAsset {
  asset_id: string
  upload_id: string
  owner_user_id: string
  namespace: string
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
  namespace: string
  file_name: string
  content_type: string
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
  namespace?: string
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
