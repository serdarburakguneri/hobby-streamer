package queue

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type SQSConsumer struct {
	client   *sqs.Client
	queueURL string
	logger   *logger.Logger
}

func NewSQSConsumer(ctx context.Context, queueURL string) (QueueConsumer, error) {
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
	return &SQSConsumer{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
		logger:   logger.WithService("sqs-consumer"),
	}, nil
}

func (c *SQSConsumer) Start(ctx context.Context, handle func(QueueMessage) error) {
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
			var msg QueueMessage
			if err := json.Unmarshal([]byte(*m.Body), &msg); err != nil {
				log.WithError(err).Error("Failed to unmarshal SQS message")
				continue
			}

			log.Debug("Processing message", "message_type", msg.Type, "message_id", *m.MessageId)
			if err := handle(msg); err == nil {
				_, _ = c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &c.queueURL,
					ReceiptHandle: m.ReceiptHandle,
				})
				log.Debug("Message processed successfully", "message_type", msg.Type, "message_id", *m.MessageId)
			} else {
				log.WithError(err).Error("Handler failed for message", "message_type", msg.Type, "message_id", *m.MessageId)
			}
		}
	}
}
