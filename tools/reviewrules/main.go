package main

import (
	"fmt"
	"os"
)

// Run the repository review rule checks and exit with a process status.
// 运行仓库级 code review 规则检查，并以进程状态码退出。
func main() {
	if err := Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
