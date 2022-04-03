package initialize

import (
	"container/rootfs"
	"container/utils"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// NsInitialisation 对容器里设置的一些 Namespace 进行初始化，注意以下部分都是在新的容器环境中执行的
func NsInitialisation() {
	newrootPath := os.Args[1]
	command := os.Args[2]

	if err := utils.MountProc(newrootPath); err != nil {
		log.Fatalf("Error mounting /proc - %s\n", err)
	}

	// 这一步进行的是将 root 文件系统更换为新的 rootfs 比如 busybox 等，请注意，在这一步执行前，容器的挂载内容是宿主机的副本，
	// 因此调用 exec.LookPath 方法找寻命令的绝对路径需要放到这一步后面进行，否则将会寻找宿主机里面的命令的绝对路径，这样后面在执行
	// 这个绝对路径中的 command 的时候，容器的根文件系统中未必能够根据传来的路径找到这个命令，然会就会导致报错。
	if err := rootfs.PivotRoot(newrootPath); err != nil {
		log.Fatalf("Error running pivot_root - %s\n", err)
	}

	if err := syscall.Sethostname([]byte("container")); err != nil {
		log.Fatalf("Error setting hostname - %s\n", err)
	}

	// LookPath 的功能是可以搜索简写的命令程序的绝对路径
	command, err := exec.LookPath(command)
	if err != nil {
		log.Fatalf("Error exec.LookPath - %s", err)
	}

	nsRun(command)
}

// nsRun 这个函数里面执行的才是容器中真正执行的程序
func nsRun(command string) {

	args := []string{command}
	// syscall.Exec 会执行参数指定的命令，但是不重新创建新的进程，只在当前进程空间内执行，也就是替换当前进程的执行内容，
	// 他们重用同一个进程号 PID。syscall.Exec 需要三个参数：1. 第一个参数为待执行命令的绝对路径地址。2. 第二个参数
	// 是一个 string list，其中 list[0] 是待运行程序的名称(任意，手动输入)，list[1] 是待执行命令需要的参数。
	// 3. 第三个参数是环境变量列表
	if err := syscall.Exec(command, args, os.Environ()); err != nil {
		log.Fatalf("Error running the %s command:%s\n", command, err)
	}
}
