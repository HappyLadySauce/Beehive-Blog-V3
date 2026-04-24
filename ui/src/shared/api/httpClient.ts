import { appConfig } from '@/shared/config/env';
import type { GatewayErrorResponse } from '@/shared/api/types';

export class GatewayHttpError extends Error {
  readonly status: number;
  readonly response: GatewayErrorResponse | undefined;

  constructor(status: number, message: string, response: GatewayErrorResponse | undefined) {
    super(message);
    this.name = 'GatewayHttpError';
    this.status = status;
    this.response = response;
  }
}

interface RequestOptions extends RequestInit {
  accessToken?: string;
}

function buildUrl(path: string): string {
  if (path.startsWith('http')) {
    return path;
  }
  return `${appConfig.gatewayBaseUrl}${path}`;
}

async function parseGatewayError(response: Response): Promise<GatewayErrorResponse | undefined> {
  try {
    return (await response.json()) as GatewayErrorResponse;
  } catch {
    return undefined;
  }
}

export async function requestJson<TResponse>(path: string, options: RequestOptions = {}): Promise<TResponse> {
  const headers = new Headers(options.headers);
  headers.set('Accept', 'application/json');

  if (options.body !== undefined && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }
  if (options.accessToken) {
    headers.set('Authorization', `Bearer ${options.accessToken}`);
  }

  const response = await fetch(buildUrl(path), {
    ...options,
    headers,
  });

  if (!response.ok) {
    const parsed = await parseGatewayError(response);
    throw new GatewayHttpError(response.status, parsed?.message ?? response.statusText, parsed);
  }

  if (response.status === 204) {
    return undefined as TResponse;
  }

  return (await response.json()) as TResponse;
}
