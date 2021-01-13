package handlers

import (
	"fmt"

	"github.com/Not-Cyrus/GoGuardian/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/valyala/fastjson"
)

func BanHandler(s *discordgo.Session, ban *discordgo.GuildBanAdd) {
	_, configData := utils.FindConfig(ban.GuildID)
	if !configData.GetBool("Config", "BanProtection") {
		return // Why you would EVER turn this off? Who knows.
	}
	bannedAnyone := readAudits(s, ban.GuildID, 22)
	if len(bannedAnyone) != 0 {
		utils.SendMessage(s, fmt.Sprintf("Banned <@!%s> who was trying to ban everyone", bannedAnyone), utils.GetGuildOwner(s, ban.GuildID))
	}
}

func ChannelCreate(s *discordgo.Session, channel *discordgo.ChannelCreate) {
	_, configData := utils.FindConfig(channel.GuildID)
	if !configData.GetBool("Config", "ChannelSpamProtection") || len(channel.GuildID) == 0 {
		return
	}
	bannedAnyone := readAudits(s, channel.GuildID, 10)
	if len(bannedAnyone) != 0 {
		utils.SendMessage(s, fmt.Sprintf("Banned <@!%s> who was trying to mass generate channels", bannedAnyone), utils.GetGuildOwner(s, channel.GuildID))
	}
}

func ChannelRemove(s *discordgo.Session, channel *discordgo.ChannelDelete) {
	_, configData := utils.FindConfig(channel.GuildID)
	if !configData.GetBool("Config", "ChannelNukeProtection") || len(channel.GuildID) == 0 {
		return
	}
	bannedAnyone := readAudits(s, channel.GuildID, 12)
	if len(bannedAnyone) != 0 {
		utils.SendMessage(s, fmt.Sprintf("Banned <@!%s> who was trying to remove all channels", bannedAnyone), utils.GetGuildOwner(s, channel.GuildID))
	}
}

func GuildCreate(s *discordgo.Session, guild *discordgo.GuildCreate) {
	DGUser, err = s.User("@me") // this apparently fires before "ready" so it'll go here I guess lol
	originalData, configData := utils.FindConfig(guild.ID)
	if configData == nil {
		return
	}

	auditEntry := findAudit(s, guild.ID, DGUser.ID, 28)
	if auditEntry == nil {
		return
	}
	inArray, _ := utils.InArray(guild.ID, "WhitelistedIDs", originalData, auditEntry.UserID)
	inArray2, _ := utils.InArray(guild.ID, "WhitelistedIDs", originalData, guild.OwnerID)

	if !inArray {

		guildArray := originalData.GetArray("Guilds", guild.ID, "WhitelistedIDs")
		originalData.Get("Guilds", guild.ID, "WhitelistedIDs").SetArrayItem(len(guildArray), fastjson.MustParse(fmt.Sprintf(`"%s"`, auditEntry.UserID)))

		if auditEntry.UserID != guild.OwnerID && !inArray2 {
			originalData.Get("Guilds", guild.ID, "WhitelistedIDs").SetArrayItem(len(guildArray)+1, fastjson.MustParse(fmt.Sprintf(`"%s"`, guild.OwnerID)))
		}

		utils.SaveJSON(nil, nil, originalData, "")
	}

}

func KickHandler(s *discordgo.Session, channel *discordgo.GuildMemberRemove) {
	_, configData := utils.FindConfig(channel.GuildID)
	if !configData.GetBool("Config", "KickProtection") {
		return // Again, Why would you turn this off???
	}
	bannedAnyone := readAudits(s, channel.GuildID, 20)
	if len(bannedAnyone) != 0 {
		utils.SendMessage(s, fmt.Sprintf("Banned <@!%s> who was trying to kick everyone", bannedAnyone), utils.GetGuildOwner(s, channel.GuildID))
	}
}

func MemberAdded(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
	var err error
	parsedData, configData := utils.FindConfig(member.GuildID)
	if !configData.GetBool("Config", "AntiBotProtection") || !member.User.Bot {
		return
	}
	auditEntry := findAudit(s, member.GuildID, member.User.ID, 28)
	if auditEntry == nil {
		return
	}
	inArray, _ := utils.InArray(member.GuildID, "WhitelistedIDs", parsedData, auditEntry.UserID)
	if inArray {
		return
	}
	err = s.GuildBanCreateWithReason(member.GuildID, member.User.ID, "Banned for being a bot that was invited by someone not whitelisted. - https://github.com/Not-Cyrus/GoGuardian", 0)
	err = s.GuildBanCreateWithReason(member.GuildID, auditEntry.UserID, "Banned for trying to invite a bot while not whitelisted. - https://github.com/Not-Cyrus/GoGuardian", 0)
	if err != nil {
		utils.SendMessage(s, fmt.Sprintf("Couldn't ban <@!%s> or <@!%s> (Bot Check): %s", member.User.ID, auditEntry.UserID, err.Error()), utils.GetGuildOwner(s, member.GuildID))
		return
	}
	utils.SendMessage(s, fmt.Sprintf("<@!%s> tried to invite <@!%s> (A bot) and got banned.", auditEntry.UserID, member.User.ID), utils.GetGuildOwner(s, member.GuildID))
}

