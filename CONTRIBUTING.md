# Contributing Guide

感谢你为 Beehive Blog v3 做出贡献。

## Development Workflow

1. Fork 或创建功能分支（建议命名：`feat/*`、`fix/*`、`chore/*`）
2. 小步提交，保证每次提交可构建、可测试
3. 发起 Pull Request，并补齐模板中的测试与风险说明
4. 通过 Code Review 与 CI 后合并

## Commit & PR Standards

- 提交信息建议遵循 Conventional Commits（如 `feat: ...`、`fix: ...`）
- 一个 PR 聚焦一个主题，避免混合无关改动
- PR 需要说明：
  - 变更背景与目标
  - 核心改动点
  - 测试证明（命令与结果）
  - 潜在风险与回滚方式

## Code Quality Requirements

- 保持分层职责清晰，避免跨层耦合
- 新增逻辑必须包含错误处理与边界检查
- 复杂逻辑请补充必要注释
- 日志统一使用英文，避免泄露敏感信息
- 涉及注释的代码改动，遵循英文在上、中文在下的双语注释约定

## Testing Requirements

- 至少执行受影响模块测试
- 对关键路径与回归风险补充测试用例
- 推荐在提交前运行：

```bash
go test ./...
```

## Security Requirements

- 不要提交密钥、令牌、账号密码等敏感信息
- 不要在日志中输出敏感字段明文
- 发现安全问题请按 `SECURITY.md` 的流程私下披露

## Communication

- 讨论尽量聚焦事实、日志与可复现步骤
- 对评审意见保持可追踪闭环（已修复/暂缓/拒绝及理由）
