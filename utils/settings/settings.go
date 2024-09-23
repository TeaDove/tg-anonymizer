package settings

import (
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/must_utils"
	"github.com/teadove/teasutils/utils/settings_utils"
)

type yc struct {
	AccessKeyId     string `env:"access_key_id,required"     json:"accessKeyId"`
	SecretAccessKey string `env:"secret_access_key,required" json:"secretAccessKey"`

	PartitionId string `env:"partition_id" json:"partitionId" envDefault:"yc"`
	Region      string `env:"region"       json:"region"      envDefault:"ru-central1"`
}

type s3 struct {
	Url    string `env:"url"    json:"url"    envDefault:"https://storage.yandexcloud.net"`
	Bucket string `env:"bucket" json:"bucket" envDefault:"prod-tg-anonymizer-storage"`
}

type sqs struct {
	Url string `env:"queue" json:"queue" envDefault:"https://message-queue.api.cloud.yandex.net/b1g15gt835j53bc2hir0/dj6000000022nftk037s/prod-tg-anonymizer-queue"`
}

type ydb struct {
	FromInside bool   `env:"ydb_from_inside" json:"ydbFromInside" endDefault:"true"`
	Url        string `env:"ydb_url"         json:"ydbUrl"                          envDefault:"grpcs://ydb.serverless.yandexcloud.net:2135/ru-central1/b1g15gt835j53bc2hir0/etnkid49o4gf60c6o88j"`
}

type tg struct {
	Token string `env:"token,required" json:"token"`
}

type baseSettings struct {
	Tg  tg  `env:"tg"  json:"tg"  envPrefix:"tg__"`
	YDB ydb `env:"ydb" json:"ydb" envPrefix:"ydb__"`
	YC  yc  `env:"yc"  json:"yc"  envPrefix:"yc__"`
	SQS sqs `env:"sqs" json:"sqs" envPrefix:"sqs__"`
	S3  s3  `env:"s3"  json:"s3"  envPrefix:"s3__"`
}

func init() {
	ctx := logger_utils.NewLoggedCtx()

	var err error
	Settings, err = settings_utils.InitSetting[baseSettings](
		ctx,
		"tg.token",
		"yc.accessKeyId",
		"yc.secretAccessKey",
	)
	if err != nil {
		must_utils.FancyPanic(ctx, errors.Wrap(err, "settings init failed"))
	}
}

var Settings baseSettings
