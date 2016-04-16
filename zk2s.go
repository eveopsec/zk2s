package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
	"github.com/spf13/viper"
	"github.com/vivace-io/evelib/zkill"
	"github.com/vivace-io/zk2s/util"
)

var config *viper.Viper
var bot *slack.Client

func main() {
	var err error
	// 1 - Load Configuration file
	config, err = util.LoadConfig()
	if err != nil {
		log.Fatalf("Unable to read zk2s.config.json with error %v", err)
		os.Exit(1)
	}
	// 2 - Setup a new Slack Bot
	bot = slack.New(config.GetString("botToken"))
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
	zClient.UserAgent = config.GetString("userAgent")
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
				util.PostKill(&kill, bot, config)
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
