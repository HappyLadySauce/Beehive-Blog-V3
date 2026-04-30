import { requestJson } from '@/shared/api/httpClient'
import { appConfig } from '@/shared/config/env'

import {
  completeMockUpload,
  createMockUpload,
  deleteMockAsset,
  getMockAsset,
  getMockUpload,
  listMockAssets,
  setMockUploadPublicUrl,
} from './mockStore'

import type {
  FileAssetDeleteResponse,
  FileAssetDetailResponse,
  FileAssetListParams,
  FileAssetListResponse,
  FileAssetResponse,
  FileUploadCreateRequest,
  FileUploadCreateResponse,
} from './types'

const fileBasePath = '/api/v3/files'
let activeCompletedMockObjectUrl = ''

// ---- Upload operations ----

export async function createFileUpload(
  payload: FileUploadCreateRequest,
  options: { accessToken?: string } = {},
): Promise<FileUploadCreateResponse> {
  if (appConfig.apiMode === 'mock') {
    const now = Math.floor(Date.now() / 1000)
    const asset = createMockUpload(payload)
    return {
      asset,
      upload_url: '',
      headers: {},
      expires_at: now + 300,
      max_bytes: 5 * 1024 * 1024,
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
    const existingState = getMockUpload(uploadId)
    const asset = completeMockUpload(uploadId)
    if (existingState) {
      replaceCompletedMockObjectUrl(asset.public_url)
    }
    return { asset }
  }

  requireLiveAccessToken(options.accessToken)
  return requestJson<FileAssetResponse>(`${fileBasePath}/uploads/${encodeURIComponent(uploadId)}/complete`, {
    method: 'POST',
    accessToken: options.accessToken,
  })
}

export async function putFileUploadObject(upload: FileUploadCreateResponse, file: File): Promise<void> {
  if (appConfig.apiMode === 'mock') {
    setMockUploadPublicUrl(upload.asset.upload_id, createPreviewUrl(file, upload.asset.upload_id))
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

// ---- Asset management operations ----

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

// ---- Helpers ----

function requireLiveAccessToken(accessToken?: string): asserts accessToken is string {
  if (!accessToken) {
    throw new Error('Sign in to upload files.')
  }
}

function createPreviewUrl(file: File, uploadId: string): string {
  if (typeof URL.createObjectURL === 'function') {
    return URL.createObjectURL(file)
  }
  return `data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='1' height='1'%3E%3C/svg%3E#${encodeURIComponent(uploadId)}`
}

function replaceCompletedMockObjectUrl(nextUrl: string): void {
  if (
    activeCompletedMockObjectUrl &&
    activeCompletedMockObjectUrl !== nextUrl &&
    typeof URL.revokeObjectURL === 'function'
  ) {
    URL.revokeObjectURL(activeCompletedMockObjectUrl)
  }
  activeCompletedMockObjectUrl = isObjectUrl(nextUrl) ? nextUrl : ''
}

function isObjectUrl(url: string): boolean {
  return url.startsWith('blob:')
}
