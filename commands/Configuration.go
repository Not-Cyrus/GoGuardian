package commands

import (
	"fmt"
	"strconv"

	"github.com/Not-Cyrus/GoGuardian/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/valyala/fastjson"
)

func (cmd *Commands) AddWhitelist(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	if len(message.Mentions) == 0 {
		s.ChannelMessageSend(message.ChannelID, "Mention someone to whitelist them")
		return
	}

	originalData, parsedData := utils.FindConfig(message.GuildID)
	if parsedData == nil {
		s.ChannelMessageSend(message.ChannelID, "Something happened and this couldn't be completed.")
		return
	}

	inArray, _ := utils.InArray(message.GuildID, "WhitelistedIDs", originalData, message.Mentions[0].ID)
	if inArray {
		s.ChannelMessageSend(message.ChannelID, "They're already whitelisted..?")
		return
	}

	guildArray := originalData.GetArray("Guilds", message.GuildID, "WhitelistedIDs")

	originalData.Get("Guilds", message.GuildID, "WhitelistedIDs").SetArrayItem(len(guildArray), fastjson.MustParse(fmt.Sprintf(`"%s"`, message.Mentions[0].ID)))
	utils.SaveJSON(s, message, originalData, "Added their whitelist.")
}

func (cmd *Commands) Config(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	originalData, _ := utils.FindConfig(message.GuildID)
	object, err := originalData.Get("Guilds", message.GuildID, "Config").Object()
	if err != nil {
		s.ChannelMessageSend(message.ChannelID, "Couldn't make the json data an object")
	}

	parseStr := validArg(ctx.Fields[0])
	if parseStr == "Failed" {
		s.ChannelMessageSend(message.ChannelID, "Not a valid argument.")
		return
	}

	boolSet := strconv.FormatBool(!object.Get(parseStr).GetBool())
	// AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA (hi)
	object.Set(parseStr, fastjson.MustParse(boolSet))

	utils.SaveJSON(s, message, originalData, fmt.Sprintf("%s has been set to %s", parseStr, boolSet))
}

func (cmd *Commands) RemoveWhitelist(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	if len(message.Mentions) == 0 {
		s.ChannelMessageSend(message.ChannelID, "Mention someone to whitelist them")
		return
	}
	originalData, parsedData := utils.FindConfig(message.GuildID)
	if parsedData == nil {
		s.ChannelMessageSend(message.ChannelID, "Something happened and this couldn't be completed.")
		return
	}

	inArray, index := utils.InArray(message.GuildID, "WhitelistedIDs", originalData, message.Mentions[0].ID)
	if inArray {
		originalData.Get("Guilds", message.GuildID, "WhitelistedIDs").Del(fmt.Sprint(index))
		utils.SaveJSON(s, message, originalData, "removed their whitelist.")
		return
	}
	s.ChannelMessageSend(message.ChannelID, "They weren't whitelisted..?")
}

func validArg(arg string) string {
	var parse string
	switch arg {
	case "antiadminrole":
		parse = "RoleUpdateProtection"
	case "antiban":
		parse = "BanProtection"
	case "antibots":
		parse = "AntiBotProtection"
	case "antichannelspam":
		parse = "ChannelSpamProtection"
	case "antichannelnuke":
		parse = "ChannelNukeProtection"
	case "antihijack":
		parse = "HijackProtection"
	case "antikick":
		parse = "KickProtection"
	case "antimemberadmin":
		parse = "MemberRoleUpdateProtection"
	case "antirolespam":
		parse = "RoleSpamProtection"
	case "antirolenuke":
		parse = "RoleNukeProtection"
	default:
		parse = "Failed" // please help my sanity feels like I am a valve dev working on CS:GO
	}
	return parse
}

var (
	defaultConfig = `{"WhitelistedIDs": [],"Config": {"Threshold":2,"Seconds":2,"BanProtection":true,"KickProtection":true,"HijackProtection":true,"AntiBotProtection":true,"RoleSpamProtection":true,"RoleNukeProtection":true,"RoleUpdateProtection":true,"ChannelSpamProtection":true,"ChannelNukeProtection":true,"MemberRoleUpdateProtection":true}}`
	parser        fastjson.Parser
)
