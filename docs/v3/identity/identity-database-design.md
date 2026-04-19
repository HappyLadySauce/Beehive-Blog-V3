# Beehive Blog v3 Identity 数据库设计

## 1. 目标

本文件用于定义 `v3` 第一阶段 `identity` 服务的数据库设计基线。

这份设计需要直接支撑以下能力：

- 本地账号注册与登录
- 本地账号与联邦身份并存
- QQ / 微信 / GitHub 第三方登录
- 多会话管理
- refresh token 轮换与吊销
- 认证审计与安全追踪

本文件不是最终 migration SQL，但应足够指导：

- 表结构设计
- 索引设计
- 唯一约束设计
- `identity.proto` 的落地实现
- `services/identity` 的 model 层实现

## 2. 设计原则

### 2.1 user 是统一主体

所有人类用户都统一落到 `users`。

本地账号与第三方登录都只是在“身份来源”层不同，不应拆成两套用户真相。

### 2.2 凭证、联邦身份、会话分离

数据库必须明确区分：

- 主体：`users`
- 本地凭证：`credential_locals`
- 第三方联邦身份：`federated_identities`
- 登录会话：`user_sessions`
- 刷新凭证：`refresh_tokens`

### 2.3 支持 provider 差异，不写死单一 subject 模型

QQ、微信、GitHub 的主体标识并不完全一致：

- QQ 以 `openid` 为核心主体标识
- 微信优先以 `unionid` 归并主体，没有时退回 `openid`
- GitHub 应优先使用 GitHub 用户唯一 ID，不使用用户名 `login` 作为主键

因此联邦身份表不能只设计成一个简单的 `provider + openid` 结构，而必须支持：

- 不同 provider 的主体标识类型
- 不同 provider 的用户资料快照
- 同一个 `user` 绑定多个外部身份

### 2.4 安全优先

第一阶段必须满足：

- 密码仅存哈希
- refresh token 仅存哈希
- 登录 `state` 具备过期与消费状态
- 敏感认证行为可审计
- 账号状态与 token 有联动失效能力

## 3. 统一字段约定

大多数表建议包含：

- `id BIGSERIAL PRIMARY KEY`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`

如表天然不需要 `updated_at`，可按需省略。

布尔或状态字段尽量避免语义重叠，优先使用枚举字符串或受控字符串值。

时间统一采用 `TIMESTAMPTZ`。

## 4. 枚举建议

## 4.1 role

- `member`
- `admin`

说明：

- `guest` 是匿名主体，不落库进 `users`

## 4.2 account_status

- `pending`
- `active`
- `disabled`
- `locked`

## 4.3 auth_source

- `local`
- `sso`

## 4.4 session_status

- `active`
- `revoked`
- `expired`

## 4.5 provider

- `qq`
- `wechat`
- `github`

## 4.6 provider_subject_type

- `openid`
- `unionid`
- `github_user_id`

## 4.7 audit_result

- `success`
- `failure`

## 5. 核心表设计

## 5.1 users

作用：

- 平台内人类用户主体真相

```sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(64) NOT NULL,
  email VARCHAR(320) NULL,
  nickname VARCHAR(128) NULL,
  avatar_url TEXT NULL,
  role VARCHAR(32) NOT NULL,
  status VARCHAR(32) NOT NULL,
  last_login_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

建议约束：

- `UNIQUE (username)`
- 对 `email` 建唯一约束，但允许为空
- `role` 仅允许 `member/admin`
- `status` 仅允许 `pending/active/disabled/locked`

建议索引：

- `UNIQUE (username)`
- `UNIQUE (email)` 或等价唯一索引
- `INDEX (role, status)`

默认值建议：

- `role = member`
- `status = active`

## 5.2 credential_locals

作用：

- 本地账号凭证存储

