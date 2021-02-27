package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Not-Cyrus/GoGuardian/api"
	"github.com/bwmarrin/discordgo"
)

func main() {

	fmt.Print("Enter your token: ")
	fmt.Scan(&token)

	fmt.Print("Enter the amount of shards you want: ")
	fmt.Scan(&dshard)

	req, _ := http.NewRequest("GET", "https://discord.com/api/v8/gateway/bot", nil) // I would use fasthttp but does speed really matter that much here?
	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", token))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("[Sharding Error]: %s\n", err.Error())
		return
	}
	defer res.Body.Close()

	gresponse := &discordgo.GatewayBotResponse{}

	json.NewDecoder(res.Body).Decode(&gresponse)
	if err != nil {
		fmt.Printf("[Decode Error]: %s\n", err.Error())
		return
	}

	var shardCount = gresponse.Shards

	if shardCount < 2 {
		shardCount = dshard
	}

	bot.Sessions = make([]*discordgo.Session, shardCount)

	for s := 0; s < shardCount; s++ {
		bot.Shard(token, shardCount, s)
		bot.Run()
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Stop()
}

var (
	bot    = api.Bot{}
	dshard int
	token  string
)
