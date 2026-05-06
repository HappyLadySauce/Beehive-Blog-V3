export const DEFAULT_FILE_CATEGORY_KEY = 'default'
export const DEFAULT_FILE_MAX_UPLOAD_BYTES = 2 * 1024 * 1024 * 1024
export const IMAGE_FILE_EXTENSIONS = ['.png', '.jpg', '.jpeg', '.webp', '.avif'] as const
export const DEFAULT_FILE_CATEGORY_EXTENSIONS = [...IMAGE_FILE_EXTENSIONS, '.pdf'] as const

export interface FileExtensionOption {
  value: string
  label: string
}

export const FILE_EXTENSION_OPTIONS: FileExtensionOption[] = [
  { value: '.png', label: 'PNG' },
  { value: '.jpg', label: 'JPG' },
  { value: '.jpeg', label: 'JPEG' },
  { value: '.webp', label: 'WebP' },
  { value: '.avif', label: 'AVIF' },
  { value: '.pdf', label: 'PDF' },
]

export function normalizeFileExtension(fileName: string): string {
  const extension = fileName.split('.').pop()?.trim().toLowerCase() ?? ''
  return extension ? `.${extension}` : ''
}

export function buildAcceptAttribute(extensions: readonly string[]): string {
  return extensions.join(',')
}
