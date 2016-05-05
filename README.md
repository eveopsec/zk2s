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

## Configuring and Customizing

[Please see the wiki for details this.](https://github.com/vivace-io/zk2s/wiki)
