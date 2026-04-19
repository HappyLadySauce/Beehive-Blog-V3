---
name: goctl-workflow
description: Scaffold and maintain Beehive Blog v3 go-zero services with goctl. Use when creating or updating the gateway API service or pure RPC backend services, generating code from `v3/api/*.api` or `v3/proto/*.proto`, or explaining how generated files should map to handwritten business code in this repository.
---

# goctl Workflow

## Overview

Use `goctl` only after the service boundary and contract are already clear.

For this repository:

- `gateway` is the only HTTP / WebSocket service
- backend business services are pure RPC
- gateway contracts live under `v3/api/`
- backend proto contracts live under `v3/proto/`
- RPC generation must include `--client=true`
- `gateway` should stay thin and transport-oriented
- business orchestration belongs in the owning domain service
- `etcd` continues through go-zero / `zrpc`
- database access uses GORM, not go-zero model/sqlx
- Redis access uses go-redis
- MQ uses third-party RabbitMQ through `pkg/mq`

## Required Workflow

1. Confirm the service boundary first.
2. Write or update the `.api` or `.proto` contract first.
3. Run `goctl` to generate the service skeleton or update transport code.
4. Edit generated code and handwritten code to match the contract.
5. Verify compile and wiring before continuing.

## Gateway API Service

For HTTP-facing gateway changes:

- keep the contract in `v3/api/gateway.api`
- generate into `services/gateway`
- treat generated `types` as transport DTOs only
- validate the `.api` file before generation
- regenerate Swagger after contract updates

Standard command:

```powershell
goctl api validate -api .\v3\api\gateway.api
goctl api go -api .\v3\api\gateway.api -dir services\gateway
goctl api swagger --api .\v3\api\gateway.api --dir .\v3\api --filename gateway --yaml
```

The gateway should:

- expose HTTP / WS routes
- perform transport adaptation
- call backend RPC services

The gateway should not become the place for business orchestration.

Read [references/goctl.md](references/goctl.md) for concrete commands.

## RPC Services

For internal service-to-service calls:

- keep contracts in `v3/proto/<service>.proto`
- generate RPC code into `services/<service>`
- include `--client=true`
- treat proto as the source of truth

**Proto `go_package`**

- Use a path **relative to the repository root**, for example:  
  `option go_package = "services/<service>/pb";`
- Do **not** use the full Go module path inside `go_package` (for example `github.com/HappyLadySauce/Beehive-Blog-V3/services/...`).  
  goctl/protoc output can mirror those segments as nested directories and land generated code in the wrong place.

Expose RPC for stable business capabilities, not every table operation.

If a capability has no clear existing service owner, create a new service contract instead of overloading gateway.

## Generated Code Boundaries

Treat `goctl` output as transport and wiring code.

- `handler/`: request entry only
- `logic/`: transport-side use case glue
- `svc/`: dependency wiring
- `types/`: request and response DTOs
- generated RPC server/client stubs: contract transport layer

Keep business rules outside request parsing handlers.

For this repository, prefer handwritten business logic in:

- `services/*/internal/logic/`
- `services/*/internal/server/`
- `services/*/internal/svc/`
- `services/*/internal/model/` for service-private data access
- `pkg/*` for cross-service shared helpers and infrastructure wrappers

Current implementation boundaries:

- `internal/config`: config structs only
- `internal/svc`: dependency wiring only
- `internal/model`: GORM entities / repositories / queries
- `internal/logic`: business use case orchestration
- `pkg/mq`: RabbitMQ publisher / consumer abstraction

Do not use go-zero built-in database model/sqlx as the default path in this repository.

## Regeneration Rules

- Regenerate from `v3/api/*.api` or `v3/proto/*.proto`, not from edited generated files
- Contract change first, generation second, handwritten adjustments third
- For gateway API changes, validate first and regenerate Swagger after code generation
- If generation changes names or signatures, update surrounding code to match the new contract cleanly
- Do not preserve outdated generated glue just for compatibility during v3 development

## Monorepo Fit For This Repo

Use the current repository layout:

- `v3/api/`
- `v3/proto/`
- `services/`
- `docs/v3/`
- `.codex/skills/`

Do not follow older examples that assume:

- `api/` instead of `v3/api/`
- `proto/` instead of `v3/proto/`
- RPC generation without `--client=true`
- every service has both API and RPC
- gateway contains business aggregation layers by default

## Validation

After generation:

1. Verify the edited contract is the actual source of truth.
2. For gateway API changes, run `goctl api validate -api .\v3\api\gateway.api`.
3. Verify `goctl` output lands in the correct service directory.
4. Verify Swagger output is regenerated from the same gateway API contract.
5. Verify gateway still acts as transport only.
6. Verify backend services remain pure RPC.
7. Build or test the affected services before moving on.
