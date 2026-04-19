# Beehive Blog v3 Identity 领域模型

## 1. 目标

本文件定义 `identity` 服务第一阶段的核心领域对象、枚举模型和实体关系，用于指导：

- `identity.proto` 设计
- 数据库表设计
- `services/identity` 代码实现
- `gateway` 的认证上下文对接

## 2. 设计原则

### 2.1 user 是统一主体

平台中的人类用户统一抽象为 `user`。

本地账号与 SSO 身份都归属于同一个 `user`。

### 2.2 凭证与主体分离

`user` 是主体，凭证是附属对象。

因此：

- 本地密码不直接塞到 `user` 聚合根里
- SSO 联邦身份也不直接等同于 `user`

### 2.3 会话与 token 分离

会话是登录上下文，token 是会话凭证。

因此：

- `user_session` 代表一次客户端登录上下文
- `refresh_token` 绑定会话
- access token 是会话的短期派生凭证

## 3. 核心枚举

## 3.1 role

- `guest`
- `member`
- `admin`

说明：

- `guest` 为未登录主体
- `member` 为默认注册用户
- `admin` 为平台管理用户

## 3.2 account_status

- `pending`
- `active`
- `disabled`
- `locked`

说明：

- `pending`：已创建未完成激活
- `active`：正常可用
- `disabled`：停用
- `locked`：临时锁定

## 3.3 auth_source

- `local`
- `sso`

说明：

- `local`：本地账号登录
- `sso`：联邦身份登录

## 3.4 session_status

- `active`
- `revoked`
- `expired`

## 4. 核心实体

## 4.1 user

作用：

- 表示平台内的人类主体

建议字段：

- `id`
- `username`
- `email`
- `nickname`
- `avatar_url`
- `role`
- `status`
- `created_at`
- `updated_at`
- `last_login_at`

约束建议：

- `username` 全局唯一
- `email` 全局唯一，可为空但一旦设置必须唯一
- `role` 默认 `member`
- `status` 默认 `active`

说明：

- 第一阶段不引入复杂用户资料中心
- 扩展资料后续可拆到独立服务或 profile 模块

## 4.2 credential_local

作用：

- 表示本地登录凭证

建议字段：

- `id`
- `user_id`
- `password_hash`
- `password_updated_at`
- `created_at`
- `updated_at`

说明：

- 一个 `user` 在第一阶段最多对应一组本地凭证
- 使用 `username 或 email` 作为登录标识，但密码哈希单独存放

## 4.3 federated_identity

作用：

- 表示来自外部身份提供方的联邦身份

建议字段：

- `id`
- `user_id`
- `provider`
- `subject`
- `email`
- `display_name`
- `raw_profile`
- `created_at`
- `updated_at`
- `last_login_at`

约束建议：

- `(provider, subject)` 必须唯一
- 一个 `user` 可绑定多个联邦身份

## 4.4 user_session

作用：

- 表示一次登录后的客户端会话

建议字段：

- `id`
- `user_id`
- `auth_source`
- `client_type`
- `device_id`
- `device_name`
- `ip_address`
- `user_agent`
- `status`
- `last_seen_at`
- `expires_at`
- `revoked_at`
- `created_at`
- `updated_at`

说明：

- 第一阶段支持多会话并存
- 一个浏览器、Studio 客户端、移动端都可以形成独立会话

## 4.5 refresh_token

作用：

- 表示可持久化、可轮换、可吊销的刷新凭证

建议字段：

- `id`
- `session_id`
- `token_hash`
- `issued_at`
- `expires_at`
- `rotated_from_token_id`
- `revoked_at`
- `created_at`

说明：

- refresh token 应存储哈希值而非明文
- 允许轮换链路
- 按会话或按 token 单独吊销

## 4.6 identity_audit

作用：

- 记录敏感认证行为和身份事件

建议字段：

- `id`
- `user_id`
- `session_id`
- `event_type`
- `auth_source`
- `client_ip`
- `user_agent`
- `result`
- `detail`
- `created_at`

典型事件：

- 注册
- 登录成功
- 登录失败
- 刷新成功
- 刷新失败
- 登出
- 账号禁用
- 账号锁定

## 5. 实体关系

```text
user
 ├─ 1 : 0..1 -> credential_local
 ├─ 1 : N    -> federated_identity
 └─ 1 : N    -> user_session
                  └─ 1 : N -> refresh_token
```

## 6. 默认值建议

新注册本地用户默认：

- `role = member`
- `status = active`

新会话默认：

- `status = active`

新 refresh token 默认：

- 有固定过期时间
- 未吊销

## 7. 模型边界说明

第一阶段不进入 `identity` 的内容：

- 用户隐私设置
- 用户偏好设置
- 用户关注关系
- 用户收藏与订阅
- 用户公开主页扩展资料
- 复杂组织 / 租户模型

这些内容后续如果出现真实需求，再拆到独立服务或独立聚合。

## 8. 当前结论

`identity` 第一阶段的领域模型可以收口为：

**一个 `user` 聚合根，外加本地凭证、联邦身份、用户会话、刷新凭证和身份审计五类核心附属对象。**
