# v3 编码与配置规范

本目录用于沉淀 `v3` 的实现级规范，指导后续服务开发、代码评审和 Agent 自动生成行为。

当前文档包括：

- [编码规范](./coding-conventions.md)
- [配置规范](./configuration-conventions.md)
- [测试规范](./testing-conventions.md)
- [错误码规范](./error-code-specification.md)
- [日志规范](./logging-conventions.md)
- [代码评审清单](./review-checklist.md)

## Skill 体系

当前仓库的 `.codex/skills/` 建议按“阶段型 + 能力型 + 规范型”三层使用：

- 阶段型：
  - `v3-start-task`
  - `v3-during-task`
  - `v3-finish-task`
- 能力型：
  - `v3-contract-first`
  - `v3-goctl`
  - `v3-review-rules`
- 规范型：
  - `v3-coding-standards`
  - `v3-error-and-logging`
  - `v3-testing`

旧 skill：

- `goctl-workflow`
- `contract-first-goctl`

当前只作为过渡入口保留，后续应优先使用新的拆分 skill，而不是继续把所有规则一次性加载。

当前约定：

- 本目录是 `v3` 编码与配置规范的唯一基线
- 架构边界继续以 `docs/v3/contracts`、`docs/v3/gateway`、`docs/v3/identity` 为准
- 实现层目录职责、技术栈边界、配置写法优先以本目录为准
- 新服务默认参考 `identity` 的配置结构与校验模式，不再继续使用“全量 optional”风格
- 代码注释统一采用“英文在上、中文在下”的双语格式
- 运行时日志统一使用英文
- 错误统一采用 `pkg/errs`，日志统一采用 `pkg/logs`
- 后续手写业务代码默认先通过 `go run ./tools/reviewrules` 做轻量规则检查
- GitHub Actions 默认执行 review rule 检查，确保新代码不回退到旧错误/日志写法
