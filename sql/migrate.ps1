# Beehive-Blog 数据库迁移入口（Windows PowerShell）
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
#   .\sql\migrate.ps1 -Mode adaptive -Verbose
#   .\sql\migrate.ps1 -Dsn 'postgres://user:pass@host:5432/dbname?sslmode=disable'

param(
    [ValidateSet('versioned', 'adaptive')][string]$Mode = 'versioned',
    [string]$Dsn = '',
    [switch]$Verbose
)

$ErrorActionPreference = 'Stop'
$RepoRoot = Split-Path -Parent $PSScriptRoot
$MigrationsDir = Join-Path $RepoRoot 'sql\migrations'

if (-not $Dsn) {
    $Dsn = $env:DB_DSN
}
if (-not $Dsn) {
    $Dsn = 'postgres://Beehive-Blog:Beehive-Blog@127.0.0.1:5432/Beehive-Blog?sslmode=disable'
}

$goArgs = @(
    'run', './scripts/db/migrate/main.go',
    '-dsn', $Dsn,
    '-dir', $MigrationsDir,
    '-mode', $Mode
)
if ($Verbose) {
    $goArgs += '-v'
}

Push-Location $RepoRoot
try {
    & go @goArgs
} finally {
    Pop-Location
}
