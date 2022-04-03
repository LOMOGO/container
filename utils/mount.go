package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// MountProc 对容器的新环境挂载 proc 文件系统
func MountProc(newroot string) error {
	source := "proc"
	target := filepath.Join(newroot, "/proc")
	fstype := "proc"
	flags := 0
	data := ""

	if err := os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("mountProc error: %v", err)
	}
	if err := syscall.Mount(
		source,
		target,
		fstype,
		uintptr(flags),
		data,
	); err != nil {
		return fmt.Errorf("mountProc error: %v", err)
	}

	return nil
}

// CreateUnionFS 以 overlay2 文件系统为基础，创建一个联合文件系统，返回联合文件系统的挂载点路径。
func CreateUnionFS(imageFilePath string) (string, error) {
	if isFileExist, err := PathExist(imageFilePath); err != nil {
		return "", fmt.Errorf("create union file system error: %s", err)
	} else if !isFileExist {
		return "", fmt.Errorf("create union file system error: image file:%s dose not exist", imageFilePath)
	}
	currentPath := filepath.Dir(imageFilePath)

	lowerDir := filepath.Join(currentPath, "lower")
	upperDir := filepath.Join(currentPath, "upper")
	mergeDir := filepath.Join(currentPath, "merge")
	workDir := filepath.Join(currentPath, "work")

	// 判断创建 overlay 联合文件系统所需要的文件目录是否存在，如果不存在的话就创建。
	mkdirPathList := []string{lowerDir, upperDir, mergeDir, workDir}
	for _, path := range mkdirPathList {
		if isFileExist, _ := PathExist(path); !isFileExist {
			err := os.MkdirAll(path, 0744)
			if err != nil {
				return "", fmt.Errorf("create union file system error: %s", err)
			}
		}
	}

	// 将给定的镜像文件解压缩，解压到 lower 目录中
	if _, err := exec.Command("tar", "-xvf", imageFilePath, "-C", lowerDir).CombinedOutput(); err != nil {
		return "", fmt.Errorf("create union file system error: excute 'tar -xvf %s, -C %s' error: %s", imageFilePath, lowerDir, err)
	}

	args := fmt.Sprintf("-t overlay overlay -o lowerdir=%s -o upperdir=%s -o workdir=%s %s", lowerDir, upperDir, workDir, mergeDir)

	mountCmd := exec.Command("mount", strings.Split(args, " ")...)
	if info, err := mountCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("create union file system error: execute command 'mount %s' error:%s", args, info)
	}
	return mergeDir, nil
}

// RemoveUnionFS 卸载联合文件系统
func UnmountUnionFS(imageFilePath string) error {
	currentPath := filepath.Dir(imageFilePath)

	// lowerDir := filepath.Join(currentPath, "lower")
	// upperDir := filepath.Join(currentPath, "upper")
	mergeDir := filepath.Join(currentPath, "merge")
	// workDir := filepath.Join(currentPath, "work")

	if err := syscall.Unmount(mergeDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("remove union fs error: %s", err)
	}

	// mkdirPathList := []string{lowerDir, upperDir, mergeDir, workDir}
	// for _, path := range mkdirPathList {
	// 	if err := os.RemoveAll(path); err != nil {
	// 		fmt.Println(err)
	// 		return fmt.Errorf("remove union fs error: %s", err)
	// 	}
	// }
	return nil
}
