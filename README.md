# GoGuardian ™️ the only way to truly protect your server against people like me

# GoGuardian features

Anti-hijack setting (makes GoGuardian leave the server if it has been doing some not so nice things!)

bans the "wizz" bot after a certain threshold has been met in a certain amount of seconds (changed with gconfig seconds number,gconfig threshold number)

faster than higher-level languages like js,python etc

covers all endpoints (except emojis)

# GoGuardian cons

that it uses JSON for a database (not really much of an issue unless you have like 10000 guilds or something)

# General help

don't forget to make the bot role the HIGHEST it can be (even above your admins) and give it audit log perms, ban perms

so basically just give it admin it's easier

you can either put your token in Config.json (to save it) or enter it when it asks (wont save)

Ask the server owner to whitelist you (gwhitelist @user) *if you are the server owner you do not have to whitelist yourself but you can*

# Command help:

Prefix: g | use ghelp

# Want to build it yourself?

Dependencies:

https://golang.org/dl/ (obviously)

```go get github.com/bwmarrin/discordgo``` (run in command line)

# How to build the program

```go build -ldflags "-s -w"``` (cd to the current path where main.go is and run this in the command line)

