package commands

import (
	"fmt"
	"strconv"

	"github.com/Not-Cyrus/GoGuardian/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/valyala/fastjson"
)

func (cmd *Commands) Config(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	parsedData, err := parser.Parse(utils.ReadFile("config.json"))
	if err != nil {
		fmt.Printf("Couldn't parse json data: %s\n", err.Error())
		return
	}
	object, err := parsedData.Object()
	if err != nil {
		panic("Couldn't make the json data an object")
	}
	parseStr := validArg(ctx.Fields[0])
	if parseStr == "Failed" {
		s.ChannelMessageSend(message.ChannelID, "Not a valid argument.")
		return
	}
	object.Set(parseStr, fastjson.MustParse(strconv.FormatBool(!parsedData.GetBool(parseStr))))
	utils.Writefile("config.json", string(object.MarshalTo(nil)))
}

func (cmd *Commands) AddWhitelist(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	if len(message.Mentions) == 0 {
		s.ChannelMessageSend(message.ChannelID, "Mention someone to whitelist them")
		return
	}
	parsedData, err := parser.Parse(utils.ReadFile("config.json"))
	if err != nil {
		fmt.Printf("Couldn't parse json data: %s\n", err.Error())
		return
	}
	inArray, _ := utils.InArray("WhitelistedIDs", parsedData, message, message.Mentions[0].ID)
	if inArray {
		s.ChannelMessageSend(message.ChannelID, "They're already whitelisted..?")
		return
	}
	parsedData.Get("WhitelistedIDs").SetArrayItem(len(parsedData.GetArray("WhitelistedIDs")), fastjson.MustParse(fmt.Sprintf(`"%s"`, message.Mentions[0].ID))) // for some reason it converts it to an int..?
	s.ChannelMessageSend(message.ChannelID, "Added their whitelist.")
	utils.Writefile("config.json", string(parsedData.MarshalTo(nil)))
}

func (cmd *Commands) RemoveWhitelist(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	if len(message.Mentions) == 0 {
		s.ChannelMessageSend(message.ChannelID, "Mention someone to whitelist them")
		return
	}
	parsedData, err := parser.Parse(utils.ReadFile("config.json"))
	if err != nil {
		fmt.Printf("Couldn't parse json data: %s\n", err.Error())
		return
	}
	inArray, index := utils.InArray("WhitelistedIDs", parsedData, message, message.Mentions[0].ID)
	if inArray {
		parsedData.Get("WhitelistedIDs").Del(fmt.Sprint(index))
		utils.Writefile("config.json", string(parsedData.MarshalTo(nil)))
		s.ChannelMessageSend(message.ChannelID, "removed their whitelist.")
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
		parse = "Failed"
	}
	return parse
}

var (
	parser fastjson.Parser
)
