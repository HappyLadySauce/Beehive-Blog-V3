export type ApiMode = 'mock' | 'live';

function readApiMode(value: string | undefined): ApiMode {
  return value === 'live' ? 'live' : 'mock';
}

export const appConfig = {
  apiMode: readApiMode(import.meta.env.VITE_API_MODE),
  gatewayBaseUrl: import.meta.env.VITE_GATEWAY_BASE_URL ?? 'http://127.0.0.1:8888',
};
