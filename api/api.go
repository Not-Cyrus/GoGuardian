package api

import (
	"github.com/Not-Cyrus/GoGuardian/config"
	"github.com/Not-Cyrus/GoGuardian/handlers"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) Setup() {
	b.DS, err = discordgo.New(config.Config.Token)
	if err != nil {
		panic("Couldn't use said token")
	}
	b.BU, err = b.DS.User("@me")
	if err != nil {
		panic("Couldn't find a local user???")
	}
	b.DS.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	handlerNames := []interface{}{handlers.BanHandler, handlers.ChannelCreate, handlers.ChannelRemove, handlers.KickHandler, handlers.RoleCreate, handlers.RoleRemove}
	for _, handler := range handlerNames {
		b.DS.AddHandler(handler)
	}
}

func (b *Bot) Run() error {
	err := b.DS.Open()
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) Stop() {
	b.DS.Close()
}

type (
	Bot struct {
		DS *discordgo.Session
		BU *discordgo.User
	}
)

var (
	err error
)
