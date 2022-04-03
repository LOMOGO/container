package utils

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// findCgroupRoot 返回该 subsystem cgroup 在系统中的绝对路径
func findCgroupRoot(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fileds := strings.Split(txt, " ")
		for _, opt := range strings.Split(fileds[len(fileds)-1], ",") {
			if opt == subsystem {
				return fileds[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

// SetCgroupPath 在指定的 subsystem Cgroup 文件树中设置指定的 cgroupPath 目录，并返回设置的 cgroup 在系统中的绝对路径。如果指定的目录已经存在，会返回
func SetCgroupPath(subsystem string, cgroupPath string) (string, error) {
	cgroupRoot := findCgroupRoot(subsystem)
	resultPath := path.Join(cgroupRoot, cgroupPath)
	_, err := os.Stat(resultPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(resultPath, 0755)
		if err != nil {
			return "", fmt.Errorf("set cgroup path error: %v", err)
		}
		return resultPath, nil
	}
	if err != nil {
		return "", fmt.Errorf("set cgroup path error: %v", err)
	}
	return resultPath, nil
}

// GetCgroupPath 查询某 subsystem Cgourp 中目录为 cgroupPath 的 cgroup 在系统中的绝对路径，如果不存在，就报错
func GetCgroupPath(subsystem string, cgroupPath string) (string, error) {
	cgroupRoot := findCgroupRoot(subsystem)
	resultPath := path.Join(cgroupRoot, cgroupPath)
	_, err := os.Stat(resultPath)
	if os.IsNotExist(err) {
		return "", err
	}
	return resultPath, nil
}
