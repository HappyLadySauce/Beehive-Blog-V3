---
name: contract-first-goctl
description: Enforce the Beehive Blog v3 contract-first workflow. Use when adding or modifying gateway HTTP contracts or internal RPC contracts so work always follows this order: edit `v3/api/*.api` or `v3/proto/*.proto`, regenerate skeletons with `goctl`, then adjust generated and handwritten code.
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

## v3 Repository Standard

This repository currently uses the following contract locations and generation style:

- gateway HTTP contract: `v3/api/gateway.api`
- backend RPC contracts: `v3/proto/*.proto`
- RPC generation target: `services/<service>`
- RPC generation command must include `--client=true`

Treat this as the repository standard unless the repository is explicitly changed later.

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

- edit `v3/api/gateway.api` first

Use `gateway.api` only for routes that `gateway` really exposes.

Under the current v3 design, this means:

- public HTTP routes
- studio HTTP routes
- gateway-owned health routes
- websocket entry routes

Do not treat `gateway.api` as a place for business orchestration design.  
It is the external HTTP contract source of truth.

### Internal RPC changes

For service-to-service changes:

- edit `v3/proto/*.proto` first

**Proto `go_package` (goctl and this repo layout)**

- In `v3/proto/*.proto`, set `go_package` **relative to the repository root**, for example:  
  `option go_package = "services/<service>/pb";` (identity: `services/identity/pb`).
- Do **not** put the full Go module path in `go_package` (for example `github.com/HappyLadySauce/Beehive-Blog-V3/services/...`).  
  Generators may treat those segments as directory names and write files under the wrong tree.

Examples:

- `v3/proto/identity.proto`
- `v3/proto/content.proto`
- `v3/proto/search.proto`

If a new business capability has no clear existing owner, create a new dedicated service contract first:

- `v3/proto/dashboard.proto`
- `v3/proto/query.proto`
- `v3/proto/composite.proto`

### Required sequencing

When an HTTP endpoint needs a new backend capability:

1. edit `v3/proto/<service>.proto`
2. regenerate that RPC service
3. edit `v3/api/gateway.api`
4. regenerate gateway skeleton
5. finish handwritten logic and wiring

## Regeneration Commands

Use `goctl` directly and follow the repository standard exactly.

### Gateway API validation

Before generation, validate the API contract:

```powershell
goctl api validate -api .\v3\api\gateway.api
```

### Gateway API generation

```powershell
goctl api go -api .\v3\api\gateway.api -dir services\gateway
```

### Gateway Swagger generation

After gateway API contract changes, also regenerate Swagger:

```powershell
goctl api swagger --api .\v3\api\gateway.api --dir .\v3\api --filename gateway --yaml
```

Notes:

- `goctl api swagger` requires a sufficiently new `goctl` version
- Swagger writes **`v3/api/gateway.yaml`** next to **`v3/api/gateway.api`** (same directory; do not use a nested `swagger/` folder)
- keep Swagger generation tied to the same `v3/api/gateway.api` contract
- if the repository later prefers json instead of yaml, update the repository standard consistently

### RPC generation

Example pattern:

```powershell
goctl rpc protoc .\v3\proto\identity.proto --go_out=. --go-grpc_out=. --zrpc_out=services\identity --client=true
```

Repeat per service:

```powershell
goctl rpc protoc .\v3\proto\content.proto --go_out=. --go-grpc_out=. --zrpc_out=services\content --client=true
goctl rpc protoc .\v3\proto\search.proto --go_out=. --go-grpc_out=. --zrpc_out=services\search --client=true
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

For gateway contract changes, the full sequence becomes:

1. edit `v3/api/gateway.api`
2. run `goctl api validate`
3. run `goctl api go`
4. run `goctl api swagger`
5. edit generated and handwritten code

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

1. Confirm `v3/api/gateway.api` passes `goctl api validate`.
2. Confirm `v3/api/gateway.api` matches the intended public HTTP surface.
3. Confirm the related `v3/proto/*.proto` files match the required backend RPC capability.
4. Confirm regenerated gateway handlers/types compile against current code.
5. Confirm regenerated RPC stubs compile against current service code.
6. Confirm `v3/api/gateway.yaml` is regenerated from the same `v3/api/gateway.api`.
7. Confirm `gateway` remains a transport layer and does not gain new business orchestration.
8. Run compile or test checks where possible, for example:
   - `go test ./services/...`
   - or scoped package tests for the touched services

If local environment limitations prevent full verification, still validate contract shape, generation success, and the affected service compile scope as far as possible.
