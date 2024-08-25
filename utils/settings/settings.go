package settings

import (
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/must_utils"
	"github.com/teadove/teasutils/utils/settings_utils"
)

type baseSettings struct {
	TgToken string `env:"tg_token,required"`
}

func init() {
	ctx := logger_utils.NewLoggedCtx()

	var err error
	Settings, err = settings_utils.InitSetting[baseSettings](ctx)
	if err != nil {
		must_utils.FancyPanic(ctx, errors.Wrap(err, "settings init failed"))
	}
}

var Settings baseSettings
