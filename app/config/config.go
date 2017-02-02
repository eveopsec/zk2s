package config

import (
	"github.com/nlopes/slack"
	"github.com/vivace-io/gonfig"
)

func ReadConfig(filepath string) (cfg *Configuration, err error) {
	cfg = new(Configuration)
	cfg.filepath = filepath
	if err := gonfig.Load(cfg); err != nil {
		return nil, err
	}
	return
}

// Configuration holds the configuration for zk2s
type Configuration struct {
	filepath     string  `json:"-"`
	TemplateFile string  `json:"template_file"`
	UserAgent    string  `json:"userAgent"`
	Teams        []*Team `json:"teams"`
}

// File returns the file name/path for gonfig interface
func (cfg *Configuration) File() string {
	return cfg.filepath
}

// Save the configuration file
func (cfg *Configuration) Save() error {
	return gonfig.Save(cfg)
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
