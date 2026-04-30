# Beehive Blog v3 UI

Vue 3 frontend for Beehive Blog v3. The first phase connects only to the gateway HTTP contract for identity and authentication.

## Commands

```powershell
pnpm install
pnpm dev
pnpm typecheck
pnpm test -- --run
pnpm test:e2e
```

## API Mode

- **Default (unset or any value other than `mock`):** `live` — Studio and auth call the gateway at `/api/v3/*` (e.g. `/api/v3/auth/*`).
- **`VITE_API_MODE=mock`:** in-memory auth/studio/uploads for UI-only work or tests that inject this explicitly (unit tests, Playwright `webServer`).
- `VITE_GATEWAY_BASE_URL=` stays empty in local dev so Vite proxies `/api`, `/healthz`, and `/readyz` to `http://127.0.0.1:8888`.

The UI must not call identity RPC or any service RPC directly.
