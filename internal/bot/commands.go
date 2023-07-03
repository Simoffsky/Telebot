package bot

import (
	"fmt"
	"log"
	"os"
	"time"
	"vladOS/internal/gpt"
	"vladOS/pkg/telebot"
)

func (bot *Bot) RegisterCommands() {
	bot.telebot.HandleMessage(commandStart(bot.telebot), "start")
	bot.telebot.HandleMessage(commandGPT(bot.telebot, bot.gptClient), "gpt")

	go bot.weatherForecastWatcher()
}

func commandStart(bot *telebot.Telebot) func(telebot.Message) {
	return func(message telebot.Message) {
		bot.SendMessage(message.Chat.Id, "Привет, я бот VladOS")
	}
}

func commandGPT(bot *telebot.Telebot, gpt *gpt.GptApiClient) func(telebot.Message) {
	return func(message telebot.Message) {
		text := message.From.Username + ": " + message.Text[4:]
		makeLog(text) //TODO: Вынести куда-то отсюда

		botMessage := bot.SendMessage(message.Chat.Id, "...")

		ans, err := gpt.SendGroupMessage(text)
		println(ans)
		if err != nil {
			bot.SendMessage(message.Chat.Id, "Ошибка: "+err.Error())
			return
		}

		botMessage = bot.EditMessage(botMessage, ans)

		if len(botMessage.Text) >= 4 && botMessage.Text[0:4] == "/gpt" {
			time.Sleep(10 * time.Second)
			botMessage.From.Username = "VladOS"
			commandGPT(bot, gpt)(botMessage)
		}
	}
}

func makeLog(text string) {
	file, err := os.OpenFile("logfile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logString := fmt.Sprintf("[%s]%s", time.Now().Format("2006-01-02 15:04:05"), text)

	if _, err := file.WriteString(logString + "\n"); err != nil {
		log.Fatal(err)
	}
}
