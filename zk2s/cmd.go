package zk2s

import (
	"github.com/urfave/cli"
)

var CMD_Run = cli.Command{
	Name:   "run",
	Usage:  "run the zk2s application",
	Action: Run,
	Flags:  []cli.Flag{Flag_Template},
}

var Flag_Template = cli.StringFlag{
	Name:  "template, t",
	Usage: "set the path to the template file (defaults to response.tmpl in working directory)",
	Value: "response.tmpl",
}
