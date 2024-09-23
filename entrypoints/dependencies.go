package entrypoints

import (
	"context"

	"tg-anonymizer/infrastructure/aws_infrastructure"
	"tg-anonymizer/infrastructure/ydb_inrastructure"
	"tg-anonymizer/repositories/user_chat_repository"
	"tg-anonymizer/services/tg_service"
	"tg-anonymizer/suppliers/s3_supplier"
	"tg-anonymizer/suppliers/sqs_supplier"
	"tg-anonymizer/utils/settings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func MakeTgService(ctx context.Context) (*tg_service.Service, error) {
	bot, err := tgbotapi.NewBotAPI(settings.Settings.Tg.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run bot")
	}

	bot.Debug = false
	zerolog.Ctx(ctx).Info().Str("username", bot.Self.UserName).Msg("bot.client.created")

	ydbDriver, err := ydb_inrastructure.GetGormConnect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ydb driver")
	}

	userToChatRepository, err := user_chat_repository.NewRepository(ctx, ydbDriver)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init user_to_chat_repository")
	}

	awsConfig, err := aws_infrastructure.NewConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init aws_config")
	}

	s3Supplier, err := s3_supplier.NewSupplier(ctx, &awsConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init s3_supplier")
	}

	sqsSupplier, err := sqs_supplier.NewSupplier(ctx, &awsConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init sqs_supplier")
	}

	tgService, err := tg_service.NewService(ctx, bot, userToChatRepository, s3Supplier, sqsSupplier)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make service")
	}

	return tgService, nil
}
