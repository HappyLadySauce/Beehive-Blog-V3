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
VITE_GATEWAY_BASE_URL=
```

`VITE_API_MODE=live` 优先用于 auth gateway 联调；内容预览也具备 live adapter，但服务未就绪时页面会降级展示。
开发环境默认留空 `VITE_GATEWAY_BASE_URL`，让 `/api` 请求走 Vite proxy；部署到已配置 CORS 或同源网关时再填写绝对地址。
