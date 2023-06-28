package bot

import (
	"log"
	"vladOS/internal/config"
	"vladOS/internal/gpt"
	"vladOS/internal/opendota"
	"vladOS/pkg/telebot"
)

type Bot struct {
	config    *config.Config
	telebot   *telebot.Telebot
	gptClient *gpt.GptApiClient
}

func NewBot(config *config.Config) *Bot {
	return &Bot{
		config:  config,
		telebot: telebot.NewTelebot(config.BotToken),
		gptClient: gpt.NewGptApiClient(
			config.GPTToken,
			config.GPTPrompt,
			config.GPTStoredMessages,
		),
	}
}

func (bot *Bot) Start() {
	bot.RegisterCommands()
	go opendota.CheckPlayedMatch(bot.telebot, bot.config.TelegramGroupId)
	err := bot.telebot.LongPolling()
	if err != nil {
		log.Fatal(err) // FIXME:
	}
}
