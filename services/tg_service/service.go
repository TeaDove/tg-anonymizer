package tg_service

import (
	"context"
	"fmt"

	"tg-anonymizer/repositories/user_chat_repository"

	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	bot                  *tgbotapi.BotAPI
	userToChatRepository *user_chat_repository.Repository
}

func NewService(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	userToChatRepository *user_chat_repository.Repository,
) (*Service, error) {
	return &Service{bot: bot, userToChatRepository: userToChatRepository}, nil
}

func (r *Service) reply(update *tgbotapi.Update, format string, a ...any) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(format, a...))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	_, err := r.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
