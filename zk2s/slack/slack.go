// Package slack handles the posting of kills to slack
package slack

import (
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
	app   config.Application
)

func Init(c *cli.Context) error {
	var err error
	app, err = config.CONFIG.Get()
	if err != nil {
		return err
	}
	manage()
	return nil
}

// Recieve kills from RedisQ on here.
func Recieve(kill zkill.Kill) {
	for _, t := range app.Teams {
		for _, c := range t.Channels {
			if filter.Within(kill, *c) {
				params := format(kill, *c)
				if !t.FailedAuth {
					log.Printf("Posting kill %v in channel %v", kill.KillID, c.Name)
					post(t, *c, params)
				}
			}
		}
	}
}

func post(team config.Team, channel config.Channel, messageParams slacklib.PostMessageParameters) {
	if team.FailedAuth {
		return
	}
	switch {
	case <-clear:
		_, _, err := team.Bot.PostMessage(channel.Name, "", messageParams)
		if err != nil {
			team.FailedAuth = true
			log.Println("Failed team for failure to post to Slack (token %v)", team.BotToken)
		}
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