```sql
CREATE TABLE credential_locals (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  password_hash VARCHAR(255) NOT NULL,
  password_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

建议约束：

- `UNIQUE (user_id)`

说明：

- 一个 `user` 第一阶段最多对应一组本地凭证
- `username` 与 `email` 放在 `users`，密码单独放在凭证表
- 不额外存密码明文、可逆密文或弱摘要

## 5.3 federated_identities

作用：

- 存储第三方登录身份绑定

这是整个第三方登录设计里最关键的一张表。

```sql
CREATE TABLE federated_identities (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  provider VARCHAR(32) NOT NULL,
  provider_subject VARCHAR(255) NOT NULL,
  provider_subject_type VARCHAR(64) NOT NULL,
  unionid VARCHAR(255) NULL,
  openid VARCHAR(255) NULL,
  provider_login VARCHAR(255) NULL,
  provider_email VARCHAR(320) NULL,
  provider_display_name VARCHAR(255) NULL,
  avatar_url TEXT NULL,
  app_id_or_client_id VARCHAR(128) NULL,
  access_scope TEXT NULL,
  raw_profile JSONB NULL,
  last_login_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

建议约束：

- `UNIQUE (provider, provider_subject)`

建议索引：

- `INDEX (user_id)`
- `INDEX (provider)`
- 微信场景可额外增加 `INDEX (provider, unionid)` 或唯一策略

字段语义建议：

- `provider_subject`
  - QQ：存 `openid`
  - 微信：优先存 `unionid`，拿不到时存 `openid`
  - GitHub：存 GitHub 用户唯一 ID
- `provider_subject_type`
  - 用于明确 `provider_subject` 的语义
- `unionid`
  - 主要给微信使用
- `openid`
  - 给 QQ / 微信保留
- `provider_login`
  - 例如 GitHub `login`
- `provider_email`
  - 第三方返回的邮箱，不能直接替代平台主邮箱真相
- `raw_profile`
  - 保存 provider 原始资料快照，便于排障和后续扩展

### 5.3.1 provider 特殊约束建议

QQ：

- `provider = qq`
- `provider_subject = openid`
- `provider_subject_type = openid`

微信：

- 如果拿到 `unionid`
  - `provider_subject = unionid`
  - `provider_subject_type = unionid`
- 如果只拿到 `openid`
  - `provider_subject = openid`
  - `provider_subject_type = openid`
- 表结构必须预留后续从 `openid` 归并到 `unionid` 的能力

GitHub：

- `provider = github`
- `provider_subject = GitHub 用户唯一 ID`
- `provider_subject_type = github_user_id`
- 不要用 `provider_login` 作为唯一主键

## 5.4 oauth_login_states

作用：

- 存储第三方登录发起态，用于回调校验、防重放与一次性消费

```sql
CREATE TABLE oauth_login_states (
  id BIGSERIAL PRIMARY KEY,
  provider VARCHAR(32) NOT NULL,
  state VARCHAR(512) NOT NULL,
  redirect_uri TEXT NOT NULL,
  client_type VARCHAR(32) NULL,
  device_id VARCHAR(128) NULL,
  code_verifier VARCHAR(255) NULL,
  requested_scopes TEXT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  consumed_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

建议约束：

- `UNIQUE (provider, state)`

建议索引：

- `INDEX (expires_at)`
- `INDEX (consumed_at)`

说明：

- 不建议只把 `state` 放在内存里
- 这张表对 QQ / 微信 / GitHub 都有价值
- `code_verifier` 用于将来接需要 PKCE 的 provider 或客户端

## 5.5 user_sessions

作用：

- 存储用户登录会话

```sql
CREATE TABLE user_sessions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  auth_source VARCHAR(32) NOT NULL,
  client_type VARCHAR(32) NULL,
  device_id VARCHAR(128) NULL,
  device_name VARCHAR(128) NULL,
  ip_address INET NULL,
  user_agent TEXT NULL,
  status VARCHAR(32) NOT NULL,
  last_seen_at TIMESTAMPTZ NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

建议索引：

- `INDEX (user_id, status)`
- `INDEX (expires_at)`
- `INDEX (device_id)`

说明：

- 第一阶段支持多会话并存
- 不强制单设备唯一
- `ip_address` 建议存服务端识别出的真实客户端 IP，不信任客户端自报

## 5.6 refresh_tokens

作用：

- 存储 refresh token 轮换链

```sql
CREATE TABLE refresh_tokens (
  id BIGSERIAL PRIMARY KEY,
  session_id BIGINT NOT NULL REFERENCES user_sessions(id),
  token_hash VARCHAR(255) NOT NULL,
  issued_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  rotated_from_token_id BIGINT NULL REFERENCES refresh_tokens(id),
  revoked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

建议约束：

- `UNIQUE (token_hash)`

建议索引：

- `INDEX (session_id)`
- `INDEX (expires_at)`
- `INDEX (revoked_at)`
- `INDEX (rotated_from_token_id)`

说明：

- 只存哈希，不存明文 refresh token
- 支持轮换链追踪
- 支持按 token 或按 session 吊销

## 5.7 identity_audits

作用：

- 存储敏感身份行为审计日志

```sql
CREATE TABLE identity_audits (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NULL REFERENCES users(id),
  session_id BIGINT NULL REFERENCES user_sessions(id),
  provider VARCHAR(32) NULL,
  auth_source VARCHAR(32) NULL,
  event_type VARCHAR(64) NOT NULL,
  result VARCHAR(32) NOT NULL,
  client_ip INET NULL,
  user_agent TEXT NULL,
  detail JSONB NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

建议索引：

- `INDEX (user_id, created_at DESC)`
- `INDEX (session_id, created_at DESC)`
- `INDEX (event_type, created_at DESC)`
- `INDEX (provider, created_at DESC)`

典型事件：

- `register_local_success`
- `register_local_failed`
- `login_local_success`
- `login_local_failed`
- `login_sso_success`
- `login_sso_failed`
- `refresh_success`
- `refresh_failed`
- `logout_success`
- `introspect_failed`
- `account_disabled`
- `account_locked`

## 6. 表关系

```text
users
 ├─ 1 : 0..1 -> credential_locals
 ├─ 1 : N    -> federated_identities
 ├─ 1 : N    -> user_sessions
 │                 └─ 1 : N -> refresh_tokens
 └─ 1 : N    -> identity_audits

oauth_login_states
  独立存在，用于 SSO 登录过程态管理
```

## 7. 第一阶段索引与约束最小集合

必须具备的唯一约束：

- `users.username`
- `users.email`
- `credential_locals.user_id`
- `federated_identities(provider, provider_subject)`
- `oauth_login_states(provider, state)`
- `refresh_tokens.token_hash`

必须具备的关键索引：

- `user_sessions(user_id, status)`
- `refresh_tokens(session_id)`
- `identity_audits(user_id, created_at DESC)`
- `identity_audits(event_type, created_at DESC)`

## 8. 与 identity.proto 的映射关系

`RegisterLocalUser`
- 写入：
  - `users`
  - `credential_locals`
  - `user_sessions`
  - `refresh_tokens`
  - `identity_audits`

`LoginLocalUser`
- 读取：
  - `users`
  - `credential_locals`
- 写入：
  - `user_sessions`
  - `refresh_tokens`
  - `identity_audits`

`StartSsoLogin`
- 写入：
  - `oauth_login_states`

`FinishSsoLogin`
- 读取：
  - `oauth_login_states`
  - `federated_identities`
  - `users`
- 写入：
  - `users`（首次登录时可能新建）
  - `federated_identities`
  - `user_sessions`
  - `refresh_tokens`
  - `identity_audits`

`RefreshSessionToken`
- 读取：
  - `refresh_tokens`
  - `user_sessions`
  - `users`
- 写入：
  - 新 `refresh_tokens`
  - 旧 token 吊销信息
  - `identity_audits`

`LogoutSession`
- 写入：
  - `user_sessions.revoked_at/status`
  - 对应 `refresh_tokens.revoked_at`
  - `identity_audits`

`GetCurrentUser`
- 读取：
  - `users`

`IntrospectAccessToken`
- 读取：
  - `users`
  - `user_sessions`

## 9. 默认值与实现建议

默认值建议：

- 新注册用户：
  - `role = member`
  - `status = active`
- 新会话：
  - `status = active`

实现建议：

- 使用数据库约束兜底唯一性，不只靠应用层校验
- 邮箱如允许为空，唯一约束要兼容空值策略
- `raw_profile` 与 `detail` 使用 `JSONB`
- `ip_address` 使用 `INET`
- `refresh_tokens` 与 `oauth_login_states` 要做定期清理

## 10. 第一阶段故意不做的数据库能力

当前阶段不纳入：

- 多租户身份隔离
- 用户组织 / 团队模型
- 第三方 provider token 长期持久化策略
- provider scope 精细授权矩阵
- “登出全部设备” 的专用批处理表
- OAuth provider 配置中心数据库化

这些能力后续再按真实需求扩展，不在第一阶段把表结构做重。

## 11. 当前结论

`identity` 第一阶段的数据库可以收口为：

**以 `users` 为统一主体，以本地凭证、联邦身份、OAuth 登录状态、用户会话、刷新凭证和审计日志为六类核心支撑表。**  
其中针对 QQ、微信、GitHub 的差异，重点通过 `federated_identities` 的 `provider + provider_subject + provider_subject_type` 模型来吸收，而不是为每个 provider 单独建一套用户表。
