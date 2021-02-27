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
	Prefix            string
	Name              string
	Help              string
	AdvancedHelp      []string
	RequiresArgs      bool
	RequiresWhitelist bool
	Run               handler
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

func (cmds *Commands) Add(name, helpInfo string, advancedHelp []string, fnc handler, requiresWhitelist, requiresArgs bool) *command {
	cmd := command{}
	cmd.Prefix = prefix
	cmd.Name = name
	cmd.Help = helpInfo
	cmd.AdvancedHelp = advancedHelp
	cmd.Run = fnc
	cmd.RequiresWhitelist = requiresWhitelist
	cmd.RequiresArgs = requiresArgs
	cmds.Commands = append(cmds.Commands, &cmd)
	return &cmd
}

func (cmds *Commands) Match(m string) (*command, []string) {
	content := strings.Fields(m)
	if len(content) == 0 {
		return nil, nil
	}
	var c *command
	var commandKey int

	for commandKey, commandName := range content {
		for _, commandValue := range cmds.Commands {
			if commandValue.Prefix+commandValue.Name == commandName {
				return commandValue, content[commandKey:]
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

	ctx := &Context{
		Content: strings.TrimSpace(m.Content),
	}
	cmd, fields := cmds.Match(ctx.Content)
	if cmd != nil {
		ctx.Fields = fields[1:]
		switch {
		case cmd.RequiresWhitelist && !inArray:
			s.ChannelMessageSend(m.ChannelID, "You need to be whitelisted to use this command. To get whitelisted ask the server owner to whitelist you. If you do not know how, type ghelp whitelist")
			return
		case len(ctx.Fields) == 0 && cmd.RequiresArgs:
			s.ChannelMessageSend(m.ChannelID, "You need to have proper arguments.")
			return
		case len(m.GuildID) == 0:
			s.ChannelMessageSend(m.ChannelID, "You need to type this command in a guild.")
			return
		}
		cmd.Run(s, m.Message, ctx)
	}
}
