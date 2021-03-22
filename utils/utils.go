package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Not-Cyrus/GoGuardian/database"
	"github.com/bwmarrin/discordgo"
)

func FindAudit(s *discordgo.Session, guildID string, auditType int) (*discordgo.AuditLogEntry, interface{}, error) {
	if !HasPerms(s, guildID, s.State.User.ID, discordgo.PermissionViewAuditLogs) {
		return nil, nil, fmt.Errorf("No perms in %s", guildID) // useless as we don't have error handling but eh
	}

	audit, err := s.GuildAuditLog(guildID, "", "", auditType, 25)

	if err != nil || len(audit.AuditLogEntries) == 0 {
		return nil, nil, err
	}

	auditLog := audit.AuditLogEntries[0]

	if whitelisted := database.Database.IsWhitelisted(guildID, auditLog.UserID); whitelisted {
		return nil, nil, err
	}

	current := time.Now()
	entryTime, err := discordgo.SnowflakeTimestamp(auditLog.ID)
	if err != nil {
		return nil, nil, err
	}

	if current.Sub(entryTime).Round(1*time.Second).Seconds() > 2 {
		return nil, nil, err
	}

	if len(auditLog.Changes) == 0 {
		return auditLog, []interface{}{}, nil
	}

	return auditLog, auditLog.Changes[0].NewValue, nil
}

func FindInSlice(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func GetGuildOwner(s *discordgo.Session, guildID string) string {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return ""
	}
	return guild.OwnerID
}

func HasPerms(s *discordgo.Session, guildID, userID string, permissions ...int) bool {
	if GetGuildOwner(s, guildID) == userID {
		return true
	}

	guild, err := s.State.Guild(guildID)
	if err != nil || len(guild.Channels) == 0 {
		return false
	}

	perms, err := s.UserChannelPermissions(userID, guild.Channels[0].ID)
	if err != nil {
		return false
	}

	for _, perm := range permissions {
		if perms&perm != perm {
			return false
		}
	}

	return true
}

func HighestRole(s *discordgo.Session, guildID string, member *discordgo.Member) *discordgo.Role {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil
	}

	var highest *discordgo.Role
	for _, roleID := range member.Roles {
		for _, role := range guild.Roles {
			if roleID != role.ID {
				continue
			}
			if highest == nil || IsAbove(role, highest) {
				highest = role
			}
			break
		}
	}
	if highest == nil {
		defaultRole, _ := s.State.Role(guildID, guildID)
		return defaultRole
	}
	return highest
}

func IsAbove(r, r2 *discordgo.Role) bool {

	switch {
	case r.Position != r2.Position:
		return r.Position > r2.Position
	case r.ID == r2.ID:
		return false
	}

	return r.Position < r2.Position
}

func LogChannel(s *discordgo.Session, guildID, postData string) {
	data, err := database.Database.FindData(guildID)
	if err != nil {
		return
	}

	if data["log-channel"] == "nil" {
		return
	}

	s.ChannelMessageSend(data["log-channel"].(string), postData)
}

func MakeRequest(url string) (resBody []byte, err error) {
	// only really for simple GET requests
	var res *http.Response

	res, err = http.Get(url)
	if err != nil {
		return
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	defer res.Body.Close()

	return
}

func ReadAudit(s *discordgo.Session, guildID, reason string, auditType int) {
	if !HasPerms(s, guildID, s.State.User.ID, discordgo.PermissionViewAuditLogs) {
		return
	}

	var (
		audits, err  = s.GuildAuditLog(guildID, "", "", auditType, 25)
		auditEntry   *discordgo.AuditLogEntry
		selfMember   *discordgo.Member
		targetMember *discordgo.Member
	)

	if err != nil || len(audits.AuditLogEntries) == 0 {
		return
	}

	auditEntry = audits.AuditLogEntries[0]

	if whitelisted := database.Database.IsWhitelisted(guildID, auditEntry.UserID); whitelisted {
		return
	}

	current := time.Now()
	entryTime, err := discordgo.SnowflakeTimestamp(auditEntry.ID)
	if err != nil {
		return
	}

	if current.Sub(entryTime).Round(1*time.Second).Seconds() > 2 {
		return
	}

	selfMember, err = s.GuildMember(guildID, s.State.User.ID)
	if err != nil {
		return
	}

	targetMember, err = s.GuildMember(guildID, auditEntry.UserID)
	if err != nil {
		return
	}

	targetHighest := HighestRole(s, guildID, targetMember)
	selfHighest := HighestRole(s, guildID, selfMember)

	if targetHighest == nil || selfHighest == nil {
		return
	}

	if !IsAbove(selfHighest, targetHighest) || !HasPerms(s, guildID, auditEntry.UserID, discordgo.PermissionBanMembers) {
		return
	}

	err = s.GuildBanCreateWithReason(guildID, auditEntry.UserID, fmt.Sprintf("%s | %s", s.State.User.Username, reason), 0)
	if err != nil {
		return
	}

	LogChannel(s, guildID, fmt.Sprintf("<@%s> %s", auditEntry.UserID, reason))
}

func RemoveFromSlice(slice []string, item string) []string {
	returnItems := []string{}
	for _, i := range slice {
		if i == item {
			continue
		}
		returnItems = append(returnItems, i)
	}
	return returnItems
}
