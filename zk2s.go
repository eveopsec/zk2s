package main

import (
	"os"

	"github.com/eveopsec/zk2s/zk2s"
	"github.com/eveopsec/zk2s/zk2s/config"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Authors = zk2s.CONTRIBUTORS
	app.Version = zk2s.VERSION
	app.Name = "zk2s"
	app.Usage = "A Slack bot for posting kills from zKillboard to slack in near-real time."
	app.Commands = []cli.Command{
		zk2s.CMD_Start,
		config.CMD_Config,
	}
	app.Run(os.Args)
}
