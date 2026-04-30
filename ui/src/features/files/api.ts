import { requestJson } from '@/shared/api/httpClient'
import { appConfig } from '@/shared/config/env'

import { deleteMockAsset, getMockAsset, listMockAssets } from '@/features/uploads/mockStore'

import type { FileAssetDeleteResponse, FileAssetDetailResponse, FileAssetListParams, FileAssetListResponse } from './types'

const fileBasePath = '/api/v3/files'

function withQuery(path: string, params?: Record<string, string | number | boolean | undefined>): string {
  const searchParams = new URLSearchParams()
  for (const [key, value] of Object.entries(params ?? {})) {
    if (value !== undefined && value !== '') {
      searchParams.set(key, String(value))
    }
  }
  const query = searchParams.toString()
  return query ? `${path}?${query}` : path
}

export async function listFileAssets(
  params: FileAssetListParams = {},
  options: { accessToken?: string } = {},
): Promise<FileAssetListResponse> {
  if (appConfig.apiMode === 'mock') {
    return listMockAssets(params)
  }
  return requestJson<FileAssetListResponse>(withQuery(`${fileBasePath}/assets`, params), {
    method: 'GET',
    accessToken: options.accessToken,
  })
}

export async function getFileAsset(assetId: string, options: { accessToken?: string } = {}): Promise<FileAssetDetailResponse> {
  if (appConfig.apiMode === 'mock') {
    const asset = getMockAsset(assetId)
    if (!asset) {
      throw new Error('File asset not found.')
    }
    return { asset }
  }
  return requestJson<FileAssetDetailResponse>(`${fileBasePath}/assets/${encodeURIComponent(assetId)}`, {
    method: 'GET',
    accessToken: options.accessToken,
  })
}

export async function deleteFileAsset(assetId: string, options: { accessToken?: string } = {}): Promise<FileAssetDeleteResponse> {
  if (appConfig.apiMode === 'mock') {
    const asset = deleteMockAsset(assetId)
    if (!asset) {
      throw new Error('File asset not found.')
    }
    return { ok: true }
  }
  return requestJson<FileAssetDeleteResponse>(`${fileBasePath}/assets/${encodeURIComponent(assetId)}`, {
    method: 'DELETE',
    accessToken: options.accessToken,
  })
}
