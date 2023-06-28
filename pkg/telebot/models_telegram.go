package telebot

// https://core.telegram.org/bots/api#update

type Update struct {
	UpdateId int64   `json:"update_id"`
	Message  Message `json:"message,omitempty"` // text, photo, sticker, etc
}

type Message struct {
	MessageId int64           `json:"message_id"`
	Chat      Chat            `json:"chat"`
	Text      string          `json:"text,omitempty"`
	From      User            `json:"from,omitempty"` // if empty then sender was channel
	Date      int             `json:"date"`
	Sticker   Sticker         `json:"sticker,omitempty"`
	Entities  []MessageEntity `json:"entities,omitempty"`
}


type MessageEntity struct {
	Type string `json:"type"`
	User User   `json:"user,omitempty"` // for type "text_mention"
}

type Sticker struct {
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Emoji   string `json:"emoji"`
	SetName string `json:"set_name"` // name of set that sticker belongs to
}

type Chat struct {
	Id       int64  `json:"id"`
	TypeChat string `json:"type"`
}

type User struct {
	Id       int64  `json:"id"`
	IsBot    bool   `json:"is_bot"`
	Username string `json:"username"`
}
