package handlers

import (
	"fmt"
	"time"

	"github.com/Not-Cyrus/GoGuardian/config"
	"github.com/bwmarrin/discordgo"
)

func readAudits(s *discordgo.Session, guildID string, auditType int) bool {
	auditMap := make(map[string]string)
	userMap := make(map[string]int)
	audits, err := s.GuildAuditLog(guildID, "", "", auditType, 25)
	if err != nil {
		panic("I can't read audits (wtf am I supposed to do???)")
	}
	for _, entry := range audits.AuditLogEntries {
		if userMap[entry.UserID] >= config.Config.Threshold {
			err := s.GuildBanCreateWithReason(guildID, entry.UserID, "You just got destroyed by https://github.com/Not-Cyrus/GoGuardian", 0)
			if err != nil {
				fmt.Println(fmt.Sprintf("I have no perms to ban: %s", err.Error()))
				return false
			}
			return true
		}
		current := time.Now()
		entryTime, err := discordgo.SnowflakeTimestamp(entry.ID)
		if err != nil {
			fmt.Println(fmt.Sprintf("how the fuck did this happen: %s", err.Error()))
			return false
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
	return false
}

func findAudit(s *discordgo.Session, guildID, targetID string, auditType int) *discordgo.AuditLogEntry {
	audits, err := s.GuildAuditLog(guildID, "", "", auditType, 10) // we really don't need 25 here so we'll use 10 instead (I could probably just use one but whatever)
	if err != nil {
		panic("I can't read audits (wtf am I supposed to do???)")
	}
	for _, entry := range audits.AuditLogEntries {
		if entry.TargetID == targetID {
			return entry
		}
	}
	return nil
}
