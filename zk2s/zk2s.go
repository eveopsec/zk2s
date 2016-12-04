package zk2s

import (
	"errors"
	"log"

	"github.com/eveopsec/zk2s/zk2s/config"
	"github.com/eveopsec/zk2s/zk2s/slack"

	"github.com/eveopsec/zk2s/zk2s/tmpl"
	"github.com/urfave/cli"
	"github.com/vivace-io/evelib/redisq"
)

var CMD_Start = cli.Command{
	Name:   "start",
	Usage:  "start the zk2s application",
	Action: start,
	Flags:  []cli.Flag{flag_Template},
}

var flag_Template = cli.StringFlag{
	Name:  "template, t",
	Usage: "set the path to the template file (defaults to response.tmpl in working directory)",
	Value: "response.tmpl",
}

func start(c *cli.Context) error {
	log.Printf("%v version %v", c.App.Name, c.App.Version)

	// [1] - Init config
	if err := config.Init(c); err != nil {
		log.Printf("[FATAL] Unable to load configuration: %v", err)
		return err
	}

	// [2] - Init tmpl
	if err := tmpl.Init(c); err != nil {
		log.Printf("[FATAL] Unable to load templates: %v", err)
		return err
	}

	// [3] - Init slack
	if err := slack.Init(c); err != nil {
		log.Printf("[FATAL] Unable create slack connections: %v", err)
		return err
	}

	// [4] - Listen to redisq
	if err := listen(); err != nil {
		log.Printf("[FATAL] Unable to listen on RedisQ: %v", err)
		return err
	}

	select {}
}

func listen() error {
	if config.CONFIG == nil {
		return errors.New("app configuration was nil")
	}
	client, err := redisq.NewClient(nil)
	if err != nil {
		return err
	}
	client.AddFunc(slack.Recieve)
	client.Listen()
	return nil
}
