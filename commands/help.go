package commands

import (
	"fmt"

	"github.com/Not-Cyrus/GoGuardian/utils"
	"github.com/bwmarrin/discordgo"
)

// we are putting this in it's own seperate file for easy access.

func (cmd *Commands) Help(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if len(ctx.Fields) == 0 {

		defaultHelp.Description = fmt.Sprintf("Below are all the commands categories, type `%shelp [category]` to view all commands in the category.", ctx.Prefix)
		defaultHelp.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)}
		defaultHelp.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: utils.BotUser.AvatarURL("500")}

		s.ChannelMessageSendEmbed(m.ChannelID, defaultHelp)
		return
	}

	certainHelp.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)}
	certainHelp.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: utils.BotUser.AvatarURL("500")}

	switch ctx.Fields[0] {

	case "information":
		certainHelp.Title = "Information Commands"
		certainHelp.Description = fmt.Sprintf("`%sserverinfo` | Returns information about the current server \n `%sbotinfo` | Shows information about the bot \n `%suserinfo [@user]` | Shows informati on about the mentioned user \n `%savatar [@user]` | Returns the mentioned users avatar \n `%smembercount` | Returns the server's member count \n `%sbanner` | Returns the server banner \n `%sservericon` | Returns the server icon", ctx.Prefix, ctx.Prefix, ctx.Prefix, ctx.Prefix, ctx.Prefix, ctx.Prefix, ctx.Prefix) //yeah ummm we don't talk about this..

	case "anti":
		certainHelp.Title = "Anti Commands"
		certainHelp.Description = fmt.Sprintf("`%swhitelist [@user] (server owner)` | Whitelists the mentioned user \n `%sunwhitelist [@user] (server owner)` | Dewhitelists the mentioned user \n `%swhitelisted` | Shows whitelisted list", ctx.Prefix, ctx.Prefix, ctx.Prefix)

	case "moderation":
		certainHelp.Title = "Moderation Commands"
		certainHelp.Description = fmt.Sprintf("`%sban [@user]` | Bans the mentioned user \n `%skick [@user]` | Kicks the mentioned user \n `%spurge [amount]` | Purges entered amount of messages\n`%slock` | Locks the channel\n `%sunlock` | Unlocks the channel\n`%sslowmode [time]` | Sets the channel slowmode to that time\n`%sunslowmode` | Disables slow mode\n`%smassunban` | Unbans all members in the server", ctx.Prefix, ctx.Prefix, ctx.Prefix, ctx.Prefix, ctx.Prefix, ctx.Prefix, ctx.Prefix, ctx.Prefix)

	case "settings":
		certainHelp.Title = "Settings"
		certainHelp.Description = fmt.Sprintf("`%sprefix [prefix]` | Sets the bot prefix\n`%slogchannel (server owner)` | Sets the log channel for all notifications relating to the anti-nuke.\n`%santiinvite [on/off]` | Enables/Disables the anti invite system", ctx.Prefix, ctx.Prefix, ctx.Prefix)
	default:

		defaultHelp.Title = fmt.Sprintf("%s anti-nuke", utils.BotUser.Username)
		defaultHelp.Description = fmt.Sprintf("Below are all the commands categories, type `%shelp [category]` to view all commands in the category.", ctx.Prefix)
		defaultHelp.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)}
		defaultHelp.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: utils.BotUser.AvatarURL("500")}

		s.ChannelMessageSendEmbed(m.ChannelID, defaultHelp)
		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, certainHelp)
}

func (cmd *Commands) Setup(s *discordgo.Session, m *discordgo.Message, ctx *Context) {

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Setup",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "How to setup:", Value: fmt.Sprintf("%s comes ready to protect your server on join, so you don't need to do anything besides **MAKE SURE** the bot is above every other role. The bot needs to be above everyone else to be able to ban.\n", utils.BotUser.Username), Inline: true},
			{Name: "Whitelist:", Value: fmt.Sprintf("To exclude someone from being banned, type `%swhitelist [@user]`. Only the server owner may whitelist and unwhitelist members.\n", ctx.Prefix)},
			{Name: "Logs:", Value: fmt.Sprintf("With %s, you can set a channel to log whenever %s bans someone for attempting to nuke. To set the log channel, type `%slog`\n", utils.BotUser.Username, utils.BotUser.Username, ctx.Prefix)},
			{Name: "Support:", Value: "If you need any help, click [here](https://github.com/Not-Cyrus/GoGuardian/issues)"},
		},
		Color: 0x36393F,
	})

}

var (
	certainHelp = &discordgo.MessageEmbed{
		Color: 0x36393F,
	}

	defaultHelp = &discordgo.MessageEmbed{

		Fields: []*discordgo.MessageEmbedField{
			{Name: "**Information**", Value: "*`Displays information related commands`*"},
			{Name: "**Anti**", Value: "*`Displays anti-nuke related commands`*"},
			{Name: "**Moderation**", Value: "*`Displays Moderation related commands.`*"},
			{Name: "**Settings**", Value: "‎‎‎‎‎‏‏‎*`Displays settings releated commands`*"},
		},

		Color: 0x36393F,
	}
)
