package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"telebot/models"
)

// https://api.telegram.org/bot5724856011:AAFXQzE-JSYIb_OH6Hm0lxC8wAE3Ohyun3Y/getUpdates
const (
	API_URL   = "https://api.telegram.org/"
	BOT_TOKEN = "5724856011:AAFXQzE-JSYIb_OH6Hm0lxC8wAE3Ohyun3Y"
	BOT_URL   = API_URL + "bot" + BOT_TOKEN
	KAKICH_ID = 688410426
)

func checkDie(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func deleteMessage(chat_id int64, message_id int64) {
	Payload := struct {
		ChatId    int64 `json:"chat_id"`
		MessageId int64 `json:"message_id"`
	}{
		-1001338826524,
		107053,
	}
	jsonPayload, err := json.Marshal(Payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", BOT_URL+"/deleteMessage", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	print(string(respBody))
}

func SendMessage(chat_id int64, message string) {
	Payload := struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}{
		chat_id,
		message,
	}

	jsonPayload, err := json.Marshal(Payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", BOT_URL+"/sendMessage", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	print(string(respBody))
}

func pong(chat models.Chat) {
	SendMessage(chat.Id, "pong")
}

func HandleCommand(message models.Message) {
	if message.Text == "/ping" {
		pong(message.Chat)
	}
}
func HandleUpdates(updates []models.Update) int64 {
	var offset int64

	for _, update := range updates {
		offset = update.UpdateId

		if len(update.Message.Entities) != 0 {
			for _, entity := range update.Message.Entities {
				if entity.Type == "bot_command" {
					HandleCommand(update.Message)
					continue
				}
			}
		}

	}
	return offset + 1
}

func LongPolling() {
	type UpdateResponse struct {
		Ok     bool            `json:"ok"`
		Result []models.Update `json:"result"`
	}
	var offset int64 = 0

	for {
		println("Some loong polling")
		resp, err := http.Get(BOT_URL + "/getUpdates?timeout=5&offset=" + fmt.Sprint(offset))
		checkDie(err)

		req := &UpdateResponse{}
		err = json.NewDecoder(resp.Body).Decode(req)
		resp.Body.Close()
		checkDie(err)

		if !req.Ok {
			log.Fatal("bad request")
		}

		offset = HandleUpdates(req.Result)
	}
}
func main() {
	LongPolling()
	//deleteMessage(1, 1)
}
