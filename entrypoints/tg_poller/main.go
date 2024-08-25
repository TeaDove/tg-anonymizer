package main

import (
	"tg-anonymizer/entrypoints"

	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/must_utils"
)

func main() {
	ctx := logger_utils.NewLoggedCtx()

	tgService, err := entrypoints.MakeTgService(ctx)
	if err != nil {
		must_utils.FancyPanic(ctx, errors.Wrap(err, "failed to make service"))
	}

	tgService.Run(ctx)
}
