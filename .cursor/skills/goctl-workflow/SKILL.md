---
name: goctl-workflow
description: Scaffold and maintain Beehive Blog v3 go-zero services with goctl. Use when creating or updating the gateway API service or pure RPC backend services, generating code from `.api` or `.proto`, or explaining how generated files should map to handwritten business code in this repository.
---

# goctl Workflow

## Overview

Use `goctl` only after the service boundary and contract are already clear.

For this repository:

- `gateway` is the only HTTP / WebSocket service
- backend business services are pure RPC
- `gateway` should stay thin and transport-oriented
- business orchestration belongs in the owning domain service

## Required Workflow

1. Confirm the service boundary first.
2. Write or update the `.api` or `.proto` contract first.
3. Run `goctl` to generate the service skeleton or update transport code.
4. Edit generated code and handwritten code to match the contract.
5. Verify compile and wiring before continuing.

## Gateway API Service

For HTTP-facing gateway changes:

- keep the contract in `api/gateway.api`
- generate into `services/gateway`
- treat generated `types` as transport DTOs only

The gateway should:

- expose HTTP / WS routes
- perform transport adaptation
- call backend RPC services

The gateway should not become the place for business orchestration.

Read [references/goctl.md](references/goctl.md) for concrete commands.

## RPC Services

For internal service-to-service calls:

- keep contracts in `proto/<service>.proto`
- generate RPC code into `services/<service>`
- treat proto as the source of truth

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

Do not introduce `domain/` / `repository/` as a default rule unless the repository later adopts that structure intentionally.

## Regeneration Rules

- Regenerate from `.api` or `.proto`, not from edited generated files
- Contract change first, generation second, handwritten adjustments third
- If generation changes names or signatures, update surrounding code to match the new contract cleanly
- Do not preserve outdated generated glue just for compatibility during v3 development

## Monorepo Fit For This Repo

Use the current repository layout:

- `api/`
- `proto/`
- `services/`
- `docs/v3/`
- `.codex/skills/`

Do not follow older examples that assume:

- `rpc/` instead of `proto/`
- every service has both API and RPC
- gateway contains business aggregation layers by default

## Validation

After generation:

1. Verify the edited contract is the actual source of truth.
2. Verify `goctl` output lands in the correct service directory.
3. Verify gateway still acts as transport only.
4. Verify backend services remain pure RPC.
5. Build or test the affected services before moving on.
