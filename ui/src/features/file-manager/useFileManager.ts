import { computed, reactive, shallowRef, watch } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { useConfirm, useProgressiveQuery, useToast } from '@/shared/composables'

import { deleteFileAsset, getFileAsset, listFileAssets } from './api'
import type { FileAssetListParams, FileAssetSummary, FileUploadNamespace } from './types'
import { useFileUpload } from './useFileUpload'

export function useFileManager() {
  const authStore = useAuthStore()
  const { confirm } = useConfirm()
  const { pushToast } = useToast()
  const { isUploading, errorMessage, uploadFile } = useFileUpload()

  const filters = reactive<FileAssetListParams>({
    keyword: '',
    namespace: '',
    status: 'uploaded',
    visibility: '',
    owner_user_id: '',
    page: 1,
    page_size: 20,
  })
  const selectedAssetId = shallowRef('')
  const isDeleting = shallowRef(false)
  const uploadNamespace = shallowRef<FileUploadNamespace>('content_image')

  const listQuery = useProgressiveQuery({
    queryKey: computed(() => ['studio-files', { ...filters }]),
    queryFn: () => listFileAssets(filters, { accessToken: authStore.accessToken }),
  })

  const detailQuery = useProgressiveQuery({
    queryKey: computed(() => ['studio-file-detail', selectedAssetId.value]),
    queryFn: () => getFileAsset(selectedAssetId.value, { accessToken: authStore.accessToken }),
    enabled: computed(() => Boolean(selectedAssetId.value)),
  })

  const items = computed(() => listQuery.data.value?.items ?? [])
  const total = computed(() => listQuery.data.value?.total ?? 0)
  const selectedAsset = computed(() => detailQuery.data.value?.asset ?? null)
  const page = computed(() => listQuery.data.value?.page ?? Number(filters.page ?? 1))
  const pageSize = computed(() => listQuery.data.value?.page_size ?? Number(filters.page_size ?? 20))

  watch(
    () => [filters.keyword, filters.namespace, filters.status, filters.visibility, filters.owner_user_id],
    () => {
      filters.page = 1
    },
  )

  function setPage(nextPage: number): void {
    filters.page = nextPage
  }

  function setPageSize(nextPageSize: number): void {
    filters.page_size = nextPageSize
    filters.page = 1
  }

  function openAsset(asset: FileAssetSummary): void {
    selectedAssetId.value = asset.asset_id
  }

  function closeAsset(): void {
    selectedAssetId.value = ''
  }

  async function uploadSelectedFile(file: File): Promise<void> {
    const namespace = uploadNamespace.value
    await uploadFile(file, authStore.accessToken, namespace)
    pushToast({
      tone: 'success',
      title: 'File uploaded',
      message: file.name,
    })
    await listQuery.refetch()
  }

  async function removeAsset(asset: FileAssetSummary): Promise<void> {
    const approved = await confirm({
      title: 'Delete file asset?',
      message: asset.file_name,
      confirmText: 'Delete',
      tone: 'danger',
    })
    if (!approved) {
      return
    }

    isDeleting.value = true
    try {
      await deleteFileAsset(asset.asset_id, { accessToken: authStore.accessToken })
      if (selectedAssetId.value === asset.asset_id) {
        closeAsset()
      }
      await listQuery.refetch()
      pushToast({
        tone: 'success',
        title: 'File deleted',
        message: asset.file_name,
      })
    } catch (error) {
      pushToast({
        tone: 'danger',
        title: 'Operation failed',
        message: error instanceof Error ? error.message : 'Unable to delete file asset.',
      })
      throw error
    } finally {
      isDeleting.value = false
    }
  }

  return {
    filters,
    items,
    total,
    page,
    pageSize,
    selectedAssetId,
    selectedAsset,
    uploadNamespace,
    isUploading,
    uploadErrorMessage: errorMessage,
    isDeleting,
    listQuery,
    detailQuery,
    setPage,
    setPageSize,
    openAsset,
    closeAsset,
    uploadSelectedFile,
    removeAsset,
  }
}
