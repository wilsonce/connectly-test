package service

import (
	"context"
	"github.com/wilsonce/connectly-test/dao"
	"github.com/wilsonce/connectly-test/model"
	"time"
)

type BotService struct {
}

func NewBotService() *BotService {
	return &BotService{}
}

func (b *BotService) Add(ctx context.Context, bot *model.WBot) error {
	now := time.Now()
	currentTime := now.Format("2006-01-02 15:04:05")
	bot.CreatedAt = currentTime
	bot.UpdatedAt = currentTime
	return dao.Q.WBot.WithContext(ctx).Create(bot)
}
