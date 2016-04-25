package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/nlopes/slack"
	"github.com/vivace-io/evelib/zkill"
	"github.com/vivace-io/gonfig"
	"github.com/vivace-io/zk2s/util"
)

/* zk2s.go
 * Main entrypoint and controller for zk2s
 */

const (
	AppAuthor  = "Nathan \"Vivace Naaris\" Morley"
	AppVersion = "0.2"
)

var config *util.Configuration
var bot *slack.Client

func main() {
	app := cli.NewApp()
	app.Author = AppAuthor
	app.Version = AppVersion
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "start",
			Usage:  "start zk2s application",
			Action: Run,
		},
		cli.Command{
			Name:   "configure",
			Usage:  "configure zk2s application to be run",
			Action: util.RunConfigure,
		},
	}
	app.Run(os.Args)
}

// Run zk2s
func Run(c *cli.Context) {
	log.Printf("%v version %v", c.App.Name, c.App.Version)
	var err error

	// 1 - Load Configuration file
	config = new(util.Configuration)
	config.FileName = util.ConfigFileName
	err = gonfig.Load(config)
	if err != nil {
		log.Fatalf("Unable to read %v with error %v", util.ConfigFileName, err)
		os.Exit(1)
	}
	// 2 - Setup a new Slack Bot
	bot = slack.New(config.BotToken)
	authResp, err := bot.AuthTest()
	if err != nil {
		log.Fatalf("Unable to authenticate with Slack - %v", err)
		os.Exit(1)
	}
	log.Printf("Connected to Slack Team %v as user %v", authResp.Team, authResp.User)

	// 3 - Watch for new kills and log errors
	errc := make(chan error, 5)
	killc := make(chan zkill.ZKill, 10)
	zClient := zkill.New()
	zClient.UserAgent = config.UserAgent
	zClient.FetchKillmails(killc, errc)
	handleKills(killc)
	handleErrors(errc)
	select {}
}

// handleKills sends the kill to be filtered/processed before posting to slack.
func handleKills(killChan chan zkill.ZKill) {
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
