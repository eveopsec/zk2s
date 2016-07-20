package zk2s

import (
	"log"

	"github.com/eveopsec/zk2s/zk2s/config"
	"github.com/nlopes/slack"
	"github.com/urfave/cli"
	"github.com/vivace-io/evelib/zkill"
	"github.com/vivace-io/gonfig"
)

var cfg *config.Configuration
var bot *slack.Client

// Run zk2s
func Run(c *cli.Context) error {
	log.Printf("%v version %v", c.App.Name, c.App.Version)
	var err error

	// 1 - Load Configuration file
	cfg = new(config.Configuration)
	err = gonfig.Load(cfg)
	if err != nil {
		log.Fatalf("Unable to read config with error %v", err)
		return err
	}
	// 2 - Setup a new Slack Bot
	bot = slack.New(cfg.BotToken)
	authResp, err := bot.AuthTest()
	if err != nil {
		log.Fatalf("Unable to authenticate with Slack - %v", err)
		return err
	}
	log.Printf("Connected to Slack Team %v as user %v", authResp.Team, authResp.User)

	// 3 - Load templates from Flag
	err = TemplateFromPath(c.String("template"))
	if err != nil {
		log.Fatalf("Unable to load template file with error: %v", err)
		return err
	}

	// 4 - Watch for new kills and log errors
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
