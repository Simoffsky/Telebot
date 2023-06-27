package main

import (
	"flag"
	"log"
	"vladOS/internal/bot"
	"vladOS/internal/config"
)

var (
	configPath    string
	gptPromptPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.json", "path to the config file")
	flag.StringVar(&gptPromptPath, "gpt-path", "configs/gpt-api-prompt.txt", "path to gpt-api init prompt")
}

func main() {
	config := parseConfig()

	bot := bot.NewBot(config)
	bot.Start()
}

func parseConfig() *config.Config {
	flag.Parse()
	config, err := config.NewConfig(configPath, gptPromptPath)

	if err != nil {
		log.Fatal(err)
	}

	return config
}
