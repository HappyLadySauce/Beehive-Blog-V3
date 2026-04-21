---
name: v3-error-and-logging
description: "Apply Beehive Blog v3 error and logging standards. Use when touching `pkg/errs`, `pkg/logs`, HTTP/gRPC error mapping, or any business code that logs or matches errors."
---

# v3 Error And Logging

## Overview

Use this skill when a task touches errors or logs.
在任务涉及错误处理或日志时使用本 skill。

It covers `pkg/errs`, `pkg/logs`, `errors.Is`, `errors.As`, `errors.Join`, error-response boundaries, and sensitive log handling.
它负责 `pkg/errs`、`pkg/logs`、`errors.Is`、`errors.As`、`errors.Join`、错误响应边界与敏感日志治理。

## Use When

- You are changing domain errors.
- You are changing HTTP or gRPC error mapping.
- You are adding or modifying logs in business code.
- You are reviewing whether a change leaks internal error details.

## Do

- Use `pkg/errs` as the domain error truth source.
- Prefer `errors.Is(err, errs.E(...))` for business matching.
- Use `errs.Parse(err)` or `errors.As` when details are needed.
- Use `pkg/logs` as the only business logging entrypoint.
- Keep client responses and server logs separate.

## Do Not

- Do not use `logx` directly in business code.
- Do not use `errs.IsCode(...)` as the primary style in new code.
- Do not branch on `err.Error()`, gRPC message text, or SQL error text.
- Do not return `errors.Join` directly as the client-facing main error.

## Hand-off

- If the task is still being implemented, return to `$v3-during-task`.
- If you need rule enforcement, continue with `$v3-review-rules`.
- If you are finishing the task, continue with `$v3-finish-task`.

## References

- [错误码规范](../../docs/v3/development/error-code-specification.md)
- [日志规范](../../docs/v3/development/logging-conventions.md)
- [代码评审清单](../../docs/v3/development/review-checklist.md)
