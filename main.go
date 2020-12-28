package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Not-Cyrus/GoGuardian/api"
)

func main() {
	bot := api.Bot{}
	bot.Setup()
	err := bot.Run()
	if err != nil {
		fmt.Printf("Could not start a discord session: %s", err.Error())
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bot.Stop()
}

// Note to self, I should really add comments to my code more often
