package tg_service

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/must_utils"
)

func (r *Service) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := r.bot.GetUpdatesChan(u)

	var wg sync.WaitGroup

	for update := range updates {
		wg.Add(1)
		go must_utils.DoOrLog(
			func(ctx context.Context) error {
				return r.processUpdate(ctx, &wg, &update)
			},
			"error.during.update.process",
		)(ctx)
	}

	wg.Wait()
}

func (r *Service) processUpdate(
	ctx context.Context,
	wg *sync.WaitGroup,
	update *tgbotapi.Update,
) error {
	defer wg.Done()

	if update.Message != nil {
		if update.Message.From != nil {
			ctx = logger_utils.WithStrContextLog(ctx, "from.username", update.Message.From.UserName)
		}

		if update.Message.Chat != nil {
			ctx = logger_utils.WithStrContextLog(ctx, "from.chat", update.Message.Chat.Title)
		}
	}

	zerolog.Ctx(ctx).Debug().Interface("update", update).Msg("tg.new.update")

	if update.Message != nil && update.Message.Chat != nil &&
		update.Message.Chat.Type == "private" {
		wg.Add(1)
		go must_utils.DoOrLog(
			func(ctx context.Context) error {
				return r.handlePrivateMessage(ctx, wg, update)
			},
			"error.during.processing.private.message",
		)(ctx)
	}

	return nil
}

func (r *Service) handlePrivateMessage(
	ctx context.Context,
	wg *sync.WaitGroup,
	update *tgbotapi.Update,
) error {
	defer wg.Done()

	zerolog.Ctx(ctx).Debug().Interface("update", update.Message).Msg("tg.handlePrivateMessage")

	msg := tgbotapi.NewMessage(-1001178533048, update.Message.Text)
	_, err := r.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
