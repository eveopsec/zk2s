package config

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var CMD_Config = cli.Command{
	Name:  "configure",
	Usage: "configure zk2s application to be run",
	Subcommands: []cli.Command{
		cmd_Assistant,
		cmd_New,
	},
}

var cmd_New = cli.Command{
	Name:      "new",
	Usage:     "create new configuration (will ask before overwriting pre-existing)",
	UsageText: "create new configuration (will ask before overwriting pre-existing)",
}

var cmd_Assistant = cli.Command{
	Name:      "assistant",
	Usage:     "setup assistant to configure the slackbot",
	UsageText: "setup assistant to configure the slackbot",
	Action:    RunAssistant,
}

func new(c *cli.Context) error {
	fmt.Println("Creating new configuration file...")
	cfg, err := LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			cfg = &Configuration{}
			err = cfg.Save()
			if err != nil {
				fmt.Printf("ERROR - unable to save new configuration file with error: %v\n", err)
			}
			return err
		}
	}
	return nil
}
