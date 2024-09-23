package s3_supplier

import (
	"context"
	"fmt"
	"io"
	"time"

	"tg-anonymizer/utils/settings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

type PutObjectInput struct {
	Key         string
	Body        io.Reader
	Expires     time.Time
	ContentType string
}

func (r *Supplier) PutObject(ctx context.Context, input *PutObjectInput) (string, error) {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(settings.Settings.S3.Bucket),
		Key:         aws.String(input.Key),
		Body:        input.Body,
		Expires:     aws.Time(input.Expires),
		ContentType: aws.String(input.ContentType),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to put object")
	}

	return fmt.Sprintf(
		"%s/%s/%s",
		settings.Settings.S3.Url,
		settings.Settings.S3.Bucket,
		input.Key,
	), nil
}
