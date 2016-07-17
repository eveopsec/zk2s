package main

import (
	"log"
	"os"

	"github.com/eveopsec/zk2s/config"
	"github.com/nlopes/slack"
	"github.com/urfave/cli"
	"github.com/vivace-io/evelib/zkill"
	"github.com/vivace-io/gonfig"
)

/* zk2s.go
 * Main entrypoint and controller for zk2s
 */

const VERSION = "0.5"

var CONTRIBUTORS = []cli.Author{
	cli.Author{
		Name: "Nathan \"Vivace Naaris\" Morley",
	},
	cli.Author{
		Name: "\"Zuke\"",
	},
}

var cfg *config.Configuration
var bot *slack.Client

func main() {
	app := cli.NewApp()
	app.Authors = CONTRIBUTORS
	app.Version = VERSION
	app.Name = "zk2s"
	app.Usage = "A Slack bot for posting kills from zKillboard to slack in near-real time."
	app.Commands = []cli.Command{
		cmdRun,
		config.CMD_Config,
	}
	app.Run(os.Args)
}

var cmdRun = cli.Command{
	Name:   "run",
	Usage:  "run the zk2s application",
	Action: Run,
}

// Run zk2s
func Run(c *cli.Context) error {
	log.Printf("%v version %v", c.App.Name, c.App.Version)
	var err error

	// 1 - Load Configuration file
	cfg = new(config.Configuration)
	err = gonfig.Load(cfg)
	if err != nil {
		log.Fatalf("Unable to read config with error %v", err)
		os.Exit(1)
	}
	// 2 - Setup a new Slack Bot
	bot = slack.New(cfg.BotToken)
	authResp, err := bot.AuthTest()
	if err != nil {
		log.Fatalf("Unable to authenticate with Slack - %v", err)
		os.Exit(1)
	}
	log.Printf("Connected to Slack Team %v as user %v", authResp.Team, authResp.User)

	// 3 - Watch for new kills and log errors
	errc := make(chan error, 5)
	killc := make(chan zkill.Kill, 10)
	zClient := zkill.NewRedisQ()
	zClient.UserAgent = cfg.UserAgent
	zClient.FetchKillmails(killc, errc)
	handleKills(killc)
	handleErrors(errc)
	select {}
}

// handleKills sends the kill to be filtered/processed before posting to slack.
func handleKills(killChan chan zkill.Kill) {
	go func() {
		for {
			select {
			case kill := <-killChan:
				PostKill(&kill)
			}
		}
	}()
}

// handleErrors logs errors returned in Zkillboard queries
func handleErrors(errChan chan error) {
	go func() {
		for {
			select {
			case err := <-errChan:
				log.Printf("ERROR - %v", err.Error())
			}
		}
	}()
}
