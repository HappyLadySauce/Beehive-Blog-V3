# Beehive Blog v3 UI

Vue3 前端工程首版默认使用 mock API，可在不启动后端服务的情况下开发页面壳。

## Commands

```powershell
pnpm install
pnpm dev
pnpm typecheck
pnpm lint
pnpm test
pnpm build
pnpm test:e2e
```

## Environment

复制 `.env.example` 为 `.env.local` 后按需调整：

```text
VITE_API_MODE=mock
VITE_GATEWAY_BASE_URL=http://127.0.0.1:8888
```

`VITE_API_MODE=live` 仅用于 auth gateway 联调；content 首版仍使用 mock。
