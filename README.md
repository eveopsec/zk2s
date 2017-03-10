# zk2s

Application to post kills/losses from zKillboard to Slack in near real time by using zKillboard's RedisQ endpoint.

This application is independent of the OpSec Project and can be run without requiring other services

Previously hosted at vivace-io/zk2s, this project has been moved to this organization to group it with other Eve Online tools and applications.

## Version 2.2

**Changes:**

Un-checked changes have not been completed, or are not fully tested.

 - [x] Updated default template to include links to attacker corporations.
 - [ ] Added `includeSystems` to monitor specific systems by ID.

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

**From Binaries:**
 1. [Download the binary from Github](https://github.com/eveopsec/zk2s/releases ) for your OS distribution and extract it in to a folder.
 2. Setup your configuration file as you see fit. [See the wiki for details.](https://github.com/eveopsec/zk2s/wiki)

**From Source:**

 *BEFORE YOU BEGIN*: This project uses `vivace-io/evelib` as one of its core dependencies, which is under active development as my time allows. As such, the API may change and suddenly this application does not compile. If it comes to this, please feel free
 to create a new issue and I can work on a quick fix, or present your own fix in a pull request.

 Regardless, below are the instructions to compile from source.

  1. Install or have installed the latest version of Go installed, with the environment properly configured. See [this document](https://golang.org/doc/install) for more information on that process.
  2. Run `go get -u github.com/vivace-io/zk2s` to retrieve the source and its dependencies.
  3. Run `zk2s configure assistant` to run the configuration setup
  4. Run `zk2s start` to run the application.

** From Docker:**

While there is a Dockerfile in this repository, I've found issues where the container is not running properly, or unexpectedly hangs/crashes. I have not had the time to work on a solution, so currently Docker is not actively supported. Should you be able to replicate this problem, and/or present a fix, please open a new issue/pull request and we can get it in to the code base ASAP!

## Configuring and Customizing

[Please see the wiki for details.](https://github.com/vivace-io/zk2s/wiki)
