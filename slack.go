package main

import (
	"html/template"
	"time"

	"github.com/nlopes/slack"
	"github.com/vivace-io/evelib/zkill"
	"github.com/vivace-io/zk2s/util"
)

/* slack.go
 * Defines functions for posting kills to Slack
 */

var t = template.Must(template.ParseGlob("response.tmpl"))

// PostKill applys the filter(s) to the kill, and posts the kill to slack
// only if the kill is within the configured filters.
func PostKill(kill *zkill.ZKill) {
	// For each filter defined in configuration,
	for c := range config.Channels {
		if util.WithinFilter(kill, config.Channels[c]) {
			params := format(kill)
			post(config.Channels[c].Name, params)
		}
	}
}

// format loads the formatting template and applies formatting
// rules from the Configuration object.
func format(kill *zkill.ZKill) (messageParams slack.PostMessageParameters) {
	return
}

// post finally sends the kill to slack
func post(channel string, messageParams slack.PostMessageParameters) {
	messageParams.AsUser = true
	bot.PostMessage(channel, "", messageParams)
	// Throttle posting rates to Slack.
	time.Sleep(1 * time.Second)
}
