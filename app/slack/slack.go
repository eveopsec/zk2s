// Package slack handles the posting of kills to slack
package slack

import (
	"log"
	"time"

	"github.com/eveopsec/zk2s/app/config"
	"github.com/eveopsec/zk2s/app/filter"
	slacklib "github.com/nlopes/slack"
	"github.com/vivace-io/evelib/redisq"
)

var (
	clear chan bool
)

// Client operates as a client to all slack teams and posts killmails to the
// configured team and channel/group.
type Client struct {
	cfg *config.Configuration
}

func NewClient(cfg *config.Configuration) (client *Client, err error) {
	client = new(Client)
	for _, t := range cfg.Teams {
		t.Bot = slacklib.New(t.BotToken)
		_, err = t.Bot.AuthTest()
		if err != nil {
			log.Printf("[FATAL] Invalid token: %v", t.BotToken)
			client = nil
			return
		}
	}
	client.cfg = cfg
	client.manage()
	return
}

// Recieve kills from RedisQ on here.
func (client *Client) Recieve(payload redisq.Payload) {
	// TODO - Handle bulk!
	for _, t := range client.cfg.Teams {
		for _, c := range t.Channels {
			if filter.Within(payload, c) {
				params := format(payload, c)
				if !t.FailedAuth {
					log.Printf("[NOTIFY] Posting kill %v in channel %v", payload.KillID, c.Name)
					client.post(t, c, params)
				}
			}
		}
	}
}

func (client *Client) post(team *config.Team, channel config.Channel, messageParams slacklib.PostMessageParameters) {
	if team.FailedAuth {
		return
	}
	switch {
	case <-clear:
		var err error
		var channels []slacklib.Channel
		channels, err = team.Bot.GetChannels(true)
		if err != nil {
			team.FailedAuth = true
			log.Printf("[ERRPR] Failed to retrieve list of channels for team with error %v", err)
		}
		for _, c := range channels {
			if c.ID == channel.Name || c.Name == channel.Name {
				_, _, err = team.Bot.PostMessage(channel.Name, "", messageParams)
				if err != nil {
					team.FailedAuth = true
					log.Printf("[ERROR] Failed team for failure to post to Slack with error: %v (token %v)", err, team.BotToken)
				}
			}
		}
		var groups []slacklib.Group
		groups, err = team.Bot.GetGroups(true)
		if err != nil {
			team.FailedAuth = true
			log.Printf("[ERROR] Failed to retrieve list of groups for team with error %v", err)
		}
		for _, g := range groups {
			if g.ID == channel.Name || g.Name == channel.Name {
				_, _, err := team.Bot.PostMessage(channel.Name, "", messageParams)
				if err != nil {
					team.FailedAuth = true
					log.Printf("[ERROR] Failed team for failure to post to Slack with error: %v (token %v)", err, team.BotToken)
				}
			}
		}
	}
}

func (client *Client) manage() {
	clear = make(chan bool, 1)
	go func() {
		for {
			clear <- true
			time.Sleep(1 * time.Second)
		}
	}()
}
