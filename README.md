# zk2s

Application to post kills/losses from zKillboard to Slack in near real time by using zKillboard's RedisQ endpoint.

This is still a very rough implementation, so the post templating features are not as full as I might like, and some of the filters may or may not work. Please report any issues you may have.

Feedback and contributions are always welcome. Please create a new issue or pull request on this repository for either, or contact "Vivace Naaris" in game to talk to me!

Read the Installing/Configuration section below for help in setting up the application.

Todo:
 - [ ] Verify filters work in various configurations.
 - [ ] Develop some method of testing without having to explode myself.
 - [ ] Add a character filter.
 - [ ] Possibly remove intermediate `postData` structure and directly use killmail.
 - [ ] Improve template features and documentation.
 - [ ] Make everything look better and more organized.

**NOTE:** *Using this application on the same server/IP as another that also uses RedisQ can and will result in kills missing and not being posted.*

## Installing

To install, you can either install from source, or download the binary from releases.

**From Source:** You will need to:
 1. Install or have installed the latest version of Go installed, with the environment properly configured. (see [this document](https://golang.org/doc/install) for more information on that process).
 2. Run `go get -u github.com/vivace-io/zk2s` to retrieve the source and its dependencies.
 3. Edit `zk2s.config.json` from `$GOPATH/src/github.com/vivace-io/zk2s/` according to the configuring section below.
 4. Execute `go run zk2s.go` to run the application. You can also execute `go install github.com/vivace-io/zk2s` to install a binary in to `$GOPATH/bin`, but you will need to copy `response.tmpl` and `zk2s.config.json` to the same directory as the executable.
 5. I personally use Supervisor on my server to run this application as a service.

## Configuration
There are two files to configure the application.
- `zk2s.config.json` is the configuration file, and initially contains values to include kills from the Opsec. Corporation and Wrong Hole alliance, and excludes kills where the victim is in an Ibis.
- `response.tmpl` is the template used to determine the format of the post to Slack.

**Configuring:**
 1. Set `userAgent` to your name and contact information. The application will refuse to run unless you set this. This is sent with requests to zKillboard's servers, per the API's terms of use.
 2. On Slack, create a new Slackbot integration and generate a token. Insert the this token in the `botToken` value.
 3. Set `channelName` to the name of the channel you wish to post new kills to.
 4. Configure your filters.

**Filters:**
 - `iskMinimum` (type int) defines a minimum value threshold for the kill to be posted. Default value is 0.
 - `excludeLosses` (type bool) will not post losses if set to true.
 - `excludeShips` (type []string) string - will not post any kills where the victim was in a ship specified in here. You can put either a name or type ID here.
 - `includeCorps` (type []string) can contain either the corporation ID or name.
 - `includeAlliances`(type []string) can contain either the alliance ID or name.

`includeCorps` and `includeAlliances` can be either a name or ID. If your spelling is off, it will not post. I recommend that you use the ID of those entities. To get the ID, you can search for the corporation or alliance on ZKillboard and pull the ID from the URL.

If a filter is not specified or contains no value(s), the filter is not applied. So, if you have nothing specified `includeCorps` and `includeAlliances`, then no kills will be posted. The exception is `excludeShips`, where if no values are specified, all ships that fall within the other filters *will* be posted.


**Post Template:**

I still need to document and expand on this a bit more. I will probably end up just using the killmail data structure directly when executing the template, but I am short on time at the moment.

However, for those who want to dive in before then, all you need to know to get started is that `response.tmpl` is executed by Go's `text/template` package with variables passed in to it from the `postData` structure defined in `util/util.go`. You can feel free to change any of the text in here to customize your messages without fear of breaking anything. Anything within double brackets (`{{}}`) is where corresponding values in `postData` will appear.