func MemberRoleUpdate(s *discordgo.Session, member *discordgo.GuildMemberUpdate) {
	parsedData, configData := utils.FindConfig(member.GuildID)
	if !configData.GetBool("Config", "MemberRoleUpdateProtection") {
		return
	}
	auditEntry := findAudit(s, member.GuildID, member.User.ID, 25)
	if auditEntry == nil {
		return
	}
	inArray, _ := utils.InArray(member.GuildID, "WhitelistedIDs", parsedData, auditEntry.UserID)
	if inArray {
		return
	}
	for _, change := range auditEntry.Changes {
		roleID := change.NewValue.([]interface{})[0].(map[string]interface{})["id"].(string)
		guildRole, err := s.State.Role(member.GuildID, roleID)
		if err != nil {
			utils.SendMessage(s, fmt.Sprintf("Couldn't find the role: %s", err.Error()), "")
			return
		}
		if guildRole.Permissions&0x8 == 0x8 {
			err = s.GuildMemberRoleRemove(member.GuildID, auditEntry.TargetID, roleID)
			err = s.GuildBanCreateWithReason(member.GuildID, auditEntry.UserID, "Banned for trying to give a role admin while not whitelisted. - https://github.com/Not-Cyrus/GoGuardian", 0)
			if err != nil {
				utils.SendMessage(s, fmt.Sprintf("Couldn't ban <@!%s> (Member Admin Role check): %s", member.User.ID, err.Error()), utils.GetGuildOwner(s, member.GuildID))
				return
			}
			utils.SendMessage(s, fmt.Sprintf("Banned <@!%s> who tried to give people admin roles without being whitelisted", auditEntry.UserID), utils.GetGuildOwner(s, member.GuildID))
		}
	}
}

func RoleCreate(s *discordgo.Session, role *discordgo.GuildRoleCreate) {
	_, configData := utils.FindConfig(role.GuildID)
	if !configData.GetBool("Config", "RoleSpamProtection") {
		return
	}
	bannedAnyone := readAudits(s, role.GuildID, 30)
	if len(bannedAnyone) != 0 {
		utils.SendMessage(s, fmt.Sprintf("Banned <@!%s> who was trying to mass generate roles", bannedAnyone), utils.GetGuildOwner(s, role.GuildID))
	}
}

func RoleRemove(s *discordgo.Session, role *discordgo.GuildRoleDelete) {
	_, configData := utils.FindConfig(role.GuildID)
	if !configData.GetBool("Config", "RoleNukeProtection") {
		return
	}
	bannedAnyone := readAudits(s, role.GuildID, 32)
	if len(bannedAnyone) != 0 {
		utils.SendMessage(s, fmt.Sprintf("Banned <@!%s> who was trying to remove all roles", bannedAnyone), utils.GetGuildOwner(s, role.GuildID))
	}
}

func RoleUpdate(s *discordgo.Session, role *discordgo.GuildRoleUpdate) {
	parsedData, configData := utils.FindConfig(role.GuildID)
	if !configData.GetBool("Config", "RoleUpdateProtection") {
		return
	}
	auditEntry := findAudit(s, role.GuildID, role.Role.ID, 31)
	if auditEntry == nil {
		return
	}
	inArray, _ := utils.InArray(role.GuildID, "WhitelistedIDs", parsedData, auditEntry.UserID)
	if inArray {
		return
	}
	guildRole, err := s.State.Role(role.GuildID, role.Role.ID)
	if err != nil {
		utils.SendMessage(s, fmt.Sprintf("Couldn't find the role: %s", err.Error()), "")
		return
	}
	if guildRole.Permissions&0x8 == 0x8 {
		err = s.GuildRoleDelete(role.GuildID, role.Role.ID)
		err = s.GuildBanCreateWithReason(role.GuildID, auditEntry.UserID, "Banned for trying to give a role admin while not whitelisted. - https://github.com/Not-Cyrus/GoGuardian", 0)
		if err != nil {
			utils.SendMessage(s, fmt.Sprintf("Couldn't ban <@!%s> (Create Admin Role check): %s", auditEntry.UserID, err.Error()), utils.GetGuildOwner(s, role.GuildID))
			return
		}
		utils.SendMessage(s, fmt.Sprintf("Banned <@!%s> who was trying to create administrator roles without being whitelisted", auditEntry.UserID), utils.GetGuildOwner(s, role.GuildID))
	}
}

func WebhookCreate(s *discordgo.Session, webhook *discordgo.WebhooksUpdate) {
	var err error
	parsedData, configData := utils.FindConfig(webhook.GuildID)
	if !configData.GetBool("Config", "WebhookProtection") {
		return
	}

	webhooks, err := s.ChannelWebhooks(webhook.ChannelID)
	if err != nil {
		utils.SendMessage(s, fmt.Sprintf("Couldn't fetch the webhooks: %s", err.Error()), utils.GetGuildOwner(s, webhook.GuildID))
		return
	}

	for _, webhook := range webhooks {
		inArray, _ := utils.InArray(webhook.GuildID, "WhitelistedIDs", parsedData, webhook.User.ID)
		if inArray {
			continue
		}

		err = s.WebhookDelete(webhook.ID)
		err = s.GuildBanCreateWithReason(webhook.GuildID, webhook.User.ID, "Banned for trying to create webhooks roles without being whitelisted. - https://github.com/Not-Cyrus/GoGuardian", 0)
		if err != nil {
			utils.SendMessage(s, fmt.Sprintf("Couldn't take moderation action against the webhook (or person who made): %s", err.Error()), utils.GetGuildOwner(s, webhook.GuildID))
			return
		}
		utils.SendMessage(s, fmt.Sprintf("Banned <@!%s> who was trying to create webhooks roles without being whitelisted", webhook.User.ID), utils.GetGuildOwner(s, webhook.GuildID))
	}
}

var (
	DGUser *discordgo.User
	err    error
)
