package main

import (
	"fmt"
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
		panic(fmt.Sprintf("Could not start a discord session: %s", err.Error()))
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bot.Stop()
}

// Note to self, I should really add comments to my code more often
