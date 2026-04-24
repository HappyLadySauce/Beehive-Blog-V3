---
name: v3-goctl
description: "Run allowed goctl generation in Beehive Blog v3. Use when contracts are already clear and you need Swagger or RPC generation boundaries."
---

# v3 goctl

## Overview

Use this skill when allowed `goctl` generation is required.
在需要执行允许范围内的 `goctl` 生成时使用本 skill。

It focuses only on Swagger generation, RPC generation, output locations, and generated-code handling.
它只关注 Swagger 生成、RPC 生成、输出位置与 generated code 的处理边界。

## Use When

- Gateway Swagger generation is required from `v3/api/gateway.api`.
- RPC code generation is required from `v3/proto/*.proto`.
- You need to know which generated files are safe to refresh.

## Do

- Run `goctl` only after contracts are already clear.
- Use repository-standard commands and output directories.
- Regenerate Swagger after gateway API contract changes.
- Regenerate RPC artifacts after backend proto contract changes.
- Check regenerated Swagger output before delivery:
  - top-level `info`
  - route `summary` / `tags`
  - request field examples and options
  - key business error descriptions
- Keep Swagger enhancement in the contract source; do not use generated-file patching as the primary workflow.

## Do Not

- Do not redesign service boundaries here.
- Do not treat `v3/api/gateway.yaml` as an editing entrypoint.
- Do not run `goctl api go` for gateway.
- Do not generate or overwrite `services/gateway` Go code from `v3/api/gateway.api`.
- Do not overwrite handwritten gateway route wiring, especially `services/gateway/internal/handler/routes.go`.
- Do not repeat project-wide logging or error-handling rules here.
- Do not treat generated glue as the long-term source of truth.

## Hand-off

- After allowed generation, continue with `$v3-during-task`.
- Before delivery, continue with `$v3-finish-task`.

## References

- [编码规范](../../docs/v3/development/coding-conventions.md)
