package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func (bot *Bot) weatherForecastWatcher() {
	for {
		if time.Now().Hour() == 7 || time.Now().Hour() == 19 {
			weather := bot.getWeatherForecast()
			answer, _ := bot.gptClient.SendMessage("Притворись что ты ведущий шоу прогоза погоды который любит жесткий троллинг. Вот сам прогноз погоды:\n" + weather)
			bot.telebot.SendMessage(bot.config.TelegramGroupId, answer)
			time.Sleep(time.Hour * 2)
		}
		time.Sleep(time.Second * 10)
	}

}

func (bot *Bot) getWeatherForecast() string {
	jsonStruct := struct {
		Ok      bool   `json:"status"`
		Message string `json:"message"`
	}{}

	jsonFile, err := os.Open(bot.config.WeatherPath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &jsonStruct)
	return jsonStruct.Message
}
