# Beehive Blog v3 测试规范

## 目标

统一 `v3` 的测试分层、测试入口、资源生命周期管理与集成测试方案，避免后续在单元测试、服务测试和基础设施测试之间反复选型。

## 分层规则

### `internal/auth`

- 作为纯单元测试入口
- 优先覆盖：
  - 输入规范化
  - token 生成与解析
  - 密码哈希与校验
  - provider registry / client 的纯行为

### `internal/service`

- 作为主业务测试入口
- 优先覆盖：
  - 注册
  - 登录
  - refresh
  - logout
  - current user
  - introspect
  - SSO start / finish

### `internal/model/repo`

- 作为持久化集成测试入口
- 优先覆盖：
  - 唯一约束
  - 锁查询
  - consume / revoke / touchLogin
  - 事务中的状态流转

### `internal/logic`

- 只保留少量 smoke tests
- 重点验证：
  - metadata 提取
  - service 调用转发
  - `service error -> gRPC status` 映射

## 测试组织方式

- 默认使用包内 `_test.go`
- 默认使用表驱动测试 + `t.Run`
- 默认使用 `t.Helper()` 抽离测试辅助
- 使用 `t.Cleanup()` 管理资源释放
- 使用 `t.Context()` 管理长生命周期资源
- 仅在纯单元测试里使用 `t.Parallel()`
- 不在并行测试里使用 `t.Setenv()`

## 集成测试方案

- `identity` 的 PostgreSQL / Redis 集成测试默认使用 `Testcontainers`
- `internal/testkit` 负责：
  - 容器启动
  - 迁移执行
  - 测试数据清理
  - 测试 `ServiceContext` 装配
- 环境变量模式只作为 fallback，不作为默认路径

## 当前约束

- `logic` 不是主要业务测试入口
- `service` 是后续身份域测试的核心入口
- provider client 必须实例化注入，测试不依赖包级 HTTP/OAuth 钩子
- `QQ/WeChat` 当前不开放入口，`GitHub` 是唯一完整 provider
