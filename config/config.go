package config

import (
	"bufio"
	"html/template"
	"os"

	"github.com/vivace-io/gonfig"
)

var t = template.Must(template.ParseGlob("response.tmpl"))
var input = bufio.NewReader(os.Stdin)

// LoadConfig reads the configuration file and returns it,
// marshalled in to Config
func LoadConfig() (*Configuration, error) {
	cfg := &Configuration{}
	err := gonfig.Load(cfg)
	return cfg, err
}

// Configuration defines zk2s' configuration
type Configuration struct {
	UserAgent string    `json:"userAgent"`
	BotToken  string    `json:"botToken"`
	Channels  []Channel `json:"channels"`
}

// File returns the file name/path for gonfig interface
func (c *Configuration) File() string {
	return "cfg.zk2s.json"
}

// Save the configuration file
func (c *Configuration) Save() error {
	return gonfig.Save(c)
}

// Channel defines the configuration for a slack channel, including its filters
type Channel struct {
	Name                string   `json:"channelName"`
	MinimumValue        int      `json:"minimumValue"`
	MaximumValue        int      `json:"maximumValue"`
	IncludeCharacters   []string `json:"includeCharacters"`
	IncludeCorporations []string `json:"includeCorporations"`
	IncludeAlliances    []string `json:"includeAlliance"`
	ExcludedShips       []string `json:"excludedShips"`
}
