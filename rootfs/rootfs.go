package rootfs

import (
	"os"
	"path/filepath"
	"syscall"
)

// PivotRoot 可以为当前调用的进程设置一个新的根文件系统，通过 syscall.PivotRoot 来实现，Pivot 必须在新的 Mount Namespace 中进行，
// 否则会改变宿主机的根文件目录 “/”
func PivotRoot(newroot string) error {
	// pivot_root 操作需要两个目录，分别是 newroot 和 put_old，newroot 是新 rootfs 的路径，put_old 是一个目录的路径，pivot_root
	// 操作对这两个目录有一些限制要求：1. 他们都必须是目录。2. 他们与当前的根目录不能在同一个文件系统上。3. put_old 目录必须在 newroot
	// 目录下面。4. 其他文件系统不能挂载到 putold 上
	putold := filepath.Join(newroot, ".pivot_root")

	// 为了满足 pivot_root 调用的第二点要求，下面的操作将会使 newroot 目录通过 syscall.MS_BIND 绑定在自身。
	if err := syscall.Mount(
		newroot,
		newroot,
		"bind",
		syscall.MS_BIND|syscall.MS_REC,
		"",
	); err != nil {
		return err
	}

	// create putold directory, 这里需要注意一点，因为需要在宿主机的 root 用户下执行 container 这个二进制文件（因为 cgroup 这个
	// 操作必须在宿主机的 /sys/fs/cgroup 下面执行，所以必须用到 root 权限），如果新的文件系统文件夹所属的用户不是 root 的话，就会报
	// 因无权限无法创建文件夹的错误。
	if err := os.MkdirAll(putold, 0777); err != nil {
		return err
	}

	// 调用 pivot_root,将旧的 root 文件系统移动到 putold 中，然后将 newroot 中的新的 root 文件系统挂载进来
	if err := syscall.PivotRoot(newroot, putold); err != nil {
		return err
	}

	// os.Chdir 命令会将当前工作目录更改为指定的目录，这一操作将确保当前的工作目录被正确设置为 newroot
	if err := os.Chdir("/"); err != nil {
		return err
	}

	// umount putold, putold 现在位于 .pivot_root 中
	putold = ".pivot_root"
	if err := syscall.Unmount(putold, syscall.MNT_DETACH); err != nil {
		return err
	}

	// remove putold
	if err := os.RemoveAll(putold); err != nil {
		return err
	}
	return nil
}
