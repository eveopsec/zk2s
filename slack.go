package main

import (
	"bytes"
	"html/template"
	"log"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/vivace-io/evelib/crest"

	"github.com/nlopes/slack"
	"github.com/vivace-io/evelib/zkill"
	"github.com/vivace-io/zk2s/util"
)

/* slack.go
 * Defines functions for posting kills to Slack
 */

var t = template.Must(template.ParseGlob("response.tmpl"))

// data is passed to template objects for defining how a slack post appears.
type data struct {
	Killmail   crest.Killmail
	TotalValue string
	IsLoss     bool
	IsSolo     bool
}

// PostKill applys the filter(s) to the kill, and posts the kill to slack
// only if the kill is within the configured filters.
func PostKill(kill *zkill.ZKill) {
	// For each filter defined in configuration,
	for c := range config.Channels {
		if util.WithinFilter(kill, config.Channels[c]) {
			params := format(kill, config.Channels[c])

			post(config.Channels[c].Name, params)
		}
	}
}

// format loads the formatting template and applies formatting
// rules from the Configuration object.
func format(kill *zkill.ZKill, channel *util.Channel) (messageParams slack.PostMessageParameters) {
	title := new(bytes.Buffer)
	body := new(bytes.Buffer)
	var err error

	// define post data for templates
	d := new(data)
	d.Killmail = kill.Killmail
	d.TotalValue = humanize.Comma(int64(kill.Zkb.TotalValue))
	d.IsLoss = util.IsLoss(kill, channel)
	if len(kill.Killmail.Attackers) == 1 {
		d.IsSolo = true
	} else {
		d.IsLoss = false
	}

	// Execute templates
	err = t.ExecuteTemplate(title, "killtitle", d)
	if err != nil {
		log.Println(err)
	}
	err = t.ExecuteTemplate(body, "killbody", d)
	if err != nil {
		log.Println(err)
	}

	attch := slack.Attachment{}
	attch.MarkdownIn = []string{"pretext", "text"}
	attch.Title = title.String()
	attch.TitleLink = "https://zkillboard.com/kill/" + strconv.Itoa(kill.KillID) + "/"
	attch.ThumbURL = "http://image.eveonline.com/render/" + strconv.Itoa(kill.Killmail.Victim.ShipType.ID) + "_64.png"
	attch.Text = body.String()
	if util.IsLoss(kill, channel) {
		attch.Color = "danger"
	} else {
		attch.Color = "good"
	}
	messageParams.Attachments = []slack.Attachment{attch}
	return
}

// post finally sends the kill to slack
func post(channel string, messageParams slack.PostMessageParameters) {
	messageParams.AsUser = true
	bot.PostMessage(channel, "", messageParams)
	// Throttle posting rates to Slack.
	time.Sleep(1 * time.Second)
}
