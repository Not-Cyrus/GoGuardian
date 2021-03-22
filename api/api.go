package api

import (
	"fmt"
	"sync"

	"github.com/Not-Cyrus/GoGuardian/commands"
	"github.com/Not-Cyrus/GoGuardian/events"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) Run() {
	for _, session := range b.Sessions {
		if session != nil {
			session.Open()
		}
	}
}

func (b *Bot) Shard(token string, shardCount, shardID int) {

	// ** Setup session ** //

	s, err := discordgo.New(fmt.Sprintf("Bot %s", token))

	if err != nil {
		fmt.Printf("[Error creating a shard]: %s", err.Error())
	}

	s.ShardCount = shardCount
	s.ShardID = shardID

	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers)

	// ** Handlers ** //

	handlers := []interface{}{
		events.AntiInvite,
		events.BanHandler,
		events.ChannelCreate,
		events.ChannelRemove,
		events.CreateGuild,
		events.DeleteGuild,
		events.KickHandler,
		events.MemberJoin,
		events.MemberLeave,
		events.MemberRoleUpdate,
		commandRoute.MessageCreate,
		events.Ready,
		events.RoleCreate,
		events.RoleRemove,
		events.WebhookCreate,
	}

	for _, handler := range handlers {
		s.AddHandler(handler)
	}

	b.Sessions[shardID] = s
}

func (b *Bot) Stop() {
	for _, session := range b.Sessions {
		session.Close()
	}
}

type (
	Bot struct {
		Sessions []*discordgo.Session
	}
)

var (
	commandRoute = &commands.Commands{
		Cooldown: &commands.CommandCooldown{
			Cooldowns: make(map[string][]string),
			Mutex:     &sync.RWMutex{},
		},
	}
	err error
)
