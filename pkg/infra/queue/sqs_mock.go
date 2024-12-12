package queue

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"github.com/moazedy/todo/internal/domain/model"
)

type mockSQS struct {
	client map[string]map[string]model.Message
}

func newMockSQS() SQSClient {
  ms := mockSQS{}
  ms.client = make(map[string]map[string]model.Message)

  return &ms
}

func (mq *mockSQS) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	theQueue, exist := mq.client[*params.QueueUrl]

	id := uuid.New().String()
	if !exist {
		mq.client[*params.QueueUrl] = map[string]model.Message{
			id: {ID: id, Body: *params.MessageBody},
		}
	} else {
		theQueue[id] = model.Message{ID: id, Body: *params.MessageBody}
		mq.client[*params.QueueUrl] = theQueue
	}

	return &sqs.SendMessageOutput{MessageId: &id}, nil
}

func (mq *mockSQS) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	out := make([]types.Message, 0)

	theQueue, exist := mq.client[*params.QueueUrl]
	if !exist {
		return nil, errors.New("queue url not found")
	}
	var idx int32
	for _, message := range theQueue {
		if idx >= params.MaxNumberOfMessages {
			break
		}

		out = append(out, types.Message{
			MessageId: &message.ID,
			Body:      &message.Body,
		})
	}

	return &sqs.ReceiveMessageOutput{Messages: out}, nil
}

func (mq *mockSQS) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	theQueue, exist := mq.client[*params.QueueUrl]
	if !exist {
		return nil, errors.New("queue url not found")
	}

	_, messageExist := theQueue[*params.ReceiptHandle]
	if !messageExist {
		return nil, errors.New("message not found")
	}

	delete(mq.client[*params.QueueUrl], *params.ReceiptHandle)

	return nil, nil
}
