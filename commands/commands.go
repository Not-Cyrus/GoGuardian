package commands

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Not-Cyrus/GoGuardian/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Not-Cyrus/GoGuardian/database"
	"github.com/bwmarrin/discordgo"
)

func (cmds *Commands) Add(name string, function handler, config *Config) *Command {
	cmd := Command{}

	cmd.Name = name
	cmd.Run = function
	cmd.Config = config

	cmds.Commands = append(cmds.Commands, &cmd)

	return &cmd
}

func (cmds *Commands) addCooldown(userID string, command string, cooldown int) {
	cmds.Cooldown.Mutex.Lock()
	defer cmds.Cooldown.Mutex.Unlock()
	cmds.Cooldown.Cooldowns[userID] = append(cmds.Cooldown.Cooldowns[userID], command)

	time.AfterFunc(time.Duration(cooldown)*time.Second, func() {
		cmds.Cooldown.Mutex.Lock()
		defer cmds.Cooldown.Mutex.Unlock()

		cmds.Cooldown.Cooldowns[userID] = utils.RemoveFromSlice(cmds.Cooldown.Cooldowns[userID], command)
	})
}

func (cmds *Commands) hasCooldown(userID string, command string) bool {
	cmds.Cooldown.Mutex.RLock()
	defer cmds.Cooldown.Mutex.RUnlock()

	return utils.FindInSlice(cmds.Cooldown.Cooldowns[userID], command)
}

func (cmds *Commands) Match(s *discordgo.Session, raw *discordgo.Message, context *Context) (*Command, []string) {
	var (
		collection bson.M
		err        error
		failure    string
		fields     = strings.Fields(context.Content)
	)

	if len(fields) == 0 {
		return nil, nil // some how hit this in testing...?
	}

	collection, err = database.Database.FindData(raw.GuildID)
	if err != nil {
		s.ChannelMessageSend(raw.ChannelID, fmt.Sprintf("```Failed to get the database collection: %s```", err.Error()))
		return nil, nil
	}

	context.Prefix = collection["prefix"].(string)

	if !strings.HasPrefix(fields[0], context.Prefix) {
		return nil, nil
	}

	fields[0] = strings.TrimPrefix(fields[0], context.Prefix)

	for _, command := range cmds.Commands {

		if fields[0] != command.Name && !utils.FindInSlice(command.Config.Alias, fields[0]) {
			continue
		}

		switch {

		case cmds.hasCooldown(raw.Author.ID, command.Name):
			return nil, nil

		case !utils.HasPerms(s, raw.GuildID, raw.Author.ID, command.Config.Perms):
			s.ChannelMessageSend(raw.ChannelID, "You do not have the required permissions to use this command.")
			return nil, nil // we want this before the failures.

		case command.Config.RequiresArgs && len(fields) < 2:
			failure = "You need args to use this command."

		case command.Config.RequiresMention && len(raw.Mentions) == 0:
			failure = "You have to mention someone to use this command."

		case command.Config.WhitelistedOnly && !database.Database.IsWhitelisted(raw.GuildID, raw.Author.ID):
			failure = "You have to be whitelisted to use this command."

		case command.Config.OwnerOnly && utils.GetGuildOwner(s, raw.GuildID) != raw.Author.ID:
			failure = "You have to be the guild owner to use this command."
		}

		if len(failure) > 0 {
			s.ChannelMessageSend(raw.ChannelID, failure)
			return nil, nil
		}

		return command, fields[0:]
	}

	return nil, nil
}

func (cmds *Commands) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.Bot || len(m.GuildID) == 0 {
		return
	}

	ctx := &Context{
		Content: strings.TrimSpace(m.Content),
	}

	cmd, fields := cmds.Match(s, m.Message, ctx)
	if cmd == nil {
		return
	}

	ctx.Fields = fields[1:]
	cmd.Run(s, m.Message, ctx)

	cmds.addCooldown(m.Author.ID, cmd.Name, cmd.Config.Cooldown)
}

type (
	CommandCooldown struct {
		Mutex     *sync.RWMutex
		Cooldowns map[string][]string
	}

	Command struct {
		Name   string
		Run    handler
		Config *Config
	}

	Commands struct {
		Commands []*Command
		Cooldown *CommandCooldown
	}

	Config struct {
		Alias           []string
		Cooldown        int
		OwnerOnly       bool
		Perms           int
		RequiresArgs    bool
		RequiresMention bool
		WhitelistedOnly bool
	}

	Context struct {
		Content string
		Prefix  string
		Fields  []string
	}

	handler func(*discordgo.Session, *discordgo.Message, *Context)
)
