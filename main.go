package main

import (
	"os"

	"github.com/newrushbolt/AppleSeal/logger"
	"github.com/newrushbolt/AppleSeal/messages"
	"gopkg.in/telegram-bot-api.v4"
)

var ()

func main() {
	logger.Init("")

	tgToken := os.Getenv("SEAL_TG_TOKEN")

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	logger.Logger.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		messages.ParseMessage(bot, update.Message)
	}
}
