package api

import (
	"sync"

	"github.com/Not-Cyrus/GoGuardian/commands"
	"github.com/bwmarrin/discordgo"
)

func init() {

	CommandRoute.Add("antiinvite", CommandRoute.AntiInvite, &commands.Config{
		Alias:        []string{"antiinv", "noinvite"},
		Cooldown:     3,
		OwnerOnly:    true,
		RequiresArgs: true,
	})

	CommandRoute.Add("avatar", CommandRoute.Avatar, &commands.Config{
		Alias:           []string{"av", "pfp", "icon"},
		Cooldown:        1,
		RequiresMention: true,
	})

	CommandRoute.Add("ban", CommandRoute.Ban, &commands.Config{
		Cooldown:        2,
		RequiresMention: true,
		Perms:           discordgo.PermissionBanMembers,
	})

	CommandRoute.Add("banner", CommandRoute.ServerBanner, &commands.Config{
		Alias:    []string{"serverbanner", "sbanner"},
		Cooldown: 1,
	})

	CommandRoute.Add("botinfo", CommandRoute.BotInfo, &commands.Config{
		Cooldown: 4,
	})

	CommandRoute.Add("credits", CommandRoute.Credits, &commands.Config{
		Cooldown: 1,
	})

	CommandRoute.Add("fox", CommandRoute.Fox, &commands.Config{
		Cooldown: 3,
	})

	CommandRoute.Add("help", CommandRoute.Help, &commands.Config{
		Cooldown: 1,
	})

	CommandRoute.Add("invite", CommandRoute.Invite, &commands.Config{
		Cooldown: 1,
	})

	CommandRoute.Add("kick", CommandRoute.Kick, &commands.Config{
		Cooldown:        2,
		RequiresMention: true,
		Perms:           discordgo.PermissionBanMembers,
	})

	CommandRoute.Add("lockdown", CommandRoute.Lockdown, &commands.Config{
		Alias:    []string{"lock"},
		Cooldown: 1,
		Perms:    discordgo.PermissionManageChannels,
	})

	CommandRoute.Add("logchannel", CommandRoute.LoggingChannel, &commands.Config{
		Alias:     []string{"setlogs", "log"},
		Cooldown:  5,
		OwnerOnly: true,
	})

	CommandRoute.Add("massunban", CommandRoute.Unban, &commands.Config{
		Alias:    []string{"unbanall"},
		Cooldown: 30,
		Perms:    discordgo.PermissionBanMembers,
	})

	CommandRoute.Add("membercount", CommandRoute.MemberCount, &commands.Config{
		Alias:    []string{"mc", "members"},
		Cooldown: 1,
	})

	CommandRoute.Add("nuke", CommandRoute.Nuke, &commands.Config{
		Alias:    []string{"nk"},
		Cooldown: 30,
		Perms:    discordgo.PermissionManageChannels,
	})

	CommandRoute.Add("ping", CommandRoute.Ping, &commands.Config{
		Alias:    []string{"b"},
		Cooldown: 5,
	})

	CommandRoute.Add("prefix", CommandRoute.Prefix, &commands.Config{
		Alias:        []string{"setprefix"},
		Cooldown:     3,
		OwnerOnly:    true,
		RequiresArgs: true,
	})

	CommandRoute.Add("servericon", CommandRoute.ServerIcon, &commands.Config{
		Alias:    []string{"serverpfp", "sicon", "serverpic"},
		Cooldown: 1,
	})

	CommandRoute.Add("serverinfo", CommandRoute.ServerInfo, &commands.Config{
		Cooldown: 1,
	})

	CommandRoute.Add("settings", CommandRoute.Settings, &commands.Config{
		Cooldown: 1,
	})

	CommandRoute.Add("setup", CommandRoute.Setup, &commands.Config{
		Cooldown: 1,
	})

	CommandRoute.Add("slowmode", CommandRoute.SlowMode, &commands.Config{
		Cooldown:     1,
		RequiresArgs: true,
		Perms:        discordgo.PermissionManageChannels,
	})

	CommandRoute.Add("unlockdown", CommandRoute.UnLockdown, &commands.Config{
		Alias:    []string{"unlock"},
		Cooldown: 1,
		Perms:    discordgo.PermissionManageChannels,
	})

	CommandRoute.Add("unslowmode", CommandRoute.UnSlowMode, &commands.Config{
		Cooldown: 1,
		Perms:    discordgo.PermissionManageChannels,
	})

	CommandRoute.Add("unwhitelist", CommandRoute.Unwhitelist, &commands.Config{
		Alias:           []string{"delwl", "removewl", "dewhitelist"},
		Cooldown:        3,
		OwnerOnly:       true,
		RequiresMention: true,
	})

	CommandRoute.Add("userinfo", CommandRoute.UserInfo, &commands.Config{
		Alias:           []string{"whois"},
		Cooldown:        1,
		RequiresMention: true,
	})

	CommandRoute.Add("whitelist", CommandRoute.Whitelist, &commands.Config{
		Alias:           []string{"wl", "addwhitelist", "bypass"},
		Cooldown:        3,
		OwnerOnly:       true,
		RequiresMention: true,
	})

	CommandRoute.Add("whitelisted", CommandRoute.ViewWhitelisted, &commands.Config{
		Cooldown:        3,
		WhitelistedOnly: true,
	})

}

var (
	CommandRoute = &commands.Commands{
		Cooldown: &commands.CommandCooldown{
			Cooldowns: make(map[string][]string),
			Mutex:     &sync.RWMutex{},
		},
	}
)
