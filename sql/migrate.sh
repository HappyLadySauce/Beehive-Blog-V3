#!/usr/bin/env bash
# Beehive-Blog 数据库迁移入口（Unix shell）
#
# 全覆盖（默认）：MODE=versioned
# 适应：MODE=adaptive  （详见 scripts/db/migrate/main.go 头部注释）
#
# 用法:
#   ./sql/migrate.sh
#   MODE=adaptive VERBOSE=1 ./sql/migrate.sh
#   DB_DSN='postgres://...' ./sql/migrate.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATIONS="${ROOT}/sql/migrations"
MODE="${MODE:-versioned}"
DSN="${DB_DSN:-postgres://Beehive-Blog:Beehive-Blog@127.0.0.1:5432/Beehive-Blog?sslmode=disable}"

GO_ARGS=(run ./scripts/db/migrate/main.go -dsn "$DSN" -dir "$MIGRATIONS" -mode "$MODE")
if [[ "${VERBOSE:-}" == "1" ]]; then
  GO_ARGS+=(-v)
fi

cd "$ROOT"
exec go "${GO_ARGS[@]}"
