package sqs

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Consumer struct {
	client   *sqs.Client
	queueURL string
	logger   *logger.Logger
}

func NewConsumer(ctx context.Context, queueURL string) (*Consumer, error) {
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
	return &Consumer{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
		logger:   logger.WithService("sqs-consumer"),
	}, nil
}

func (c *Consumer) Start(ctx context.Context, handle func(Message) error) {
	const maxConsecutiveFailures = 10
	consecutiveFailures := 0
	backoff := time.Second

	log := c.logger.WithContext(ctx)
	log.Info("Starting SQS consumer", "queue_url", c.queueURL)

	for {
		out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &c.queueURL,
			MaxNumberOfMessages: 5,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			consecutiveFailures++
			log.WithError(err).Error("SQS receive error", "attempt", consecutiveFailures)
			if consecutiveFailures >= maxConsecutiveFailures {
				log.Error("Too many consecutive SQS receive failures, exiting", "max_failures", maxConsecutiveFailures)
				return
			}
			time.Sleep(backoff)
			if backoff < 30*time.Second {
				backoff *= 2
			}
			continue
		}
		consecutiveFailures = 0
		backoff = time.Second

		if len(out.Messages) > 0 {
			log.Debug("Received messages", "count", len(out.Messages))
		}

		for _, m := range out.Messages {
			log.Debug("Raw SQS message body", "body", *m.Body, "message_id", *m.MessageId)
			var msg Message
			if err := json.Unmarshal([]byte(*m.Body), &msg); err != nil {
				log.WithError(err).Error("Failed to unmarshal SQS message")
				continue
			}

			log.Debug("Processing message", "message_type", msg.Type, "message_id", *m.MessageId, "payload_length", len(msg.Payload))
			if err := handle(msg); err == nil {
				_, _ = c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &c.queueURL,
					ReceiptHandle: m.ReceiptHandle,
				})
				log.Debug("Message processed successfully", "message_type", msg.Type, "message_id", *m.MessageId)
			} else {
				log.WithError(err).Error("Handler failed for message", "message_type", msg.Type, "message_id", *m.MessageId)

				if isValidationError(err) {
					log.Info("Deleting message due to validation error (non-retryable)", "message_type", msg.Type, "message_id", *m.MessageId)
					_, _ = c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
						QueueUrl:      &c.queueURL,
						ReceiptHandle: m.ReceiptHandle,
					})
				}
			}
		}
	}
}

func isValidationError(err error) bool {
	if err == nil {
		return false
	}

	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.Type == apperrors.ErrorTypeValidation
	}

	return false
}
