package main

import (
	"log"
	"os"

	"github.com/eveopsec/zk2s/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Authors = cmd.CONTRIBUTORS
	app.Version = cmd.VERSION
	app.Name = "zk2s"
	app.Usage = "A Slack bot for posting kills from zKillboard to slack in near-real time."
	app.UsageText = "To start zk2s, run 'zk2s start' and pass '--config' if the configuration file is not in the default location."
	app.Commands = []cli.Command{
		cmd.CMD_Start,
	}
	if err := app.Run(os.Args); err != nil {
		log.Println("[FATAL] Application exited with error.")
	}
}
