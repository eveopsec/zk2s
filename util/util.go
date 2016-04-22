// Package util contains definitions for filtering kills and loading configuration.
package util

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/vivace-io/gonfig"
)

/* util/util.go
 * Defines functions application configuration
 */

const ConfigFileName = "cfg.zk2s.json"

var t = template.Must(template.ParseGlob("response.tmpl"))
var config *Configuration

// LoadConfig reads the configuration file and returns it,
// marshalled in to Config
func LoadConfig() (*Configuration, error) {
	c := new(Configuration)
	c.FileName = ConfigFileName
	err := gonfig.Load(c)
	return c, err
}

// Configuration defines zk2s' configuration
type Configuration struct {
	FileName  string
	UserAgent string     `json:"userAgent"`
	BotToken  string     `json:"botToken"`
	Channels  []*Channel `json:"channels"`
}

// File returns the file name/path for gonfig interface
func (c *Configuration) File() string {
	return c.FileName
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

func RunConfigure(c *cli.Context) {
	fmt.Printf("Configure %v version %v", c.App.Name, c.App.Version)
	fmt.Println("What would you like to do?")
	fmt.Println("1 - Edit/Create configuration")
	fmt.Println("2 - Verify configuration")
	fmt.Println("0 - Exit")
	option := getOptionInt(0, 2)
	switch option {
	case 0:
		return
	case 1:
		configure(c)
	case 2:
		verifyConfig(c)
	}
}

func configure(c *cli.Context) {
	var err error
	config = new(Configuration)
	fmt.Println("---------------------------------------")
	fmt.Println("CONFIGURATION")
	fmt.Println("---------------------------------------")
	config.FileName = ConfigFileName
	err = gonfig.Load(config)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Printf("Unable to read/write to %v due to permission errors.", config.File())
			fmt.Println("Check permissions and try again.")
			return
		} else if os.IsNotExist(err) {
			fmt.Println("File does not exist, creating a new file...")
			err = gonfig.Save(config)
			if err != nil {
				fmt.Printf("Unable to create configuration - %v", err)
				return
			}
		} else {
			fmt.Printf("Error - %v", err)
		}
	}
	fmt.Println("Enter a UserAgent Name/E-mail (i.e. your/admin name). CANNOT be empty")
	fmt.Scanln(&config.UserAgent)
	fmt.Println("Enter the auth token for Slack. This can be either a bot token(recommended) or user token.")
	fmt.Scanln(&config.BotToken)
	configureChannels(c)
}

func configureChannels(c *cli.Context) {
	var choice int
	fmt.Println("---------------------------------------")
	fmt.Println("CONFIGURE CHANNELS")
	fmt.Println("---------------------------------------")
	if len(config.Channels) > 0 {
		fmt.Println("You have channels already configured:")
		for c := range config.Channels {
			fmt.Printf("%v - %v\n", c, config.Channels[c])
		}
		fmt.Printf("%v - New Channel", len(config.Channels)+1)
		fmt.Println("0 - Continue")
		fmt.Println("Select a channel to edit it or another option: ")
		choice = getOptionInt(0, len(config.Channels)+1)
		switch choice {
		case 0:
			return
		case 1:
			newChannel(c)
		default:
			editChannel(c, config.Channels[choice])
		}
	}
	newChannel(c)
	fmt.Println("Saving configuration...")
	err := gonfig.Save(config)
	if err != nil {
		fmt.Printf("Unable to save configuration to file - %v", err)
		return
	}
	fmt.Println("Done. Configuration complete, zk2s is now configured to run.")
}

func newChannel(c *cli.Context) {
	fmt.Println("---------------------------------------")
	fmt.Println("CONFIGURE NEW CHANNEL")
	fmt.Println("---------------------------------------")
	channel := new(Channel)
	fmt.Println("Enter the name of the channel you wish to post to: ")
	fmt.Scanln(&channel.Name)
	fmt.Println("ISK Values -- enter the following as an integer")
	fmt.Println("Minimum ISK value of the kill/loss for it to be posted:")
	fmt.Scanln(&channel.MinimumValue)
	fmt.Println("Maximum ISK value of the kill/loss for it to be posted:")
	fmt.Println("Note: value of 0 means no maximum is set")
	fmt.Scanln(&channel.MaximumValue)

	// Ships
	fmt.Println("Exclude any ships? Y/N")
	if yesOrNo() {
		var ok = false
		for !ok {
			var ship string
			fmt.Println("Enter Ship name or TypeID (caps sensitive/must be exact)")
			fmt.Scanln(&ship)
			channel.ExcludedShips = append(channel.ExcludedShips, ship)
			fmt.Println("Add another? Y/N")
			if !yesOrNo() {
				ok = true
			}
		}
	}

	// Alliances
	fmt.Println("Specify Alliance(s) to watch? Y/N")
	if yesOrNo() {
		var ok = false
		for !ok {
			var alliance string
			fmt.Println("Enter Alliance name or ID (caps sensitive/must be exact)")
			fmt.Scanln(alliance)
			channel.IncludeAlliances = append(channel.IncludeAlliances, alliance)
			fmt.Println("Add another? Y/N")
			if !yesOrNo() {
				ok = true
			}
		}
	}

	// Corporations
	fmt.Println("Specify Corporation(s) to watch? Y/N")
	if yesOrNo() {
		var ok = false
		for !ok {
			var corporation string
			fmt.Println("Enter Corporation name or ID (caps sensitive/must be exact)")
			fmt.Scanln(&corporation)
			channel.IncludeCorporations = append(channel.IncludeCorporations, corporation)
			fmt.Println("Add another? Y/N")
			if !yesOrNo() {
				ok = true
			}
		}
	}

	// Characters
	fmt.Println("Specify Character(s) to watch? Y/N")
	if yesOrNo() {
		var ok = false
		for !ok {
			var character string
			fmt.Println("Enter Character name or ID (caps sensitive/must be exact)")
			fmt.Scanln(&character)
			channel.IncludeCharacters = append(channel.IncludeCharacters, character)
			fmt.Println("Add another? Y/N")
			if !yesOrNo() {
				ok = true
			}
		}
	}
	config.Channels = append(config.Channels, channel)
}

func editChannel(c *cli.Context, channel *Channel) {

}

func verifyConfig(c *cli.Context) {
	fmt.Println("---------------------------------------")
	fmt.Println("VERIFY CONFIGURATION")
	fmt.Println("---------------------------------------")
}

func getOptionInt(lower int, upper int) int {
	var option int
	var ok = false
	for !ok {
		fmt.Scanln(&option)
		if !((option >= lower) && (option <= upper)) {
			fmt.Printf("Invalid option - please choose a number between %v and %v", lower, upper)
		} else {
			ok = true
		}
	}
	return option
}

// returns true for yes, false for no
func yesOrNo() bool {
	var option string
	fmt.Scanln(&option)
	strings.ToLower(option)
	if len(option) == 0 {
		fmt.Println("Please enter yes(y) or no(n)")
		return yesOrNo()
	} else if option[0] == []byte("y")[0] {
		return true
	} else if option[0] == []byte("n")[0] {
		return false
	}
	fmt.Println("Please enter yes(y) or no(n)")
	return yesOrNo()
}
