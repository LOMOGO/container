package action

import (
	"container/cgroups"
	"container/cgroups/subsystem"
	"container/utils"
	"flag"
	"github.com/docker/docker/pkg/reexec"
	"os"
	"syscall"

	"log"
)

func Run(tty bool, command string, resCfg *subsystem.ResourceConfig) {
	var imageFilePath string
	flag.StringVar(&imageFilePath, "imageFilePath", "/home/lomogo/playground/unionfs_ubuntu/ubuntu.tar", "镜像文件地址")

	// 创建联合系统
	newrootPath, err := utils.CreateUnionFS(imageFilePath)
	if err != nil {
		log.Fatalf("Error Run: %s\n", err)
	}
	defer utils.UnmountUnionFS(imageFilePath)
	// reexec.Command 这个函数会返回一个 *exec.Cmd 结构体，但是这个结构体要执行的程序被指向自身，因此这个进程在 clone 之后会
	// 进行一个自己调用自己的动作，这在 exec.Cmd{Path:/proc/self/exe} 中声明，然后这个函数需要添加一个参数，这个参数是在
	// reexec.Register 中注册的函数的名称
	cmd := reexec.Command("nsInitialisation", newrootPath, command)
	// 之所以使用 reexec 方法重新调用自身是为了使 nsRun 中运行的进程不可能获取到宿主机的信息。因为这个 nsRun 中运行的程序总是在
	// 创建好了 Namespace 之后运行的

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	/*
		| Namespace | 用途                       |
		| --------- | -------------------------- |
		| Cgroup    | 隔离 Cgroup 目录           |
		| IPC       | 隔离进程间通信（IPC）资源  |
		| Network   | 隔离网络接口               |
		| Mount     | 隔离文件系统挂载点         |
		| PID       | 隔离进程 ID (PID) 号码空间 |
		| Time      | 隔离单调时钟、启动时钟     |
		| User      | 隔离 UID/GID 号码空间      |
		| UTS       | 隔离主机名和域名           |
	*/
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		// 下面的代码将启动的容器在宿主机上的 uid 和 gid 与 root 用户的 uid 和 gid 进行映射，使得容器认为自己获得了 root 权限
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("error starting the reexec.Command: %s", err)
	}

	cm := cgroups.NewCgroupManager(cmd.Process.Pid, "container-test", resCfg)
	defer cm.Remove()
	err = cm.SetCgroup()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("error waiting for the reexec.Command: %s", err)
	}
}
