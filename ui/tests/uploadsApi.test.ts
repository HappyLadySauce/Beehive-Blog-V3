import { afterEach, describe, expect, it, vi } from 'vitest'

import type { FileUploadCreateResponse } from '@/features/uploads/types'

function imageFile(): File {
  return new File(['data'], 'avatar.png', { type: 'image/png' })
}

function uploadResponse(): FileUploadCreateResponse {
  return {
    asset: {
      asset_id: 'asset_1',
      upload_id: 'upload_1',
      owner_user_id: 'user_1',
      scope: 'avatar',
      visibility: 'public',
      status: 'pending',
      bucket: 'local',
      object_key: 'avatars/user_1/avatar.png',
      public_url: '',
      file_name: 'avatar.png',
      content_type: 'image/png',
      byte_size: 4,
      created_at: 1,
      expires_at: 2,
    },
    upload_url: 'http://127.0.0.1:8084/files/uploads/upload_1',
    headers: {
      'Content-Type': 'image/png',
      'X-Upload-Token': 'signed-token',
    },
    expires_at: 2,
    max_bytes: 2 * 1024 * 1024,
  }
}

describe('uploads api', () => {
  afterEach(() => {
    vi.useRealTimers()
    vi.unstubAllEnvs()
    vi.unstubAllGlobals()
    vi.resetModules()
  })

  it('mock avatar upload does not require a token or call the data plane', async () => {
    vi.stubEnv('VITE_API_MODE', 'mock')
    const fetcher = vi.fn<typeof fetch>()
    vi.stubGlobal('fetch', fetcher)
    Object.defineProperty(URL, 'createObjectURL', {
      configurable: true,
      value: vi.fn(() => 'blob:mock-avatar'),
    })
    const { useAvatarUpload } = await import('@/features/uploads/useAvatarUpload')

    const { uploadAvatar } = useAvatarUpload()
    const publicUrl = await uploadAvatar(imageFile())

    expect(publicUrl).toBe('blob:mock-avatar')
    expect(fetcher).not.toHaveBeenCalled()
  })

  it('mock upload clears completed state and revokes the previous completed preview URL', async () => {
    vi.stubEnv('VITE_API_MODE', 'mock')
    const createObjectURL = vi.fn()
      .mockReturnValueOnce('blob:first-avatar')
      .mockReturnValueOnce('blob:second-avatar')
    const revokeObjectURL = vi.fn()
    Object.defineProperty(URL, 'createObjectURL', {
      configurable: true,
      value: createObjectURL,
    })
    Object.defineProperty(URL, 'revokeObjectURL', {
      configurable: true,
      value: revokeObjectURL,
    })
    const { completeFileUpload, createFileUpload, putFileUploadObject } = await import('@/features/uploads/api')

    const firstUpload = await createFileUpload({
      scope: 'avatar',
      file_name: 'first.png',
      content_type: 'image/png',
      byte_size: 4,
      visibility: 'public',
    })
    await putFileUploadObject(firstUpload, imageFile())
    const firstCompleted = await completeFileUpload(firstUpload.asset.upload_id)
    const repeatedFirst = await completeFileUpload(firstUpload.asset.upload_id)

    expect(firstCompleted.asset.public_url).toBe('blob:first-avatar')
    expect(repeatedFirst.asset.public_url).not.toBe('blob:first-avatar')

    const secondUpload = await createFileUpload({
      scope: 'avatar',
      file_name: 'second.png',
      content_type: 'image/png',
      byte_size: 4,
      visibility: 'public',
    })
    await putFileUploadObject(secondUpload, imageFile())
    expect(revokeObjectURL).not.toHaveBeenCalled()
    const secondCompleted = await completeFileUpload(secondUpload.asset.upload_id)

    expect(secondCompleted.asset.public_url).toBe('blob:second-avatar')
    expect(createObjectURL).toHaveBeenCalledTimes(2)
    expect(revokeObjectURL).toHaveBeenCalledTimes(1)
    expect(revokeObjectURL).toHaveBeenCalledWith('blob:first-avatar')
  })

  it('live object upload forwards presigned headers', async () => {
    vi.stubEnv('VITE_API_MODE', 'live')
    const fetcher = vi.fn<typeof fetch>().mockResolvedValue(new Response(null, { status: 204 }))
    vi.stubGlobal('fetch', fetcher)
    const { putFileUploadObject } = await import('@/features/uploads/api')

    await putFileUploadObject(uploadResponse(), imageFile())

    expect(fetcher).toHaveBeenCalledTimes(1)
    const [url, init] = fetcher.mock.calls[0]!
    expect(url).toBe('http://127.0.0.1:8084/files/uploads/upload_1')
    expect(init?.method).toBe('PUT')
    expect(init?.headers).toEqual({
      'Content-Type': 'image/png',
      'X-Upload-Token': 'signed-token',
    })
  })

  it('live object upload surfaces a timeout when the data plane does not respond', async () => {
    vi.useFakeTimers()
    vi.stubEnv('VITE_API_MODE', 'live')
    const fetcher = vi.fn<typeof fetch>((_input, init) => new Promise((_, reject) => {
      init?.signal?.addEventListener('abort', () => {
        reject(new DOMException('Aborted', 'AbortError'))
      })
    }))
    vi.stubGlobal('fetch', fetcher)
    const { putFileUploadObject } = await import('@/features/uploads/api')

    const uploadPromise = expect(putFileUploadObject(uploadResponse(), imageFile())).rejects.toThrow('File upload timed out.')
    await vi.advanceTimersByTimeAsync(60_000)

    await uploadPromise
    expect(fetcher).toHaveBeenCalledTimes(1)
  })
})
