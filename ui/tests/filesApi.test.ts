import { afterEach, describe, expect, it, vi } from 'vitest'

function imageFile(name = 'studio.png', size = 4, type = 'image/png'): File {
  return new File([new Uint8Array(size)], name, { type })
}

describe('files api', () => {
  afterEach(() => {
    vi.unstubAllEnvs()
    vi.unstubAllGlobals()
    vi.resetModules()
  })

  it('live list/get/delete requests target gateway file asset routes', async () => {
    vi.stubEnv('VITE_API_MODE', 'live')
    const fetcher = vi.fn<typeof fetch>()
      .mockResolvedValueOnce(new Response(JSON.stringify({ items: [], total: 0, page: 1, page_size: 20 }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ asset: { asset_id: 'asset_1' } }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ ok: true }), { status: 200 }))
    vi.stubGlobal('fetch', fetcher)
    const { deleteFileAsset, getFileAsset, listFileAssets } = await import('@/features/files/api')

    await listFileAssets({ keyword: 'avatar', status: 'uploaded', page: 2, page_size: 10 }, { accessToken: 'access-token' })
    await getFileAsset('asset_1', { accessToken: 'access-token' })
    await deleteFileAsset('asset_1', { accessToken: 'access-token' })

    expect(fetcher).toHaveBeenCalledTimes(3)
    expect(fetcher.mock.calls[0]?.[0]).toBe('/api/v3/files/assets?keyword=avatar&status=uploaded&page=2&page_size=10')
    expect(fetcher.mock.calls[1]?.[0]).toBe('/api/v3/files/assets/asset_1')
    expect(fetcher.mock.calls[2]?.[0]).toBe('/api/v3/files/assets/asset_1')
  })

  it('mock uploads become queryable file assets and can be soft deleted from the default uploaded view', async () => {
    vi.stubEnv('VITE_API_MODE', 'mock')
    Object.defineProperty(URL, 'createObjectURL', {
      configurable: true,
      value: vi.fn(() => 'blob:studio-file'),
    })

    const { useFileUpload } = await import('@/features/uploads/useFileUpload')
    const { deleteFileAsset, getFileAsset, listFileAssets } = await import('@/features/files/api')
    const { uploadFile } = useFileUpload()

    await uploadFile(imageFile(), undefined, 'content_image')
    const uploaded = await listFileAssets({ status: 'uploaded' })

    expect(uploaded.total).toBe(1)
    expect(uploaded.items[0]?.file_name).toBe('studio.png')

    const detail = await getFileAsset(uploaded.items[0]!.asset_id)
    expect(detail.asset.public_url).toBe('blob:studio-file')

    await deleteFileAsset(uploaded.items[0]!.asset_id)

    const uploadedAfterDelete = await listFileAssets({ status: 'uploaded' })
    const deleted = await listFileAssets({ status: 'deleted' })
    expect(uploadedAfterDelete.total).toBe(0)
    expect(deleted.total).toBe(1)
    expect(deleted.items[0]?.status).toBe('deleted')
  })
})
