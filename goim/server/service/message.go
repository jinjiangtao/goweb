package service

import (
	"goim/server/model"
	"goim/server/storage"
	"time"

	"github.com/google/uuid"
)

func SaveMessage(msg *model.WSMessage) {
	dbMsg := &model.Message{
		ID:           uuid.New().String(),
		SenderID:     msg.From,
		ReceiverID:   msg.To,
		ReceiverType: msg.ToType,
		Content:      msg.Content,
		Type:         msg.MsgType,
		Status:       0,
		CreatedAt:    time.Now(),
	}

	err := storage.CreateMessage(dbMsg)
	if err != nil {
		return
	}

	msg.ID = dbMsg.ID
	msg.Timestamp = dbMsg.CreatedAt.Unix()
}

func GetMessageHistory(userID, targetID string, targetType int, limit, offset int) ([]*model.Message, error) {
	if targetType == 0 {
		return storage.GetMessagesByUser(userID, limit, offset)
	}
	return storage.GetGroupMessages(targetID, limit, offset)
}

func MarkMessagesRead(userID, targetID string) error {
	return nil
}

func GetUnreadMessages(userID string) ([]*model.Message, error) {
	return storage.GetUnreadMessages(userID)
}