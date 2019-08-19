package models

type SendMessage struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	MessageID int32  `json:"reply_to_message_id,omitempty"`
}
