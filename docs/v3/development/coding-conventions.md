# Beehive Blog v3 编码规范

## 目标

统一 `v3` 的服务分层、目录职责、基础设施边界与公共代码组织方式，降低实现分歧和技术债务。

## 目录职责

### `services/<svc>/internal/config`

- 只放配置结构定义
- 只负责表达配置字段、tag 约束和 `Validate()` 校验入口
- 不在这里初始化数据库、Redis、MQ、RPC client

### `services/<svc>/internal/svc`

- 只放依赖初始化与装配
- 统一持有：
  - GORM DB
  - Redis client
  - RPC client
  - MQ publisher / consumer
  - JWT helper
  - OAuth provider client
- 不在这里写业务规则

### `services/<svc>/internal/model`

- 只放服务私有数据访问层
- 默认进一步拆分为：
  - `internal/model/entity`
  - `internal/model/repo`
- 不放 HTTP DTO、RPC DTO、handler 参数结构

#### `services/<svc>/internal/model/entity`

- 只放 GORM 表结构映射
- 一个表结构一个文件优先

#### `services/<svc>/internal/model/repo`

- 只放 repository、store、query 封装、事务辅助
- `service` 通过 `repo.Store` 访问数据库
- 不在 `logic` 或 `server` 中直接出现主要 CRUD

### `services/<svc>/internal/service`

- 只放核心用例编排
- service 负责：
  - 输入规范化与业务校验
  - 事务边界
  - 调用 `repo`
  - 调用 `auth` helper
  - 写审计
- service 不负责：
  - gRPC status
  - HTTP / RPC DTO 绑定
  - transport metadata 读取

### `services/<svc>/internal/logic`

- 只放 transport 适配逻辑
- 一条 RPC 或一个 HTTP 用例对应一个 logic
- logic 只消费 `svc.ServiceContext` 暴露的 `Services`
- logic 负责：
  - 请求参数适配
  - metadata 提取
  - 调用 `service`
  - 将 `service` 错误映射到 transport 错误
- 不直接 new 数据库连接、Redis 连接或第三方 provider client
- 不直接执行主要数据库 CRUD 或事务编排

### `services/<svc>/internal/handler` / `server`

- 只做 transport 适配
- 负责请求绑定、调用 logic、返回结果
- 不写业务编排和数据访问逻辑

### `pkg/`

- 只放跨服务共享能力
- 例如：
  - `pkg/auth`
  - `pkg/errs`
  - `pkg/ctxmeta`
  - `pkg/xgorm`
  - `pkg/xredis`
  - `pkg/mq`
- 不放服务私有 repository 或服务专属业务规则

## 技术栈边界

### 保留 go-zero 的范围

- HTTP / RPC 骨架
- `zrpc`
- `etcd` 注册发现集成
- 日志
- 错误码与服务基础壳
- 配置加载

### 默认不使用 go-zero 的范围

- 数据库访问层
- 内置 model / sqlx
- 内置 MQ 方案

## 基础设施标准

- PostgreSQL：`GORM`
- Redis：`go-redis`
- etcd：`go-zero / zrpc`
- MQ：`RabbitMQ + pkg/mq`

## 服务级约束

### `identity`

- 使用 `internal/model/entity` 承载账户、会话、refresh token、SSO 绑定、审计等表结构映射
- 使用 `internal/model/repo` 承载账户、会话、refresh token、SSO 绑定、审计等持久化访问
- 使用 `internal/service` 承载认证、token、会话、SSO 用例编排
- `internal/logic` 只做 RPC transport 适配
- 配置结构和配置校验放在 `internal/config`
- 第三阶段实际开放的 SSO provider 只有 `GitHub`
- provider client 必须实例化注入，不允许继续依赖包级全局 HTTP/OAuth 钩子

### `gateway`

- 保持 transport 层定位
- 不直接访问业务数据库
- 通过 RPC 调用后端服务
- 如需要鉴权中间件，统一收口到 `internal/middleware`

## 注释与日志规范

### 注释规范

- 代码注释统一采用中英双语
- 英文注释写在上方，中文注释写在下方
- 函数、方法、结构体、接口、关键业务分支、复杂数据流都应补充双语注释
- 简单赋值、明显语义的单行代码不强制加注释，避免噪音

推荐格式：

```go
// Create a new session and persist the refresh token hash.
// 创建新会话并持久化 refresh token 哈希。
func CreateSession() {}
```

### 日志规范

- 运行时日志统一使用英文
- 日志内容应面向排障和监控，不写中英文混合消息
- 错误日志优先包含：
  - 行为
  - 关键标识
  - 失败原因
- 不在日志中输出密码、token 明文、客户端密钥等敏感信息

推荐格式：

```go
logx.Errorf("failed to refresh session token, user_id=%s session_id=%s: %v", userID, sessionID, err)
```

## 禁止项

- 不在 `logic` 中直接初始化 DB / Redis / MQ 连接
- 不在 `logic` 中直接执行主要数据库 CRUD 或事务编排
- 不在 `gateway` 中直接访问 PostgreSQL
- 不把服务私有 repository 放到 `pkg`
- 不把业务编排写进 `handler`
- 不使用 go-zero 内置 model/sqlx 作为默认实现路径
