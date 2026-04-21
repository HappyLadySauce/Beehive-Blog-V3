---
name: v3-review-rules
description: "Apply Beehive Blog v3 review-rule enforcement. Use when you need to check or explain the automated rules around `pkg/logs`, `pkg/errs`, and forbidden legacy patterns."
---

# v3 Review Rules

## Overview

Use this skill for review-rule enforcement.
在需要执行或解释 review 规则时使用本 skill。

It explains the automated rule gate and how to run it locally.
它解释自动门禁规则及其本地执行方式。

## Use When

- You are reviewing handwritten business code.
- You need to run `tools/reviewrules`.
- A task touches error matching or logging style.
- You need to explain why CI blocks a PR on review rules.

## Do

- Enforce these rules:
  - no direct `logx` imports or calls outside `pkg/logs`
  - prefer `errors.Is(err, errs.E(...))`
  - do not use `errs.IsCode(...)` as the primary matching style in new code
  - do not branch on `err.Error()` or similar string matching
- Run:
  - `go run ./tools/reviewrules`
  - or `./scripts/check-review-rules.ps1`
- Check whitelist behavior before changing the scanner.

## Do Not

- Do not restate all coding conventions here.
- Do not treat this skill as the only source of truth; the docs remain authoritative.
- Do not bypass the rule gate for handwritten business code.

## Hand-off

- If the task is still being implemented, return to `$v3-during-task`.
- If the task is being finalized, continue with `$v3-finish-task`.

## References

- [代码评审清单](../../docs/v3/development/review-checklist.md)
- [日志规范](../../docs/v3/development/logging-conventions.md)
- [错误码规范](../../docs/v3/development/error-code-specification.md)
