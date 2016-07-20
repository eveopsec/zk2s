package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
	"github.com/vivace-io/gonfig"
)

var config *Configuration

// RunAssistant runs configuration CLI process
func RunAssistant(c *cli.Context) error {
	// TODO - implment config validation and option
	fmt.Println("***************************************")
	fmt.Printf("CONFIGURE %v VERSION %v\n", c.App.Name, c.App.Version)
	fmt.Println("***************************************")
	fmt.Println("1 - Edit/Create configuration")
	fmt.Println("0 - Exit")
	fmt.Println("Select an option:")
	fmt.Println("---------------------------------------")
	option := getOptionInt(0, 1)
	if option == 1 {
		configure(c)
	}
	return nil
}

func configure(c *cli.Context) {
	var err error
	config = &Configuration{}
	err = gonfig.Load(config)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Printf("Unable to read/write to %v due to permission errors.\n", config.File())
			fmt.Println("Check permissions and try again.")
			return
		} else if os.IsNotExist(err) {
			fmt.Println("File does not exist, creating a new file...")
			err = gonfig.Save(config)
			if err != nil {
				fmt.Printf("Unable to create configuration - %v\na", err)
				return
			}
			fmt.Printf("New configuration file created!\n\n")
			configureInfo(c)
			return
		} else {
			fmt.Printf("Error - %v\n", err)
		}
	}
	fmt.Println("A configuration file already exists. Editing may overwrite these settings.")
	fmt.Println("Continue anyways? Yes/No")
	if !yesOrNo() {
		return
	}
	configureInfo(c)
}

func configureInfo(c *cli.Context) {
	fmt.Println("***************************************")
	fmt.Println("CONFIGURATION - INFO")
	fmt.Println("***************************************")
	if config != nil {
		if config.UserAgent != "" || config.BotToken != "" {
			fmt.Printf("UserAgent: %v\n", config.UserAgent)
			fmt.Printf("BotToken: %v\n", config.BotToken)
			fmt.Println("Overwrite these values? Y/n")
			if !yesOrNo() {
				configureChannels(c)
			}
		}
	}
	fmt.Println("Enter a UserAgent Name/E-mail (i.e. your/admin name). CANNOT be empty")
	config.UserAgent = getInputString()
	fmt.Println("Enter the auth token for Slack. This can be either a bot token(recommended) or user token.")
	config.BotToken = getInputString()
	configureChannels(c)
}

func configureChannels(c *cli.Context) {
	fmt.Println("***************************************")
	fmt.Println("CHANNEL CONFIGURATIONS")
	fmt.Println("***************************************")
	if len(config.Channels) == 0 {
		fmt.Printf("You have no channels configured. Please create a new channel.\n\n")
		newChannel(c)
	} else {
		fmt.Println("Channels already configured: ")
		for c := range config.Channels {
			fmt.Printf("%v - %v\n", c+1, config.Channels[c].Name)
		}
		fmt.Println("---------------------------------------")
		fmt.Println("Add new channel to configuration? Yes/No")
		if yesOrNo() {
			newChannel(c)
		}
	}
	fmt.Println("Saving configuration...")
	err := gonfig.Save(config)
	if err != nil {
		fmt.Printf("Unable to save configuration to file - %v\n", err)
		return
	}
	fmt.Println("Done. Configuration complete, zk2s is now configured to run.")
	os.Exit(0)
}

func newChannel(c *cli.Context) {
	fmt.Println("***************************************")
	fmt.Println("NEW CHANNEL")
	fmt.Println("***************************************")
	channel := Channel{}
	fmt.Println("Enter the name of the channel you wish to post to: ")
	channel.Name = getInputString()
	fmt.Println("---------------------------------------")
	fmt.Println("ISK Values -- enter the following as an integer")
	fmt.Println("Minimum ISK value of the kill/loss for it to be posted:")
	fmt.Println("---------------------------------------")
	fmt.Scanln(&channel.MinimumValue)
	fmt.Println("Maximum ISK value of the kill/loss for it to be posted:")
	fmt.Println("Note: value of 0 means no maximum is set")
	fmt.Println("---------------------------------------")
	fmt.Scanln(&channel.MaximumValue)
	fmt.Println("---------------------------------------")

	// Ships
	fmt.Println("Exclude any ships? Yes/No")
	if yesOrNo() {
		var ok = false
		for !ok {
			fmt.Println("Enter Ship name(caps sensitive/must be exact)")
			ship := getInputString()
			channel.ExcludedShips = append(channel.ExcludedShips, ship)
			fmt.Println("Add another? Yes/No")
			if !yesOrNo() {
				ok = true
			}
		}
	}
	fmt.Println("---------------------------------------")

	// Alliances
	fmt.Println("Specify Alliance(s) to watch? Yes/No")
	if yesOrNo() {
		var ok = false
		for !ok {
			fmt.Println("Enter Alliance name (caps sensitive/must be exact)")
			alliance := getInputString()
			channel.IncludeAlliances = append(channel.IncludeAlliances, alliance)
			fmt.Println("Add another? Yes/No")
			if !yesOrNo() {
				ok = true
			}
		}
	}
	fmt.Println("---------------------------------------")

	// Corporations
	fmt.Println("Specify Corporation(s) to watch? Yes/No")
	if yesOrNo() {
		var ok = false
		for !ok {
			fmt.Println("Enter Corporation name (caps sensitive/must be exact)")
			corporation := getInputString()
			channel.IncludeCorporations = append(channel.IncludeCorporations, corporation)
			fmt.Println("Add another? Yes/No")
			if !yesOrNo() {
				ok = true
			}
		}
	}
	fmt.Println("---------------------------------------")

	// Characters
	fmt.Println("Specify Character(s) to watch? Yes/No")
	if yesOrNo() {
		var ok = false
		for !ok {
			fmt.Println("Enter Character name (caps sensitive/must be exact)")
			character := getInputString()
			channel.IncludeCharacters = append(channel.IncludeCharacters, character)
			fmt.Println("Add another? Yes/No")
			if !yesOrNo() {
				ok = true
			}
		}
	}
	fmt.Println("---------------------------------------")
	config.Channels = append(config.Channels, channel)
	fmt.Println("Configure another channel? Yes/No")
	if yesOrNo() {
		newChannel(c)
	}
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
	option := getInputString()
	option = strings.ToLower(option)
	if len(option) == 0 {
		fmt.Println("enter yes(y) or no(n)")
		return yesOrNo()
	} else if option[0] == []byte("y")[0] {
		return true
	} else if option[0] == []byte("n")[0] {
		return false
	}
	fmt.Println("enter yes(y) or no(n)")
	return yesOrNo()
}

// returns input as a string, prevents empty input
func getInputString() string {
	s, _, err := input.ReadLine()
	if err != nil {
		fmt.Printf("Error - %v\n", err)
		return getInputString()
	}
	if len(s) == 0 {
		fmt.Println("input cannot be empty")
		return getInputString()
	}
	return string(s)
}
