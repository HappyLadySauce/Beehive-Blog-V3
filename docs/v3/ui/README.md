# Beehive Blog v3 UI 设计基线

## 1. 目标

`ui/` 是 Beehive Blog v3 的 Web UI 客户端工程，承载两个产品面：

- `Public Web`：面向访客的公开内容、项目、经历与个人表达入口。
- `Studio`：面向创作者的内容生产、管理、审阅与发布工作台。

首版目标是先完成良好工程化骨架、页面壳、响应式布局和 auth API client。`content` 服务仍在开发中，前端只保留契约形态类型与 mock 数据，不进行真实 content 联调。

## 2. 技术栈

- Vue 3 + Vite + TypeScript
- Vue Router：Public/Auth/Studio 路由分区
- Pinia：认证会话与后续跨页状态
- UnoCSS：设计令牌、原子样式和响应式工具
- Vitest：单元测试
- Playwright：关键页面响应式 smoke test
- pnpm：前端包管理

## 3. 模块边界

```text
ui/
  src/
    app/                # 应用入口、router、全局样式
    pages/              # 路由页面壳
      public/
      auth/
      studio/
    features/           # 业务能力模块
      auth/
      content-preview/
      navigation/
    shared/             # 共享基础设施
      api/
      components/
      config/
      layouts/
      storage/
```

约束：

- 页面只编排布局和交互，不直接写底层 fetch。
- `features/*` 承载业务域状态、组合组件和领域用例。
- `shared/api` 是 gateway HTTP 契约的唯一前端访问入口。
- `shared/components` 只放可复用基础 UI，不放业务逻辑。
- `content-preview` 首版只读 mock 数据，不访问真实 content 服务。

## 4. 运行方式

```powershell
pnpm --dir ui install
pnpm --dir ui dev
```

默认 mock：

```text
VITE_API_MODE=mock
```

auth 实联调：

```text
VITE_API_MODE=live
VITE_GATEWAY_BASE_URL=http://127.0.0.1:8888
```

live 模式只要求 `gateway` 和 `identity` auth 链路可用；不要求 content ready。
