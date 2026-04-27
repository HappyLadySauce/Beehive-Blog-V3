import { fileURLToPath, URL } from 'node:url'
import UnoCSS from 'unocss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [vue(), UnoCSS()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    host: '127.0.0.1',
    port: 5173,
    strictPort: false,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8888',
        changeOrigin: true,
      },
      '/healthz': {
        target: 'http://127.0.0.1:8888',
        changeOrigin: true,
      },
      '/readyz': {
        target: 'http://127.0.0.1:8888',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://127.0.0.1:8888',
        changeOrigin: true,
        ws: true,
      },
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./tests/setup.ts'],
    exclude: ['e2e/**', 'node_modules/**', 'dist/**'],
  },
})
