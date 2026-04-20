# Beehive Blog v3 日志规范

## 目标

统一 `v3` 项目的日志入口、字段规范、敏感信息治理与错误关联方式，确保：

- 业务代码不直接依赖底层日志实现
- 日志可稳定关联 `request_id`、业务错误码和关键标识
- 客户端不看到内部细节，服务端日志保留足够排障信息
- 后续 `content`、`search`、`realtime` 可以直接复用

## 统一入口

- 全项目业务代码统一使用 `pkg/logs`
- `pkg/logs` 当前底层基于 go-zero `logx`
- `service`、`logic`、`middleware`、`svc` 中不再直接调用 `logx`
- `logx` 只允许保留在 `pkg/logs` 或明确批准的底层适配层中

## 基本写法

推荐写法：

```go
logs.Ctx(ctx).Info(
    "identity_register_local_user_succeeded",
    logs.Int64("user_id", user.ID),
    logs.String("username", user.Username),
)
```

错误日志写法：

```go
logs.Ctx(ctx).Error(
    "auth_introspect",
    err,
    logs.String("route", r.URL.Path),
    logs.String("upstream_code", upstreamCode),
)
```

## action 规范

- 每条日志必须有明确 `action`
- `action` 统一使用英文、小写、下划线风格
- 建议结构：
  - `<domain>_<operation>_<result>`
  - 例如：
    - `identity_register_local_user_succeeded`
    - `auth_introspect`
    - `readyz_check`

## 字段规范

优先使用这些字段名：

- `request_id`
- `user_id`
- `session_id`
- `route`
- `provider`
- `code`
- `dependency`
- `reason`

字段值统一使用英文或机器可读值，不写中英文混合说明。

## 错误与日志的关系

- 错误的机器语义以 `pkg/errs` 业务错误码为准
- 日志通过 `pkg/logs` 自动附带或显式写入业务错误码
- 客户端响应与服务端日志必须分离：
  - 客户端只看 `code/message/reference/request_id`
  - 服务端日志可以看 `cause`

## 敏感信息黑名单

以下信息禁止直接写入日志：

- `password`
- `access token`
- `refresh token`
- `authorization`
- `cookie`
- `provider secret`
- 第三方敏感响应体
- 完整 SQL 原文

如果字段名包含上述敏感关键词，必须被掩码。

## Review 与自动检查

- 后续手写业务代码必须通过 `pkg/logs` 写日志
- 提交前建议执行：
  - `go run ./tools/reviewrules`
  - 或 `./scripts/check-review-rules.ps1`
- 自动检查会拦截：
  - 直接导入 `logx`
  - 直接调用 `logx.*`
  - 新增 `errs.IsCode(`
  - 基于错误字符串的业务分支判断
- GitHub Actions 默认执行同一套规则检查，防止后续 PR 引入直接 `logx` 依赖或旧式错误判断

## 级别约定

- `Info`
  - 正常流程关键节点
  - 生命周期日志
  - 成功完成的业务动作
- `Warn`
  - 可恢复异常
  - 降级、回退、非致命异常
- `Error`
  - 当前请求或当前动作失败
  - 依赖不可用
  - 状态不一致或严重异常
- `Debug`
  - 只用于本地开发或临时排查
  - 不应成为线上主路径依赖
