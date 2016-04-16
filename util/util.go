// Package util contains definitions for filtering and posting kills to Slack from zKillboard.
package util

import (
	"bytes"
	"log"
	"strconv"
	"text/template"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"
	"github.com/vivace-io/evelib/zkill"
)

/* util/util.go
 * Defines functions for loading the configuration file, as well as formating
 * and posting kills to Slack.
 */

var t = template.Must(template.ParseGlob("response.tmpl"))

// postData is passed to the executed template for formatting.
type postData struct {
	SoloKill       bool
	AttackerName   string
	VictimName     string
	VictimShipName string
	DamageTaken    int
	TotalValue     string
	NumberInvolved int
}

// LoadConfig reads the configuration file and returns it,
// marshalled in to Config
func LoadConfig() (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile("zk2s.config.json")
	err := v.ReadInConfig()
	return v, err
}

// PostKill applys the filter(s) to the kill, and posts the kill to slack
// only if the kill is within the configured filters.
func PostKill(kill *zkill.ZKill, bot *slack.Client, config *viper.Viper) {
	if isWithinFilters(kill, config) {
		format(kill, bot, config)
	}
}

// format loads the formatting template and applies formatting
// rules from the Configuration object.
func format(kill *zkill.ZKill, bot *slack.Client, config *viper.Viper) {
	title := new(bytes.Buffer)
	body := new(bytes.Buffer)
	var err error

	data := postData{}
	getData(kill, &data)
	err = t.ExecuteTemplate(title, "killtitle", data)
	if err != nil {
		log.Println(err)
	}
	err = t.ExecuteTemplate(body, "killbody", data)
	if err != nil {
		log.Println(err)
	}

	attch := slack.Attachment{}
	attch.MarkdownIn = []string{"pretext", "text"}
	attch.Title = title.String()
	attch.TitleLink = "https://zkillboard.com/kills/" + strconv.Itoa(kill.KillID) + "/"
	attch.ThumbURL = "http://image.eveonline.com/render/" + strconv.Itoa(kill.Killmail.Victim.ShipType.ID) + "_64.png"
	attch.Text = body.String()
	if withinCorpFilter(kill.Killmail.Victim.Corporation.ID, config) {
		attch.Color = "danger"
	} else if withinAllianceFilter(kill.Killmail.Victim.Alliance.ID, config) {
		attch.Color = "danger"
	} else {
		attch.Color = "good"
	}
	messageParams := slack.PostMessageParameters{}
	messageParams.Attachments = []slack.Attachment{attch}
	post(bot, messageParams, config)
}

// getData takes a kill and builds a postData object for use
func getData(kill *zkill.ZKill, data *postData) {
	data.VictimName = kill.Killmail.Victim.Character.Name
	data.VictimShipName = kill.Killmail.Victim.ShipType.Name
	data.TotalValue = humanize.Comma(int64(kill.Zkb.TotalValue))
	if len(kill.Killmail.Attackers) == 1 {
		data.SoloKill = true
		data.AttackerName = kill.Killmail.Attackers[0].Character.Name
	}
	for a := range kill.Killmail.Attackers {
		data.DamageTaken += kill.Killmail.Attackers[a].DamageDone
		data.NumberInvolved++
	}
}

// post finally sends the kill to slack
func post(bot *slack.Client, messageParams slack.PostMessageParameters, config *viper.Viper) {
	messageParams.AsUser = true
	bot.PostMessage(config.GetString("channelName"), "", messageParams)
	// Throttle posting rates to Slack.
	time.Sleep(1 * time.Second)
}
