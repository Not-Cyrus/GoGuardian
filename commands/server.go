package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (cmd *Commands) MemberCount(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: guild.Name,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "Members", Value: fmt.Sprint(guild.MemberCount), Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
		Color:  0x36393F,
	})

}

func (cmd *Commands) Nuke(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}

	_, err = s.ChannelDelete(channel.ID)
	if err != nil {
		return
	}

	channel, err = s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
		Name:                 channel.Name,
		Type:                 channel.Type,
		Topic:                channel.Topic,
		RateLimitPerUser:     channel.RateLimitPerUser,
		Position:             channel.Position,
		PermissionOverwrites: channel.PermissionOverwrites,
		ParentID:             channel.ParentID,
		NSFW:                 channel.NSFW,
	})
	if err != nil {
		return
	}

	s.ChannelMessageSend(channel.ID, fmt.Sprintf("http://gph.is/1s201Ez\n Channel has been nuked by: %s", m.Author.Username))
}

func (cmd *Commands) ServerBanner(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return
	}

	if len(guild.Banner) == 0 {
		s.ChannelMessageSend(m.ChannelID, "There is no guild banner.")
		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s's server banner", guild.Name),
		Image: &discordgo.MessageEmbedImage{URL: guild.Banner},

		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
		Color:  0x36393F,
	})

}

func (cmd *Commands) ServerIcon(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return
	}

	if len(guild.IconURL()) == 0 {
		s.ChannelMessageSend(m.ChannelID, "There is no guild icon.")
		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s's server icon", guild.Name),
		Image: &discordgo.MessageEmbedImage{
			URL: guild.IconURL(),
		},
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
		Color:  0x36393F,
	})

}

func (cmd *Commands) ServerInfo(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return
	}

	guildTime, _ := discordgo.SnowflakeTimestamp(guild.ID)

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s | Server Info", guild.Name),
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "Members:", Value: fmt.Sprint(guild.MemberCount), Inline: true},
			&discordgo.MessageEmbedField{Name: "Roles:", Value: fmt.Sprint(len(guild.Roles)), Inline: true},
			&discordgo.MessageEmbedField{Name: "Channels:", Value: fmt.Sprint(len(guild.Channels)), Inline: true},
			&discordgo.MessageEmbedField{Name: "Server Created on", Value: guildTime.Format("01/02/2006"), Inline: true},
			&discordgo.MessageEmbedField{Name: "Server Owner", Value: fmt.Sprintf("<@%s>", guild.OwnerID), Inline: true},
			&discordgo.MessageEmbedField{Name: "Server Region", Value: guild.Region, Inline: true},
		},
		Footer:    &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: guild.IconURL()},
		Color:     0x36393F,
	})

}
