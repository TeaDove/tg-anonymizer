package tg_service

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	bot *tgbotapi.BotAPI
}

func NewService(ctx context.Context, bot *tgbotapi.BotAPI) (*Service, error) {
	return &Service{bot: bot}, nil
}
