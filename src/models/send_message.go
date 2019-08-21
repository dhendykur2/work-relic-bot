package models

type SendMessage struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
	MessageID int32  `json:"reply_to_message_id,omitempty"`
}
