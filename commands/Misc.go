package commands

import "github.com/Not-Cyrus/GoGuardian/discord"

func (cmd *Commands) Invite(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	s.ChannelMessageSend(message.ChannelID, "https://top.gg/bot/775890268364210196")
}
