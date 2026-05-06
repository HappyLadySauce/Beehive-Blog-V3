import { shallowRef } from 'vue'

import { i18n } from '@/shared/i18n'

import { completeFileUpload, createFileUpload, putFileUploadObject } from './api'
import {
  DEFAULT_FILE_CATEGORY_EXTENSIONS,
  DEFAULT_FILE_CATEGORY_KEY,
  DEFAULT_FILE_MAX_UPLOAD_BYTES,
  normalizeFileExtension,
} from './constants'
import type { FileCategoryKey, FileUploadCreateRequest, FileUploadVisibility } from './types'

const contentTypesByExtension: Record<string, string> = {
  '.avif': 'image/avif',
  '.jpeg': 'image/jpeg',
  '.jpg': 'image/jpeg',
  '.pdf': 'application/pdf',
  '.png': 'image/png',
  '.webp': 'image/webp',
}

export interface UploadFileOptions {
  allowedExtensions?: readonly string[]
  maxUploadBytes?: number
  visibility?: FileUploadVisibility
}

export function useFileUpload() {
  const isUploading = shallowRef(false)
  const errorMessage = shallowRef('')

  async function uploadFile(
    file: File,
    accessToken: string | undefined,
    categoryKey: FileCategoryKey = DEFAULT_FILE_CATEGORY_KEY,
    options: UploadFileOptions = {},
  ): Promise<string> {
    errorMessage.value = ''
    isUploading.value = true
    try {
      const contentType = normalizeContentType(file)
      validateFile(file, {
        allowedExtensions: options.allowedExtensions ?? DEFAULT_FILE_CATEGORY_EXTENSIONS,
        maxUploadBytes: options.maxUploadBytes ?? DEFAULT_FILE_MAX_UPLOAD_BYTES,
      })
      const upload = await createFileUpload(createUploadPayload(file, categoryKey, contentType, options.visibility), { accessToken })
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
    return uploadFile(file, accessToken, DEFAULT_FILE_CATEGORY_KEY, {
      allowedExtensions: DEFAULT_FILE_CATEGORY_EXTENSIONS,
    })
  }

  async function uploadImage(
    file: File,
    accessToken: string | undefined,
    categoryKey: FileCategoryKey = DEFAULT_FILE_CATEGORY_KEY,
    options: UploadFileOptions = {},
  ): Promise<string> {
    return uploadFile(file, accessToken, categoryKey, options)
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
  return contentTypesByExtension[normalizeFileExtension(file.name)] ?? 'application/octet-stream'
}

function validateFile(
  file: File,
  options: Required<Pick<UploadFileOptions, 'allowedExtensions' | 'maxUploadBytes'>>,
): void {
  if (file.name.trim() === '') {
    throw new Error(String(i18n.global.t('uploads.fileNameRequired')))
  }

  const extension = normalizeFileExtension(file.name)
  if (extension === '' || !options.allowedExtensions.map((item) => item.toLowerCase()).includes(extension)) {
    throw new Error(String(i18n.global.t('uploads.fileTypeUnsupported')))
  }

  if (file.size <= 0) {
    throw new Error(String(i18n.global.t('uploads.fileEmpty')))
  }

  if (file.size > options.maxUploadBytes) {
    throw new Error(String(i18n.global.t('uploads.fileTooLarge')))
  }
}

function createUploadPayload(
  file: File,
  categoryKey: FileCategoryKey,
  contentType: string,
  visibility: FileUploadVisibility = 'public',
): FileUploadCreateRequest {
  return {
    category_key: categoryKey,
    file_name: file.name,
    content_type: contentType,
    byte_size: file.size,
    visibility,
  }
}
