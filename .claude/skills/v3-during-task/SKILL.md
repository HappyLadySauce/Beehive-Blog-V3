---
name: v3-during-task
description: "Guide Beehive Blog v3 implementation work while coding. Use when handwritten business code is being added or changed and you need stage-specific implementation guidance without reloading every project rule."
---

# v3 During Task

## Overview

Use this skill during implementation.
在编码实现阶段使用本 skill。

It focuses on what the current stage should do, not on contract sequencing or final verification.
它关注当前实现阶段该做什么，不负责契约顺序和最终验收。

## Use When

- You are editing handwritten business code.
- You need implementation-stage guidance on layer boundaries.
- You need to know whether the current change belongs in transport, service, repository, config, or tests.

## Do

- Keep work inside the correct layer:
  - transport -> `handler/server/logic`
  - business orchestration -> `service`
  - persistence -> `internal/model/repo`
  - shared cross-service helpers -> `pkg`
- Follow repository-level conventions by loading:
  - `$v3-coding-standards`
  - `$v3-error-and-logging` when touching errors or logs
  - `$v3-testing` when adding or updating tests
- Use the thinnest valid layer for the change.

## Do Not

- Do not make contract-order decisions here.
- Do not duplicate `goctl` command details here.
- Do not skip service/repo boundaries by convenience.
- Do not turn `gateway` into a business orchestration layer.

## Hand-off

- If you discover contract changes are required, switch to `$v3-contract-first`.
- If code generation becomes necessary, switch to `$v3-goctl`.
- Before finishing, switch to `$v3-finish-task`.

## References

- [编码规范](../../docs/v3/development/coding-conventions.md)
- [错误码规范](../../docs/v3/development/error-code-specification.md)
- [日志规范](../../docs/v3/development/logging-conventions.md)
- [测试规范](../../docs/v3/development/testing-conventions.md)
