import { shallowRef } from 'vue'

import { requestJson } from '@/shared/api/httpClient'

import type {
  FileAssetResponse,
  FileUploadCreateRequest,
  FileUploadCreateResponse,
  FileUploadScope,
} from './types'

const uploadBasePath = '/api/v3/files/uploads'
const defaultMaxBytesByScope: Record<FileUploadScope, number> = {
  avatar: 2 * 1024 * 1024,
  content_cover: 5 * 1024 * 1024,
  content_image: 5 * 1024 * 1024,
  attachment: 20 * 1024 * 1024,
}

const contentTypesByExtension: Record<string, string> = {
  avif: 'image/avif',
  jpeg: 'image/jpeg',
  jpg: 'image/jpeg',
  pdf: 'application/pdf',
  png: 'image/png',
  webp: 'image/webp',
}

export function useAvatarUpload() {
  const isUploading = shallowRef(false)
  const errorMessage = shallowRef('')

  async function uploadAvatar(file: File, accessToken?: string): Promise<string> {
    return uploadImage(file, accessToken, 'avatar')
  }

  async function uploadImage(file: File, accessToken: string | undefined, scope: FileUploadScope): Promise<string> {
    errorMessage.value = ''
    isUploading.value = true
    try {
      const contentType = normalizeContentType(file)
      validateFile(file, scope, contentType)
      const upload = await createUpload(file, scope, contentType, accessToken)
      await putObject(upload.upload_url, upload.headers, file)
      const completed = await completeUpload(upload.asset.upload_id, accessToken)
      return completed.asset.public_url || upload.asset.public_url
    } catch (error) {
      errorMessage.value = error instanceof Error ? error.message : 'Unable to upload file.'
      throw error
    } finally {
      isUploading.value = false
    }
  }

  return {
    isUploading,
    errorMessage,
    uploadAvatar,
    uploadImage,
  }
}

function normalizeContentType(file: File): string {
  const explicitType = file.type.split(';')[0]?.trim().toLowerCase()
  if (explicitType) {
    return explicitType
  }
  const extension = file.name.split('.').pop()?.trim().toLowerCase() ?? ''
  return contentTypesByExtension[extension] ?? ''
}

function validateFile(file: File, scope: FileUploadScope, contentType: string): void {
  if (file.name.trim() === '') {
    throw new Error('File name is required.')
  }
  if (contentType === '') {
    throw new Error('File type is not supported.')
  }
  if (file.size <= 0) {
    throw new Error('File is empty.')
  }
  if (file.size > defaultMaxBytesByScope[scope]) {
    throw new Error('File is too large.')
  }
}

async function createUpload(
  file: File,
  scope: FileUploadScope,
  contentType: string,
  accessToken?: string,
): Promise<FileUploadCreateResponse> {
  requireAccessToken(accessToken)
  const payload: FileUploadCreateRequest = {
    scope,
    file_name: file.name,
    content_type: contentType,
    byte_size: file.size,
    visibility: 'public',
  }
  return requestJson<FileUploadCreateResponse>(uploadBasePath, {
    method: 'POST',
    body: JSON.stringify(payload),
    accessToken,
  })
}

async function completeUpload(uploadId: string, accessToken?: string): Promise<FileAssetResponse> {
  requireAccessToken(accessToken)
  return requestJson<FileAssetResponse>(`${uploadBasePath}/${encodeURIComponent(uploadId)}/complete`, {
    method: 'POST',
    accessToken,
  })
}

async function putObject(uploadURL: string, headers: Record<string, string>, file: File): Promise<void> {
  const controller = new AbortController()
  const timeoutId = window.setTimeout(() => controller.abort(), 60_000)
  try {
    const response = await fetch(uploadURL, {
      method: 'PUT',
      body: file,
      headers,
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

function requireAccessToken(accessToken?: string): asserts accessToken is string {
  if (!accessToken) {
    throw new Error('Sign in to upload files.')
  }
}
