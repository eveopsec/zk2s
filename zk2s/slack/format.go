package slack

import (
	"bytes"
	"log"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/eveopsec/zk2s/zk2s/config"
	"github.com/eveopsec/zk2s/zk2s/filter"
	"github.com/eveopsec/zk2s/zk2s/tmpl"
	slacklib "github.com/nlopes/slack"
	"github.com/vivace-io/evelib/crest"
	"github.com/vivace-io/evelib/zkill"
)

// data is passed to templates for defining how a slacklib post appears.
type data struct {
	Killmail       crest.Killmail
	TotalValue     string
	IsLoss         bool
	IsSolo         bool
	InAlli         bool
	LosingCorp     string
	LosingAlli     string
	CorpsInvolved  []string
	AlliInvolved   []string
	PilotInvolved  []string
	FinalBlowPilot []string
	FinalBlowCorp  []string
	FinalBlowAlli  []string
	TotalCorp      []string
	TotalAlli      []string
}

// format loads the formatting template and applies formatting
// rules from the Configuration object.
func format(kill zkill.Kill, channel config.Channel) (messageParams slacklib.PostMessageParameters) {
	title := new(bytes.Buffer)
	body := new(bytes.Buffer)
	var err error

	// define post data for templates
	d := new(data)
	d.Killmail = *kill.Killmail
	d.TotalValue = humanize.Comma(int64(kill.Zkb.TotalValue))
	d.IsLoss = filter.IsLoss(kill, channel)
	//Solo kill testing
	if len(kill.Killmail.Attackers) == 1 {
		d.IsSolo = true
	} else {
		d.IsLoss = false
	}
	//Testing to see if the victim is in an alliance

	if kill.Killmail.Victim.Alliance.Name != "" {
		d.InAlli = true
		d.LosingAlli = kill.Killmail.Victim.Alliance.Name
	}
	d.LosingCorp = kill.Killmail.Victim.Corporation.Name

	// Compile list of pilots involved, if not final blow
	for a := range kill.Killmail.Attackers {
		okToAdd := true
		if kill.Killmail.Attackers[a].FinalBlow == true {
			okToAdd = false
		}
		if kill.Killmail.Attackers[a].Character.Name == "" {
			okToAdd = false
		}
		if okToAdd {
			d.PilotInvolved = append(d.PilotInvolved, kill.Killmail.Attackers[a].Character.Name)
		}
	}

	//Compile the list for the final blow pilot, mainly use for formatting commas on the post
	for a := range kill.Killmail.Attackers {
		if kill.Killmail.Attackers[a].FinalBlow == true {
			okToAdd := true
			if kill.Killmail.Attackers[a].Character.Name == "" {
				okToAdd = false
			}
			if okToAdd {
				d.FinalBlowPilot = append(d.FinalBlowPilot, kill.Killmail.Attackers[a].Character.Name)
				d.FinalBlowCorp = append(d.FinalBlowCorp, kill.Killmail.Attackers[a].Corporation.Name)
				d.FinalBlowAlli = append(d.FinalBlowAlli, kill.Killmail.Attackers[a].Alliance.Name)
				d.TotalCorp = append(d.TotalCorp, kill.Killmail.Attackers[a].Corporation.Name)
				d.TotalAlli = append(d.TotalAlli, kill.Killmail.Attackers[a].Alliance.Name)
			}

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
			if kill.Killmail.Attackers[a].Corporation.Name == d.FinalBlowCorp[c] {
				okToAdd = false
				break
			}
			if okToAdd {
				d.CorpsInvolved = append(d.CorpsInvolved, kill.Killmail.Attackers[a].Corporation.Name)
				d.TotalCorp = append(d.TotalCorp, kill.Killmail.Attackers[a].Corporation.Name)
			}
		}
	}
	// Compile list of alliances involved from attackers, ignoring duplicates
	for a := range kill.Killmail.Attackers {

		okToAdd := true

		for c := range d.AlliInvolved {

			// Do not add blank alliances (corp is not in an alliance)
			if kill.Killmail.Attackers[a].Alliance.Name == "" {
				okToAdd = false
				break
			}
			if kill.Killmail.Attackers[a].Alliance.Name == d.AlliInvolved[c] {
				okToAdd = false
				d.InAlli = true
				break
			}
			if kill.Killmail.Attackers[a].Alliance.Name == d.FinalBlowAlli[c] {
				okToAdd = false
				d.InAlli = true
				break
			}
			if okToAdd {
				d.AlliInvolved = append(d.AlliInvolved, kill.Killmail.Attackers[a].Alliance.Name)
				d.TotalAlli = append(d.TotalAlli, kill.Killmail.Attackers[a].Alliance.Name)
				d.InAlli = true
			}
		}
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
	attch.TitleLink = "https://zkillboard.com/kill/" + strconv.Itoa(kill.KillID) + "/"
	attch.ThumbURL = "http://image.eveonline.com/render/" + strconv.Itoa(kill.Killmail.Victim.ShipType.ID) + "_64.png"
	attch.Text = body.String()
	//Color Coding
	if filter.IsLoss(kill, channel) {
		attch.Color = "danger"
	} else {
		attch.Color = "good"
	}
	messageParams.Attachments = []slacklib.Attachment{attch}
	return
}
