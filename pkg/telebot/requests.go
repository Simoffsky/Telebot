package telebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// client for making api bot requests
type TelegramClient struct {
	apiUrl string
	offset int64 // id of the last update
}

func NewTelegramClient(apiKey string) *TelegramClient {
	return &TelegramClient{
		apiUrl: "https://api.telegram.org/bot" + apiKey,
		offset: 0,
	}
}

func (tc *TelegramClient) GetUpdates() ([]Update, error) {
	type UpdateResponse struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	resp, err := http.Get(tc.apiUrl + "/getUpdates?timeout=60&offset=" + fmt.Sprint(tc.offset))
	if err != nil {
		return nil, err
	}

	req := &UpdateResponse{}
	err = json.NewDecoder(resp.Body).Decode(req)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if !req.Ok {
		return nil, err
	}
	if len(req.Result) != 0 {
		tc.offset = req.Result[len(req.Result)-1].UpdateId + 1
	}

	return req.Result, nil
}

func (tc *TelegramClient) SendMessage(chatId int64, text string) Message {
	Payload := struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}{
		chatId,
		text,
	}

	jsonPayload, err := json.Marshal(Payload)
	checkDie(err)

	resp := tc.makePostRequest(jsonPayload, "/sendMessage")
	message := tc.parseMessageResponse(resp)

	return message
}

func (tc *TelegramClient) EditMessage(message Message, editedText string) Message {
	Payload := struct {
		ChatId    int64  `json:"chat_id"`
		MessageId int64  `json:"message_id"`
		Text      string `json:"text"`
	}{
		message.Chat.Id,
		message.MessageId,
		editedText,
	}

	jsonPayload, err := json.Marshal(Payload)
	checkDie(err)
	resp := tc.makePostRequest(jsonPayload, "/editMessageText")
	editedMessage := tc.parseMessageResponse(resp)

	return editedMessage
}

func (tc *TelegramClient) makePostRequest(jsonPayload []byte, endpoint string) *http.Response {
	req, err := http.NewRequest("POST", tc.apiUrl+endpoint, bytes.NewBuffer(jsonPayload))
	checkDie(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	checkDie(err)
	return resp
}

func (tc *TelegramClient) parseMessageResponse(resp *http.Response) Message {
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
	if !message.Ok {
		log.Fatal("Messed up, parseMessageResponse bad response")
	}
	return message.Result
}
