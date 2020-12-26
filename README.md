# GoGuardian ™️ the only way to truly protect your server against people like me

# GoGuardian features

Anti-hijack setting (makes GoGuardian leave the server if it has doing some not so nice things!)

bans the "wizz" bot very fast

faster than higher-level languages like js,python etc

covers all endpoints (except emojis)

# GoGuardian cons

that it uses JSON for a database (not really much of an issue unless you have like 10000 guilds or something)

there is no way to currently change threshold/seconds outside of config.json file (will probably be in by like tomorrow)

# General help

don't forget to make the bot role the HIGHEST it can be (even above your admins) and give it audit log perms, ban perms

so basically just give it admin it's easier

you can either put your token in Configs.json (to save it) or enter it when it asks (wont save)

Ask the server owner to whitelist you (gwhitelist @user) *if you are the server owner you do not have to whitelist yourself but you can*


# Command help:

Prefix: g | layout: gcommand

Commands: 

config antiadminrole - blocks people/bots from making roles have Administator permissions

config antiban - blocks wizz/nuke/destroyer bots from banning everyone

config antibots - bans any bot that gets invited (and the person who invited it)

config antichannelspam - stops people/bots from spam creating channels (mainly bots)

config antichannelnuke - stops people/bots from deleting lots of channels (mainly bots)

config antihijack - if GoGuardian ever gets comprimised and starts "nuking" the server it will leave.

config antikick - blocks wizz/nuke/destroyer bots from kicking everyone

config antimemberadmin - blocks people/bots from giving a member a role that has Administator permissions

config antirolespam - stops people/bots from spamm creating roles (mainly bots)

config antirolenuke - stops people/bots from deleting lots of roles (mainly bots)

whitelist @user - whitelists a user so that they do not get affected by any protection settings.

unwhitelist @user - unwhitelists a user so that they are affected by any toggled protection settings (default for all users.)

# Want to build it yourself?

Dependencies:

https://golang.org/dl/ (obviously)

```go get github.com/bwmarrin/discordgo``` (run in command line)

# How to build the program

```go build -ldflags "-s -w"``` (cd to the current path where main.go is and run this in the command line)

