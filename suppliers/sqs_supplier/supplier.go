package sqs_supplier

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Supplier struct {
	client *sqs.Client
}

func NewSupplier(ctx context.Context, cfg *aws.Config) (*Supplier, error) {
	return &Supplier{sqs.NewFromConfig(*cfg)}, nil
}
