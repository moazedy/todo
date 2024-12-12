package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClient interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

type sqsClient struct {
	client *sqs.Client
}

func newSQSClient(client *sqs.Client) SQSClient {
	return &sqsClient{
		client: client,
	}
}

func (sc sqsClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	return sc.ReceiveMessage(ctx, params, optFns...)
}

func (sc sqsClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	return sc.client.SendMessage(ctx, params, optFns...)
}

func (sc sqsClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	return sc.client.DeleteMessage(ctx, params, optFns...)
}
