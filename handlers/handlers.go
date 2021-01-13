package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/Not-Cyrus/GoGuardian/utils"

	"github.com/bwmarrin/discordgo"
)

func readAudits(s *discordgo.Session, guildID string, auditType int) string {
	parsedData, configData := utils.FindConfig(guildID)

	auditMap := make(map[string]string)
	userMap := make(map[string]int)
	audits, err := s.GuildAuditLog(guildID, "", "", auditType, 25)
	if err != nil {
		utils.SendMessage(s, fmt.Sprintf("I can't read audits : %s", err.Error()), utils.GetGuildOwner(s, guildID))
		return ""
	}
	for _, entry := range audits.AuditLogEntries {
		if userMap[entry.UserID] >= configData.GetInt("Config", "Threshold") {
			if entry.UserID == DGUser.ID && configData.GetBool("Config", "AntiHijackEnabled") {
				if strings.Contains(entry.Reason, "https://github.com/Not-Cyrus/GoGuardian") {
					return "" // lazy as fuck, but it'll do the trick to stop false flags. I'll probably make this better soon idk
				}
				utils.SendMessage(s, "The bot has been comprimised. I have left the guild for your safety.", utils.GetGuildOwner(s, guildID)) // this is an important message so we'll DM the owner.
				s.GuildLeave(guildID)
			}
			err := s.GuildBanCreateWithReason(guildID, entry.UserID, "You just got destroyed by https://github.com/Not-Cyrus/GoGuardian", 0)
			if err != nil {
				utils.SendMessage(s, fmt.Sprintf("I have no perms to ban <@!%s>: %s", entry.UserID, err.Error()), utils.GetGuildOwner(s, guildID)) // this is an important message so we'll DM the owner.
				return ""
			}
			return entry.UserID
		}
		current := time.Now()
		entryTime, err := discordgo.SnowflakeTimestamp(entry.ID)
		if err != nil {
			utils.SendMessage(s, fmt.Sprintf("how the fuck did this happen: %s", err.Error()), "")
			return ""
		}
		if current.Sub(entryTime).Round(1*time.Second).Seconds() <= configData.GetFloat64("Config", "Seconds") {
			if _, ok := auditMap[entry.ID]; !ok {
				inArray, _ := utils.InArray(guildID, "WhitelistedIDs", parsedData, entry.UserID)
				if !inArray {
					auditMap[entry.ID] = entry.ID
					userMap[entry.UserID]++
				}
			}
		}
	}
	return ""
}

func findAudit(s *discordgo.Session, guildID, targetID string, auditType int) *discordgo.AuditLogEntry {
	audits, err := s.GuildAuditLog(guildID, "", "", auditType, 10) // we really don't need 25 here so we'll use 10 instead (I could probably just use one but whatever)
	if err != nil {
		utils.SendMessage(s, fmt.Sprintf("I can't read audits: %s | if you think this is a mistake make an issue at https://github.com/Not-Cyrus/GoGuardian/issues", err.Error()), utils.GetGuildOwner(s, guildID)) // scrap this last comment as it's now going to be a public bot.
		return nil
	}
	for _, entry := range audits.AuditLogEntries {
		if entry.TargetID == targetID {
			return entry
		}
	}
	return nil
}
