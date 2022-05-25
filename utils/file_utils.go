package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// PathExist 判断文件是否存在。函数首先会判读对该目录是否具有访问权限，如果可以访问的话就判断文件是否存在，否则直接 return,注意：只有在无权限的情况下才会报错
func PathExist(Path string) (bool, error) {
	// syscall.Access 可以判断当前对目录是否有相应的访问权限 syscall.O_RDWR 这一参数是尝试当前程序是否可以读写该目录，F_OK 参数判断
	// 该目录是否存在,可以参考 https://man7.org/linux/man-pages/man2/access.2.html, golang 用例参考 https://golang.hotexamples.com/examples/syscall/-/Access/golang-access-function-examples.html
	err := syscall.Access(filepath.Dir(Path), syscall.O_RDWR)
	if err != nil {
		return false, fmt.Errorf("path exist error: can't judge path: %s, reason: %s", Path, err)
	}
	if _, err = os.Stat(Path); os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}
