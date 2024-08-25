package entrypoints

import (
	"context"

	"tg-anonymizer/services/tg_service"
	"tg-anonymizer/utils/settings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func MakeTgService(ctx context.Context) (*tg_service.Service, error) {
	bot, err := tgbotapi.NewBotAPI(settings.Settings.TgToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run bot")
	}

	bot.Debug = false
	zerolog.Ctx(ctx).Info().Str("username", bot.Self.UserName).Msg("bot.client.created")

	tgService, err := tg_service.NewService(ctx, bot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make service")
	}

	return tgService, nil
}
