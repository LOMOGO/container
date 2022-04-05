package cgroups

import (
	"container/cgroups/subsystem"
)

type CgroupManager struct {
	Pid        int
	CgroupPath string // Cgroup 的目录名
	resCfg     *subsystem.ResourceConfig
}

var subSysIns = []subsystem.Subsystem{
	&subsystem.MemorySubSystem{},
	&subsystem.CPUShareSubsystem{},
	&subsystem.CPUSetSubsystem{},
}

func NewCgroupManager(pid int, cgroupPath string, resCfg *subsystem.ResourceConfig) *CgroupManager {
	return &CgroupManager{
		Pid:        pid,
		CgroupPath: cgroupPath,
		resCfg:     resCfg,
	}
}

// SetCgroup 设置各种 cgroup
func (cm *CgroupManager) SetCgroup() error {
	for _, s := range subSysIns {

		err := s.Set(cm.CgroupPath, cm.resCfg, cm.Pid)
		if err != nil {
			return err
		}
	}
	return nil
}

// Remove 删除各种 cgroup
func (cm *CgroupManager) Remove() error {
	for _, s := range subSysIns {
		err := s.Remove(cm.CgroupPath, cm.resCfg)
		if err != nil {
			return err
		}
	}
	return nil
}
