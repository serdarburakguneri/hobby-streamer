package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	for {
		out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &c.queueURL,
			MaxNumberOfMessages: 5,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			log.Printf("receive error: %v", err)
			//handle error appropriately, retry instead of log
			continue
		}

		for _, m := range out.Messages {
			var msg QueueMessage
			if err := json.Unmarshal([]byte(*m.Body), &msg); err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}
			if err := handle(msg); err == nil {
				_, _ = c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &c.queueURL,
					ReceiptHandle: m.ReceiptHandle,
				})
			} else {
				log.Printf("handler error: %v", err)
			}
		}
	}
}