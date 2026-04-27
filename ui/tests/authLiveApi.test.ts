import { afterEach, describe, expect, it, vi } from 'vitest'
import { GatewayHttpError } from '../src/shared/api/httpClient'
import { createLiveAuthApi } from '../src/features/auth/api/liveAuthApi'

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

describe('liveAuthApi', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('posts login to the gateway auth contract', async () => {
    const fetcher = vi.fn<typeof fetch>().mockResolvedValue(jsonResponse({
      access_token: 'access',
      refresh_token: 'refresh',
      expires_in: 3600,
      token_type: 'Bearer',
      session_id: 'session',
      user: {
        user_id: 'user',
        username: 'alice',
        email: 'alice@example.com',
        role: 'user',
        status: 'active',
      },
      session: {
        session_id: 'session',
        user_id: 'user',
        auth_source: 'password',
        status: 'active',
      },
    }))
    vi.stubGlobal('fetch', fetcher)
    const api = createLiveAuthApi()

    const response = await api.login({
      login_identifier: 'alice@example.com',
      password: 'password',
      client_type: 'web',
    })

    expect(response.access_token).toBe('access')
    const [url, init] = fetcher.mock.calls[0]!
    const headers = init?.headers as Headers

    expect(url).toBe('/api/v3/auth/login')
    expect(init?.method).toBe('POST')
    expect(init?.body).toBe(JSON.stringify({
      login_identifier: 'alice@example.com',
      password: 'password',
      client_type: 'web',
    }))
    expect(headers.get('Accept')).toBe('application/json')
    expect(headers.get('Content-Type')).toBe('application/json')
  })

  it('sends bearer authorization for protected endpoints', async () => {
    const fetcher = vi.fn<typeof fetch>().mockResolvedValue(jsonResponse({
      user: {
        user_id: 'user',
        username: 'alice',
        email: 'alice@example.com',
        role: 'admin',
        status: 'active',
      },
    }))
    vi.stubGlobal('fetch', fetcher)
    const api = createLiveAuthApi()

    await api.me({ accessToken: 'access-token' })

    const [url, init] = fetcher.mock.calls[0]!
    const headers = init?.headers as Headers

    expect(url).toBe('/api/v3/auth/me')
    expect(init?.method).toBe('GET')
    expect(headers.get('Accept')).toBe('application/json')
    expect(headers.get('Authorization')).toBe('Bearer access-token')
  })

  it('normalizes gateway error bodies', async () => {
    const fetcher = vi.fn<typeof fetch>().mockResolvedValue(jsonResponse({
      code: 110201,
      message: 'invalid credentials',
      reference: 'auth.login',
      request_id: 'req_001',
    }, 401))
    vi.stubGlobal('fetch', fetcher)
    const api = createLiveAuthApi()

    await expect(api.login({
      login_identifier: 'alice@example.com',
      password: 'bad-password',
    })).rejects.toMatchObject<Partial<GatewayHttpError>>({
      message: 'invalid credentials',
      status: 401,
      response: {
        code: 110201,
        message: 'invalid credentials',
        reference: 'auth.login',
        request_id: 'req_001',
      },
    })
  })
})
