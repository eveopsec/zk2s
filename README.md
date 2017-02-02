# zk2s

Application to post kills/losses from zKillboard to Slack in near real time by using zKillboard's RedisQ endpoint.

This application is independent of the OpSec Project and can be run without requiring other services

Previously hosted at vivace-io/zk2s, this project has been moved to this organization to group it with other Eve Online tools and applications.

## Version 2.0

**Changes:**

 - [x] Template file path is now specified in configuration. If the file path is not specified, defaults to `response.tmpl` of the working directory.
 - [x] Improved API such that it can be used in other applications.
 - [x] More code cleanup... I'm cringing at my own code.
 - [x] Fixed posting to private groups (hopefully).

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

**From Docker:**
 1. Pull the image from DockerHub: `docker pull eveopsec/zk2s`
 2. Create your desired configuration and templates, and store them in a directory of your choosing. (If you do not have a template, feel free to use the provided default template!)
 3. Create your container from the image, like so: `docker create --name eve-zk2s -v /host/path/to/config/dir/zk2s:/zk2s zk2s:1.1`
 4. Star the docker image and you're good to go!

 **From Source:**

 *THIS IS NOT RECOMMENDED*
 Mostly because the `evelib` project API that ZK2S uses is in a constant of flux, half because I'm always changing my mind as to how to program it, the other half because my time is currently limited!

 If you know your way around the code, then by all means. You will get several errors when Go tries to compile it, generally for methods not being found or wrong return types (literals vs pointers etc.).

 I'll try to get back to consistency as soon as possible!

 You will need to:
  1. Install or have installed the latest version of Go installed, with the environment properly configured. (see [this document](https://golang.org/doc/install) for more information on that process).
  2. Run `go get -u github.com/vivace-io/zk2s` to retrieve the source and its dependencies.
  3. Run `zk2s configure assistant` to run the configuration setup
  4. Run `zk2s start` to run the application.

## Configuring and Customizing

[Please see the wiki for details.](https://github.com/vivace-io/zk2s/wiki)
