package config

import "github.com/urfave/cli"

var CMD_Config = cli.Command{
	Name:        "configure",
	Usage:       "configure zk2s application to be run",
	Subcommands: []cli.Command{},
}
