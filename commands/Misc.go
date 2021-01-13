package commands

import "github.com/bwmarrin/discordgo"

func (cmd *Commands) Invite(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	s.ChannelMessageSend(message.ChannelID, "https://discord.com/api/oauth2/authorize?client_id=775890268364210196&permissions=8&scope=bot")
}
