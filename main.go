package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type (
	configData struct {
		Token        string   `json:"token"`
		WhitelistIDs []string `json:"Whitelisted,omitempty"`
		BanCount     int      `json:"BanCount"`
		BanSeconds   float64  `json:"BanSeconds"`
	}
)

var (
	auditLog       = make(map[string]int)
	bans           = make(map[string]string)
	config         = configData{}
	whitelistedIDs = map[string]string{}
)

func main() {
	dg, _ := discordgo.New(config.Token)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildBans | discordgo.IntentsGuildMembers)
	dg.AddHandler(banCreate)
	dg.AddHandler(memberAdd)
	dg.Open()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func init() {
	file, err := os.Open("Config.json")
	if err != nil {
		panic("HELLO, DO YOU KNOW HOW TO MOVE FILES??")
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("I couldn't read the data")
	}
	json.Unmarshal([]byte(data), &config)
	// pretty sure indexing a key value is faster than searching through an array resulting in this code
	for _, v := range config.WhitelistIDs {
		whitelistedIDs[v] = "https://github.com/Not-Cyrus is pretty cool"
	}
}

func memberAdd(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
	if !member.User.Bot {
		return
	}
	memberName := strings.ToLower(member.User.Username)
	// you see I could clean up this code (I made it before making the JSON part), but I am lazy.
	if strings.Contains(memberName, "dyno") || strings.Contains(memberName, "mee6") {
		if member.User.ID != "155149108183695360" && member.User.ID != "159985870458322944" {
			err := s.GuildBanCreateWithReason(member.GuildID, member.User.ID, "Fake Bot", 0)
			if err != nil {
				fmt.Println(fmt.Sprintf("I have no perms to ban: %s#%s", member.User.Username, member.User.Discriminator))
				return
			}
		}
	}
}

func banCreate(s *discordgo.Session, ban *discordgo.GuildBanAdd) {
	audits, err := s.GuildAuditLog(ban.GuildID, "", "", int(discordgo.AuditLogActionMemberBanAdd), 25)
	if err != nil {
		panic("I have no perms to view audit logs")
	}
	for _, entry := range audits.AuditLogEntries {
		if auditLog[entry.UserID] >= config.BanCount {
			err := s.GuildBanCreateWithReason(ban.GuildID, entry.UserID, "You suck at nuking!!!", 0)
			if err != nil {
				fmt.Println(fmt.Sprintf("I have no perms to ban: %s", err.Error()))
				return
			}
			return
		}
		current := time.Now()
		entryTime, err := discordgo.SnowflakeTimestamp(entry.ID)
		if err != nil {
			fmt.Println(fmt.Sprintf("how the fuck did this happen: %s", err.Error()))
			return
		}
		if current.Sub(entryTime).Round(1*time.Second).Seconds() <= config.BanSeconds {
			if _, ok := bans[entry.ID]; !ok {
				if _, whitelisted := whitelistedIDs[entry.UserID]; whitelisted {
					return
				}
				auditLog[entry.UserID]++
				bans[entry.ID] = "lol this was useless oops"
			}
		}
	}
}
