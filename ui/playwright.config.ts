import { dirname } from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig, devices } from '@playwright/test'

const currentDir = dirname(fileURLToPath(import.meta.url))
const localNoProxy = '127.0.0.1,localhost,::1'

process.env.NO_PROXY = [process.env.NO_PROXY, localNoProxy].filter(Boolean).join(',')
process.env.no_proxy = [process.env.no_proxy, localNoProxy].filter(Boolean).join(',')

export default defineConfig({
  testDir: './e2e',
  timeout: 30_000,
  workers: 1,
  expect: {
    timeout: 5_000,
  },
  use: {
    baseURL: 'http://127.0.0.1:5174',
    trace: 'on-first-retry',
  },
  webServer: {
    command: 'pnpm.cmd dev --host 127.0.0.1 --port 5174',
    cwd: currentDir,
    url: 'http://127.0.0.1:5174',
    reuseExistingServer: false,
    timeout: 120_000,
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'mobile-chrome',
      use: { ...devices['Pixel 5'] },
    },
  ],
})
