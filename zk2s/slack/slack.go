// Package slack handles the posting of kills to slack
package slack

import (
	"errors"
	"log"
	"time"

	"github.com/eveopsec/zk2s/zk2s/config"
	"github.com/eveopsec/zk2s/zk2s/filter"
	slacklib "github.com/nlopes/slack"
	"github.com/urfave/cli"
	"github.com/vivace-io/evelib/zkill"
)

var (
	clear chan bool
	app   *config.Application
)

func Init(c *cli.Context) error {
	if config.CONFIG == nil {
		return errors.New("app configuration was nil")
	}
	app = config.CONFIG

	for _, t := range app.Teams {
		t.Bot = slacklib.New(t.BotToken)
		_, err := t.Bot.AuthTest()
		if err != nil {
			log.Printf("[WARNING] - token %v is invalid and will not be posted to!", t.BotToken)
		}
	}

	manage()

	return nil
}

// Recieve kills from RedisQ on here.
func Recieve(kill zkill.Kill) {
	// TODO - Handle bulk!
	for _, t := range app.Teams {
		for _, c := range t.Channels {
			if filter.Within(kill, c) {
				params := format(kill, c)
				if !t.FailedAuth {
					log.Printf("Posting kill %v in channel %v", kill.KillID, c.Name)
					post(t, c, params)
				}
			}
		}
	}
}

func post(team *config.Team, channel config.Channel, messageParams slacklib.PostMessageParameters) {
	if team.FailedAuth {
		return
	}
	switch {
	case <-clear:
		channels, err := team.Bot.GetChannels(true)
		if err != nil {
			team.FailedAuth = true
			log.Printf("Failed to retrieve list of channels for team with error %v", err)
		}
		for _, c := range channels {
			if c.ID == channel.Name || c.Name == channel.Name {
				_, _, err := team.Bot.PostMessage(channel.Name, "", messageParams)
				if err != nil {
					team.FailedAuth = true
					log.Printf("Failed team for failure to post to Slack (token %v) with error: %v", team.BotToken, err)
				}
				return
			}
		}
		team.FailedAuth = true
		log.Printf("Unable to find channel %v associated with token %v", channel.Name, team.BotToken)
	}
}

func manage() {
	clear = make(chan bool, 1)
	go func() {
		for {
			clear <- true
			time.Sleep(1 * time.Second)
		}
	}()
}
