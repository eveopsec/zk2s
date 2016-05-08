package main

import (
	"bytes"
	"log"
	"strconv"
	"text/template"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/vivace-io/evelib/crest"

	"github.com/nlopes/slack"
	"github.com/vivace-io/evelib/zkill"
	"github.com/vivace-io/zk2s/util"
)

/* slack.go
 * Defines functions for formatting/posting kills to Slack
 */

var t = template.Must(template.ParseGlob("response.tmpl"))

// data is passed to template objects for defining how a slack post appears.
type data struct {
	Killmail      crest.Killmail
	TotalValue    string
	IsLoss        bool
	IsSolo        bool
	CorpsInvolved []string
	AlliInvolved  []string
}

// PostKill applys the filter(s) to the kill, and posts the kill to slack
// only if the kill is within the configured filters.
func PostKill(kill *zkill.Kill) {
	// For each filter defined in configuration,
	for c := range config.Channels {
		if util.WithinFilter(kill, config.Channels[c]) {
			params := format(kill, config.Channels[c])
			log.Printf("Posting kill %v to channel %v", kill.KillID, config.Channels[c].Name)
			post(config.Channels[c].Name, params)
		}
	}
}

// format loads the formatting template and applies formatting
// rules from the Configuration object.
func format(kill *zkill.Kill, channel util.Channel) (messageParams slack.PostMessageParameters) {
	title := new(bytes.Buffer)
	body := new(bytes.Buffer)
	var err error

	// define post data for templates
	d := new(data)
	d.Killmail = kill.Killmail
	d.TotalValue = humanize.Comma(int64(kill.Zkb.TotalValue))
	d.IsLoss = util.IsLoss(kill, channel)
	//Solo kill testing
	if len(kill.Killmail.Attackers) == 1 {
		d.IsSolo = true
	} else {
		d.IsLoss = false
	}
	
	// Compile list of pilots involved, with special text formatting
	for a := range kill.Killmail.Attackers {
		//Special Formatting
		if {{kill.Killmail.attacker.finalBlow}} == 1
			{{kill.Killmail.attacker.Character.Name}}, //Regular Attacker
		else
			*{{kill.Killmail.attacker.Character.Name}}*, //Final Blow
		{{end}}
		//End Special formatting
		okToAdd := true
		for c := range d.PilotInvolved[c] {
			// Do not add blank pilots
			if kill.Killmail.Attackers[a].Character.Name == " ,"||","||" " {
				okToAdd = false
				break
			}
			if kill.Killmail.Attackers[a].Character.Name == d.PilotInvolved[c] {
				okToAdd = false
				break
			}
		}
		if okToAdd {
			d.PilotInvolved= append(d.PilotInvolved, kill.Killmail.Attackers[a].Character.Name)
		}
	}
	
	// Compile list of corporations involved from attackers, ignoring duplicates
	for a := range kill.Killmail.Attackers {
		okToAdd := true
		for c := range d.CorpsInvolved {
			if kill.Killmail.Attackers[a].Corporation.Name == d.CorpsInvolved[c] {
				okToAdd = false
				break
			}
		}
		if okToAdd {
			d.CorpsInvolved = append(d.CorpsInvolved, kill.Killmail.Attackers[a].Corporation.Name)
		}
	}

	// Compile list of alliances involved from attackers, ignoring duplicates
	for a := range kill.Killmail.Attackers {
		okToAdd := true
		for c := range d.AlliInvolved {
			// Do not add blank alliances (corp is not in an alliance)
			if kill.Killmail.Attackers[a].Alliance.Name == " ,"||","||" " {
				okToAdd = false
				break
			}
			if kill.Killmail.Attackers[a].Alliance.Name == d.AlliInvolved[c] {
				okToAdd = false
				break
			}
		}
		if okToAdd {
			d.AlliInvolved = append(d.AlliInvolved, kill.Killmail.Attackers[a].Alliance.Name)
		}
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
	//Color Coding
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
