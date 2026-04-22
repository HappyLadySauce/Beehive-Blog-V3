---
name: v3-goctl
description: "Run goctl correctly in Beehive Blog v3. Use when contracts are already clear and you need the exact generation commands, output locations, and generated-code boundaries."
---

# v3 goctl

## Overview

Use this skill when `goctl` generation is required.
在需要执行 `goctl` 生成时使用本 skill。

It focuses only on command usage, output locations, and generated-code handling.
它只关注命令用法、输出位置与 generated code 的处理边界。

## Use When

- Gateway HTTP contract generation is required.
- RPC code generation is required from `v3/proto/*.proto`.
- You need to know where generated transport code stops and handwritten code begins.

## Do

- Run `goctl` only after contracts are already clear.
- Use repository-standard commands and output directories.
- Regenerate Swagger after gateway API contract changes.
- Treat generated files as transport and wiring scaffolding.
- Check regenerated Swagger output before delivery:
  - top-level `info`
  - route `summary` / `tags`
  - request field examples and options
  - key business error descriptions
- Keep Swagger enhancement in the contract source; do not use generated-file patching as the primary workflow.

## Do Not

- Do not redesign service boundaries here.
- Do not treat `v3/api/gateway.yaml` as an editing entrypoint.
- Do not repeat project-wide logging or error-handling rules here.
- Do not treat generated glue as the long-term source of truth.

## Hand-off

- After generation, continue with `$v3-during-task`.
- Before delivery, continue with `$v3-finish-task`.

## References

- [goctl 参考](../goctl-workflow/references/goctl.md)
- [编码规范](../../docs/v3/development/coding-conventions.md)
