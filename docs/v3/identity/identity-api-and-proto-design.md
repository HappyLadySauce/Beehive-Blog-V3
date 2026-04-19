# Beehive Blog v3 Identity API 与 Proto 设计

## 1. 目标

本文件定义 `identity` 第一阶段的 RPC 契约目标，并据此反推 `gateway.api` 第一批 `/auth/*` 接口。

当前第三阶段已经收口为：

- 本地认证主链路完整实现
- SSO 继续采用统一抽象
- 实际对外开放的 SSO provider 仅为 `GitHub`
- `QQ/WeChat` 当前仅保留建模、配置与授权地址辅助，不开放登录入口

## 2. 设计原则

### 2.1 proto 只描述 identity 的真能力

`identity.proto` 只定义 `identity` 自己负责的能力：

- 注册
- 登录
- SSO 登录
- 刷新
- 登出
- 当前用户
- token introspection

### 2.2 gateway.api 只表达对外 HTTP 契约

`gateway.api` 直接映射 `identity` 提供的能力。

`gateway` 负责：

- 请求参数绑定
- token 提取
- RPC 调用
- 错误包装

`gateway` 不负责：

- 授权裁决
- 业务编排

### 2.3 SSO 采用两步抽象

proto 层继续抽象为：

- `StartSsoLogin`
- `FinishSsoLogin`

这样后续接其他 provider 时，无需修改 proto 结构。

### 2.4 第三阶段只开放 GitHub SSO

当前实现状态固定为：

- `GitHub`：完整实现并开放
- `QQ`：仅保留抽象、配置与授权地址辅助，不开放登录入口
- `WeChat`：仅保留抽象、配置与授权地址辅助，不开放登录入口

实现行为固定为：

- `StartSsoLogin` 对 `QQ/WeChat` 返回 `sso_provider_not_ready`
- `FinishSsoLogin` 当前只完成 `GitHub` 的回调交换与平台会话建立

## 3. 第一版 RPC 集合

### 3.1 `RegisterLocalUser`

作用：

- 创建本地账号用户
- 注册成功后直接建立平台会话

输入语义：

- `username`
- `email`
- `password`
- `nickname`

输出语义：

- `current_user`
- `token_pair`
- `session_info`

### 3.2 `LoginLocalUser`

作用：

- 本地账号登录

输入语义：

- `login_identifier`
- `password`
- `client_type`
- `device_id`
- `device_name`
- `user_agent`

说明：

- `login_identifier` 允许传 `username` 或 `email`
- 终端真实 IP 由 `gateway` 从连接 / 可信代理头解析后，通过 gRPC metadata 传给 `identity`

输出语义：

- `token_pair`
- `current_user`
- `session_info`

### 3.3 `StartSsoLogin`

作用：

- 发起 SSO 登录流程
- 第三阶段当前仅真正开放 `GitHub`

输入语义：

- `provider`
- `redirect_uri`
- `state`

输出语义：

- `provider`
- `auth_url`
- `state`

### 3.4 `FinishSsoLogin`

作用：

- 完成 SSO 回调交换与用户会话建立
- 第三阶段当前仅真正完成 `GitHub`

输入语义：

- `provider`
- `code`
- `state`
- `redirect_uri`
- `client_type`
- `device_id`
- `device_name`
- `user_agent`

说明：

- 终端真实 IP 由 `gateway` 从连接 / 可信代理头解析后，通过 gRPC metadata 传给 `identity`

输出语义：

- `token_pair`
- `current_user`
- `session_info`

### 3.5 `RefreshSessionToken`

作用：

- 使用 refresh token 刷新 access token

输入语义：

- `refresh_token`
- `user_agent`

输出语义：

- 新的 `token_pair`
- `session_info`

约束：

- refresh token 采用轮换策略
- 若会话已吊销、账号不可用或 token 过期，应返回明确错误

### 3.6 `LogoutSession`

作用：

- 登出并吊销当前会话

输入语义：

- `session_id`
- 可选 `refresh_token`

输出语义：

- 操作结果

说明：

- `session_id` 须由可信调用方在调用 RPC 前从 access token 上下文中解析并填入

### 3.7 `GetCurrentUser`

作用：

- 获取当前登录用户基础信息

输入语义：

- `user_id`

输出语义：

- `current_user`

### 3.8 `IntrospectAccessToken`

作用：

- 统一解析并校验 access token

调用方：

- `gateway`
- 其他内部服务

输入语义：

- `access_token`

输出语义：

- `active`
- `user_id`
- `role`
- `account_status`
- `session_id`
- `auth_source`
- `expires_at`

说明：

- 这是内部认证基础能力，不直接作为公网 HTTP 接口暴露
- `gateway` 依赖该 RPC 完成 access token 的标准化校验

## 4. 第一版核心消息模型

### 4.1 `Role`

- `ROLE_GUEST`
- `ROLE_MEMBER`
- `ROLE_ADMIN`

### 4.2 `AccountStatus`

- `ACCOUNT_STATUS_PENDING`
- `ACCOUNT_STATUS_ACTIVE`
- `ACCOUNT_STATUS_DISABLED`
- `ACCOUNT_STATUS_LOCKED`

### 4.3 `AuthSource`

- `AUTH_SOURCE_LOCAL`
- `AUTH_SOURCE_SSO`

### 4.4 `TokenPair`

字段：

- `access_token`
- `refresh_token`
- `expires_in`
- `session_id`
- `token_type`

### 4.5 `CurrentUser`

字段：

- `user_id`
- `username`
- `email`
- `nickname`
- `avatar_url`
- `role`
- `status`

### 4.6 `SessionInfo`

字段：

- `session_id`
- `user_id`
- `auth_source`
- `client_type`
- `device_id`
- `device_name`
- `status`
- `last_seen_at`
- `expires_at`

### 4.7 `FederatedIdentity`

字段：

- `provider`
- `subject`
- `email`
- `display_name`

## 5. 错误语义建议

第一阶段至少统一以下错误语义：

- `invalid_credentials`
- `account_disabled`
- `account_locked`
- `account_pending`
- `invalid_token`
- `token_expired`
- `session_revoked`
- `refresh_token_expired`
- `sso_provider_not_supported`
- `sso_provider_not_ready`
- `sso_state_invalid`

## 6. 反推 gateway.api

基于 `identity` 第一版 RPC，`gateway.api` 第一批应暴露：

- `POST /api/v3/auth/register`
- `POST /api/v3/auth/login`
- `POST /api/v3/auth/sso/start`
- `POST /api/v3/auth/sso/callback`
- `POST /api/v3/auth/refresh`
- `POST /api/v3/auth/logout`
- `GET /api/v3/auth/me`

映射关系：

- `/auth/register` -> `RegisterLocalUser`
- `/auth/login` -> `LoginLocalUser`
- `/auth/sso/start` -> `StartSsoLogin`
- `/auth/sso/callback` -> `FinishSsoLogin`
- `/auth/refresh` -> `RefreshSessionToken`
- `/auth/logout` -> `LogoutSession`
- `/auth/me` -> `GetCurrentUser`

`IntrospectAccessToken` 不直接映射公网路由，由 `gateway` 中间件或内部 handler 使用。

## 7. 当前结论

`identity.proto` 第一版围绕“本地账号 + SSO 并存、多会话、可吊销 refresh token、标准化 token introspection”来设计。  
第三阶段实际对外开放的 SSO provider 为 `GitHub`，`QQ/WeChat` 仅保留抽象与模型准备。  
`gateway.api` 第一版则只需要把这组能力清晰地映射成 `/api/v3/auth/*` 对外接口。
