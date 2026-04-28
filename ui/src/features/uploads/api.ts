import { requestJson } from '@/shared/api/httpClient'
import { appConfig } from '@/shared/config/env'

import type { FileAsset, FileAssetResponse, FileUploadCreateRequest, FileUploadCreateResponse } from './types'

const fileBasePath = '/api/v3/files'
const mockUploads = new Map<string, { asset: FileAsset, publicUrl: string }>()

export async function createFileUpload(
  payload: FileUploadCreateRequest,
  options: { accessToken?: string } = {},
): Promise<FileUploadCreateResponse> {
  if (appConfig.apiMode === 'mock') {
    const now = Math.floor(Date.now() / 1000)
    const uploadId = `mock_upload_${Date.now()}`
    const asset: FileAsset = {
      asset_id: `mock_asset_${Date.now()}`,
      upload_id: uploadId,
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
      created_at: now,
      expires_at: now + 300,
    }
    mockUploads.set(uploadId, { asset, publicUrl: '' })
    return {
      asset,
      upload_url: '',
      headers: {},
      expires_at: now + 300,
      max_bytes: payload.scope === 'avatar' ? 2 * 1024 * 1024 : 5 * 1024 * 1024,
    }
  }

  requireLiveAccessToken(options.accessToken)
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
    const now = Math.floor(Date.now() / 1000)
    const state = mockUploads.get(uploadId)
    const asset = state?.asset ?? {
      asset_id: `mock_asset_${uploadId}`,
      upload_id: uploadId,
      owner_user_id: 'mock',
      scope: 'avatar',
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
    return {
      asset: {
        ...asset,
        status: 'uploaded',
        public_url: state?.publicUrl || asset.public_url || mockPublicUrl(uploadId),
        uploaded_at: now,
      },
    }
  }

  requireLiveAccessToken(options.accessToken)
  return requestJson<FileAssetResponse>(`${fileBasePath}/uploads/${encodeURIComponent(uploadId)}/complete`, {
    method: 'POST',
    accessToken: options.accessToken,
  })
}

export async function putFileUploadObject(upload: FileUploadCreateResponse, file: File): Promise<void> {
  if (appConfig.apiMode === 'mock') {
    const state = mockUploads.get(upload.asset.upload_id)
    if (state) {
      state.publicUrl = createPreviewUrl(file, upload.asset.upload_id)
    }
    return
  }
  if (!upload.upload_url) {
    throw new Error('Upload URL is missing.')
  }

  const controller = new AbortController()
  const timeoutId = window.setTimeout(() => controller.abort(), 60_000)
  try {
    const response = await fetch(upload.upload_url, {
      method: 'PUT',
      body: file,
      headers: upload.headers,
      signal: controller.signal,
    })
    if (!response.ok) {
      throw new Error(response.statusText || 'Unable to upload file.')
    }
  } catch (error) {
    if (error instanceof DOMException && error.name === 'AbortError') {
      throw new Error('File upload timed out.')
    }
    throw error
  } finally {
    window.clearTimeout(timeoutId)
  }
}

function requireLiveAccessToken(accessToken?: string): asserts accessToken is string {
  if (!accessToken) {
    throw new Error('Sign in to upload files.')
  }
}

function createPreviewUrl(file: File, uploadId: string): string {
  if (typeof URL.createObjectURL === 'function') {
    return URL.createObjectURL(file)
  }
  return mockPublicUrl(uploadId)
}

function mockPublicUrl(uploadId: string): string {
  return `data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='1' height='1'%3E%3C/svg%3E#${encodeURIComponent(uploadId)}`
}
