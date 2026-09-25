package models

import (
	"fmt"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type TGLastRequest struct {
	UserRequest map[int]string
}
type TelegramDeleteList struct {
	List map[string]MessageDeleteList
}
type MessageDeleteList struct {
	Messages  []*tb.Message
	ExpiresAt *time.Time
}

type TGNotification struct {
	Message    interface{}
	Keyboard   *tb.ReplyMarkup
	TelegramID uint64
}

type StoredMessage struct {
	MessageID string `sql:"message_id" json:"message_id"`
	ChatID    int64  `sql:"chat_id" json:"chat_id"`
}

func (x StoredMessage) MessageSig() (string, int64) {
	return x.MessageID, x.ChatID
}

func (sm *StoredMessage) DeleteMessage(b *tb.Bot, msg *tb.Message) {
	sm = &StoredMessage{
		MessageID: fmt.Sprintf("%v", msg.ID),
		ChatID:    msg.Chat.ID,
	}
	e := b.Delete(sm)
	if e != nil {
		log.Println("[DeleteMessage] error deleting message:", e)
	}
}
