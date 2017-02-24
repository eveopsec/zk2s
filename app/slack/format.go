package slack

import (
	"bytes"
	"log"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/eveopsec/zk2s/app/config"
	"github.com/eveopsec/zk2s/app/filter"
	"github.com/eveopsec/zk2s/app/tmpl"
	slacklib "github.com/nlopes/slack"
	"github.com/vivace-io/evelib/crest"
	"github.com/vivace-io/evelib/redisq"
)

// TemplateData is passed to templates for defining how a slacklib post appears.
type TemplateData struct {
	// The raw CREST killmail model.
	Killmail crest.Killmail
	// The total value of the kill, as a formatted string with appropriate abbreviations.
	TotalValue string
	// IsLoss is true when the kill is a loss to the filter, false otherwise.
	IsLoss bool
	// IsSolo is true when the killmail was a solo kill/loss.
	IsSolo bool
}

// format loads the formatting template and applies formatting
// rules from the Configuration object.
func format(payload redisq.Payload, channel config.Channel) (messageParams slacklib.PostMessageParameters) {
	title := new(bytes.Buffer)
	body := new(bytes.Buffer)
	var err error

	// define post data for templates
	d := new(TemplateData)
	d.Killmail = payload.Killmail
	d.TotalValue = humanize.Comma(int64(payload.Zkb.TotalValue))
	d.IsLoss = filter.IsLoss(payload, channel)
	// Check if the kill was solo (NOTE: NPCs count)
	if len(payload.Killmail.Attackers) == 1 {
		d.IsSolo = true
	} else {
		d.IsLoss = false
	}

	// Execute templates
	err = tmpl.T.ExecuteTemplate(title, "kill-title", d)
	if err != nil {
		log.Println(err)
	}
	err = tmpl.T.ExecuteTemplate(body, "kill-body", d)
	if err != nil {
		log.Println(err)
	}

	attch := slacklib.Attachment{}
	attch.MarkdownIn = []string{"pretext", "text"}
	attch.Title = title.String()
	attch.TitleLink = "https://zkillboard.com/kill/" + strconv.Itoa(payload.KillID) + "/"
	attch.ThumbURL = "http://image.eveonline.com/render/" + strconv.Itoa(payload.Killmail.Victim.ShipType.ID) + "_64.png"
	attch.Text = body.String()
	//Color Coding
	if filter.IsLoss(payload, channel) {
		attch.Color = "danger"
	} else {
		attch.Color = "good"
	}
	messageParams.Attachments = []slacklib.Attachment{attch}
	messageParams.AsUser = true
	return
}
