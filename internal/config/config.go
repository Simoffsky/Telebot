package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	BotToken          string `json:"telegram_bot_token"`
	TelegramGroupId   int64  `json:"telegram_group_id"`
	GPTPrompt         string `json:"-"`
	GPTStoredMessages int    `json:"gpt_stored_messages"`
	GPTToken          string `json:"gpt_token"`

	WeatherPath string `json:"weather_json_path"`
}

func NewConfig(configPath, gptPromptPath string) (*Config, error) {
	config := &Config{}
	err := config.LoadConfig(configPath, gptPromptPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (config *Config) LoadConfig(configPath, promptPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return errors.New("could not open config file: " + err.Error())
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return errors.New("could not parse config file: " + err.Error())
	}

	data, err = os.ReadFile(promptPath)
	if err != nil {
		return errors.New("could not open gpt prompt file: " + err.Error())
	}
	config.GPTPrompt = string(data)

	if config.BotToken == "" {
		return errors.New("error while setup config: bot token is empty")
	}
	return nil

}
