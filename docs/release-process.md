# Release Process

本文档定义 Beehive Blog v3 的标准发布流程，确保版本可追踪、可回滚、可审计。

## Versioning

- 使用 SemVer：`MAJOR.MINOR.PATCH`
- 建议约定：
  - `MAJOR`：不兼容变更
  - `MINOR`：向后兼容功能新增
  - `PATCH`：向后兼容缺陷修复

## Pre-release Checklist

- 主分支 CI 全绿
- 关键路径测试通过（单测/集成测试）
- 安全相关变更已评审
- `CHANGELOG.md` 已更新
- 文档与配置变更已同步

## Release Steps

1. 从主分支创建发布分支（可选）
2. 确认版本号并更新 `CHANGELOG.md`
3. 打 Tag（例如 `v3.1.0`）
4. 推送 Tag 到远端
5. 在 GitHub 创建 Release，附上变更说明与升级注意事项

## Rollback Strategy

- 快速回滚到上一个稳定 Tag
- 如涉及数据库变更，必须执行预先准备的回滚脚本
- 回滚后立即发布事故说明与后续修复计划

## Post-release

- 观察日志、错误率、延迟、资源占用等关键指标
- 收集用户反馈并归档到 Issue
- 若出现高优先级问题，发布补丁版本
