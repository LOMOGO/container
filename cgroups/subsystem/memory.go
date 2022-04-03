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

type MemorySubSystem struct{}

func (mss *MemorySubSystem) Name() string {
	return "memory"
}

func (mss *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig, pid int) error {
	if res.MemoryLimit == "" {
		return nil
	}
	fmt.Println(res.MemoryLimit)
	subSysCgroupPath, err := utils.SetCgroupPath(mss.Name(), cgroupPath)
	if err != nil {
		return fmt.Errorf("set memory cgroup path: %v fail %v", subSysCgroupPath, err)
	}

	err = ioutil.WriteFile(path.Join(subSysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644)
	if err != nil {
		return fmt.Errorf("set memory cgroup path: %v fail %v", subSysCgroupPath, err)
	}

	// 通过将 pid 写入到 Cgroup path 下面的 tasks 文件中来达到将某进程加入到某 cgroup 中的效果
	if err = ioutil.WriteFile(path.Join(subSysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("apply pid:%v join memory cgroup path:%v error: %v", pid, subSysCgroupPath, err)
	}
	log.Printf("pid : %v 已经加入到 memory cgroup path: %v 中\n", pid, subSysCgroupPath)

	return nil
}

func (mss *MemorySubSystem) Remove(cgroupPath string, res *ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}
	subSysCgroupPath, err := utils.GetCgroupPath(mss.Name(), cgroupPath)
	if err != nil {
		return fmt.Errorf("remove memory cgroup path: %v error:%v", subSysCgroupPath, err)
	}
	err = os.RemoveAll(subSysCgroupPath)
	if err != nil {
		return fmt.Errorf("remove memory cgroup path: %v error:%v", subSysCgroupPath, err)
	}
	return nil
}
