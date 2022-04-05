package command

import (
	"container/cgroups/subsystem"
	"container/cmd/action"
	"fmt"
	"github.com/urfave/cli"
)

var RunCommand = cli.Command{
	Name:  "run",
	Usage: `运行经过 Namespace 隔离和 cgroup 资源限制后的容器`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "设置参数以限制容器进程内存使用量",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "设置参数以限制容器进程的 cpu 时间片",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "设置参数以限制容器进程所能使用的 cpu 数量",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		cmd := ctx.Args().Get(0)
		tty := ctx.Bool("it")

		mL := ctx.String("m")
		cShare := ctx.String("cpushare")
		cSet := ctx.String("cpuset")
		resCfg := subsystem.NewResourceConfig(mL, cShare, cSet)

		action.Run(tty, cmd, resCfg)
		return nil
	},
}
