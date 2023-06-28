package telebot

import (
	"log"
	"strings"
)

func checkDie(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Telebot struct {
	apiClient       *TelegramClient
	commandHandlers map[string]func(Message) // calls if command is found
	commonHandlers  []func(Message)
}

func NewTelebot(apiKey string) *Telebot {
	return &Telebot{
		apiClient:       NewTelegramClient(apiKey),
		commandHandlers: make(map[string]func(Message), 0),
	}
}

// infinite cycle of getting updates
// TODO: skip updates while offline
func (t *Telebot) LongPolling() error {
	println("Long polling is started")
	for {
		updates, err := t.apiClient.GetUpdates()
		if err != nil {
			return err
		}
		t.handleUpdates(updates)
	}
}

func (t *Telebot) handleUpdates(updates []Update) {
	for _, update := range updates {
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
	return t.apiClient.SendMessage(chatId, text)
}

// TODO: implement
func (t *Telebot) ReplyTo(message Message, text string) {
	panic("not implemented")
}

func (t *Telebot) EditMessage(message Message, updatedText string) Message {
	return t.apiClient.EditMessage(message, updatedText)
}

func (t *Telebot) HandleMessage(commandFunc func(Message), command string) {
	if command == "" {
		t.commonHandlers = append(t.commonHandlers, commandFunc)
		return
	}
	t.commandHandlers[command] = commandFunc
}
