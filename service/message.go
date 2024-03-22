package service

import (
	"context"
	"github.com/wilsonce/connectly-test/dao"
	"github.com/wilsonce/connectly-test/model"
	"time"
)

type MessageService struct {
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (m *MessageService) Add(ctx context.Context, message *model.WMessage) error {
	now := time.Now()
	currentTime := now.Format("2006-01-02 15:04:05")
	message.CreatedAt = currentTime
	return dao.Q.WMessage.WithContext(ctx).Create(message)
}
