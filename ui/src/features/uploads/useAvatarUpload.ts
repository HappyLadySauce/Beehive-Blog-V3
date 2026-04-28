import { shallowRef } from 'vue'

import { completeFileUpload, createFileUpload, putFileUploadObject } from './api'
import type {
  FileUploadCreateRequest,
  FileUploadScope,
} from './types'

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
      const upload = await createFileUpload(createUploadPayload(file, scope, contentType), { accessToken })
      await putFileUploadObject(upload, file)
      const completed = await completeFileUpload(upload.asset.upload_id, { accessToken })
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

function createUploadPayload(
  file: File,
  scope: FileUploadScope,
  contentType: string,
): FileUploadCreateRequest {
  return {
    scope,
    file_name: file.name,
    content_type: contentType,
    byte_size: file.size,
    visibility: 'public',
  }
}
