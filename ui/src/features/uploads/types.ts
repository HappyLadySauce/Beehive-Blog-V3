export type FileUploadScope = 'avatar' | 'content_cover' | 'content_image' | 'attachment'
export type FileUploadVisibility = 'public' | 'private'
export type FileAssetStatus = 'pending' | 'uploaded' | 'deleted'

export interface FileAsset {
  asset_id: string
  upload_id: string
  owner_user_id: string
  scope: FileUploadScope
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

export interface FileUploadCreateRequest {
  scope: FileUploadScope
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
