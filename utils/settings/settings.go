package settings

import (
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/must_utils"
	"github.com/teadove/teasutils/utils/settings_utils"
)

type baseSettings struct {
	TgToken string `env:"tg_token,required" json:"tgToken"`

	YdbFromInside bool   `env:"ydb_from_inside" json:"ydbFromInside" endDefault:"true"`
	YdbUrl        string `env:"ydb_url"         json:"ydbUrl"                          envDefault:"grpcs://ydb.serverless.yandexcloud.net:2135/ru-central1/b1g15gt835j53bc2hir0/etnkid49o4gf60c6o88j"`
}

func init() {
	ctx := logger_utils.NewLoggedCtx()

	var err error
	Settings, err = settings_utils.InitSetting[baseSettings](ctx, "tgToken")
	if err != nil {
		must_utils.FancyPanic(ctx, errors.Wrap(err, "settings init failed"))
	}
}

var Settings baseSettings
