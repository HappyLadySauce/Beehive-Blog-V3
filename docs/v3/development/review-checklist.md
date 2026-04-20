# Beehive Blog v3 代码评审清单

## 目标

统一 `v3` 项目的 code review 关注点，确保后续服务开发在错误处理、日志、安全与分层上保持一致。

## 错误与日志专项必查项

- 是否统一使用 `pkg/errs` 构造领域错误，而不是自造私有错误模型
- 是否优先使用 `errors.Is(err, errs.E(...))` 判断业务错误
- 是否避免把 `errs.IsCode(err, ...)` 作为新增代码的首选写法
- 是否避免使用 `err.Error()`、`strings.Contains(err.Error(), ...)`、gRPC message、SQL message 做业务分支判断
- 是否统一使用 `pkg/logs` 写业务日志，而不是直接依赖 `logx`
- 是否避免向客户端暴露底层 `cause`、SQL 原文、gRPC 原始错误文本
- 是否避免在日志中输出密码、token、secret、raw SQL、第三方敏感响应体
- 是否在新增错误码时同步更新 `pkg/errs` 常量、错误码文档和对应测试

## 分层与边界必查项

- `logic` 是否只做 transport 适配，而不是直接做事务编排和主要 CRUD
- `service` 是否只返回领域错误，而不是直接拼 HTTP / gRPC transport 错误
- `gateway` 是否继续保持薄接入层，不直接访问业务数据库
- `pkg/` 中是否只放跨服务共享能力，而没有混入服务私有规则

## 测试与验证必查项

- 是否补了与新增错误分支对应的单元测试或服务测试
- 是否补了 HTTP / gRPC 错误映射测试
- 是否在必要时补了日志与敏感信息保护相关测试
- 是否执行了当前服务的最小验证命令
- 是否执行了 `go run ./tools/reviewrules` 或 `./scripts/check-review-rules.ps1`
- 是否确认 GitHub Actions 的 `review-rules` 工作流会覆盖本次改动范围
