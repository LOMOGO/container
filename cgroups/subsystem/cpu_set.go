package subsystem

import (
	"container/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

type CPUSetSubsystem struct{}

func (css *CPUSetSubsystem) Name() string {
	return "cpuset"
}

func (css *CPUSetSubsystem) Set(cgroupPath string, res *ResourceConfig, pid int) error {
	if res.CPUSet == "" {
		return nil
	}
	subSysCgroupPath, err := utils.SetCgroupPath(css.Name(), cgroupPath)
	if err != nil {
		return fmt.Errorf("set cpuset cgroup path: %v, fail: %v", subSysCgroupPath, err)
	}

	err = ioutil.WriteFile(path.Join(subSysCgroupPath, "cpuset.cpus"), []byte(res.CPUSet), 0644)
	if err != nil {
		return fmt.Errorf("set cpuset cgroup path: %v, fail: %v", subSysCgroupPath, err)
	}

	if err = ioutil.WriteFile(path.Join(subSysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("apply cpuset cgroup path: %v, fail: %v", subSysCgroupPath, err)
	}

	log.Printf("pid : %v 已经加入到 cpuset cgroup path: %v 中\n", pid, subSysCgroupPath)
	return nil
}

func (css *CPUSetSubsystem) Remove(cgroupPath string, res *ResourceConfig) error {
	if res.CPUSet == "" {
		return nil
	}
	subSysCgroupPath, err := utils.GetCgroupPath(css.Name(), cgroupPath)
	if err != nil {
		return fmt.Errorf("remove cpuset cgroup path: %v error: %v", subSysCgroupPath, err)
	}
	err = os.RemoveAll(subSysCgroupPath)
	if err != nil {
		return fmt.Errorf("remove cpuset cgroup path: %v error: %v", subSysCgroupPath, err)
	}
	return nil
}
