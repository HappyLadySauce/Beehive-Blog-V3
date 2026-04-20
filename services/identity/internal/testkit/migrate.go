package testkit

import "os"

func osReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
