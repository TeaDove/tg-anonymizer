package sqs_supplier

import (
	"context"

	"tg-anonymizer/utils/settings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pkg/errors"
)

func (r *Supplier) SendMessage(ctx context.Context, message string) error {
	_, err := r.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(settings.Settings.SQS.Url),
	})
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (r *Supplier) ReceiveAndDeleteMessage(ctx context.Context) ([]string, error) {
	received, err := r.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(settings.Settings.SQS.Url),
		MaxNumberOfMessages: 5,
		WaitTimeSeconds:     2,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to receive message")
	}

	messages := make([]string, 0, len(received.Messages))
	for _, message := range received.Messages {
		messages = append(messages, *message.Body)

		_, err = r.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(settings.Settings.SQS.Url),
			ReceiptHandle: message.ReceiptHandle,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to delete message")
		}
	}

	return messages, nil
}
