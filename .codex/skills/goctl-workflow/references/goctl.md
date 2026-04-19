# goctl Reference

## Standard gateway API generation

Validate and generate gateway server code from `v3/api/gateway.api`:

```powershell
goctl api validate -api .\v3\api\gateway.api
goctl api go -api .\v3\api\gateway.api -dir services\gateway
```

Use this when gateway HTTP contracts change.

## Standard gateway Swagger generation

Generate Swagger from `v3/api/gateway.api`:

```powershell
goctl api swagger --api .\v3\api\gateway.api --dir .\v3\api\swagger --filename gateway --yaml
```

Use this after gateway contract changes so API docs stay in sync.

## Proto `go_package`

In each `v3/proto/<service>.proto`, set for example:

```protobuf
option go_package = "services/<service>/pb";
```

Use paths relative to the repo root. Avoid full module paths in `go_package` (they can cause goctl to emit files under unintended `github.com/...` directories).

## Standard RPC generation

Generate RPC code from `v3/proto/*.proto` and always include `--client=true`:

```powershell
goctl rpc protoc .\v3\proto\identity.proto --go_out=. --go-grpc_out=. --zrpc_out=services\identity --client=true
goctl rpc protoc .\v3\proto\content.proto --go_out=. --go-grpc_out=. --zrpc_out=services\content --client=true
goctl rpc protoc .\v3\proto\search.proto --go_out=. --go-grpc_out=. --zrpc_out=services\search --client=true
```

Replace the proto path and target service directory as needed.

## Recommended v3 workflow

1. Decide which service owns the capability.
2. If it is a backend capability, edit `v3/proto/<service>.proto` first.
3. If it is an external HTTP route, edit `v3/api/gateway.api`.
4. Run `goctl api validate` for gateway contract changes.
5. Run `goctl`.
6. Regenerate Swagger for gateway contract changes.
7. Adjust generated code and handwritten business code.
8. Verify compile and wiring.

## Service ownership guidance

- `gateway`: only HTTP / WebSocket entry and transport glue
- `identity`: RPC
- `content`: RPC
- `search`: RPC
- additional no-clear-owner capabilities: create a dedicated RPC service such as `dashboard`, `query`, or `composite`

## Rules

- Do not use `goctl` as a substitute for service design
- Do not put business orchestration into `gateway`
- Do not manually edit contracts after generation and skip regeneration
- Do not use outdated repo assumptions such as root-level `api/` / `proto/`, RPC generation without `--client=true`, or mixed API+RPC services unless the repository is explicitly changed later
