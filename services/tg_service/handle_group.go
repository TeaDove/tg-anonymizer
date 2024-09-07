package tg_service

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func (r *Service) handleGroupMessage(
	ctx context.Context,
	wg *sync.WaitGroup,
	update *tgbotapi.Update,
) error {
	defer wg.Done()

	err := r.userToChatRepository.PutCommonChat(ctx, update.Message.From.ID, update.Message.Chat)
	if err != nil {
		return errors.Wrap(err, "failed to put common chat")
	}

	return nil
}
