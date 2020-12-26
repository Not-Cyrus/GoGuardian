package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Not-Cyrus/GoGuardian/api"
)

func main() {
	bot := api.Bot{}
	bot.Setup()
	err := bot.Run()
	if err != nil {
		panic("Could not start a discord session")
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bot.Stop()
}

// Note to self, I should really add comments to my code more often
