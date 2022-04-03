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

type CPUShareSubsystem struct{}

func (css *CPUShareSubsystem) Name() string {
	return "cpu"
}

func (css *CPUShareSubsystem) Set(cgroupPath string, res *ResourceConfig, pid int) error {
	if res.CPUShare == "" {
		return nil
	}
	subSysCgroupPath, err := utils.SetCgroupPath(css.Name(), cgroupPath)
	if err != nil {
		return fmt.Errorf("set cpu cgroup path: %v fail: %v", subSysCgroupPath, err)
	}

	err = ioutil.WriteFile(path.Join(subSysCgroupPath, "cpu.shares"), []byte(res.CPUShare), 0644)
	if err != nil {
		return fmt.Errorf("set cpu cgroup: %v fail: %v", subSysCgroupPath, err)
	}

	if err = ioutil.WriteFile(path.Join(subSysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("apply pid: %v join cpu cgroup path: %v error:%v", pid, subSysCgroupPath, err)
	}

	log.Printf("pid : %v 已经加入到 cpu cgroup path: %v 中\n", pid, subSysCgroupPath)
	return nil
}

func (css *CPUShareSubsystem) Remove(cgroupPath string, res *ResourceConfig) error {
	if res.CPUShare == "" {
		return nil
	}
	subSysCgroupPath, err := utils.GetCgroupPath(css.Name(), cgroupPath)
	if err != nil {
		return fmt.Errorf("remove cpu cgroup path: %v error:%v", subSysCgroupPath, err)
	}
	err = os.RemoveAll(subSysCgroupPath)
	if err != nil {
		return fmt.Errorf("remove cpu cgroup path: %v error:%v", subSysCgroupPath, err)
	}
	return nil
}
