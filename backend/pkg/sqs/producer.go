package sqs

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Producer struct {
	client   *sqs.Client
	queueURL string
	logger   *logger.Logger
}

func NewProducer(ctx context.Context, queueURL string) (*Producer, error) {
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint == "" {
		awsEndpoint = "http://localstack:4566"
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           awsEndpoint,
			SigningRegion: region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, err
	}
	return &Producer{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
		logger:   logger.WithService("sqs-producer"),
	}, nil
}

func (p *Producer) SendMessage(ctx context.Context, messageType string, payload interface{}) error {
	log := p.logger.WithContext(ctx)

	messageBody, err := json.Marshal(map[string]interface{}{
		"type":    messageType,
		"payload": payload,
	})
	if err != nil {
		log.WithError(err).Error("Failed to marshal message")
		return err
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    &p.queueURL,
		MessageBody: aws.String(string(messageBody)),
	}

	_, err = p.client.SendMessage(ctx, input)
	if err != nil {
		log.WithError(err).Error("Failed to send SQS message", "message_type", messageType)
		return err
	}

	log.Info("SQS message sent successfully", "message_type", messageType)
	return nil
}
