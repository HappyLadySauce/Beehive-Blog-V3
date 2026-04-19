---
name: contract-first-goctl
description: Enforce the Beehive Blog v3 contract-first workflow. Use when adding or modifying gateway HTTP contracts or internal RPC contracts so work always follows this order: edit `api/*.api` or `proto/*.proto`, regenerate skeletons with `goctl`, then adjust generated and handwritten code.
---

# Contract-First goctl Workflow

Apply this workflow for every interface change in v3.

## Core Rule

Always follow this order:

1. Edit contracts first.
2. Regenerate code from contracts.
3. Adjust generated and handwritten code.
4. Verify compile and transport wiring.

Never start by editing generated transport files first.

## v3 Architecture Assumptions

This workflow assumes the current v3 architecture:

- `gateway` is the only HTTP / WebSocket entry service
- other business services are pure RPC
- gateway does not carry business orchestration
- business orchestration belongs in the primary domain service
- if an interface has no clear domain owner, create a dedicated service instead of pushing orchestration into `gateway`

## Contract Update Rules

### Gateway HTTP changes

For gateway-facing HTTP changes:

- edit `api/gateway.api` first

Use `gateway.api` only for routes that `gateway` really exposes.

Under the current v3 design, this means:

- public HTTP routes
- studio HTTP routes
- gateway-owned health / websocket entry routes

Do not treat `gateway.api` as a place for business orchestration design.  
It is the external HTTP contract source of truth.

### Internal RPC changes

For service-to-service changes:

- edit `proto/*.proto` first

Examples:

- `proto/identity.proto`
- `proto/content.proto`
- `proto/search.proto`

If a new business capability has no clear existing owner, create a new dedicated service contract first:

- `proto/dashboard.proto`
- `proto/query.proto`
- `proto/composite.proto`

### Required sequencing

When an HTTP endpoint needs a new backend capability:

1. edit `proto/<service>.proto`
2. regenerate that RPC service
3. edit `api/gateway.api`
4. regenerate gateway skeleton
5. finish handwritten logic and wiring

## Regeneration Commands

Use `goctl` directly. Do not rely on repo-local wrapper scripts unless they are later added to the repository.

### Gateway API generation

```powershell
goctl api go -api .\api\gateway.api -dir .\services\gateway
```

### RPC generation

Example pattern:

```powershell
goctl rpc protoc .\proto\identity.proto --go_out=. --go-grpc_out=. --zrpc_out=.\services\identity
```

Repeat per service:

```powershell
goctl rpc protoc .\proto\content.proto --go_out=. --go-grpc_out=. --zrpc_out=.\services\content
goctl rpc protoc .\proto\search.proto --go_out=. --go-grpc_out=. --zrpc_out=.\services\search
```

If `protoc` is not installed or unavailable, stop and fix the local toolchain first.

## Generated Code Handling

Treat goctl output as transport and wiring scaffolding.

- `handler/`: request entry
- `logic/`: transport-side use case glue
- `svc/`: dependency wiring
- `types/`: HTTP DTOs
- generated RPC server/client stubs: transport contracts only

After generation, it is expected that you will edit generated code and adjacent handwritten code.

Under this repository's workflow, the normal sequence is:

1. edit contract
2. run `goctl`
3. edit generated code as needed
4. edit handwritten business code

Do not manually keep old generated glue alive if the new contract changes names or shapes.  
Adjust the generated code and the surrounding business code to match the new contract cleanly.

## Business Code Update Rules

After regeneration, update business layers as needed:

- `services/*/internal/logic/`
- `services/*/internal/svc/`
- `services/*/internal/server/`
- `services/gateway/internal/handler/`
- `services/gateway/internal/middleware/`

Put business orchestration in the primary domain service, not in `gateway`.

Examples:

- order complete view -> `order-service`
- content complete view -> `content-service`

If no primary domain owner exists, create a new service instead of composing inside gateway.

## Verification Checklist

1. Confirm `api/gateway.api` matches the intended public HTTP surface.
2. Confirm the related `proto/*.proto` files match the required backend RPC capability.
3. Confirm regenerated gateway handlers/types compile against current code.
4. Confirm regenerated RPC stubs compile against current service code.
5. Confirm `gateway` remains a transport layer and does not gain new business orchestration.
6. Run compile or test checks where possible, for example:
   - `go test ./services/...`
   - or scoped package tests for the touched services

If local environment limitations prevent full verification, still validate contract shape, generation success, and the affected service compile scope as far as possible.
