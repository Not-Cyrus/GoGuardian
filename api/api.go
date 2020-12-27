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

	// yeah lets just hide the sheer UGLINESS of this V on one line and forget about it ok? thanks
	route.Add("config", "- Turns on/off a protection setting that you specify. (type help config for more info)", []string{"antiadminrole - blocks people/bots from making roles have Administator permissions", "antiban - blocks wizz/nuke/destroyer bots from banning everyone", "antibots - bans any bot that gets invited (and the person who invited it)", "antichannelspam - stops people/bots from spam creating channels (mainly bots)", "antichannelnuke - stops people/bots from deleting lots of channels (mainly bots)", "antihijack - if GoGuardian ever gets comprimised and starts \"nuking\" the server it will leave.", "antikick - blocks wizz/nuke/destroyer bots from kicking everyone", "antimemberadmin - blocks people/bots from giving a member a role that has Administator permissions", "antirolespam - stops people/bots from spamm creating roles (mainly bots)", "antirolenuke - stops people/bots from deleting lots of roles (mainly bots)", "seconds - how many seconds it'll take before that moderation action is no longer classified as \"malicious\" by the bot.", "threshold - changes the amount of times someone can do a certain moderation action in x seconds (see above)"}, route.Config, true, true)

	route.Add("help", "Literally just help info what else do you expect a cookie?", []string{"WOW, You got a cookie!"}, route.Help, false, false)

	route.Add("invite", "Gives you an invite link to invite GoGuardian to your server.", []string{"ANOTHER COOKIE!!! (I'm running out of ideas for what to put here)"}, route.Invite, false, false)

	route.Add("whitelist", "whitelists a user so that they do not get affected by any protection settings.", []string{"whitelist @user"}, route.AddWhitelist, true, true)

	route.Add("unwhitelist", "unwhitelists a user so that they are affected by any toggled protection settings (default for all users minus the Guild owner.)", []string{"unwhitelist @user"}, route.RemoveWhitelist, true, true)
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
