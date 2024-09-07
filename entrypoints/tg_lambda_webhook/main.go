package main

import (
	"net/http"

	"tg-anonymizer/entrypoints"
	"tg-anonymizer/services/tg_service"

	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/must_utils"
)

var service *tg_service.Service

func init() {
	ctx := logger_utils.NewLoggedCtx()
	var err error

	service, err = entrypoints.MakeTgService(ctx)
	if err != nil {
		must_utils.FancyPanic(ctx, errors.Wrap(err, "failed to make tg service"))
	}
}

func Handler(rw http.ResponseWriter, req *http.Request) {
	ctx := logger_utils.AddLoggerToCtx(req.Context())

	reqId, ok := ctx.Value("lambdaRuntimeRequestID").(string)
	if ok {
		ctx = logger_utils.WithStrContextLog(ctx, "request_id", reqId)
	}

	service.ProcessWebhook(ctx, rw, req)
}
