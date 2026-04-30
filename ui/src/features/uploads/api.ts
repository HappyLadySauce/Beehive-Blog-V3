import { requestJson } from '@/shared/api/httpClient'
import { appConfig } from '@/shared/config/env'

import type { FileAssetResponse, FileUploadCreateRequest, FileUploadCreateResponse } from './types'
import { completeMockUpload, createMockUpload, getMockUpload, setMockUploadPublicUrl } from './mockStore'

const fileBasePath = '/api/v3/files'
let activeCompletedMockObjectUrl = ''

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

function mockPublicUrl(uploadId: string): string {
  return `data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='1' height='1'%3E%3C/svg%3E#${encodeURIComponent(uploadId)}`
}
