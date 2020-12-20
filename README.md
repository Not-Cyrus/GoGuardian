# GoGuardian ™️ the only way to truly protect your server against people like me

# GoGuardian features

Anti-hijack setting (makes GoGuardian leave the server if it has doing some not so nice things!)

bans the "wizz" bot very fast

faster than higher-level languages like js,python etc

covers all endpoints (except emojis)

# GoGuardian cons

doesn't really have multi-guild support or sharding so no real scalability

# General help

don't forget to make the bot role the HIGHEST it can be (even above your admins) and give it audit log perms, ban perms
so basically just give it admin it's easier


# Help for config.json

Threshold: how much times a person/bot can do a certain action in x seconds (see below) if they reach over x limit they get banned themselves

Seconds: if someone has banned equal to or over the Threshold (see above) in x seconds

token: (your bot token duh?)

WhitelistedIDs: Whitelisted IDs that can bypass the banning

BanProtection: Keeps you safe from bots/people that mass ban

KickProtection: Keeps you safe from bots/people that mass kick

RoleSpamProtection: Keeps you safe from bots/people that mass create roles (also bans annoying admins that do it so be warned)

RoleNukeProtection: Keeps you safe from bots/people that mass delete roles

ChannelSpamProtection: Keeps you safe from bots/people that mass create channels (also bans annoying admins that do it so be warned)

ChannelNukeProtection: Keeps you safe from bots/people that mass delete channels

MemberRoleUpdateProtection: Keeps you safe from bots/people giving administrator roles without being whitelisted

AntiBotProtection: Removes all bots that get invited and bans the person who invited it

HijackProtection: Makes the bot leave the guild if it tried nuking the server itself.

# Want to build it yourself?

Dependencies:

https://golang.org/dl/ (obviously)

```go get github.com/bwmarrin/discordgo``` (run in command line)

# How to build the program

```go build -ldflags "-s -w"``` (cd to the current path where main.go is and run this in the command line)

