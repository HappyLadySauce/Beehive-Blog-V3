# v3 编码与配置规范

本目录用于沉淀 `v3` 的实现级规范，指导后续服务开发、代码评审和 Agent 自动生成行为。

当前文档包括：

- [编码规范](./coding-conventions.md)
- [配置规范](./configuration-conventions.md)
- [测试规范](./testing-conventions.md)

当前约定：

- 本目录是 `v3` 编码与配置规范的唯一基线
- 架构边界继续以 `docs/v3/contracts`、`docs/v3/gateway`、`docs/v3/identity` 为准
- 实现层目录职责、技术栈边界、配置写法优先以本目录为准
- 新服务默认参考 `identity` 的配置结构与校验模式，不再继续使用“全量 optional”风格
- 代码注释统一采用“英文在上、中文在下”的双语格式
- 运行时日志统一使用英文
