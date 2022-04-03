package cgroups

import (
	"container/cgroups/subsystem"
)

type cgroupManager struct {
	Pid        int
	CgroupPath string // Cgroup 的目录名
	resCfg     *subsystem.ResourceConfig
}

var subSysIns = []subsystem.Subsystem{
	&subsystem.MemorySubSystem{},
	&subsystem.CPUShareSubsystem{},
	&subsystem.CPUSetSubsystem{},
}

func NewCgroupManager(pid int, cgroupPath string, memoryLimit string, cpuShare string, cpuSet string) *cgroupManager {
	return &cgroupManager{
		Pid:        pid,
		CgroupPath: cgroupPath,
		resCfg: &subsystem.ResourceConfig{
			memoryLimit,
			cpuShare,
			cpuSet,
		},
	}
}

// SetCgroup 设置各种 cgroup
func (cm *cgroupManager) SetCgroup() error {
	for _, s := range subSysIns {

		err := s.Set(cm.CgroupPath, cm.resCfg, cm.Pid)
		if err != nil {
			return err
		}
	}
	return nil
}

// Remove 删除各种 cgroup
func (cm *cgroupManager) Remove() error {
	for _, s := range subSysIns {
		err := s.Remove(cm.CgroupPath, cm.resCfg)
		if err != nil {
			return err
		}
	}
	return nil
}
