package telebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func checkDie(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Telebot struct {
	apiUrl          string
	offset          int64                    //id of the last update
	commandHandlers map[string]func(Message) // calls if command is found
	commonHandlers  []func(Message)
}

func NewTelebot(apiKey string) *Telebot {
	return &Telebot{
		apiUrl:          "https://api.telegram.org/bot" + apiKey,
		offset:          0,
		commandHandlers: make(map[string]func(Message), 0),
	}
}

// infinite cycle of getting updates
func (t *Telebot) LongPolling() {
	println("Long polling is started")
	for {
		updates := t.getUpdates()
		t.handleUpdates(updates)
	}
}

func (t *Telebot) getUpdates() []Update {
	type UpdateResponse struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	resp, err := http.Get(t.apiUrl + "/getUpdates?timeout=60&offset=" + fmt.Sprint(t.offset))
	checkDie(err)

	req := &UpdateResponse{}
	err = json.NewDecoder(resp.Body).Decode(req)
	resp.Body.Close()
	checkDie(err)

	if !req.Ok {
		log.Fatal("error while getting updates from telegram server")
	}
	return req.Result
}

func (t *Telebot) handleUpdates(updates []Update) {
	for _, update := range updates {
		t.offset = update.UpdateId + 1
		t.handleUpdate(update)
	}
}

func (t *Telebot) handleUpdate(update Update) {
	// t.handleMiddlewares(update) TODO: add middlewares and pass them to MessageHandler
	for _, entity := range update.Message.Entities {
		if entity.Type == "bot_command" {
			t.handleCommand(update.Message)
			return
		}
	}

	for _, commonHandler := range t.commonHandlers {
		commonHandler(update.Message)
	}

}

// can handle message that starts with '/'
func (t *Telebot) handleCommand(message Message) {
	command := strings.Split(message.Text, " ")[0][1:]
	_, ok := t.commandHandlers[command]
	if !ok {
		log.Println("unknown command: " + command)
		return
	}

	t.commandHandlers[command](message)
}

// send string message to chat with chatId
// TODO: change chatId to chat and text to message
// Maybe not? I don't know
// Returns telegamm sended message
func (t *Telebot) SendMessage(chatId int64, text string) Message {
	Payload := struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}{
		chatId,
		text,
	}

	jsonPayload, err := json.Marshal(Payload)
	checkDie(err)

	req, err := http.NewRequest("POST", t.apiUrl+"/sendMessage", bytes.NewBuffer(jsonPayload))
	checkDie(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	checkDie(err)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error() + "; " + string(respBody))
	}
	message := struct {
		Ok     bool    `json:"ok"`
		Result Message `json:"result"`
	}{}
	err = json.Unmarshal(respBody, &message)
	checkDie(err)
	return message.Result
}

// TODO: implement
func (t *Telebot) ReplyTo(message Message, text string) {
	panic("not implemented")
}

func (t *Telebot) HandleMessage(commandFunc func(Message), command string) {
	if command == "" {
		t.commonHandlers = append(t.commonHandlers, commandFunc)
		return
	}
	t.commandHandlers[command] = commandFunc
}
