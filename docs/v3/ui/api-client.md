# UI API Client

## 1. 入口

前端只面向 `gateway` HTTP 契约，不直接调用任何 RPC 服务。

当前入口：

- `ui/src/shared/api/authApi.ts`
- `ui/src/shared/api/contentPreviewApi.ts`
- `ui/src/shared/api/httpClient.ts`
- `ui/src/shared/api/types.ts`

## 2. Mock / Live 切换

环境变量：

```text
VITE_API_MODE=mock|live
VITE_GATEWAY_BASE_URL=
```

默认 `mock`，页面不依赖后端即可运行。
开发环境默认留空 `VITE_GATEWAY_BASE_URL`，让 `/api` 请求走 Vite proxy；部署到已配置 CORS 或同源网关时再填写绝对地址。

`live` 会优先验证 auth 链路：

- `POST /api/v3/auth/register`
- `POST /api/v3/auth/login`
- `GET /api/v3/auth/me`
- `POST /api/v3/auth/logout`

内容预览接口已经保留 gateway adapter：

- `GET /api/v3/public/content/items`
- `GET /api/v3/public/content/items/:slug`
- `GET /api/v3/studio/content/items`

content 服务未就绪时，页面需要展示可理解的错误或空状态；本地开发仍可切回 `mock` 模式独立开发 UI。

## 3. 错误处理

`requestJson` 将非 2xx 响应包装为 `GatewayHttpError`：

- `status`：HTTP 状态码
- `response`：gateway 标准错误体
- `message`：优先使用 gateway message

页面层只展示用户可理解文案，不根据后端字符串做业务分支。

## 4. Token 策略

- `access_token` 保存在 Pinia 内存状态，只用于请求头。
- `refresh_token` 通过 `tokenStorage` 封装后存入 localStorage。
- 页面不直接访问 localStorage。
- logout 无论后端请求是否成功，都会清理本地会话。

后续接入 refresh 接口时，刷新逻辑应收口在 auth store 或 API client adapter 中，不散落到页面。
