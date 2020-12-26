package commands

import (
	"strings"

	"github.com/Not-Cyrus/GoGuardian/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	prefix string = "g"
)

type Commands struct {
	Commands []*command
}

type command struct {
	Name string
	Run  handler
}

type Context struct {
	Fields  []string
	Content string
}

type handler func(*discordgo.Session, *discordgo.Message, *Context)

func New() *Commands {
	c := &Commands{}
	return c
}

func (cmds *Commands) Add(name string, fnc handler) *command {
	cmd := command{}
	cmd.Name = prefix + name
	cmd.Run = fnc
	cmds.Commands = append(cmds.Commands, &cmd)
	return &cmd
}

func (cmds *Commands) Match(m string) (*command, []string) {
	content := strings.Fields(m)
	if len(content) == 0 {
		return nil, nil
	}
	var c *command
	var rank int
	var commandKey int

	for commandKey, commandName := range content {
		for _, commandValue := range cmds.Commands {
			if commandValue.Name == commandName {
				return commandValue, content[commandKey:]
			}
			if strings.HasPrefix(commandValue.Name, commandName) {
				if len(commandName) > rank {
					c = commandValue
					rank = len(commandName)
				}
			}
		}
	}
	return c, content[commandKey:]
}

func (cmds *Commands) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	originalData, _ := utils.FindConfig(m.GuildID)
	inArray, _ := utils.InArray(m.GuildID, "WhitelistedIDs", originalData, m.Author.ID)
	if !inArray && utils.GetGuildOwner(s, m.GuildID) != m.Author.ID {
		return
	}
	ctx := &Context{
		Content: strings.TrimSpace(m.Content),
	}
	cmd, fields := cmds.Match(ctx.Content)
	if cmd != nil {
		ctx.Fields = fields[1:]
		switch {
		case len(ctx.Fields) == 0:
			s.ChannelMessageSend(m.ChannelID, "You need to have proper arguments.")
			return
		case len(m.GuildID) == 0:
			s.ChannelMessageSend(m.ChannelID, "You need to type this command in a guild.")
			return
		}
		cmd.Run(s, m.Message, ctx)
	}
}
