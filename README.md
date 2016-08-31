# zk2s

Application to post kills/losses from zKillboard to Slack in near real time by using zKillboard's RedisQ endpoint.

This application is independent of the OpSec Project and can be run without requiring other services

Previously hosted at vivace-io/zk2s, this project has been moved to this organization to group it with other Eve Online tools and applications.

## Version 1.0

**Changes:**

**This release will break old configurations!**

First major release!

 - [x] **Multiple Slack Team support**
 - [ ] Configuration Changes
  - [ ] Options to post kills in bulk rather than individually.
 - [ ] Template Changes
  - [x] Templates `killtitle` and `killbody` renamed to `kill-title` and `kill-body`
  - [ ] New templates `bulk-title` and `bulk-body` introduced.
 - [ ] Lots and lots of code changes and packages created/deleted/consolidated.
  - [ ] Automated testing
  - [ ] Testing flag to run a test slack post
 - [ ] Setup Assistant improved

## Note

Feedback and contributions are always welcome. Please create a new issue or pull request on this repository for either, or contact "Vivace Naaris" in game to talk to me!

Read the Installing/Configuration section below for help in setting up the application.

Todo:
 - [ ] Develop some method of testing without having to explode myself.
 - [ ] Make everything look better and more organized.
 - [ ] Allow already configured channels to inspected/edited.
 - [ ] Command/Option for testing configuration.

**NOTE:** *Using this application on the same server/IP as another that also uses RedisQ can and will result in kills missing and not being posted.*

## Installing

To install, you can either install from source, or download the binary from releases.

**From Source:**
You will need to:
 1. Install or have installed the latest version of Go installed, with the environment properly configured. (see [this document](https://golang.org/doc/install) for more information on that process).
 2. Run `go get -u github.com/vivace-io/zk2s` to retrieve the source and its dependencies.
 3. Run `zk2s configure assistant` to run the configuration setup
 4. Run `zk2s start` to run the application.

## Configuring and Customizing

[Please see the wiki for details this.](https://github.com/vivace-io/zk2s/wiki)
