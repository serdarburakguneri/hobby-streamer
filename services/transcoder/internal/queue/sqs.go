package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSConsumer struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSConsumer(ctx context.Context, queueURL string) (QueueConsumer, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &SQSConsumer{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}

func (c *SQSConsumer) Start(ctx context.Context, handle func(QueueMessage) error) {
	const maxConsecutiveFailures = 10
	consecutiveFailures := 0
	backoff := time.Second

	for {
		out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &c.queueURL,
			MaxNumberOfMessages: 5,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			consecutiveFailures++
			log.Printf("[ERROR] SQS receive error (attempt %d): %v", consecutiveFailures, err)
			if consecutiveFailures >= maxConsecutiveFailures {
				log.Printf("[FATAL] Too many consecutive SQS receive failures (%d). Exiting.", maxConsecutiveFailures)
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

		for _, m := range out.Messages {
			var msg QueueMessage
			if err := json.Unmarshal([]byte(*m.Body), &msg); err != nil {
				log.Printf("[ERROR] Failed to unmarshal SQS message: %v", err)
				// TODO: move to DLQ here
				continue
			}
			if err := handle(msg); err == nil {
				_, _ = c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &c.queueURL,
					ReceiptHandle: m.ReceiptHandle,
				})
			} else {
				log.Printf("[ERROR] Handler failed for message (type: %s): %v", msg.Type, err)
				// TODO: move to DLQ here
			}
		}
	}
}
