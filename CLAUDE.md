# CLAUDE.md 请始终使用简体中文与我对话，并在回答时保持专业、简洁。

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Quick Reference

- **Language**: Go 1.26.1, module `github.com/HappyLadySauce/Beehive-Blog-V3`
- **Framework**: go-zero v1.10.1 (HTTP/RPC scaffold, service discovery, config loading)
- **Database access**: GORM + PostgreSQL (`github.com/jackc/pgx/v5`)
- **Cache/Session state**: go-redis v9
- **Service discovery**: etcd (via go-zero zrpc)
- **Message queue**: RabbitMQ (through `pkg/mq`)
- **Code generation**: `goctl` — scaffolds from `v3/api/gateway.api` and `v3/proto/*.proto`
- **QA suite**: Python + pytest + httpx in `qa/`
- **Commit style**: Conventional Commits (`feat:`, `fix:`, `chore:`, etc.)
- **Indentation**: 2 spaces (general), 4 spaces (Go files), tabs (Makefiles only)

## Build & Run

```bash
# Download dependencies
go mod download

# Start infrastructure (PostgreSQL, Redis, Etcd)
docker compose -f docker/Infrastructure/docker-compose.yaml up -d

# Run database migrations
./sql/migrate.sh          # Linux/macOS
./sql/migrate.ps1         # Windows PowerShell

# Start Identity service (terminal 1 — gRPC on :8080)
go run ./services/identity -f services/identity/etc/identity.yaml

# Start Gateway service (terminal 2 — HTTP on :8888)
go run ./services/gateway -f services/gateway/etc/gateway.yaml

# Verify
curl http://127.0.0.1:8888/healthz
curl http://127.0.0.1:8888/readyz
```

### Via goctl (dev mode)

```bash
goctl api go -api v3/api/gateway.api -dir services/gateway
goctl rpc protoc v3/proto/identity.proto -o services/identity --go_out=pb --go-grpc_out=pb --zrpc_out=.
```

## Test

```bash
# Full test suite
go test ./...

# Single package
go test ./services/identity/internal/service/...

# Single test
go test ./services/identity/internal/service/ -run TestRegisterLocalUser

# With race detection and coverage
go test -race -coverprofile=cover.out ./...

# Before committing, run review rules check
go run ./tools/reviewrules
# Or: ./scripts/check-review-rules.ps1

# QA suite (from qa/ directory, requires services running)
uv run pytest -v
```

## Architecture

### Service Topology

```
Client → [edge] → Gateway (HTTP :8888) → Identity (gRPC :8080)
                                    → Content  (gRPC)
                                    → File     (gRPC)
                                    → (search, indexer — future)
```

**Layer Boundary**:
- `gateway`: HTTP/WS entry, auth middleware, rate limiting, error wrapping, **no business logic, no direct DB access**
- `identity`: auth, sessions, tokens, SSO, user CRUD (gRPC-only)
- `content`: content items, revisions, tags, relations, outbox events (gRPC-only)
- `file`: file asset management, upload sessions, S3/local storage (gRPC + raw HTTP for upload)
- Backend services use go-zero zrpc with etcd discovery

### Inter-Service Authentication

Services call each other over gRPC using a **shared bearer token** scheme:
- `x-beehive-internal-auth-token`: shared secret (constant-time compare)
- `x-beehive-internal-caller`: caller service name
- `x-beehive-trusted-client-ip`, `x-beehive-user-id`, `x-beehive-session-id`, `x-beehive-user-role`: forwarded by gateway to downstream services
- Configured in each service's `etc/*.yaml` — `InternalAuthToken` must match across gateway and all downstream services

### Internal Service Layer Pattern (every service follows this)

```
main.go
internal/
  config/        — config structs + Validate(), no DB/RPC/MQ init
  svc/           — dependency assembly (DB, Redis, RPC clients, MQ), no business rules
  model/
    entity/      — GORM table structs (one table per file)
    repo/        — repositories wrapping a Store, all DB access through Store
  service/       — core use-case orchestration (tx boundaries, repo calls, audit, auth helpers)
                    ← main test entry for business logic
  logic/         — transport adaptation: extracts metadata, calls service, maps service errors → transport
                    ← smoke tests only (metadata extraction, error mapping)
  handler/       — (gateway only) HTTP handler funcs: bind request → call logic → return
  server/        — (RPC services only) gRPC server registration + internal auth interceptors
                    (gateway instead has handler/routes.go)
```

### `pkg/` Shared Packages

- `errs` — single source of truth for business error codes (6-digit numeric codes, segmented by service)
- `errs/httpx` — maps `errs` → HTTP JSON responses
- `errs/grpcx` — maps `errs` → gRPC status codes
- `logs` — structured logging wrapper (all logs via `logs.Ctx(ctx).Info/Debug/Warn/Error("action_name", fields...)`)
- `ctxmeta` — gRPC metadata keys, internal auth context, trusted proxy IP extraction, user claims forwarding
- `mq` — RabbitMQ abstraction (publish/consume)

### Error Handling Rules

- All business errors use `pkg/errs` codes (e.g., `errs.CodeIdentityUserNotFound`)
- Match errors with `errors.Is(err, errs.E(code))`, never with string matching
- Service layer returns domain errors only (no HTTP/gRPC codes)
- Logic/handler layers map domain errors to transport
- Prohibited: `err.Error()`, `strings.Contains(err.Error(), ...)`, direct `logx` imports

### Contract-First Development

- `v3/api/gateway.api` — HTTP contract + Swagger source of truth (goctl syntax v1)
- `v3/api/gateway.yaml` — generated Swagger (do not edit manually)
- `v3/proto/*.proto` — gRPC contracts
- `v3/api/file.api` + `v3/api/file.yaml` — File service contracts
- Generated code in `services/*/pb/`

## Database Migrations

Located in `sql/migrations/v3/<service>/` with sequential numbering:
- `020-028` — identity
- `030-034` — content
- `040-` — file

Run via `sql/migrate/main.go` (self-contained migration tool).

## Bilingual Comment Convention

All comments follow English-above, Chinese-below format:
```go
// Create a new session and persist the refresh token hash.
// 创建新会话并持久化 refresh token 哈希。
func CreateSession() {}
```

Runtime logs are English-only. Use `pkg/logs` for all logging (never import `logx` directly).

## Pre-Commit Checks

```bash
go run ./tools/reviewrules    # Checks for banned patterns:
                              #   direct logx imports
                              #   errs.IsCode() usage
                              #   string-based error matching
                              #   new service-layer violations
```

CI runs the same rules via `.github/workflows/review-rules.yml`.

## UI

The `ui/` directory contains a separate React/TypeScript frontend (pnpm, Vite, Playwright). The backend repo does not depend on the UI beyond the gateway serving as API server.
