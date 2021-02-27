package events

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Not-Cyrus/GoGuardianv2/database"
	"github.com/Not-Cyrus/GoGuardianv2/utils"
	"github.com/bwmarrin/discordgo"
)

func AntiInvite(s *discordgo.Session, m *discordgo.MessageCreate) {
	data, err := database.Database.FindData(m.GuildID)

	switch {
	case err != nil:
		return
	case data["anti-invite"] == "off":
		return
	case utils.HasPerms(s, m.GuildID, m.Author.ID, discordgo.PermissionManageMessages):
		return
	}

	if strings.Contains(m.Content, "discord.gg/") {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func BanHandler(s *discordgo.Session, event *discordgo.GuildBanAdd) {
	utils.ReadAudit(s, event.GuildID, "banned a member", 22)
}

func ChannelCreate(s *discordgo.Session, event *discordgo.ChannelCreate) {
	utils.ReadAudit(s, event.GuildID, "created a channel", 10)
}

func ChannelRemove(s *discordgo.Session, event *discordgo.ChannelDelete) {
	utils.ReadAudit(s, event.GuildID, "deleted a channel", 12)
}

func CreateGuild(s *discordgo.Session, event *discordgo.GuildCreate) {
	database.Database.CreateGuild(utils.BotUser, event.Guild)
	muteX.Lock()
	defer muteX.Unlock()

	guilds[event.Guild.ID] = event.Guild.MemberCount

	MemberCount += guilds[event.Guild.ID]
	GuildCount++
}

func DeleteGuild(s *discordgo.Session, event *discordgo.GuildDelete) {
	database.Database.DeleteGuild(event.Guild.ID)

	MemberCount -= guilds[event.Guild.ID]

	muteX.RLock()
	defer muteX.RUnlock()

	delete(guilds, event.Guild.ID)
	GuildCount--
}

func KickHandler(s *discordgo.Session, event *discordgo.GuildMemberRemove) {
	utils.ReadAudit(s, event.GuildID, "kicked a member", 20)
}

func MemberJoin(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	MemberCount++

	// ANTI-Bot check
	var (
		entry, _, err = utils.FindAudit(s, event.GuildID, 28)
	)

	if entry == nil || err != nil {
		return
	}

	err = s.GuildBanCreateWithReason(event.GuildID, entry.UserID, fmt.Sprintf("%s | invited a bot", utils.BotUser.Username), 0)
	err = s.GuildBanCreateWithReason(event.GuildID, event.User.ID, fmt.Sprintf("%s - invited by a non-whitelisted person (Bot)", utils.BotUser.Username), 0)

	if err != nil {
		return
	}

	utils.LogChannel(s, event.GuildID, fmt.Sprintf("<@%s> %s", entry.UserID, fmt.Sprintf("<@%s> tried inviting a bot (Both the bot, user have been banned)", entry.UserID)))
}

func MemberLeave(s *discordgo.Session, event *discordgo.GuildMemberRemove) {
	MemberCount--
}

func MemberRoleUpdate(s *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	var (
		entry, change, err = utils.FindAudit(s, event.GuildID, 25)
	)

	if err != nil || change == nil || len(change.([]interface{})) == 0 {
		return
	}

	roleID := change.([]interface{})[0].(map[string]interface{})["id"].(string)

	guildRole, err := s.State.Role(event.GuildID, roleID)
	if err != nil {
		return
	}

	if guildRole.Permissions&0x8 != 0x8 {
		return
	}

	err = s.GuildMemberRoleRemove(event.GuildID, entry.TargetID, roleID)
	err = s.GuildBanCreateWithReason(event.GuildID, entry.UserID, fmt.Sprintf("%s | gave a member an admin role", utils.BotUser.Username), 0)

	if err != nil {
		return
	}

	utils.LogChannel(s, event.GuildID, fmt.Sprintf("<@%s> gave a member an admin role", entry.UserID))
}

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStreamingStatus(2, fmt.Sprintf(">help | Shard %d", s.ShardID), "https://twitch.tv/discord")
	fmt.Printf("Connected to shard %d\n", s.ShardID)
}
func RoleCreate(s *discordgo.Session, event *discordgo.GuildRoleCreate) {
	utils.ReadAudit(s, event.GuildID, "created a role", 30)
}

func RoleRemove(s *discordgo.Session, event *discordgo.GuildRoleDelete) {
	utils.ReadAudit(s, event.GuildID, "deleted a role", 32)
}

func WebhookCreate(s *discordgo.Session, event *discordgo.WebhooksUpdate) {
	var (
		err          error
		selfMember   *discordgo.Member
		targetMember *discordgo.Member
	)

	webhooks, err := s.ChannelWebhooks(event.ChannelID)
	if err != nil {
		return
	}

	for _, webhook := range webhooks {
		whitelisted := database.Database.IsWhitelisted(event.GuildID, webhook.User.ID)
		if whitelisted {
			return
		}

		err = s.WebhookDelete(webhook.ID)

		selfMember, err = s.GuildMember(event.GuildID, utils.BotUser.ID)
		if err != nil {
			return
		}

		targetMember, err = s.GuildMember(event.GuildID, webhook.User.ID)
		if err != nil {
			return
		}

		targetHighest := utils.HighestRole(s, event.GuildID, targetMember)
		selfHighest := utils.HighestRole(s, event.GuildID, selfMember)

		if !utils.IsAbove(selfHighest, targetHighest) || !utils.HasPerms(s, event.GuildID, webhook.User.ID, discordgo.PermissionBanMembers) {
			return
		}

		err = s.GuildBanCreateWithReason(event.GuildID, webhook.User.ID, fmt.Sprintf("%s | created a webhook", utils.BotUser.Username), 0)
		if err != nil {
			utils.LogChannel(s, event.GuildID, fmt.Sprintf("<@%s> created a webhook | %s couldn't take moderation action: %s", webhook.User.ID, utils.BotUser.Username, err.Error()))
			return
		}

		utils.LogChannel(s, event.GuildID, fmt.Sprintf("<@%s> created a webhook", webhook.User.ID))
	}
}

var (
	guilds      = make(map[string]int)
	GuildCount  int
	MemberCount int
	muteX       = &sync.RWMutex{}
)
