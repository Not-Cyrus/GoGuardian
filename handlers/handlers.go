package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/Not-Cyrus/GoGuardian/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/valyala/fastjson"
)

func readAudits(s *discordgo.Session, guildID string, auditType int) string {
	parsedData, configData := utils.FindConfig(guildID)

	auditMap := make(map[string]string)
	userMap := make(map[string]int)
	audits, _ := utils.SendRequest("GET", fmt.Sprintf("https://discord.com/api/v8/guilds/%s/audit-logs?action_type=%d&limit=25", guildID, auditType), "", nil)
	err := fastjson.Validate(audits)
	if err != nil {
		utils.SendMessage(s, fmt.Sprintf("I can't read audits: %s | if you think this is a mistake make an issue at https://github.com/Not-Cyrus/GoGuardian/issues", err.Error()), utils.GetGuildOwner(s, guildID)) // scrap this last comment as it's now going to be a public bot.
		return ""
	}
	parsed, _ := parser.Parse(audits)
	auditLogEntries := parsed.GetArray("audit_log_entries")
	for _, entry := range auditLogEntries {

		var (
			auditID = string(entry.GetStringBytes("id"))
			reason  = string(entry.GetStringBytes("reason"))
			userID  = string(entry.GetStringBytes("user_id"))
		)

		if userMap[userID] >= configData.GetInt("Config", "Threshold") {
			if userID == DGUser.ID && configData.GetBool("Config", "AntiHijackEnabled") {
				if strings.Contains(reason, "https://github.com/Not-Cyrus/GoGuardian") {
					return "" // lazy as fuck, but it'll do the trick to stop false flags. I'll probably make this better soon idk
				}
				utils.SendMessage(s, "The bot has been comprimised. I have left the guild for your safety.", utils.GetGuildOwner(s, guildID)) // this is an important message so we'll DM the owner.
				s.GuildLeave(guildID)
			}
			err := utils.BanCreate(guildID, userID, "You just got destroyed by https://github.com/Not-Cyrus/GoGuardian")
			if len(err) != 0 {
				utils.SendMessage(s, fmt.Sprintf("I have no perms to ban <@!%s>: %s", userID, err), utils.GetGuildOwner(s, guildID)) // this is an important message so we'll DM the owner.
				return ""
			}
			return userID
		}

		current := time.Now()
		entryTime, err := discordgo.SnowflakeTimestamp(auditID)

		if err != nil {
			utils.SendMessage(s, fmt.Sprintf("how the fuck did this happen: %s", err.Error()), "")
			return ""
		}

		if current.Sub(entryTime).Round(1*time.Second).Seconds() <= configData.GetFloat64("Config", "Seconds") {
			if _, ok := auditMap[auditID]; !ok {
				inArray, _ := utils.InArray(guildID, "WhitelistedIDs", parsedData, userID)
				if !inArray {
					auditMap[userID] = auditID
					userMap[userID]++
				}
			}
		}
	}
	return ""
}

func findAudit(s *discordgo.Session, guildID, targetID string, auditType int) *fastjson.Value {
	audits, _ := utils.SendRequest("GET", fmt.Sprintf("https://discord.com/api/v8/guilds/%s/audit-logs?action_type=%d&limit=10", guildID, auditType), "", nil)
	err := fastjson.Validate(audits)
	if err != nil {
		utils.SendMessage(s, fmt.Sprintf("I can't read audits: %s | if you think this is a mistake make an issue at https://github.com/Not-Cyrus/GoGuardian/issues", err.Error()), utils.GetGuildOwner(s, guildID)) // scrap this last comment as it's now going to be a public bot.
		return nil
	}
	parsed, _ := parser.Parse(audits)
	auditLogEntries := parsed.GetArray("audit_log_entries")
	for _, entry := range auditLogEntries {
		var (
			auditTargetID = string(entry.GetStringBytes("target_id"))
		)
		if auditTargetID == targetID {
			return entry
		}
	}
	return nil
}

var (
	parser fastjson.Parser
)
