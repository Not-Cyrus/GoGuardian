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

	var coolInt string
	if len(ctx.Fields) > 1 {
		_, err := strconv.Atoi(ctx.Fields[1])
		if err != nil {
			s.ChannelMessageSend(message.ChannelID, "You sure that's a number?")
			return
		}
		coolInt = ctx.Fields[1]
	}

	switch parseStr {
	case "Seconds":
		object.Set(parseStr, fastjson.MustParse(coolInt))
		utils.SaveJSON(s, message, originalData, fmt.Sprintf("%s has been set to %s", parseStr, coolInt))
		return
	case "Threshold":
		object.Set(parseStr, fastjson.MustParse(coolInt)) // AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA (hi)
		utils.SaveJSON(s, message, originalData, fmt.Sprintf("%s has been set to %s", parseStr, coolInt))
		return
	default:
		set := strconv.FormatBool(!object.Get(parseStr).GetBool())
		object.Set(parseStr, fastjson.MustParse(set))
		utils.SaveJSON(s, message, originalData, fmt.Sprintf("%s has been set to %s", parseStr, set))
		return
	}

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
	case "antiwebhook":
		parse = "WebhookProtection"
	case "seconds":
		parse = "Seconds"
	case "threshold":
		parse = "Threshold"
	default:
		parse = "Failed" // please help my sanity feels like I am a valve dev working on CS:GO
	}
	return parse
}

var (
	parser fastjson.Parser
)
