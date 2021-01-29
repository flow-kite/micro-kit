package stack

import (
	"os"
	"path/filepath"
)

// 根据 os.Args[0] 获取该项目的路径
func GetRootService() string {
	path := filepath.Base(os.Args[0])
	return path
}
