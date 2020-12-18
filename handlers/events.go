package handlers

import (
	"fmt"

	"github.com/Not-Cyrus/GoGuardian/config"
	"github.com/bwmarrin/discordgo"
)

func AddMember(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
	if _, ok := config.WhitelistedIDs[member.User.ID]; !ok || !config.Config.AntiBotEnabled || !member.User.Bot {
		return
	}
	s.GuildBanCreateWithReason(member.GuildID, member.User.ID, "Bot Destroyed by https://github.com/Not-Cyrus/GoGuardian", 0)
}

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
	if !config.Config.ChannelSpamEnabled || len(channel.GuildID) == 0 {
		return
	}
	bannedAnyone := readAudits(s, channel.GuildID, 10)
	if bannedAnyone {
		fmt.Println("Banned a bot/Account that was trying to mass generate channels")
	}
}

func ChannelRemove(s *discordgo.Session, channel *discordgo.ChannelDelete) {
	if !config.Config.ChannelNukeEnabled || len(channel.GuildID) == 0 {
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

func MemberAdded(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
	var err error
	if !config.Config.AntiBotEnabled || !member.User.Bot {
		return
	}
	auditEntry := findAudit(s, member.GuildID, member.User.ID, 28)
	if auditEntry == nil {
		fmt.Println("hello")
		return
	}
	if _, ok := config.WhitelistedIDs[auditEntry.UserID]; !ok {
		return
	}
	err = s.GuildBanCreateWithReason(member.GuildID, member.User.ID, "Banned for being a bot that was invited by someone not whitelisted. - https://github.com/Not-Cyrus/GoGuardian", 0)
	err = s.GuildBanCreateWithReason(member.GuildID, auditEntry.UserID, "Banned for trying to invite a bot while not whitelisted. - https://github.com/Not-Cyrus/GoGuardian", 0)
	if err != nil {
		fmt.Println(fmt.Sprintf("Couldn't ban either the user/bot: %s", err.Error()))
		return
	}
	fmt.Println("Someone tried to invite a bot and got banned with the bot.")
}

func MemberRoleUpdate(s *discordgo.Session, member *discordgo.GuildMemberUpdate) {
	if !config.Config.MemberRoleUpdateEnabled {
		return
	}
	auditEntry := findAudit(s, member.GuildID, member.User.ID, 25)
	if auditEntry == nil {
		return
	}
	if _, ok := config.WhitelistedIDs[auditEntry.UserID]; !ok {
		return
	}
	for _, change := range auditEntry.Changes {
		roleID := change.NewValue.([]interface{})[0].(map[string]interface{})["id"].(string)
		guildRole, err := s.State.Role(member.GuildID, roleID)
		if err != nil {
			fmt.Println(fmt.Sprintf("Couldn't find the role: %s", err.Error()))
			return
		}
		if guildRole.Permissions&0x8 == 0x8 {
			err = s.GuildBanCreateWithReason(member.GuildID, auditEntry.UserID, "Banned for trying to give a role admin while not whitelisted. - https://github.com/Not-Cyrus/GoGuardian", 0)
			if err != nil {
				fmt.Println(fmt.Sprintf("Couldn't ban the person who gave a member a role without being whitelisted: %s", err.Error()))
				return
			}
			fmt.Println("Banned a user trying to give people admin roles without being whitelisted")
		}
	}
}

func RoleCreate(s *discordgo.Session, channel *discordgo.GuildRoleCreate) {
	if !config.Config.RoleSpamEnabled {
		return
	}
	bannedAnyone := readAudits(s, channel.GuildID, 30)
	if bannedAnyone {
		fmt.Println("Banned a bot/Account that was trying to mass generate roles")
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

func RoleUpdate(s *discordgo.Session, role *discordgo.GuildRoleUpdate) {
	var err error
	if !config.Config.RoleUpdateEnabled {
		return
	}
	auditEntry := findAudit(s, role.GuildID, role.Role.ID, 31)
	if auditEntry == nil {
		return
	}
	if _, ok := config.WhitelistedIDs[auditEntry.UserID]; !ok {
		return
	}
	guildRole, err := s.State.Role(role.GuildID, role.Role.ID)
	if err != nil {
		fmt.Println(fmt.Sprintf("Couldn't find the role: %s", err.Error()))
		return
	}
	if guildRole.Permissions&0x8 == 0x8 {
		err = s.GuildRoleDelete(role.GuildID, role.Role.ID)
		err = s.GuildBanCreateWithReason(role.GuildID, auditEntry.UserID, "Banned for trying to give a role admin while not whitelisted. - https://github.com/Not-Cyrus/GoGuardian", 0)
		if err != nil {
			fmt.Println(fmt.Sprintf("Couldn't properly ban the user/delete the role: %s", err.Error()))
			return
		}
		fmt.Println("Banned a user trying to create administrator roles without being whitelisted")
	}
}
