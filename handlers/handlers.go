package handlers

import (
	"fmt"
	"time"

	"github.com/Not-Cyrus/GoGuardian/utils"

	"github.com/Not-Cyrus/GoGuardian/config"
	"github.com/bwmarrin/discordgo"
)

func readAudits(s *discordgo.Session, guildID string, auditType int) string {
	auditMap := make(map[string]string)
	userMap := make(map[string]int)
	audits, err := s.GuildAuditLog(guildID, "", "", auditType, 25)
	if err != nil {
		utils.SendMessage(s, fmt.Sprintf("I can't read audits : %s", err.Error()), "")
		return ""
	}
	for _, entry := range audits.AuditLogEntries {
		if userMap[entry.UserID] >= config.Config.Threshold {
			if entry.UserID == DGUser.ID && config.Config.AntiHijackEnabled {
				utils.SendMessage(s, "The bot has been comprimised. PANIC! (Also refresh the token)", utils.GetGuildOwner(s, guildID)) // this is an important message so we'll DM the owner.
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
		if current.Sub(entryTime).Round(1*time.Second).Seconds() <= config.Config.Seconds {
			if _, ok := auditMap[entry.ID]; !ok {
				if _, whitelisted := config.WhitelistedIDs[entry.UserID]; !whitelisted {
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
		utils.SendMessage(s, fmt.Sprintf("I can't read audits: %s", err.Error()), "") // not TOO important as it can error if the connection drops.
		return nil
	}
	for _, entry := range audits.AuditLogEntries {
		if entry.TargetID == targetID {
			return entry
		}
	}
	return nil
}
