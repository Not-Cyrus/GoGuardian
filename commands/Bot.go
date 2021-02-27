package commands

import (
	"fmt"
	"time"

	"github.com/Not-Cyrus/GoGuardian/events"
	"github.com/Not-Cyrus/GoGuardian/utils"

	"github.com/bwmarrin/discordgo"
)

func (cmd *Commands) BotInfo(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Bot Info",

		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "Name:", Value: utils.BotUser.Username, Inline: true},
			&discordgo.MessageEmbedField{Name: "Server Count:", Value: fmt.Sprint(events.GuildCount), Inline: true},
			&discordgo.MessageEmbedField{Name: "User Count:", Value: fmt.Sprint(events.MemberCount), Inline: true},
			&discordgo.MessageEmbedField{Name: "Ping:", Value: fmt.Sprintf("%s", s.HeartbeatLatency().Round(1*time.Millisecond)), Inline: true},
			&discordgo.MessageEmbedField{Name: "discordgo Version", Value: "v0.22.0", Inline: true},
			&discordgo.MessageEmbedField{Name: "Shard", Value: fmt.Sprint(s.ShardID), Inline: true},
		},

		Footer:    &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: utils.BotUser.AvatarURL("500")},
		Color:     0x36393F,
	})
}

func (cmd *Commands) Credits(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Credits",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "Creators:", Value: "[!fishgang Cy](https://github.com/Not-Cyrus) - Bot developer"},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
		Color:  0x36393F,
	})
}

func (cmd *Commands) Invite(s *discordgo.Session, m *discordgo.Message, ctx *Context) {

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "Bot Invite", Value: fmt.Sprintf("[Click Here](https://discord.com/api/oauth2/authorize?client_id=%s&permissions=8&scope=bot)", utils.BotUser.ID), Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
		Color:  0x36393F,
	})
}

func (cmd *Commands) Ping(s *discordgo.Session, m *discordgo.Message, ctx *Context) {

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("Bot Ping: `%s` | on shard `%d`", s.HeartbeatLatency().Round(1*time.Millisecond), s.ShardID),
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
		Color:  0x36393F,
	})
}
