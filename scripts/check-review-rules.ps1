# Run the repository review rule checks.
# 运行仓库级 code review 规则检查。
param()

$ErrorActionPreference = "Stop"

go run ./tools/reviewrules
