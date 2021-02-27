package commands

import (
	"fmt"
	"strings"

	"github.com/Not-Cyrus/GoGuardianv2/database"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func (cmd *Commands) AntiInvite(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if !(ctx.Fields[0] == "on" || ctx.Fields[0] == "off") {
		return
	}

	if _, err := database.Database.SetData(m.GuildID, "anti-invite", ctx.Fields[0]); err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Set Anti-Invite to %s", ctx.Fields[0]))
}

func (cmd *Commands) LoggingChannel(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if set, err := database.Database.SetData(m.GuildID, "log-channel", m.ChannelID); !set {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Set the logging channel to the current channel")
}

func (cmd *Commands) Prefix(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if set, err := database.Database.SetData(m.GuildID, "prefix", ctx.Fields[0]); !set {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("Prefix has been set to `%s`", ctx.Fields[0]),
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
		Color:  0x36393F,
	})
}

func (cmd *Commands) Settings(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	data, err := database.Database.FindData(m.GuildID)
	guild, _ := s.State.Guild(m.GuildID)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var (
		embed = &discordgo.MessageEmbed{
			Title:  fmt.Sprintf("%s current settings", guild.Name),
			Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
			Color:  0x36393F,
		}
		tempValue string
	)

	for index, value := range data {
		if index == "users" || index == "_id" || index == "guild_id" {
			continue
		}

		switch value.(string) {
		case "on":
			tempValue = "<:enabled:799507631274197022>"

		case "off":
			tempValue = "<:disabled:799507673648594954>"

		case "nil":
			tempValue = "<:disabled:799507673648594954>"

		default:
			tempValue = value.(string)
			if index == "log-channel" {
				tempValue = fmt.Sprintf("<#%s>", value.(string))
			}
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   index,
			Value:  tempValue,
			Inline: false,
		})

	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func (cmd *Commands) Whitelist(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if whitelisted, err := database.Database.SetWhitelist(m.GuildID, m.Mentions[0], true); !whitelisted {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Whitelisted that user.")
}

func (cmd *Commands) Unwhitelist(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if whitelisted, err := database.Database.SetWhitelist(m.GuildID, m.Mentions[0], false); !whitelisted {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Unwhitelisted that user.")
}

func (cmd *Commands) ViewWhitelisted(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	data, err := database.Database.FindData(m.GuildID)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var whitelistedUsers []string

	for _, userID := range data["users"].(bson.A) {
		user, _ := s.User(userID.(string))
		whitelistedUsers = append(whitelistedUsers, fmt.Sprintf("📋 | %s#%s", user.Username, user.Discriminator))
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Whitelisted Members",
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s | made by https://github.com/Not-Cyrus", m.Author.Username)},
		Description: strings.Join(whitelistedUsers, "\n"),
		Color:       0x36393F,
	})
}
