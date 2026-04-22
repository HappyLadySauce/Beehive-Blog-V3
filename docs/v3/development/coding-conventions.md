# Beehive Blog v3 编码规范

## 目标

统一 `v3` 的服务分层、目录职责、基础设施边界与公共代码组织方式，降低实现分歧和技术债务。

## 目录职责

### `services/<svc>/internal/config`

- 只放配置结构定义
- 只负责表达配置字段、tag 约束和 `Validate()` 校验入口
- 启动前必须通过 `Validate()` 给出显式错误，不允许把配置错误留到 runtime panic
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
- 初始化失败必须显式返回 `error`，不允许继续保留 `Must*` 或 panic 兜底路径

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
  - `pkg/logs`
  - `pkg/ctxmeta`
  - `pkg/xgorm`
  - `pkg/xredis`
  - `pkg/mq`
- 不放服务私有 repository 或服务专属业务规则

## 技术栈边界

## 错误模型规范

- 全项目统一使用 `pkg/errs` 作为领域错误真相源
- 业务错误码使用六位整数错误码，具体规则见 `docs/v3/development/error-code-specification.md`
- `service` 或核心业务层只返回领域错误，不直接返回 HTTP 状态码或 gRPC 状态码
- HTTP 输出统一通过 `pkg/errs/httpx`
- gRPC 输出统一通过 `pkg/errs/grpcx`
- 不再为单个服务维护私有错误模型
- 业务错误匹配优先使用 `errors.Is(err, errs.E(code))`
- `errs.IsCode(err, code)` 只作为辅助函数保留
- 需要提取丰富错误信息时，使用 `errs.Parse(err)` 或 `errors.As`
- `errors.Join` 只允许用于诊断聚合，不允许直接作为客户端主错误返回
- 禁止使用 `err.Error()`、`strings.Contains(err.Error(), ...)`、gRPC message、SQL message 作为业务分支判断依据

### 保留 go-zero 的范围

- HTTP / RPC 骨架
- `zrpc`
- `etcd` 注册发现集成
- 日志后端
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

### API 契约与 Swagger 真相源

- `v3/api/gateway.api` 同时承担 HTTP 契约与 Swagger 文档真相源职责
- 接口测试需要的说明、示例、可选值与关键业务错误码描述应优先写在 `.api`
- `v3/api/gateway.yaml` 只通过 `goctl` 生成，不手工维护

### `identity`

- 使用 `internal/model/entity` 承载账户、会话、refresh token、SSO 绑定、审计等表结构映射
- 使用 `internal/model/repo` 承载账户、会话、refresh token、SSO 绑定、审计等持久化访问
- 使用 `internal/service` 承载认证、token、会话、SSO 用例编排
- `internal/logic` 只做 RPC transport 适配
- 配置结构和配置校验放在 `internal/config`
- 当前实际开放的 SSO provider 包括 `GitHub`、`QQ`、`WeChat`
- provider client 必须实例化注入，不允许继续依赖包级全局 HTTP/OAuth 钩子
- **测试分层（与 `docs/v3/development/testing-conventions.md` 一致）**：
  - `internal/service` 是身份域的**主业务测试入口**（注册、登录、refresh、logout、introspect、SSO 等闭环优先写在这里）。
  - `internal/logic` **只保留少量 smoke test**，守住 gRPC 适配与错误映射，不在此重复主业务分支。
  - PostgreSQL / Redis 集成测试默认使用 **Testcontainers**，由 `internal/testkit` 统一容器、迁移与清表；仅用环境变量直连本地实例为 **fallback**，不作为默认路径。
  - fallback 环境变量命名统一以 `BEEHIVE_TEST_` 为前缀，具体约定见 `docs/v3/development/testing-conventions.md`。
  - SSO 测试：`GitHub`、`QQ`、`WeChat` 都应覆盖 `start + finish` 主链路和关键失败分支。
  - OAuth 交互一律用 **`httptest.Server` 桩**，测试不访问真实第三方外网。

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
- 项目业务代码统一通过 `pkg/logs` 写日志，不直接调用 `logx`
- 日志内容应面向排障和监控，不写中英文混合消息
- `action` 是日志必填概念
- 错误日志优先包含：
  - 行为
  - `request_id`
  - 业务错误码
  - 关键标识
  - 失败原因
- 不在日志中输出密码、token 明文、客户端密钥等敏感信息
- 结构化字段优先使用 `logs.String`、`logs.Int64`、`logs.Any`
- `logx` 只允许出现在 `pkg/logs` 或后续明确批准的底层适配层中

### Code Review Rules

- 后续手写业务代码必须优先使用 `errors.Is(err, errs.E(...))` 进行业务错误匹配
- 后续手写业务代码必须统一使用 `pkg/logs`
- 禁止新增直接 `logx` 导入与调用
- 禁止新增 `errs.IsCode(` 作为主错误匹配写法
- 禁止新增基于错误字符串的业务分支判断
- 提交前应执行：
  - `go run ./tools/reviewrules`
  - 或 `./scripts/check-review-rules.ps1`
- GitHub Actions 默认执行同一套 review rule 检查，新增违规代码会直接失败

推荐格式：

```go
logs.Ctx(ctx).Error(
    "refresh_session_token_failed",
    err,
    logs.UserID(userID),
    logs.SessionID(sessionID),
)
```

## 禁止项

- 不在 `logic` 中直接初始化 DB / Redis / MQ 连接
- 不在 `logic` 中直接执行主要数据库 CRUD 或事务编排
- 不在 `gateway` 中直接访问 PostgreSQL
- 不把服务私有 repository 放到 `pkg`
- 不把业务编排写进 `handler`
- 不使用 go-zero 内置 model/sqlx 作为默认实现路径
