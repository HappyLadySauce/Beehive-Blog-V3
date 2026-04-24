import { fileURLToPath, URL } from 'node:url';

import { defineConfig, devices } from '@playwright/test';

const repoRoot = fileURLToPath(new URL('../', import.meta.url));
process.env.NO_PROXY = '127.0.0.1,localhost';
process.env.no_proxy = '127.0.0.1,localhost';

const pnpmCommand =
  process.platform === 'win32'
    ? 'pnpm.cmd --dir ui dev'
    : 'pnpm --dir ui dev';

export default defineConfig({
  testDir: './e2e',
  outputDir: './test-results',
  reporter: [['list'], ['html', { open: 'never' }]],
  webServer: {
    command: pnpmCommand,
    cwd: repoRoot,
    url: 'http://127.0.0.1:5173',
    reuseExistingServer: false,
    timeout: 120_000,
  },
  use: {
    baseURL: 'http://127.0.0.1:5173',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'desktop',
      use: { ...devices['Desktop Chrome'], viewport: { width: 1440, height: 900 } },
    },
    {
      name: 'tablet',
      use: { ...devices['iPad Pro 11'], viewport: { width: 768, height: 1024 } },
    },
    {
      name: 'mobile',
      use: { ...devices['Pixel 7'], viewport: { width: 390, height: 844 } },
    },
  ],
});
