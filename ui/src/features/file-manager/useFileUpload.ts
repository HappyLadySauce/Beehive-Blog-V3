import { shallowRef } from 'vue'

import { i18n } from '@/shared/i18n'

import { completeFileUpload, createFileUpload, putFileUploadObject } from './api'
import type { FileUploadCreateRequest, FileUploadNamespace } from './types'

const imageContentTypes = [
  'image/avif',
  'image/jpeg',
  'image/png',
  'image/webp',
] as const

const genericAllowedContentTypes: readonly string[] = [...imageContentTypes, 'application/pdf']

const contentTypesByExtension: Record<string, string> = {
  avif: 'image/avif',
  jpeg: 'image/jpeg',
  jpg: 'image/jpeg',
  pdf: 'application/pdf',
  png: 'image/png',
  webp: 'image/webp',
}

export function useFileUpload() {
  const isUploading = shallowRef(false)
  const errorMessage = shallowRef('')

  async function uploadFile(file: File, accessToken: string | undefined, namespace: FileUploadNamespace): Promise<string> {
    errorMessage.value = ''
    isUploading.value = true
    try {
      const contentType = normalizeContentType(file)
      validateFile(file, contentType)
      const upload = await createFileUpload(createUploadPayload(file, namespace, contentType), { accessToken })
      await putFileUploadObject(upload, file)
      const completed = await completeFileUpload(upload.asset.upload_id, { accessToken })
      return completed.asset.public_url || upload.asset.public_url
    } catch (error) {
      errorMessage.value = error instanceof Error ? error.message : String(i18n.global.t('uploads.uploadFailed'))
      throw error
    } finally {
      isUploading.value = false
    }
  }

  async function uploadAvatar(file: File, accessToken?: string): Promise<string> {
    return uploadFile(file, accessToken, 'avatar')
  }

  async function uploadImage(file: File, accessToken: string | undefined, namespace: FileUploadNamespace): Promise<string> {
    return uploadFile(file, accessToken, namespace)
  }

  return {
    isUploading,
    errorMessage,
    uploadFile,
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

function validateFile(file: File, contentType: string): void {
  if (file.name.trim() === '') {
    throw new Error(String(i18n.global.t('uploads.fileNameRequired')))
  }
  if (contentType === '') {
    throw new Error(String(i18n.global.t('uploads.fileTypeUnsupported')))
  }
  if (!genericAllowedContentTypes.includes(contentType)) {
    throw new Error(String(i18n.global.t('uploads.fileTypeUnsupported')))
  }
  if (file.size <= 0) {
    throw new Error(String(i18n.global.t('uploads.fileEmpty')))
  }
  if (file.size > 20 * 1024 * 1024) {
    throw new Error(String(i18n.global.t('uploads.fileTooLarge')))
  }
}

function createUploadPayload(file: File, namespace: FileUploadNamespace, contentType: string): FileUploadCreateRequest {
  return {
    namespace,
    file_name: file.name,
    content_type: contentType,
    byte_size: file.size,
    visibility: 'public',
  }
}
