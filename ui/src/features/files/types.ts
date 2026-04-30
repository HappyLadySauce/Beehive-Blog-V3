import type { FileAsset } from '@/features/uploads/types'

export type FileAssetSummary = FileAsset
export type FileAssetDetail = FileAsset

export interface FileAssetListParams {
  scope?: string
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
