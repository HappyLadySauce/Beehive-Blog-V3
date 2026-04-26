# Beehive-Blog-V3 development seed entrypoint.
# Beehive-Blog-V3 开发种子数据入口。
#
# This script is intended for local development after recreating the database.
# 本脚本用于本地开发删库重刷后的种子数据写入。
#
# Usage:
#   .\sql\migrate.ps1
#   .\sql\seed.ps1

param(
    [string]$Dsn = '',
    [ValidateSet('v3', 'all')][string]$SeedsScope = 'v3',
    [switch]$Force,
    [switch]$Reapply,
    [switch]$Verbose
)

$ErrorActionPreference = 'Stop'
$RepoRoot = Split-Path -Parent $PSScriptRoot
$SeedsCatalog = Join-Path $RepoRoot 'sql\seeds'
$SeedsDir = $SeedsCatalog
switch ($SeedsScope) {
    'v3' { $SeedsDir = Join-Path $SeedsDir 'v3' }
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
    '-dir', $SeedsDir,
    '-catalog', $SeedsCatalog,
    '-mode', 'versioned'
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
