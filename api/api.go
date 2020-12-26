package api

import (
	"fmt"

	"github.com/Not-Cyrus/GoGuardian/utils"
	"github.com/valyala/fastjson"

	"github.com/Not-Cyrus/GoGuardian/commands"
	"github.com/Not-Cyrus/GoGuardian/handlers"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) Setup() {
	token := getToken()
	if len(token) == 0 {
		fmt.Print("Enter your token: ")
		fmt.Scan(&token)
	}

	b.DS, err = discordgo.New("Bot " + token)
	if err != nil {
		panic("Couldn't use said token")
	}

	b.BU, err = b.DS.User("@me")
	if err != nil {
		panic("Couldn't find a local user???")
	}

	b.DS.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	handlerNames := []interface{}{handlers.BanHandler, handlers.ChannelCreate, handlers.ChannelRemove, handlers.KickHandler, handlers.MemberAdded, handlers.ReadyHandler, handlers.MemberRoleUpdate, handlers.RoleCreate, handlers.RoleRemove, handlers.RoleUpdate}
	for _, handler := range handlerNames {
		b.DS.AddHandler(handler)
	}

	b.DS.AddHandler(route.MessageCreate)
	route.Add("config", route.Config)
	route.Add("whitelist", route.AddWhitelist)
	route.Add("unwhitelist", route.RemoveWhitelist)
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

func getToken() string {
	fileContents := utils.ReadFile("Config.json")
	parsed, err := parser.Parse(fileContents)
	if err != nil {
		panic("Couldn't parse Config.json to get your Token")
	}
	if fastjson.Exists([]byte(fileContents), "Token") {
		return string(parsed.GetStringBytes("Token"))
	}
	return ""
}

type (
	Bot struct {
		DS *discordgo.Session
		BU *discordgo.User
	}
)

var (
	err    error
	token  string
	parser fastjson.Parser
	route  = commands.New()
)
