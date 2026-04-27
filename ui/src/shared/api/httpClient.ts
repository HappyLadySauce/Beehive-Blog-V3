import { appConfig } from '@/shared/config/env'

export interface GatewayErrorResponse {
  code: number
  message: string
  reference: string
  request_id: string
}

export class GatewayHttpError extends Error {
  readonly status: number
  readonly response: GatewayErrorResponse | undefined

  constructor(status: number, message: string, response: GatewayErrorResponse | undefined) {
    super(message)
    this.name = 'GatewayHttpError'
    this.status = status
    this.response = response
  }
}

export class GatewayNetworkError extends Error {
  constructor(message = 'Gateway request failed') {
    super(message)
    this.name = 'GatewayNetworkError'
  }
}

export interface RequestJsonOptions extends RequestInit {
  accessToken?: string
  timeoutMs?: number
}

function buildUrl(path: string): string {
  if (/^https?:\/\//i.test(path)) {
    return path
  }
  return `${appConfig.gatewayBaseUrl}${path}`
}

async function parseGatewayError(response: Response): Promise<GatewayErrorResponse | undefined> {
  try {
    return (await response.json()) as GatewayErrorResponse
  } catch {
    return undefined
  }
}

export async function requestJson<TResponse>(
  path: string,
  options: RequestJsonOptions = {},
): Promise<TResponse> {
  const controller = new AbortController()
  const timeoutMs = options.timeoutMs ?? 12_000
  const timeoutId = window.setTimeout(() => controller.abort(), timeoutMs)
  const headers = new Headers(options.headers)
  headers.set('Accept', 'application/json')

  if (options.body !== undefined && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }
  if (options.accessToken) {
    headers.set('Authorization', `Bearer ${options.accessToken}`)
  }

  try {
    const response = await fetch(buildUrl(path), {
      ...options,
      headers,
      signal: options.signal ?? controller.signal,
    })

    if (!response.ok) {
      const parsed = await parseGatewayError(response)
      throw new GatewayHttpError(response.status, parsed?.message ?? response.statusText, parsed)
    }

    if (response.status === 204) {
      return undefined as TResponse
    }

    return (await response.json()) as TResponse
  } catch (error) {
    if (error instanceof GatewayHttpError) {
      throw error
    }
    if (error instanceof DOMException && error.name === 'AbortError') {
      throw new GatewayNetworkError('Gateway request timed out')
    }
    if (error instanceof TypeError) {
      throw new GatewayNetworkError(error.message)
    }
    throw error
  } finally {
    window.clearTimeout(timeoutId)
  }
}
