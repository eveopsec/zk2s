package cmd

import "github.com/urfave/cli"

/* info.go
 * General information on the application/library
 */

var (
	// VERSION of zk2s
	VERSION = "2.0"

	// USAGE text for zk2s application/library
	USAGE = "A Slack bot for posting kills from zKillboard to slack in near-real time."

	// CONTRIBUTORS to zk2s
	CONTRIBUTORS = []cli.Author{
		cli.Author{
			Name: "Nathan \"Vivace Naaris\" Morley",
		},
		cli.Author{
			Name: "\"Zuke\"",
		},
	}
)
