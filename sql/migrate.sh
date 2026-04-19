#!/usr/bin/env bash
# Beehive-Blog-V3 数据库迁移入口（Unix shell）
#
# 全覆盖（默认）：MODE=versioned
# 适应：MODE=adaptive  （详见 sql/migrate/main.go 头部注释）
#
# 用法:
#   ./sql/migrate.sh
#     默认仅 sql/migrations/v3（递归），不跑 v2。
#   MIGRATIONS_SCOPE=all ./sql/migrate.sh   # v2 + v3
#   MIGRATIONS_SCOPE=v2 ./sql/migrate.sh    # 仅 v2
#   MODE=adaptive VERBOSE=1 ./sql/migrate.sh
#   DB_DSN='postgres://...' ./sql/migrate.sh
#   MIGRATION_FORCE=1 ./sql/migrate.sh
#     改过迁移 SQL 后与库 checksum 不一致时仍执行并覆盖记录。
#   MIGRATION_REAPPLY=1 MODE=adaptive ./sql/migrate.sh
#     已应用的迁移再执行一遍（多为 DML）。

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATIONS_CATALOG="${ROOT}/sql/migrations"
MIGRATIONS="${MIGRATIONS_CATALOG}"
case "${MIGRATIONS_SCOPE:-v3}" in
  v2) MIGRATIONS="${MIGRATIONS}/v2" ;;
  v3) MIGRATIONS="${MIGRATIONS}/v3" ;;
  all) ;;
  *)
    echo "migration error: unknown MIGRATIONS_SCOPE=${MIGRATIONS_SCOPE:-} (use v3, v2, or all)" >&2
    exit 1
    ;;
esac
MODE="${MODE:-versioned}"
DSN="${DB_DSN:-postgres://Beehive-Blog-V3:Beehive-Blog-V3@127.0.0.1:5432/Beehive-Blog-V3?sslmode=disable}"

GO_ARGS=(run ./sql/migrate/main.go -dsn "$DSN" -dir "$MIGRATIONS" -catalog "$MIGRATIONS_CATALOG" -mode "$MODE")
if [[ "${VERBOSE:-}" == "1" ]]; then
  GO_ARGS+=(-v)
fi
if [[ "${MIGRATION_FORCE:-}" == "1" ]]; then
  GO_ARGS+=(-force)
fi
if [[ "${MIGRATION_REAPPLY:-}" == "1" ]]; then
  GO_ARGS+=(-reapply)
fi

cd "$ROOT"
exec go "${GO_ARGS[@]}"
