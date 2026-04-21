---
name: v3-contract-first
description: "Apply the contract-first workflow in Beehive Blog v3. Use when a task changes `v3/api/*.api` or `v3/proto/*.proto` and the contract sequence must be handled correctly."
---

# v3 Contract First

## Overview

Use this skill when contracts change.
在契约发生变化时使用本 skill。

It only defines the correct source-of-truth order and contract sequencing.
它只定义正确的契约真相源与修改顺序。

## Use When

- The task changes `v3/api/gateway.api`.
- The task changes `v3/proto/*.proto`.
- A gateway HTTP capability needs a new backend RPC capability.

## Do

- Change contracts before handwritten implementation.
- Keep gateway HTTP contracts in `v3/api/gateway.api`.
- Keep backend RPC contracts in `v3/proto/*.proto`.
- Follow the sequence:
  - backend RPC contract first when HTTP needs new backend capability
  - gateway HTTP contract second
- Keep proto `go_package` relative to the repository root.

## Do Not

- Do not start from generated files.
- Do not move business orchestration into `gateway`.
- Do not duplicate project-wide coding or logging rules here.

## Hand-off

- After the contract shape is clear, continue with `$v3-goctl`.
- After generation, continue with `$v3-during-task`.

## References

- [goctl 参考](../goctl-workflow/references/goctl.md)
- [编码规范](../../docs/v3/development/coding-conventions.md)
