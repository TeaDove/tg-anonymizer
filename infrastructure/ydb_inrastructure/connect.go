package ydb_inrastructure

import (
	"context"

	"tg-anonymizer/utils/settings"

	"github.com/rs/zerolog"
	ydb2 "github.com/ydb-platform/gorm-driver"
	"gorm.io/gorm"

	"github.com/pkg/errors"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	yc "github.com/ydb-platform/ydb-go-yc"
)

func GetConnect(ctx context.Context) (db *ydb.Driver, err error) {
	if settings.Settings.YdbFromInside {
		db, err = ydb.Open(ctx,
			settings.Settings.YdbUrl,
			yc.WithInternalCA(),
			yc.WithMetadataCredentials(),
		)
	} else {
		db, err = ydb.Open(ctx,
			settings.Settings.YdbUrl,
			yc.WithInternalCA(),
			yc.WithServiceAccountKeyFileCredentials(".ydb_sa_keys.json"),
		)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to ydb")
	}

	zerolog.Ctx(ctx).Info().Str("connect", db.String()).Msg("connected.to.ydb")

	return db, nil
}

func GetGormConnect(ctx context.Context) (db *gorm.DB, err error) {
	db, err = gorm.Open(ydb2.Open(
		settings.Settings.YdbUrl,
		ydb2.With(yc.WithInternalCA()),
		ydb2.With(yc.WithServiceAccountKeyFileCredentials(".ydb_sa_keys.json")),
	),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to ydb")
	}

	return db, nil
}
