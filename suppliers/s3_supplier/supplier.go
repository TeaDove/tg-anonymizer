package s3_supplier

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Supplier struct {
	client *s3.Client
}

func NewSupplier(ctx context.Context, cfg *aws.Config) (*Supplier, error) {
	return &Supplier{s3.NewFromConfig(*cfg)}, nil
}
