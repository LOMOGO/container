package subsystem

// ResourceConfig 用于传递资源限制配置的结构体，包括内存限制、CPU 时间片设置、CPU核心数设置
type ResourceConfig struct {
	MemoryLimit string
	CPUShare    string
	CPUSet      string
}

// Subsystem 每个 Subsystem 可以实现以下四个接口
type Subsystem interface {
	// Name 返回 subsystem 的名字
	Name() string
	// Set 设置某个 cgroup 在这个 subsystem 中的资源限制, 并将当前进程加入进去
	Set(path string, res *ResourceConfig, pid int) error
	// Remove 移除某个 cgroup
	Remove(path string, res *ResourceConfig) error
}
