package aws_infrastructure

import (
	"context"

	"tg-anonymizer/utils/settings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pkg/errors"
)

func NewConfig(ctx context.Context) (aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...any) (aws.Endpoint, error) {
			if service == s3.ServiceID && region == settings.Settings.YC.Region {
				return aws.Endpoint{
					PartitionID:   settings.Settings.YC.PartitionId,
					URL:           settings.Settings.S3.Url,
					SigningRegion: settings.Settings.YC.Region,
				}, nil
			}
			if service == sqs.ServiceID && region == settings.Settings.YC.Region {
				return aws.Endpoint{
					PartitionID:   settings.Settings.YC.PartitionId,
					URL:           settings.Settings.SQS.Url,
					SigningRegion: settings.Settings.YC.Region,
				}, nil
			}
			return aws.Endpoint{}, errors.Errorf(
				"unknown endpoint requested, service=%s, region=%s",
				service,
				region,
			)
		},
	)

	options := []func(*config.LoadOptions) error{
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRegion(settings.Settings.YC.Region),
	}
	if !settings.Settings.YC.FromEnv {
		if settings.Settings.YC.AccessKeyId == "" || settings.Settings.YC.SecretAccessKey == "" {
			return aws.Config{}, errors.New("missing YC access key or secret access key")
		}

		options = append(options, config.WithCredentialsProvider(
			&credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     settings.Settings.YC.AccessKeyId,
					SecretAccessKey: settings.Settings.YC.SecretAccessKey,
				},
			},
		))
	}

	cfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return aws.Config{}, errors.Wrap(err, "failed to load default config")
	}

	return cfg, nil
}
