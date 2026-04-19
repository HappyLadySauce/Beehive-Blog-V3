# Beehive-Blog-V3 数据库迁移入口（Windows PowerShell）
#
# 全覆盖（默认）：-Mode versioned
#   每个 .sql 文件在一个事务内整段执行，schema_migrations 记录 checksum，适合空库或严格与仓库一致。
#
# 适应：-Mode adaptive
#   将单个迁移文件按语句拆分执行；遇到表/列/对象已存在等 SQLSTATE 时跳过该句并继续，
#   用于已有部分表结构、重复执行或半旧库向前对齐（仍有风险，重要环境请先备份）。
#
# 用法（在仓库根目录执行亦可）:
#   .\sql\migrate.ps1
#     默认仅 sql/migrations/v3（递归，如 v3/identity/*.sql），不跑 v2。
#   .\sql\migrate.ps1 -MigrationsScope all
#     扫描 sql/migrations 下全部版本（v2 + v3）。
#   .\sql\migrate.ps1 -MigrationsScope v2
#     仅 sql/migrations/v2。
#   .\sql\migrate.ps1 -Mode adaptive -Verbose
#   .\sql\migrate.ps1 -Dsn 'postgres://user:pass@host:5432/dbname?sslmode=disable'
#   .\sql\migrate.ps1 -Force
#     迁移文件改过、与库内 checksum 不一致时仍执行并覆盖 schema_migrations。
#   .\sql\migrate.ps1 -Reapply -Mode adaptive
#     已应用过的迁移再跑一遍（多为 DML；含 CREATE 时常需 adaptive）。

param(
    [ValidateSet('versioned', 'adaptive')][string]$Mode = 'versioned',
    [string]$Dsn = '',
    [ValidateSet('v3', 'v2', 'all')][string]$MigrationsScope = 'v3',
    [switch]$Force,
    [switch]$Reapply,
    [switch]$Verbose
)

$ErrorActionPreference = 'Stop'
$RepoRoot = Split-Path -Parent $PSScriptRoot
$MigrationsCatalog = Join-Path $RepoRoot 'sql\migrations'
$MigrationsDir = $MigrationsCatalog
switch ($MigrationsScope) {
    'v2' { $MigrationsDir = Join-Path $MigrationsDir 'v2' }
    'v3' { $MigrationsDir = Join-Path $MigrationsDir 'v3' }
    'all' { }
}

if (-not $Dsn) {
    $Dsn = $env:DB_DSN
}
if (-not $Dsn) {
    $Dsn = 'postgres://Beehive-Blog-V3:Beehive-Blog-V3@127.0.0.1:5432/Beehive-Blog-V3?sslmode=disable'
}

$goArgs = @(
    'run', './sql/migrate/main.go',
    '-dsn', $Dsn,
    '-dir', $MigrationsDir,
    '-catalog', $MigrationsCatalog,
    '-mode', $Mode
)
if ($Verbose) {
    $goArgs += '-v'
}
if ($Force) {
    $goArgs += '-force'
}
if ($Reapply) {
    $goArgs += '-reapply'
}

Push-Location $RepoRoot
try {
    & go @goArgs
} finally {
    Pop-Location
}
