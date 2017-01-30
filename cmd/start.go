package cmd

import (
	"log"

	"github.com/eveopsec/zk2s/app"
	"github.com/eveopsec/zk2s/app/config"

	"github.com/urfave/cli"
)

var (
	CMD_Start = cli.Command{
		Name:   "start",
		Usage:  "start the zk2s application",
		Action: start,
		Flags:  []cli.Flag{flagConfig},
	}
	flagConfig = cli.StringFlag{
		Name:  "config,c",
		Usage: "Set the file path to the conffiguration file. Defaults to 'cfg.zk2s.json' in working directory.",
		Value: "cfg.zk2s.json",
	}
)

func start(c *cli.Context) (err error) {
	log.Printf("%v version %v", c.App.Name, c.App.Version)
	var cfg *config.Configuration
	var application *app.App

	// Load Configuration
	if cfg, err = config.ReadConfig(c.String("config")); err != nil {
		log.Printf("[FATAL] Failed to load configuration with error: %v", err)
		return err
	}

	// Create a new App.
	if application, err = app.NewApp(cfg); err != nil {
		log.Printf("[FATAL] Failed to initialize with error: %v", err)
		return err
	}

	// Start the App
	if err = application.Start(); err != nil {
		log.Printf("[FATAL] Failed to start with error: %v", err)
		return err
	}

	log.Println("[NOTIFY] Running... ")
	select {}
}
