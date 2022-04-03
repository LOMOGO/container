package main

import (
	"container/cmd/command"
	"container/initialize"
	"github.com/docker/docker/pkg/reexec"
	// reexec 是 docker 的一个功能包，它的功能是重新调用自己，下面将分别用到这个包的 Register、Init 和 Command 包，功能介绍见下方注释
	"github.com/urfave/cli"
	"log"
	"os"
)

const usage = `简易版容器引擎，类似 docker，`

func init() {
	// reexec.Register 这个函数会将注册的函数保存到内存里面，以便后续的调用，参数的格式为 函数名标识+函数体, 这个函数貌似
	// 要在 reexec.Init 之前运行
	reexec.Register("nsInitialisation", initialize.NsInitialisation)

	// reexec.Init 这个函数的功能有两个，首先它会判断进程在重新执行自身后的进程有没有在运行，（如果返回 true 就说明当前进程是调用自身后
	// 执行的进程，false 就是第一次运行），如果返回 true 的话，这个 reexec.Init 内部会执行下面 reexec.Command 中标明的已注册到
	// reexec.Register 函数中的函数。
	if reexec.Init() {
		os.Exit(0)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "container"
	app.Usage = usage

	app.Commands = []cli.Command{
		command.RunCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
