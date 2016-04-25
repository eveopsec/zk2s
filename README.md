# zk2s

Application to post kills/losses from zKillboard to Slack in near real time by using zKillboard's RedisQ endpoint.

## Version 0.2
Get it while it's hot! See releases for release notes and binaries.

## Note

This is still a very rough implementation, so the post templating features are not as full as I might like, and some of the filters may or may not work. Please report any issues you may have.

Feedback and contributions are always welcome. Please create a new issue or pull request on this repository for either, or contact "Vivace Naaris" in game to talk to me!

Read the Installing/Configuration section below for help in setting up the application.

Todo:
 - [ ] Verify filters work in various configurations.
 - [ ] Develop some method of testing without having to explode myself.
 - [x] Add a character filter.
 - [x] Possibly remove intermediate `postData` structure and directly use killmail.
 - [x] Improve template features and documentation.
 - [ ] Make everything look better and more organized.
 - [ ] Clean up/organize `util/util.go`
 - [ ] Allow already configured channels to inspected/edited.
 - [ ] Command/Option for testing configuration.

**NOTE:** *Using this application on the same server/IP as another that also uses RedisQ can and will result in kills missing and not being posted.*

## Installing

To install, you can either install from source, or download the binary from releases.

**From Source:**
You will need to:
 1. Install or have installed the latest version of Go installed, with the environment properly configured. (see [this document](https://golang.org/doc/install) for more information on that process).
 2. Run `go get -u github.com/vivace-io/zk2s` to retrieve the source and its dependencies.
 3. Run `zk2s configure` to configure your setup
 4. Run `zk2s start` to run the application.

## Configuration
There are two files to configure the application.
- `zk2s.config.json` is the configuration file, which can be created/edited by running `zk2s configure`.
- `response.tmpl` is the template used to determine the format of the post to Slack.

**Configuring:**

At the top level of configuration, a `userAgent` must be set to your name/e-mail for getting data from zKillboard, and will refuse to exectute if this value is not set. This information is used only by zKillboard. `botToken` is the token used for posting to Slack. This may be either a bot token (recommended) or a user OAuth token from Slack([click here to set one up](https://api.slack.com/docs/oauth-test-tokens)). Whatever token you choose, that is the user posts will be made under.

Also, this new iteration allows for multiple channels on one team to receive killmail posts, each with their own filters! `channels` is an array of `Channel` objects (defined in `util/util.go`). These can be set up by hand or with `zk2s configure` utility, which also explains which each field means and how it is used.

All ships, characters, corporations, and alliances must match their exact name, and capitalization matters.

Entity IDs and Type IDs are not currently supported (but will be back soon!&trade)


**Post Template:**

This is on my TODO.

However, it can be noted that the `data` struct defined in `slack.go` is passed to *all* templates when they are executed. This means that you can access anything in that pipeline.

For example, you might try inserting these directly in to your template:
 - `{{.Killmail.Victim.Character.Name}}` inserts the victim's character name in to the template
 - `{{.Killmail.Victim.Corporation.Name}}` inserts the victim's corporation name in to the template
 - `{{.TotalValue}}` inserts the total value (formatted) in to the template
 - `{{if .IsLoss}} message A {{else}} message B {{end}}` will insert message A in to the template if it is a loss, otherwise it inserts message B.
 - `{{if .IsSolo}} message A {{else}} message B {{end}}` will insert message A in to the template if it is a solo kill/loss, otherwise it inserts message B.

`.Killmail` is a CREST Killmail model, so it can be assumed that any data available in a CREST killmail is accessible from this object.
