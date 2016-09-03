package config

import (
	"github.com/nlopes/slack"
	"github.com/urfave/cli"
	"github.com/vivace-io/gonfig"
)

var (
	CONFIG *Application
)

func Init(c *cli.Context) error {
	CONFIG = new(Application)
	return gonfig.Load(CONFIG)
}

// Application holds the configuration for zk2s
type Application struct {
	UserAgent string  `json:"userAgent"`
	Teams     []*Team `json:"teams"`
}

// File returns the file name/path for gonfig interface
func (this *Application) File() string {
	return "cfg.zk2s.json"
}

// Save the configuration file
func (this *Application) Save() error {
	return gonfig.Save(this)
}

// Team is the configuration object for a slack team.
type Team struct {
	BotToken   string        `json:"botToken"`
	Channels   []Channel     `json:"channels"`
	Bot        *slack.Client `json:"-"`
	FailedAuth bool          `json:"-"`
}

// Channel defines the configuration for a slack channel in a team, including its filters
type Channel struct {
	Name                string   `json:"channelName"`
	MinimumValue        int      `json:"minimumValue"`
	MaximumValue        int      `json:"maximumValue"`
	IncludeCharacters   []string `json:"includeCharacters"`
	IncludeCorporations []string `json:"includeCorporations"`
	IncludeAlliances    []string `json:"includeAlliances"`
	ExcludedShips       []string `json:"excludedShips"`
}
