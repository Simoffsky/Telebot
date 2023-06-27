package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const (
	URL = "https://api.openai.com/v1"
)

type GptApiClient struct {
	apiKey          string
	messages        []Message
	messagesToStore int
}

func NewGptApiClient(gptKey string, gptInitPrompt string, countStoredMessages int) *GptApiClient {
	messages := []Message{
		{
			Role:    "system",
			Content: gptInitPrompt,
		},
	}

	return &GptApiClient{
		apiKey:          gptKey,
		messages:        messages,
		messagesToStore: countStoredMessages,
	}
}

// returns answer from GPT-3.5-turbo
func (gpt *GptApiClient) SendMessage(message string) (string, error) {
	gpt.addMessage(message)

	requestData := NewRequest(gpt.messages)

	jsonPayload, err := json.Marshal(requestData)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := gpt.makeRequest(jsonPayload)
	if err != nil {
		return "", err
	}

	answer := resp.Choices[0].Message.Content
	gpt.addMessage("VladOS: " + answer)
	return answer, nil

}

func (client *GptApiClient) addMessage(messageString string) {
	message := Message{
		Role:    "user",
		Content: messageString,
	}
	if len(client.messages) > 10 {
		client.messages = append(client.messages[0:1], client.messages[2:]...) // delete second message
	}
	message.Content = "\n\n" + message.Content //forgot why I did this
	client.messages = append(client.messages, message)
}

func (gpt *GptApiClient) makeRequest(jsonRequest []byte) (*Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", URL+"/chat/completions", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Authorization", "Bearer "+gpt.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respData := &Response{}

	err = json.NewDecoder(resp.Body).Decode(respData)
	if err != nil {
		log.Fatal(err)
	}

	if respData.Error.Message != "" {
		return nil, errors.New(respData.Error.Message)
	}
	return respData, nil
}
