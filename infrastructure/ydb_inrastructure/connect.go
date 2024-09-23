package ydb_inrastructure

import (
	"context"

	"github.com/rs/zerolog"

	"tg-anonymizer/utils/settings"

	ydb2 "github.com/ydb-platform/gorm-driver"
	"gorm.io/gorm"

	"github.com/pkg/errors"
	yc "github.com/ydb-platform/ydb-go-yc"
)

func GetGormConnect(ctx context.Context) (db *gorm.DB, err error) {
	var dialector gorm.Dialector
	if settings.Settings.YDB.FromInside {
		dialector = ydb2.Open(
			settings.Settings.YDB.Url,
			ydb2.With(yc.WithInternalCA()),
			ydb2.With(yc.WithMetadataCredentials()),
		)
	} else {
		dialector = ydb2.Open(
			settings.Settings.YDB.Url,
			ydb2.With(yc.WithInternalCA()),
			ydb2.With(yc.WithServiceAccountKeyFileCredentials(".ydb_sa_keys.json")),
		)
	}

	db, err = gorm.Open(dialector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to ydb")
	}
	zerolog.Ctx(ctx).Info().Msg("connected.to.ydb")

	return db, nil
}
