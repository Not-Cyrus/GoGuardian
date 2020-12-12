package handlers

import (
	"fmt"

	"github.com/Not-Cyrus/GoGuardian/config"
	"github.com/bwmarrin/discordgo"
)

func BanHandler(s *discordgo.Session, ban *discordgo.GuildBanAdd) {
	if !config.Config.BanEnabled {
		return // Why you would EVER turn this off? Who knows.
	}
	bannedAnyone := readAudits(s, ban.GuildID, 22)
	if bannedAnyone {
		fmt.Println("Banned a bot/Account that was trying to ban everyone")
	}
}

func ChannelCreate(s *discordgo.Session, channel *discordgo.ChannelCreate) {
	if len(channel.GuildID) == 0 || !config.Config.ChannelSpamEnabled {
		return
	}
	bannedAnyone := readAudits(s, channel.GuildID, 10)
	if bannedAnyone {
		fmt.Println("Banned a bot/Account that was trying to mass generate channels")
	}
}

func ChannelRemove(s *discordgo.Session, channel *discordgo.ChannelDelete) {
	if len(channel.GuildID) == 0 || !config.Config.ChannelNukeEnabled {
		return
	}
	bannedAnyone := readAudits(s, channel.GuildID, 12)
	if bannedAnyone {
		fmt.Println("Banned a bot/Account that was trying to remove all channels")
	}
}

func KickHandler(s *discordgo.Session, channel *discordgo.GuildMemberRemove) {
	if !config.Config.KickEnabled {
		return // Again, Why would you turn this off???
	}
	bannedAnyone := readAudits(s, channel.GuildID, 20)
	if bannedAnyone {
		fmt.Println("Banned a bot/Account that was trying to kick everyone")
	}
}

func RoleCreate(s *discordgo.Session, channel *discordgo.GuildRoleCreate) {
	if !config.Config.RoleSpamEnabled {
		return
	}
	bannedAnyone := readAudits(s, channel.GuildID, 30)
	if bannedAnyone {
		fmt.Println("Banned a bot/Account that was trying to mass generate channels")
	}
}

func RoleRemove(s *discordgo.Session, channel *discordgo.GuildRoleDelete) {
	if !config.Config.RoleNukeEnabled {
		return
	}
	bannedAnyone := readAudits(s, channel.GuildID, 32)
	if bannedAnyone {
		fmt.Println("Banned a bot/Account that was trying to remove all roles")
	}
}
