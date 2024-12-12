package repoimplement

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/moazedy/todo/internal/adapter/driven/db/repoimplement/cast"
	"github.com/moazedy/todo/internal/domain/model"
	"github.com/moazedy/todo/internal/port/driven/db/repository"
	"github.com/moazedy/todo/pkg/infra/queue"
)

type queueRepo struct {
	sqsClientFactory queue.SQSClientFactory
	queueUrl         string
}

func NewQueue(sqsClientFactory queue.SQSClientFactory, queueUrl string) repository.Queue {
	return queueRepo{
		sqsClientFactory: sqsClientFactory,
		queueUrl:         queueUrl,
	}
}

func (q queueRepo) SendMessage(ctx context.Context, messageBody string) error {
	sqsClient := q.sqsClientFactory.NewSQSClient()

	_, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &q.queueUrl,
		MessageBody: &messageBody,
	})

	return err
}

func (q queueRepo) ReceiveMessage(ctx context.Context, maxNumberOfMessages int32) ([]model.Message, error) {
	sqsClient := q.sqsClientFactory.NewSQSClient()

	messages, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &q.queueUrl,
		MaxNumberOfMessages: maxNumberOfMessages,
	})
	if err != nil {
		return nil, err
	}

	return cast.ToSliceOfModelMessages(messages.Messages), nil
}

func (q queueRepo) DeleteMessage(ctx context.Context, messageId string) error {
	sqsClient := q.sqsClientFactory.NewSQSClient()

	_, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &q.queueUrl,
		ReceiptHandle: &messageId,
	})

	return err
}
