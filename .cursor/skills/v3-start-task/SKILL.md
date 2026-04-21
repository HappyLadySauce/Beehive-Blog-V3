---
name: v3-start-task
description: "Route a Beehive Blog v3 task at the start of development. Use when a task is just beginning and you need to decide whether it involves contract changes, code generation, handwritten business code, configuration, testing, or documentation updates."
---

# v3 Start Task

## Overview

Use this skill at the beginning of a development task.
在开发任务开始时使用本 skill。

Its purpose is to classify the task and decide which follow-up skills must be used next.
它的目标是为当前任务分流，并决定后续必须串联哪些 skill。

## Use When

- The task is just starting.
- You need to determine whether the work involves API/proto changes.
- You need to decide whether `goctl` generation is required.
- You need to identify whether the task is mainly transport, service, repository, config, test, or docs work.

## Do

- Confirm the owning service and boundary first.
- Identify whether the task changes:
  - `v3/api/*.api`
  - `v3/proto/*.proto`
  - handwritten business code
  - config structs or validation
  - tests
  - docs and error-code registrations
- Route follow-up work to the right skills:
  - contract changes -> `$v3-contract-first`
  - code generation -> `$v3-goctl`
  - implementation -> `$v3-during-task`
  - tests -> `$v3-testing`
  - error/log work -> `$v3-error-and-logging`

## Do Not

- Do not explain full repository conventions here.
- Do not inline all `goctl` commands here.
- Do not duplicate detailed error/log/testing rules.
- Do not start coding before the task type and service boundary are clear.

## Hand-off

- If the task changes contracts, continue with `$v3-contract-first`.
- If the task requires generation, continue with `$v3-goctl`.
- If the task is entering implementation, continue with `$v3-during-task`.

## References

- [开发规范索引](../../docs/v3/development/README.md)
- [编码规范](../../docs/v3/development/coding-conventions.md)
