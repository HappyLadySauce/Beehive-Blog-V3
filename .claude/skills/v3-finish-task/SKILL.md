---
name: v3-finish-task
description: "Close out a Beehive Blog v3 development task. Use when implementation is done and you need to verify tests, review rules, docs, and final delivery checks."
---

# v3 Finish Task

## Overview

Use this skill at the end of a development task.
在开发任务收尾阶段使用本 skill。

It is the hand-off and verification skill.
它是交付前的收尾与验收 skill。

## Use When

- Implementation is complete.
- You need to run the final verification pass.
- You need to confirm docs, tests, and review rules are all in sync.

## Do

- Run the relevant verification commands for touched services.
- Run review rule checks:
  - `go run ./tools/reviewrules`
  - or `./scripts/check-review-rules.ps1`
- Check whether new errors/log rules required docs or test updates.
- Review the final change against `review-checklist.md`.

## Do Not

- Do not re-explain implementation details here.
- Do not introduce new contract changes here unless blocked.
- Do not skip review-rule verification for handwritten business code.

## Hand-off

- If final verification fails because contracts changed, return to `$v3-contract-first` or `$v3-goctl`.
- If final verification fails because of implementation issues, return to `$v3-during-task`.

## References

- [代码评审清单](../../docs/v3/development/review-checklist.md)
- [测试规范](../../docs/v3/development/testing-conventions.md)
- [日志规范](../../docs/v3/development/logging-conventions.md)
