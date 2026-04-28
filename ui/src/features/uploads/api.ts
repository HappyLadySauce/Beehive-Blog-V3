import { requestJson } from '@/shared/api/httpClient'
import { appConfig } from '@/shared/config/env'

import type { FileAssetResponse, FileUploadCreateRequest, FileUploadCreateResponse } from './types'

const fileBasePath = '/api/v3/files'

export async function createFileUpload(
  payload: FileUploadCreateRequest,
  options: { accessToken?: string } = {},
): Promise<FileUploadCreateResponse> {
  if (appConfig.apiMode === 'mock') {
    return {
      asset: {
        asset_id: `mock_asset_${Date.now()}`,
        upload_id: `mock_upload_${Date.now()}`,
        owner_user_id: 'mock',
        scope: payload.scope,
        visibility: payload.visibility ?? 'public',
        status: 'pending',
        bucket: 'mock',
        object_key: `mock/${payload.scope}/${payload.file_name}`,
        public_url: '',
        file_name: payload.file_name,
        content_type: payload.content_type,
        byte_size: payload.byte_size,
        created_at: Math.floor(Date.now() / 1000),
        expires_at: Math.floor(Date.now() / 1000) + 300,
      },
      upload_url: '',
      headers: {},
      expires_at: Math.floor(Date.now() / 1000) + 300,
      max_bytes: payload.scope === 'avatar' ? 2 * 1024 * 1024 : 5 * 1024 * 1024,
    }
  }

  return requestJson<FileUploadCreateResponse>(`${fileBasePath}/uploads`, {
    method: 'POST',
    body: JSON.stringify(payload),
    accessToken: options.accessToken,
  })
}

export async function completeFileUpload(
  uploadId: string,
  options: { accessToken?: string } = {},
): Promise<FileAssetResponse> {
  if (appConfig.apiMode === 'mock') {
    return {
      asset: {
        asset_id: `mock_asset_${uploadId}`,
        upload_id: uploadId,
        owner_user_id: 'mock',
        scope: 'avatar',
        visibility: 'public',
        status: 'uploaded',
        bucket: 'mock',
        object_key: `mock/${uploadId}`,
        public_url: '',
        file_name: 'mock',
        content_type: 'image/png',
        byte_size: 0,
        created_at: Math.floor(Date.now() / 1000),
        expires_at: Math.floor(Date.now() / 1000) + 300,
        uploaded_at: Math.floor(Date.now() / 1000),
      },
    }
  }

  return requestJson<FileAssetResponse>(`${fileBasePath}/uploads/${encodeURIComponent(uploadId)}/complete`, {
    method: 'POST',
    accessToken: options.accessToken,
  })
}
