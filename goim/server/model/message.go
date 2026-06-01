package model

import "time"

type MessageType int

const (
	MessageText MessageType = iota
	MessageImage
	MessageSystem
)

type Message struct {
	ID           string      `json:"id"`
	SenderID     string      `json:"sender_id"`
	ReceiverID   string      `json:"receiver_id"`
	ReceiverType int         `json:"receiver_type"`
	Content      string      `json:"content"`
	Type         MessageType `json:"type"`
	Status       int         `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
}

type WSMessage struct {
	ID        string          `json:"id,omitempty"`
	Type      string          `json:"type"`
	From      string          `json:"from"`
	To        string          `json:"to"`
	ToType    int             `json:"to_type"`
	Content   string          `json:"content"`
	MsgType   MessageType     `json:"msg_type"`
	Timestamp int64           `json:"timestamp"`
}

type MessageRead struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	MessageID  string    `json:"message_id"`
	ReadAt     time.Time `json:"read_at"`
}