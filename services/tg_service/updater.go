package tg_service

import (
	"context"
	"net/http"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/must_utils"
)

func (r *Service) PollerRun(ctx context.Context) {
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

func (r *Service) ProcessWebhook(rw http.ResponseWriter, req *http.Request) {
	ctx := logger_utils.AddLoggerToCtx(req.Context())

	var wg sync.WaitGroup
	for update := range r.bot.ListenForWebhookRespReqFormat(rw, req) {
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
			ctx = logger_utils.WithStrContextLog(ctx, "from.chat", getMessageChatTitle(update))
		}
	}

	// Log only every 10th update
	if update.UpdateID%10 == 0 {
		zerolog.Ctx(ctx).Debug().Interface("update", update).Msg("tg.new.update")
	}

	if update.Message != nil && update.Message.Chat != nil {

		if update.Message.Chat.Type == "private" {
			wg.Add(1)
			go must_utils.DoOrLog(
				func(ctx context.Context) error {
					return r.handlePrivateMessage(ctx, wg, update)
				},
				"error.during.processing.private.message",
			)(ctx)
		}

		if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
			wg.Add(1)
			go must_utils.DoOrLog(
				func(ctx context.Context) error {
					return r.handleGroupMessage(ctx, wg, update)
				},
				"error.during.processing.group.message",
			)(ctx)
		}
	}

	return nil
}
