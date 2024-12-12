package queue

import "github.com/aws/aws-sdk-go-v2/service/sqs"

type SQSClientFactory interface {
	NewSQSClient() SQSClient
}

type sqsClientFactory struct {
	isMock bool
	client *sqs.Client
}

func NewSQSClientFactory(isMock bool, client *sqs.Client) SQSClientFactory {
	return sqsClientFactory{
		isMock: isMock,
		client: client,
	}
}

func (scf sqsClientFactory) NewSQSClient() SQSClient {
	if scf.isMock {
		return newMockSQS()
	} else {
		return newSQSClient(scf.client)
	}
}
