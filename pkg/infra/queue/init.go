package queue

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	myconfig "github.com/moazedy/todo/pkg/infra/config"
)

func NewSQSClient(cfg myconfig.SQS) *sqs.Client {
	if cfg.IsMock {
		return nil
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
	)
	if err != nil {
		log.Fatalf("error while creating sqs client: %s \n", err.Error())
	}

	return sqs.NewFromConfig(awsCfg)
}
