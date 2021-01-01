# GoGuardian ™️ the only way to truly protect your server against people like me

# GoGuardian pros

- completely open source

- faster than higher-level languages like js,python etc

- covers all endpoints (except emojis)


# GoGuardian cons

- that it uses JSON for a database (not really much of an issue unless you have like 10000 guilds or something)

# General help

- don't forget to make the bot role the HIGHEST it can be (even above your admins) for the best use-case

- remember to read each commands help info to see what permissions it needs.

- To toggle on/off each command just type it again.

Ask the server owner to whitelist you (gwhitelist @user) *if you are the server owner you do not have to whitelist yourself but you can*

# Command help:

Prefix: g | use ghelp

# Want to build it yourself?

Dependencies:

https://golang.org/dl/ (obviously)

```go get github.com/bwmarrin/discordgo``` (run in command line)

# How to build the program

```go build -ldflags "-s -w"``` (cd to the current path where main.go is and run this in the command line)

